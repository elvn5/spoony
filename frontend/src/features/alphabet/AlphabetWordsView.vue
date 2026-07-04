<template>
  <div class="max-w-md md:max-w-2xl mx-auto px-4 pb-24 md:pb-10 pt-5 md:pt-8 min-h-screen">
    <!-- Header -->
    <header class="flex items-center gap-3 mb-4">
      <Button variant="ghost" size="icon" class="shrink-0" @click="backToLevels">
        <ArrowLeftIcon class="h-5 w-5" />
      </Button>
      <div class="min-w-0">
        <h1 class="text-lg font-extrabold leading-none truncate">
          {{ $t('alphabet.level3Title') }}
        </h1>
      </div>
    </header>

    <WordBuildGame :key="restartKey" :items="words" @complete="finish" />

    <!-- Win overlay -->
    <Transition name="fade">
      <div v-if="won" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm px-6">
        <Card class="w-full max-w-xs p-6 text-center animate-pop">
          <div class="text-6xl mb-3">🎉</div>
          <h2 class="text-xl font-extrabold mb-1">{{ $t('alphabet.level3Win') }}</h2>
          <p class="text-muted-foreground text-sm mb-4">{{ $t('alphabet.level3WinSub') }}</p>

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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { hapticFeedback } from '../../services/telegram'
import { markAlphabetLevelCompleted } from './progress'
import { greetingWords } from './data/greetingWords'
import Button from '../../components/ui/button.vue'
import Card from '../../components/ui/card.vue'
import WordBuildGame from '../../components/games/WordBuildGame.vue'
import { ArrowLeft as ArrowLeftIcon } from 'lucide-vue-next'

const router = useRouter()

const words = greetingWords

const won = ref(false)
const restartKey = ref(0)

function finish() {
  won.value = true
  hapticFeedback('heavy')
  markAlphabetLevelCompleted(4)
}

function restart() {
  won.value = false
  restartKey.value++
}

function backToLevels() {
  router.push('/alphabet')
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
