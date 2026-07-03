<template>
  <div class="max-w-md md:max-w-2xl mx-auto px-4 pb-24 md:pb-10 pt-5 md:pt-8">
    <header class="mb-5">
      <div class="flex items-center gap-2">
        <PuzzleIcon class="h-5 w-5 text-primary" />
        <h1 class="text-xl font-extrabold">{{ $t('alphabet.title') }}</h1>
      </div>
      <p class="text-muted-foreground text-sm mt-1">🔤 {{ $t('alphabet.subtitle') }}</p>
      <p class="text-muted-foreground text-xs mt-2 leading-relaxed">{{ $t('alphabet.intro') }}</p>
    </header>

    <div class="space-y-3">
      <button
        v-for="lvl in levelCards"
        :key="lvl.id"
        class="w-full flex items-center gap-4 rounded-2xl border border-border bg-card p-4 text-left transition-transform active:scale-[0.98]"
        :class="lvl.unlocked ? '' : 'opacity-60'"
        @click="openLevel(lvl)"
      >
        <div
          class="h-14 w-14 shrink-0 rounded-2xl flex items-center justify-center text-2xl font-extrabold"
          :class="lvl.completed ? 'bg-green-500/15 text-green-600' : lvl.unlocked ? 'bg-primary/15 text-primary' : 'bg-muted text-muted-foreground'"
        >
          <span v-if="!lvl.unlocked">🔒</span>
          <span v-else-if="lvl.completed">✓</span>
          <span v-else>{{ lvl.id }}</span>
        </div>
        <div class="min-w-0 flex-1">
          <p class="font-bold truncate">{{ lvl.title }}</p>
          <p class="text-xs text-muted-foreground mt-0.5">{{ lvl.range }}</p>
        </div>
        <span class="text-sm font-semibold text-primary shrink-0" v-if="lvl.unlocked">
          {{ lvl.completed ? $t('alphabet.review') : $t('alphabet.play') }}
        </span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { hapticFeedback, showAlert } from '../services/telegram'
import { getCompletedAlphabetLevels } from '../services/alphabetProgress'
import { Puzzle as PuzzleIcon } from 'lucide-vue-next'

const router = useRouter()
const { t } = useI18n()

const levelCards = computed(() => {
  const completed = getCompletedAlphabetLevels()
  return [
    { id: 1, title: t('alphabet.level1Title'), range: 'A B C D E F G H I J K L M', completed: completed.includes(1), unlocked: true },
    { id: 2, title: t('alphabet.level2Title'), range: 'N O P Q R S T U V W X Y Z', completed: completed.includes(2), unlocked: completed.includes(1) },
  ]
})

function openLevel(lvl) {
  if (!lvl.unlocked) {
    hapticFeedback('rigid')
    showAlert(t('alphabet.locked'))
    return
  }
  hapticFeedback('light')
  router.push(`/alphabet/${lvl.id}`)
}
</script>
