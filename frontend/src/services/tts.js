// English word/sentence pronunciation via pre-generated audio files
// (frontend/public/audio/<slug>.mp3, one per vocab string in the app — see
// gen_audio.py). Telegram's mobile in-app WebView (Android/iOS) doesn't
// implement window.speechSynthesis at all, so relying on it left every game
// silent there; real audio files sidestep that entirely. speechSynthesis is
// kept only as a fallback for words that don't have a pre-generated file
// (e.g. added directly to the DB without regenerating audio).

const AUDIO_BASE = '/audio/'
const COMBINING_MARKS = new RegExp('[̀-ͯ]', 'g')

let voices = []
let unlocked = false
let currentAudio = null

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

// Must produce the exact same filename as gen_audio.py's slugify().
function slugify(text) {
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(COMBINING_MARKS, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

// Mobile WebViews only allow an <audio> element to play programmatically
// (outside a click handler) after one has already played following a real
// user gesture in the same session. Priming with a near-silent clip on the
// first tap anywhere unlocks playback for every speakWord() call after.
export function unlockSpeech() {
  if (unlocked || typeof window === 'undefined') return
  unlocked = true
  const audio = new Audio(`${AUDIO_BASE}_unlock.mp3`)
  audio.play().catch(() => {})
}

if (typeof window !== 'undefined') {
  window.addEventListener('pointerdown', unlockSpeech, { once: true, capture: true })
}

function speakWithSynthesis(word) {
  if (typeof window === 'undefined' || !window.speechSynthesis) return
  const synth = window.speechSynthesis
  synth.cancel()
  const utterance = new SpeechSynthesisUtterance(word)
  utterance.lang = 'en-US'
  utterance.rate = 0.85
  const voice = pickEnglishVoice()
  if (voice) utterance.voice = voice
  setTimeout(() => synth.speak(utterance), 0)
}

export function speakWord(word) {
  if (!word || typeof window === 'undefined') return
  stopSpeaking()
  const audio = new Audio(`${AUDIO_BASE}${slugify(word)}.mp3`)
  currentAudio = audio
  audio.addEventListener('error', () => speakWithSynthesis(word), { once: true })
  audio.play().catch(() => speakWithSynthesis(word))
}

export function stopSpeaking() {
  currentAudio?.pause()
  window.speechSynthesis?.cancel()
}
