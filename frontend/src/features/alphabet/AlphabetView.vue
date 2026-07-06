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

    <!-- Only the unlocked part of the course is listed — the next level
         appears when the previous one is completed. -->
    <div class="space-y-3">
      <button
        v-for="lvl in levelCards"
        :key="lvl.id"
        class="w-full flex items-center gap-4 rounded-2xl border border-border bg-card p-4 text-left transition-transform active:scale-[0.98]"
        @click="openLevel(lvl)"
      >
        <div
          class="h-14 w-14 shrink-0 rounded-2xl flex items-center justify-center text-2xl font-extrabold"
          :class="lvl.completed ? 'bg-green-500/15 text-green-600' : 'bg-primary/15 text-primary'"
        >
          <span v-if="lvl.completed">✓</span>
          <span v-else>{{ lvl.id }}</span>
        </div>
        <div class="min-w-0 flex-1">
          <p class="font-bold truncate">{{ lvl.title }}</p>
          <p class="text-xs text-muted-foreground mt-0.5">{{ lvl.range }}</p>
        </div>
        <span class="text-sm font-semibold text-primary shrink-0">
          {{ lvl.completed ? $t('alphabet.review') : $t('alphabet.play') }}
        </span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { hapticFeedback } from '../../services/telegram'
import { getCompletedAlphabetLevels, loadAlphabetProgress } from './progress'
import { phonicsWords, COMBO_GROUPS } from './data/phonicsWords'
import { greetingWords } from './data/greetingWords'
import { Puzzle as PuzzleIcon } from 'lucide-vue-next'

const router = useRouter()
const { t } = useI18n()

onMounted(loadAlphabetProgress)

const levelCards = computed(() => {
  const completed = getCompletedAlphabetLevels()
  const greetingRange = greetingWords.map(w => w.word_en.toUpperCase()).join(' · ')
  const cards = [
    { id: 1, title: t('alphabet.level1Title'), range: 'A B C D E F G H I J K L M', completed: completed.includes(1), unlocked: true, path: '/alphabet/1' },
    { id: 2, title: t('alphabet.level2Title'), range: 'N O P Q R S T U V W X Y Z', completed: completed.includes(2), unlocked: completed.includes(1), path: '/alphabet/2' },
    { id: 3, title: t('alphabet.matchTitle'), range: greetingRange, completed: completed.includes(3), unlocked: completed.includes(2), path: '/alphabet/match/play' },
    { id: 4, title: t('alphabet.level3Title'), range: greetingRange, completed: completed.includes(4), unlocked: completed.includes(3), path: '/alphabet/words/play' },
  ]

  COMBO_GROUPS.forEach((group, i) => {
    const id = 5 + i
    const examples = phonicsWords
      .filter(w => group.combos.includes(w.combo))
      .slice(0, 3)
      .map(w => w.word_en.toUpperCase())
    cards.push({
      id,
      title: group.combos.join(', '),
      range: examples.join(' · ') + '...',
      completed: completed.includes(id),
      unlocked: completed.includes(id - 1),
      path: `/alphabet/combos/${group.id}`,
    })
  })

  // Hide everything past the current level: the next card appears
  // only when the previous one is completed.
  return cards.filter(c => c.unlocked)
})

function openLevel(lvl) {
  hapticFeedback('light')
  router.push(lvl.path)
}
</script>
