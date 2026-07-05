// English word-of-a-card pronunciation via the browser's SpeechSynthesis API.

let voices = []
let unlocked = false

function loadVoices() {
  voices = window.speechSynthesis?.getVoices() || []
}

if (typeof window !== 'undefined' && window.speechSynthesis) {
  loadVoices()
  window.speechSynthesis.onvoiceschanged = loadVoices
}

function pickEnglishVoice() {
  return (
    voices.find(v => v.lang === 'en-US') ||
    voices.find(v => v.lang?.startsWith('en')) ||
    null
  )
}

// Telegram's in-app WebView (unlike a regular browser tab on the same
// engine) only lets speechSynthesis produce audio when speak() runs
// synchronously inside a trusted user gesture. Games call speakWord() from
// onMounted or after an API call finishes, which is always one tick past
// the gesture that opened the level — so the very first utterance of a
// session gets silently dropped. Priming the engine with a near-silent
// utterance on the first tap anywhere in the app unlocks it for every
// speak() call for the rest of the session, gesture or not.
export function unlockSpeech() {
  if (unlocked || typeof window === 'undefined' || !window.speechSynthesis) return
  unlocked = true
  const utterance = new SpeechSynthesisUtterance(' ')
  utterance.volume = 0
  window.speechSynthesis.speak(utterance)
}

if (typeof window !== 'undefined') {
  window.addEventListener('pointerdown', unlockSpeech, { once: true, capture: true })
}

export function speakWord(word) {
  if (!word || typeof window === 'undefined' || !window.speechSynthesis) return
  const synth = window.speechSynthesis
  synth.cancel()
  const utterance = new SpeechSynthesisUtterance(word)
  utterance.lang = 'en-US'
  utterance.rate = 0.85
  const voice = pickEnglishVoice()
  if (voice) utterance.voice = voice
  // Some WebViews (incl. Telegram's) silently drop a speak() call made in
  // the same tick as cancel() — queue it for the next tick instead.
  setTimeout(() => synth.speak(utterance), 0)
}

export function stopSpeaking() {
  window.speechSynthesis?.cancel()
}
