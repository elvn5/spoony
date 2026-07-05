<template>
  <div>
    <p class="text-center text-muted-foreground text-xs mb-1">{{ t('sentenceBuild.hint') }}</p>
    <p class="text-center text-xs font-semibold text-primary mb-4">
      {{ t('wordBuild.wordOf', { current: currentIndex + 1, total: items.length }) }}
    </p>

    <!-- Revealed full sentence -->
    <div v-if="revealed" class="flex flex-col items-center justify-center py-10 animate-pop">
      <div class="text-6xl mb-4">{{ currentItem.emoji }}</div>
      <p class="text-xl font-extrabold text-primary text-center px-4">{{ currentItem.word_en }}</p>
      <p class="text-base font-bold text-muted-foreground mt-1 text-center px-4">{{ currentItem.word_ru }}</p>
      <Button variant="gradient" class="mt-6" @click="next">
        {{ t('wordBuild.next') }} <ArrowRightIcon class="h-4 w-4" />
      </Button>
    </div>

    <!-- Sentence-building puzzle -->
    <template v-else>
      <div class="text-5xl text-center mb-2">{{ currentItem.emoji }}</div>
      <p class="text-center text-lg font-bold mb-6 px-4">{{ currentItem.word_ru }}</p>

      <!-- Placed words -->
      <div class="flex justify-center flex-wrap gap-2 mb-8 min-h-11 px-2">
        <div
          v-for="(blockId, i) in placed"
          :key="i"
          class="h-11 px-3 rounded-xl border-2 border-primary bg-primary/10 text-primary flex items-center justify-center text-base font-extrabold"
        >
          {{ wordBlocks.find(b => b.id === blockId).word }}
        </div>
        <div
          v-if="placed.length < targetWords.length"
          class="h-11 w-16 rounded-xl border-2 border-dashed border-border"
        />
      </div>

      <!-- Word blocks -->
      <div class="flex justify-center flex-wrap gap-2 px-2" :class="{ 'animate-shake': wrongShake }">
        <button
          v-for="block in wordBlocks"
          :key="block.id"
          class="h-11 px-4 rounded-xl flex items-center justify-center text-base font-extrabold select-none transition-all duration-150 border-2 bg-gradient-to-br from-primary to-emerald-600 border-primary/40 text-white shadow active:scale-95"
          :disabled="placed.includes(block.id)"
          :class="{ 'opacity-0 pointer-events-none': placed.includes(block.id) }"
          @click="tapBlock(block)"
        >
          {{ block.word }}
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { hapticFeedback } from '../../services/telegram'
import { speakWord } from '../../services/tts'
import Button from '../ui/button.vue'
import { ArrowRight as ArrowRightIcon } from 'lucide-vue-next'

const props = defineProps({
  items: { type: Array, required: true },
})
const emit = defineEmits(['complete'])
const { t } = useI18n()

const currentIndex = ref(0)
const mistakes = ref(0)
const revealed = ref(false)
const placed = ref([])
const wordBlocks = ref([])
const wrongShake = ref(false)

const currentItem = computed(() => props.items[currentIndex.value])
const targetWords = computed(() => currentItem.value.word_en.split(' '))

function shuffle(arr) {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[arr[i], arr[j]] = [arr[j], arr[i]]
  }
  return arr
}

function setupSentence() {
  const blocks = targetWords.value.map((word, i) => ({ id: `${currentIndex.value}-${i}`, word }))
  wordBlocks.value = shuffle(blocks)
  placed.value = []
  revealed.value = false
}

function tapBlock(block) {
  if (revealed.value || placed.value.includes(block.id)) return

  const expected = targetWords.value[placed.value.length]
  if (block.word === expected) {
    placed.value.push(block.id)
    hapticFeedback('light')
    if (placed.value.length === targetWords.value.length) {
      setTimeout(() => {
        revealed.value = true
        speakWord(currentItem.value.word_en)
        hapticFeedback('medium')
      }, 250)
    }
  } else {
    mistakes.value++
    hapticFeedback('rigid')
    wrongShake.value = true
    setTimeout(() => { wrongShake.value = false }, 400)
  }
}

function next() {
  if (currentIndex.value + 1 < props.items.length) {
    currentIndex.value++
    setupSentence()
  } else {
    const stars = mistakes.value === 0 ? 3 : mistakes.value <= props.items.length ? 2 : 1
    emit('complete', stars)
  }
}

onMounted(setupSentence)
</script>

<style scoped>
@keyframes shake {
  10%, 90% { transform: translateX(-2px); }
  20%, 80% { transform: translateX(4px); }
  30%, 50%, 70% { transform: translateX(-8px); }
  40%, 60% { transform: translateX(8px); }
}
.animate-shake { animation: shake 0.4s; }
</style>
