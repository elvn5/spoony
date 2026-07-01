<template>
  <div class="max-w-md md:max-w-3xl mx-auto px-4 pb-24 md:pb-10 pt-5 md:pt-8 min-h-screen">
    <!-- Header -->
    <header class="flex items-center gap-3 mb-4">
      <Button variant="ghost" size="icon" class="shrink-0" @click="backToMap">
        <ArrowLeftIcon class="h-5 w-5" />
      </Button>
      <div class="min-w-0">
        <h1 class="text-lg font-extrabold leading-none truncate">
          <span class="mr-1">{{ level?.emoji }}</span>{{ level?.city || $t('game.title') }}
        </h1>
        <p class="text-muted-foreground text-xs mt-0.5">
          {{ level?.title_ru }} · {{ items.length }} {{ $t('trainer.words') }}
        </p>
      </div>
    </header>

    <!-- Stats bar -->
    <div class="flex gap-3 mb-4">
      <div class="flex-1 rounded-xl bg-card border border-border py-2 text-center">
        <p class="text-[10px] uppercase tracking-wide text-muted-foreground">{{ $t('game.found') }}</p>
        <p class="font-bold text-primary">{{ matchedPairs }} / {{ items.length }}</p>
      </div>
      <div class="flex-1 rounded-xl bg-card border border-border py-2 text-center">
        <p class="text-[10px] uppercase tracking-wide text-muted-foreground">{{ $t('game.moves') }}</p>
        <p class="font-bold">{{ moves }}</p>
      </div>
    </div>

    <p class="text-center text-muted-foreground text-xs mb-4">{{ $t('game.hint') }}</p>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-6 gap-3">
      <div v-for="n in 12" :key="n" class="aspect-square rounded-2xl bg-muted animate-pulse" />
    </div>

    <!-- Card grid -->
    <div v-else class="grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-6 gap-3">
      <button
        v-for="card in cards"
        :key="card.uid"
        class="aspect-square rounded-2xl flex items-center justify-center text-center select-none transition-all duration-200 border-2"
        :class="cardClass(card)"
        :disabled="isFaceUp(card) || locking"
        @click="flip(card)"
      >
        <template v-if="isFaceUp(card)">
          <span v-if="card.type === 'image'" class="text-4xl animate-pop">{{ card.content }}</span>
          <span v-else class="px-1 font-extrabold text-base leading-tight animate-pop break-words">{{ card.content }}</span>
        </template>
        <span v-else class="text-2xl opacity-90">🥄</span>
      </button>
    </div>

    <!-- Win overlay -->
    <Transition name="fade">
      <div v-if="won" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm px-6">
        <Card class="w-full max-w-xs p-6 text-center animate-pop">
          <div class="text-6xl mb-3">🎉</div>
          <h2 class="text-xl font-extrabold mb-1">{{ $t('game.win') }}</h2>
          <p class="text-muted-foreground text-sm mb-4">{{ $t('game.winSub') }}</p>

          <div class="flex justify-center gap-1 mb-5">
            <StarIcon
              v-for="s in 3" :key="s"
              class="h-9 w-9 transition-transform"
              :class="s <= stars ? 'text-yellow-400 scale-110' : 'text-muted-foreground/25'"
              :fill="s <= stars ? 'currentColor' : 'none'"
            />
          </div>

          <div class="space-y-2">
            <Button v-if="nextLevel" variant="gradient" class="w-full" @click="goNext">
              {{ $t('game.nextCity') }} <ArrowRightIcon class="h-4 w-4" />
            </Button>
            <Button variant="secondary" class="w-full" @click="restart">
              {{ $t('game.playAgain') }}
            </Button>
            <Button variant="ghost" class="w-full" @click="backToMap">
              {{ $t('game.backToMap') }}
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
import { learnApi } from '../services/api'
import { hapticFeedback } from '../services/telegram'
import { speakWord, stopSpeaking } from '../services/tts'
import Button from '../components/ui/button.vue'
import Card from '../components/ui/card.vue'
import { ArrowLeft as ArrowLeftIcon, ArrowRight as ArrowRightIcon, Star as StarIcon } from 'lucide-vue-next'

const props = defineProps({ id: { type: [String, Number], required: true } })
const router = useRouter()

const items = ref([])
const levels = ref([])
const cards = ref([])
const flipped = ref([])   // uids currently face up (unmatched), max 2
const matched = ref([])   // uids that are matched
const wrong = ref([])     // uids briefly shown as wrong
const moves = ref(0)
const locking = ref(false)
const loading = ref(true)
const won = ref(false)

const level = computed(() => levels.value.find(l => String(l.id) === String(props.id)))
const nextLevel = computed(() => {
  if (!level.value) return null
  return levels.value.find(l => l.order_index === level.value.order_index + 1) || null
})
const matchedPairs = computed(() => matched.value.length / 2)
const stars = computed(() => {
  const pairs = items.value.length
  if (!pairs) return 1
  if (moves.value <= pairs) return 3
  if (moves.value <= pairs * 2) return 2
  return 1
})

function shuffle(arr) {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[arr[i], arr[j]] = [arr[j], arr[i]]
  }
  return arr
}

function buildDeck() {
  const deck = []
  items.value.forEach(item => {
    deck.push({ uid: `i${item.id}`, pairId: item.id, type: 'image', content: item.emoji })
    deck.push({ uid: `w${item.id}`, pairId: item.id, type: 'word', content: item.word_en })
  })
  cards.value = shuffle(deck)
}

function isFaceUp(card) {
  return flipped.value.includes(card.uid) || matched.value.includes(card.uid)
}

function cardClass(card) {
  if (matched.value.includes(card.uid)) return 'bg-green-50 border-green-400 card-matched'
  if (wrong.value.includes(card.uid)) return 'bg-rose-50 border-rose-400 card-wrong'
  if (flipped.value.includes(card.uid)) return 'bg-card border-primary shadow-md'
  return 'bg-gradient-to-br from-primary to-emerald-600 border-primary/40 text-white shadow active:scale-95'
}

function flip(card) {
  if (locking.value || isFaceUp(card) || won.value) return
  hapticFeedback('light')
  flipped.value.push(card.uid)

  const item = items.value.find(i => i.id === card.pairId)
  if (item) speakWord(item.word_en)

  if (flipped.value.length === 2) {
    moves.value++
    locking.value = true
    const [aUid, bUid] = flipped.value
    const a = cards.value.find(c => c.uid === aUid)
    const b = cards.value.find(c => c.uid === bUid)

    if (a.pairId === b.pairId && a.uid !== b.uid) {
      // match!
      setTimeout(() => {
        matched.value.push(aUid, bUid)
        flipped.value = []
        locking.value = false
        hapticFeedback('medium')
        if (matched.value.length === cards.value.length) finish()
      }, 350)
    } else {
      // wrong
      wrong.value = [aUid, bUid]
      hapticFeedback('rigid')
      setTimeout(() => {
        flipped.value = []
        wrong.value = []
        locking.value = false
      }, 800)
    }
  }
}

async function finish() {
  won.value = true
  hapticFeedback('heavy')
  try {
    await learnApi.completeLevel(props.id, stars.value)
  } catch {}
}

function restart() {
  flipped.value = []
  matched.value = []
  wrong.value = []
  moves.value = 0
  locking.value = false
  won.value = false
  buildDeck()
}

function backToMap() {
  router.push('/trainer')
}

function goNext() {
  if (nextLevel.value) {
    router.push(`/trainer/${nextLevel.value.id}`)
  } else {
    backToMap()
  }
}

async function load() {
  loading.value = true
  won.value = false
  try {
    const [cardsRes, levelsRes] = await Promise.all([
      learnApi.getCards(props.id),
      learnApi.getLevels(),
    ])
    items.value = cardsRes.data || []
    levels.value = levelsRes.data || []
    restart()
  } catch {
    items.value = []
  } finally {
    loading.value = false
  }
}

// Reload when navigating between levels (Next city) without unmounting.
watch(() => props.id, load)

onMounted(load)
onUnmounted(stopSpeaking)
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
