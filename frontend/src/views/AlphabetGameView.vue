<template>
  <div class="max-w-md md:max-w-2xl mx-auto px-4 pb-24 md:pb-10 pt-5 md:pt-8 min-h-screen">
    <!-- Header -->
    <header class="flex items-center gap-3 mb-4">
      <Button variant="ghost" size="icon" class="shrink-0" @click="backToLevels">
        <ArrowLeftIcon class="h-5 w-5" />
      </Button>
      <div class="min-w-0">
        <h1 class="text-lg font-extrabold leading-none truncate">
          {{ levelTitle }}
        </h1>
      </div>
    </header>

    <!-- Stats bar -->
    <div class="flex gap-3 mb-4">
      <div class="flex-1 rounded-xl bg-card border border-border py-2 text-center">
        <p class="text-[10px] uppercase tracking-wide text-muted-foreground">{{ $t('alphabet.moves') }}</p>
        <p class="font-bold text-primary">{{ moves }}</p>
      </div>
    </div>

    <p class="text-center text-muted-foreground text-xs mb-2">{{ $t('alphabet.hint') }}</p>
    <p class="text-center font-extrabold text-base tracking-wide mb-4 text-primary">{{ letterRange }}</p>

    <!-- Puzzle grid -->
    <div class="grid grid-cols-4 gap-2 aspect-square">
      <button
        v-for="(letter, i) in board"
        :key="i"
        class="rounded-2xl flex items-center justify-center text-2xl font-extrabold select-none transition-all duration-150 border-2"
        :class="tileClass(letter, i)"
        :disabled="won || (!letter && !isValidTarget(i))"
        @click="tap(i)"
      >
        {{ letter }}
      </button>
    </div>

    <!-- Win overlay -->
    <Transition name="fade">
      <div v-if="won" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm px-6">
        <Card class="w-full max-w-xs p-6 text-center animate-pop">
          <div class="text-6xl mb-3">🎉</div>
          <h2 class="text-xl font-extrabold mb-1">{{ $t('alphabet.win') }}</h2>
          <p class="text-muted-foreground text-sm mb-4">{{ $t('alphabet.winSub') }}</p>

          <div class="space-y-2">
            <Button v-if="hasNextLevel" variant="gradient" class="w-full" @click="goNext">
              {{ $t('alphabet.nextLevel') }} <ArrowRightIcon class="h-4 w-4" />
            </Button>
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
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { hapticFeedback } from '../services/telegram'
import { markAlphabetLevelCompleted } from '../services/alphabetProgress'
import { speakWord } from '../services/tts'
import Button from '../components/ui/button.vue'
import Card from '../components/ui/card.vue'
import { ArrowLeft as ArrowLeftIcon, ArrowRight as ArrowRightIcon } from 'lucide-vue-next'

const props = defineProps({ level: { type: [String, Number], required: true } })
const router = useRouter()
const { t } = useI18n()

const COLS = 4
const ROWS = 4
const LETTERS_PER_LEVEL = 13

const levelNum = computed(() => Number(props.level))
const hasNextLevel = computed(() => levelNum.value === 1)
const levelTitle = computed(() => t(levelNum.value === 1 ? 'alphabet.level1Title' : 'alphabet.level2Title'))
const letterRange = computed(() => letters.value.join(' '))

const letters = computed(() => {
  const startCode = 65 + (levelNum.value - 1) * LETTERS_PER_LEVEL // 'A' = 65
  return Array.from({ length: LETTERS_PER_LEVEL }, (_, i) => String.fromCharCode(startCode + i))
})

const solved = computed(() => {
  const arr = [...letters.value]
  while (arr.length < COLS * ROWS) arr.push(null)
  return arr
})

const board = ref([])
const moves = ref(0)
const won = ref(false)
const selected = ref(null)

function neighborsOf(i) {
  const row = Math.floor(i / COLS)
  const col = i % COLS
  const result = []
  if (row > 0) result.push(i - COLS)
  if (row < ROWS - 1) result.push(i + COLS)
  if (col > 0) result.push(i - 1)
  if (col < COLS - 1) result.push(i + 1)
  return result
}

function shuffle() {
  const arr = [...solved.value]
  for (let n = 0; n < 200; n++) {
    const emptyIndices = arr.reduce((acc, v, idx) => (v === null ? [...acc, idx] : acc), [])
    const empty = emptyIndices[Math.floor(Math.random() * emptyIndices.length)]
    const candidates = neighborsOf(empty)
    const from = candidates[Math.floor(Math.random() * candidates.length)]
    ;[arr[empty], arr[from]] = [arr[from], arr[empty]]
  }
  board.value = arr
}

function emptyNeighborsOf(i) {
  return neighborsOf(i).filter(n => board.value[n] === null)
}

function isValidTarget(i) {
  return selected.value !== null && board.value[i] === null && neighborsOf(selected.value).includes(i)
}

function tileClass(letter, i) {
  if (!letter) return isValidTarget(i) ? 'bg-primary/10 border-primary/40 border-dashed animate-pulse' : 'bg-transparent border-transparent'
  if (won.value) return 'bg-green-50 border-green-400 text-green-700'
  if (i === selected.value) return 'bg-gradient-to-br from-primary to-emerald-600 border-yellow-400 ring-4 ring-yellow-400/50 text-white shadow'
  return 'bg-gradient-to-br from-primary to-emerald-600 border-primary/40 text-white shadow active:scale-95'
}

function moveInto(from, to) {
  hapticFeedback('light')
  ;[board.value[from], board.value[to]] = [board.value[to], board.value[from]]
  selected.value = null
  moves.value++

  if (board.value.every((v, idx) => v === solved.value[idx])) {
    finish()
  }
}

function tap(i) {
  if (won.value) return

  // A tile is selected and this cell is one of its empty neighbors — complete the move.
  if (isValidTarget(i)) {
    moveInto(selected.value, i)
    return
  }

  if (!board.value[i]) {
    selected.value = null
    return
  }

  speakWord(board.value[i].toLowerCase())

  // Tapping the already-selected tile again deselects it.
  if (i === selected.value) {
    selected.value = null
    return
  }

  const empties = emptyNeighborsOf(i)
  if (empties.length === 0) {
    selected.value = null
    return
  }
  if (empties.length === 1) {
    moveInto(i, empties[0])
    return
  }
  // Multiple empty neighbors — let the player pick which one to slide into.
  selected.value = i
}

function finish() {
  won.value = true
  hapticFeedback('heavy')
  markAlphabetLevelCompleted(levelNum.value)
}

function restart() {
  moves.value = 0
  won.value = false
  selected.value = null
  shuffle()
}

function backToLevels() {
  router.push('/alphabet')
}

function goNext() {
  if (hasNextLevel.value) {
    router.push('/alphabet/2')
  } else {
    backToLevels()
  }
}

watch(() => props.level, restart)
onMounted(restart)
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
