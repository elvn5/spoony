<template>
  <div class="max-w-md mx-auto">
    <p class="text-center text-xs font-semibold text-primary mb-4">
      {{ t('theory.cardOf', { current: idx + 1, total: slides.length }) }}
    </p>

    <div :key="idx" class="rounded-3xl border border-border bg-card p-6 animate-pop">
      <h2 class="text-lg font-extrabold mb-3">{{ slide.title_ru }}</h2>
      <p class="text-sm leading-relaxed mb-4">{{ slide.body_ru }}</p>

      <button
        v-if="slide.example_en"
        class="w-full rounded-xl bg-secondary/60 p-3 text-left active:scale-[0.98] transition-transform"
        @click="speak"
      >
        <p class="font-bold text-primary flex items-center gap-2">
          {{ slide.example_en }}
          <Volume2Icon class="h-4 w-4 shrink-0" />
        </p>
        <p class="text-xs text-muted-foreground mt-1">{{ slide.example_ru }}</p>
      </button>
    </div>

    <!-- Progress dots -->
    <div class="flex justify-center gap-1.5 mt-4">
      <span
        v-for="(s, i) in slides" :key="i"
        class="h-1.5 rounded-full transition-all"
        :class="i === idx ? 'w-5 bg-primary' : 'w-1.5 bg-border'"
      />
    </div>

    <div class="flex justify-center gap-2 mt-5">
      <Button v-if="idx > 0" variant="secondary" @click="idx--">
        <ArrowLeftIcon class="h-4 w-4" /> {{ t('common.back') }}
      </Button>
      <Button v-if="idx + 1 < slides.length" variant="gradient" @click="idx++">
        {{ t('wordBuild.next') }} <ArrowRightIcon class="h-4 w-4" />
      </Button>
      <Button v-else variant="gradient" @click="emit('complete', 3)">
        {{ t('theory.done') }} ✅
      </Button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { hapticFeedback } from '../../services/telegram'
import { speakWord } from '../../services/tts'
import Button from '../ui/button.vue'
import { ArrowLeft as ArrowLeftIcon, ArrowRight as ArrowRightIcon, Volume2 as Volume2Icon } from 'lucide-vue-next'

const props = defineProps({
  slides: { type: Array, required: true },
})
const emit = defineEmits(['complete'])
const { t } = useI18n()

const idx = ref(0)
const slide = computed(() => props.slides[idx.value])

function speak() {
  hapticFeedback('light')
  speakWord(slide.value.example_en)
}
</script>
