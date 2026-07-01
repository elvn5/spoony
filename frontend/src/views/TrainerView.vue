<template>
  <div class="max-w-md md:max-w-3xl mx-auto px-4 pb-24 md:pb-10 pt-5 md:pt-8">
    <!-- Header -->
    <header class="mb-4">
      <div class="flex items-center gap-2">
        <MapIcon class="h-5 w-5 text-primary" />
        <h1 class="text-xl font-extrabold">{{ $t('trainer.title') }}</h1>
      </div>
      <p class="text-muted-foreground text-sm mt-1">🇬🇧 {{ $t('trainer.subtitle') }}</p>
      <p class="text-muted-foreground text-xs mt-2 leading-relaxed">{{ $t('trainer.intro') }}</p>
    </header>

    <div v-if="loading" class="space-y-3">
      <Card v-for="n in 4" :key="n" class="h-16 animate-pulse" />
    </div>

    <!-- The route map -->
    <div
      v-else
      class="relative rounded-3xl bg-gradient-to-b from-secondary/60 via-background to-secondary/40 border border-border overflow-hidden"
      :style="{ height: mapHeight + 'px' }"
    >
      <!-- winding road -->
      <svg class="absolute inset-0 w-full h-full" viewBox="0 0 100 100" preserveAspectRatio="none">
        <path
          :d="pathD"
          fill="none"
          stroke="hsl(var(--primary) / 0.25)"
          stroke-width="6"
          stroke-linecap="round"
          stroke-linejoin="round"
          vector-effect="non-scaling-stroke"
        />
        <path
          :d="pathD"
          fill="none"
          stroke="hsl(var(--primary) / 0.6)"
          stroke-width="2"
          stroke-dasharray="1 5"
          stroke-linecap="round"
          vector-effect="non-scaling-stroke"
        />
      </svg>

      <!-- decorative scenery -->
      <span class="absolute text-2xl opacity-70 select-none" style="left:78%; top:84%">🌳</span>
      <span class="absolute text-2xl opacity-70 select-none" style="left:12%; top:64%">⛅</span>
      <span class="absolute text-2xl opacity-70 select-none" style="left:80%; top:30%">🐑</span>
      <span class="absolute text-2xl opacity-70 select-none" style="left:14%; top:14%">🏁</span>

      <!-- city nodes -->
      <div
        v-for="level in levels"
        :key="level.id"
        class="absolute flex flex-col items-center"
        :style="{ left: level.pos_x + '%', top: level.pos_y + '%', transform: 'translate(-50%, -50%)' }"
      >
        <button
          class="relative h-16 w-16 rounded-full flex items-center justify-center text-3xl transition-transform active:scale-95 shadow-lg"
          :class="nodeClass(level)"
          :disabled="!level.unlocked"
          @click="openLevel(level)"
        >
          <span v-if="!level.unlocked" class="text-2xl">🔒</span>
          <span v-else>{{ level.emoji }}</span>

          <!-- completed check -->
          <span
            v-if="level.completed"
            class="absolute -top-1 -right-1 h-6 w-6 rounded-full bg-green-500 text-white flex items-center justify-center text-xs shadow"
          >✓</span>
        </button>

        <!-- label -->
        <div class="mt-1.5 text-center">
          <p class="text-xs font-bold leading-none" :class="level.unlocked ? 'text-foreground' : 'text-muted-foreground'">
            {{ level.city }}
          </p>
          <p class="text-[10px] text-muted-foreground">{{ level.title_ru }}</p>

          <!-- stars -->
          <div v-if="level.completed" class="flex justify-center gap-0.5 mt-0.5">
            <StarIcon
              v-for="s in 3" :key="s"
              class="h-3 w-3"
              :class="s <= level.stars ? 'text-yellow-400' : 'text-muted-foreground/30'"
              :fill="s <= level.stars ? 'currentColor' : 'none'"
            />
          </div>
        </div>
      </div>
    </div>

    <p v-if="!loading && !levels.length" class="text-center text-muted-foreground text-sm py-10">
      {{ $t('common.error') }}
    </p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { learnApi } from '../services/api'
import { hapticFeedback, showAlert } from '../services/telegram'
import { useI18n } from 'vue-i18n'
import Card from '../components/ui/card.vue'
import { Map as MapIcon, Star as StarIcon } from 'lucide-vue-next'

const router = useRouter()
const { t } = useI18n()
const levels = ref([])
const loading = ref(true)

const mapHeight = computed(() => Math.max(560, levels.value.length * 130))

// SVG path connecting the cities in order (coords are 0..100 percentages).
const pathD = computed(() => {
  if (!levels.value.length) return ''
  return levels.value
    .map((l, i) => `${i === 0 ? 'M' : 'L'} ${l.pos_x} ${l.pos_y}`)
    .join(' ')
})

function nodeClass(level) {
  if (level.completed) return 'bg-green-500/15 ring-4 ring-green-500/60'
  if (level.unlocked) return 'bg-primary text-primary-foreground ring-4 ring-primary/30 animate-bounce-slow'
  return 'bg-muted ring-4 ring-border opacity-70 cursor-not-allowed'
}

function openLevel(level) {
  if (!level.unlocked) {
    hapticFeedback('rigid')
    showAlert(t('trainer.locked'))
    return
  }
  hapticFeedback('light')
  router.push(`/trainer/${level.id}`)
}

async function load() {
  loading.value = true
  try {
    const res = await learnApi.getLevels()
    levels.value = res.data || []
  } catch {
    levels.value = []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
@keyframes bounce-slow {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}
.animate-bounce-slow { animation: bounce-slow 1.6s ease-in-out infinite; }
</style>
