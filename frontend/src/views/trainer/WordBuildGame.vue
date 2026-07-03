<template>
  <div>
    <p class="text-center text-muted-foreground text-xs mb-1">{{ t('wordBuild.hint') }}</p>
    <p class="text-center text-xs font-semibold text-primary mb-4">
      {{ t('wordBuild.wordOf', { current: currentIndex + 1, total: items.length }) }}
    </p>

    <!-- Revealed translation card -->
    <div v-if="revealed" class="flex flex-col items-center justify-center py-10 animate-pop">
      <div class="text-6xl mb-4">{{ currentItem.emoji }}</div>
      <p class="text-2xl font-extrabold text-primary">{{ currentItem.word_ru }}</p>
      <p class="text-lg font-bold text-muted-foreground mt-1">{{ currentItem.word_en }}</p>
      <Button variant="gradient" class="mt-6" @click="next">
        {{ t('wordBuild.next') }} <ArrowRightIcon class="h-4 w-4" />
      </Button>
    </div>

    <!-- Word-building puzzle -->
    <template v-else>
      <div class="text-5xl text-center mb-6">{{ currentItem.emoji }}</div>

      <!-- Answer slots -->
      <div class="flex justify-center flex-wrap gap-2 mb-8">
        <div
          v-for="(letter, i) in targetWord.split('')"
          :key="i"
          class="h-12 w-10 rounded-xl border-2 flex items-center justify-center text-xl font-extrabold uppercase"
          :class="placed[i] ? 'border-primary bg-primary/10 text-primary' : 'border-dashed border-border text-transparent'"
        >
          {{ placed[i] ? letterBlocks.find(b => b.id === placed[i]).letter : letter }}
        </div>
      </div>

      <!-- Letter blocks -->
      <div class="flex justify-center flex-wrap gap-2" :class="{ 'animate-shake': wrongShake }">
        <button
          v-for="block in letterBlocks"
          :key="block.id"
          class="h-12 w-12 rounded-xl flex items-center justify-center text-xl font-extrabold uppercase select-none transition-all duration-150 border-2 bg-gradient-to-br from-primary to-emerald-600 border-primary/40 text-white shadow active:scale-95"
          :disabled="placed.includes(block.id)"
          :class="{ 'opacity-0 pointer-events-none': placed.includes(block.id) }"
          @click="tapBlock(block)"
        >
          {{ block.letter }}
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
import Button from '../../components/ui/button.vue'
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
const letterBlocks = ref([])
const wrongShake = ref(false)

const currentItem = computed(() => props.items[currentIndex.value])
const targetWord = computed(() => currentItem.value.word_en.toLowerCase())

function shuffle(arr) {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[arr[i], arr[j]] = [arr[j], arr[i]]
  }
  return arr
}

function setupWord() {
  const blocks = targetWord.value.split('').map((letter, i) => ({ id: `${currentIndex.value}-${i}`, letter }))
  letterBlocks.value = shuffle(blocks)
  placed.value = []
  revealed.value = false
}

function tapBlock(block) {
  if (revealed.value || placed.value.includes(block.id)) return

  const expected = targetWord.value[placed.value.length]
  if (block.letter === expected) {
    placed.value.push(block.id)
    hapticFeedback('light')
    if (placed.value.length === targetWord.value.length) {
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
    setupWord()
  } else {
    const stars = mistakes.value === 0 ? 3 : mistakes.value <= props.items.length ? 2 : 1
    emit('complete', stars)
  }
}

onMounted(setupWord)
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
