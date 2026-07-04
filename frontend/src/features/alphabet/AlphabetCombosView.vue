<template>
  <div class="max-w-md md:max-w-2xl mx-auto px-4 pb-24 md:pb-10 pt-5 md:pt-8 min-h-screen">
    <!-- Header -->
    <header class="flex items-center gap-3 mb-4">
      <Button variant="ghost" size="icon" class="shrink-0" @click="backToLevels">
        <ArrowLeftIcon class="h-5 w-5" />
      </Button>
      <div class="min-w-0">
        <p class="text-xs text-muted-foreground leading-none mb-1">{{ $t('alphabet.combosLabel') }}</p>
        <h1 class="text-lg font-extrabold leading-none truncate">
          {{ groupCombos.join(' · ') }}
        </h1>
      </div>
    </header>

    <p class="text-center text-muted-foreground text-xs mb-1">{{ $t('alphabet.comboHint') }}</p>
    <p class="text-center text-xs font-semibold text-primary mb-4">
      {{ $t('wordBuild.wordOf', { current: currentIndex + 1, total: words.length }) }}
    </p>

    <!-- Word card -->
    <div class="rounded-3xl border border-border bg-card p-6 text-center mb-6">
      <div class="text-6xl mb-3">{{ currentItem.emoji }}</div>
      <button
        class="inline-flex items-center gap-2 text-2xl font-extrabold uppercase tracking-wide text-foreground active:scale-95 transition-transform"
        @click="replay"
      >
        {{ currentItem.word_en }}
        <Volume2Icon class="h-5 w-5 text-primary" />
      </button>

      <p v-if="revealed" class="mt-3 text-lg font-bold text-primary animate-pop">
        {{ currentItem.word_ru }}
      </p>
    </div>

    <!-- Combo options -->
    <div class="grid grid-cols-2 sm:grid-cols-3 gap-3" :class="{ 'animate-shake': wrongShake }">
      <button
        v-for="combo in options"
        :key="combo"
        class="h-14 rounded-2xl flex items-center justify-center text-xl font-extrabold uppercase select-none transition-all duration-150 border-2"
        :class="comboClass(combo)"
        :disabled="revealed || wrongOptions.includes(combo)"
        @click="tapCombo(combo)"
      >
        {{ combo }}
      </button>
    </div>

    <div v-if="revealed" class="flex justify-center mt-6">
      <Button variant="gradient" @click="next">
        {{ $t('wordBuild.next') }} <ArrowRightIcon class="h-4 w-4" />
      </Button>
    </div>

    <!-- Win overlay -->
    <Transition name="fade">
      <div v-if="won" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm px-6">
        <Card class="w-full max-w-xs p-6 text-center animate-pop">
          <div class="text-6xl mb-3">🎉</div>
          <h2 class="text-xl font-extrabold mb-1">{{ $t('alphabet.comboWin') }}</h2>
          <p class="text-muted-foreground text-sm mb-4">{{ $t('alphabet.comboWinSub') }}</p>

          <div class="space-y-2">
            <Button variant="secondary" class="w-full" @click="restart">
              {{ $t('alphabet.playAgain') }}
            </Button>
            <Button variant="ghost" class="w-full" @click="backToLevels">
              {{ $t('alphabet.backToLevels') }}
            </Button>
          </div>
        </Card>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { hapticFeedback } from '../../services/telegram'
import { speakWord, stopSpeaking } from '../../services/tts'
import { markAlphabetLevelCompleted } from './progress'
import { phonicsWords, COMBO_GROUPS } from './data/phonicsWords'
import Button from '../../components/ui/button.vue'
import Card from '../../components/ui/card.vue'
import { ArrowLeft as ArrowLeftIcon, ArrowRight as ArrowRightIcon, Volume2 as Volume2Icon } from 'lucide-vue-next'

const props = defineProps({ group: { type: String, required: true } })
const router = useRouter()

const ROUND_SIZE = 10
// The alphabet section's first 4 levels are ids 1-4; each combo group takes
// the next id in COMBO_GROUPS order (5, 6, 7, ...).
const levelId = computed(() => 5 + COMBO_GROUPS.findIndex(g => g.id === props.group))
const groupCombos = computed(() => COMBO_GROUPS.find(g => g.id === props.group)?.combos || [])
const groupWords = computed(() => phonicsWords.filter(w => groupCombos.value.includes(w.combo)))

function shuffle(arr) {
  const a = [...arr]
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[a[i], a[j]] = [a[j], a[i]]
  }
  return a
}

function pickRound() {
  return shuffle(groupWords.value).slice(0, ROUND_SIZE)
}

const words = ref(pickRound())
const currentIndex = ref(0)
const mistakes = ref(0)
const revealed = ref(false)
const wrongOptions = ref([])
const wrongShake = ref(false)
const options = ref([])
const won = ref(false)

const currentItem = computed(() => words.value[currentIndex.value])

function buildOptions(item) {
  // Only offer combos from this word's own group — and never one that's
  // also literally present in the word (e.g. "quick" contains both QU and CK).
  const candidates = groupCombos.value.filter(
    c => c === item.combo || !item.word_en.toLowerCase().includes(c.toLowerCase())
  )
  return shuffle(candidates)
}

function setupWord() {
  options.value = buildOptions(currentItem.value)
  revealed.value = false
  wrongOptions.value = []
  speakWord(currentItem.value.word_en)
}

function replay() {
  speakWord(currentItem.value.word_en)
}

function comboClass(combo) {
  if (revealed.value && combo === currentItem.value.combo) {
    return 'bg-green-50 border-green-400 text-green-700'
  }
  if (wrongOptions.value.includes(combo)) {
    return 'bg-rose-50 border-rose-400 text-rose-400'
  }
  return 'bg-gradient-to-br from-primary to-emerald-600 border-primary/40 text-white shadow active:scale-95'
}

function tapCombo(combo) {
  if (revealed.value || wrongOptions.value.includes(combo)) return

  if (combo === currentItem.value.combo) {
    revealed.value = true
    hapticFeedback('medium')
    speakWord(currentItem.value.word_en)
  } else {
    mistakes.value++
    wrongOptions.value.push(combo)
    hapticFeedback('rigid')
    wrongShake.value = true
    setTimeout(() => { wrongShake.value = false }, 400)
  }
}

function next() {
  if (currentIndex.value + 1 < words.value.length) {
    currentIndex.value++
    setupWord()
  } else {
    finish()
  }
}

function finish() {
  won.value = true
  hapticFeedback('heavy')
  markAlphabetLevelCompleted(levelId.value)
}

function restart() {
  won.value = false
  words.value = pickRound()
  currentIndex.value = 0
  mistakes.value = 0
  setupWord()
}

function backToLevels() {
  router.push('/alphabet')
}

// Reload when navigating between groups without unmounting.
watch(() => props.group, restart)

onMounted(setupWord)
onUnmounted(stopSpeaking)
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

@keyframes shake {
  10%, 90% { transform: translateX(-2px); }
  20%, 80% { transform: translateX(4px); }
  30%, 50%, 70% { transform: translateX(-8px); }
  40%, 60% { transform: translateX(8px); }
}
.animate-shake { animation: shake 0.4s; }
</style>
