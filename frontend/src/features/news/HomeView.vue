<template>
  <div class="max-w-md md:max-w-2xl mx-auto px-4 pb-24 md:pb-10 pt-5 md:pt-8">
    <!-- Brand header (hidden on desktop — the sidebar already shows it) -->
    <header class="flex items-center gap-3 mb-5 md:hidden">
      <div class="h-12 w-12 rounded-2xl bg-primary/15 flex items-center justify-center text-3xl shadow-sm">🥄</div>
      <div>
        <h1 class="text-2xl font-extrabold leading-none text-gradient">Spoony</h1>
        <p class="text-muted-foreground text-xs mt-1">{{ $t('app.tagline') }}</p>
      </div>
    </header>

    <LoginCard />

    <h2 class="text-sm font-semibold text-muted-foreground mb-3 px-1">{{ $t('home.feedHint') }}</h2>

    <!-- Loading skeletons -->
    <div v-if="loading" class="space-y-4">
      <Card v-for="n in 3" :key="n" class="p-4 animate-pulse">
        <div class="h-4 w-1/2 bg-muted rounded mb-3"></div>
        <div class="h-3 w-full bg-muted rounded mb-2"></div>
        <div class="h-3 w-2/3 bg-muted rounded"></div>
      </Card>
    </div>

    <!-- Feed -->
    <div v-else class="space-y-4">
      <Card v-for="post in posts" :key="post.id" class="overflow-hidden">
        <div class="p-4">
          <!-- author row -->
          <div class="flex items-center gap-3 mb-3">
            <div class="h-10 w-10 rounded-full bg-secondary flex items-center justify-center text-xl">{{ post.avatar }}</div>
            <div class="flex-1 min-w-0">
              <p class="font-semibold text-sm leading-tight">{{ post.author }}</p>
              <p class="text-muted-foreground text-xs">{{ formatRelative(post.created_at) }}</p>
            </div>
            <Badge variant="secondary" class="shrink-0">{{ post.category }}</Badge>
          </div>

          <h3 class="font-bold text-base mb-1">{{ post.title }}</h3>
          <p class="text-sm text-foreground/80 leading-relaxed">{{ post.body }}</p>
        </div>

        <!-- big emoji "image" banner -->
        <div v-if="post.image" class="h-32 flex items-center justify-center text-6xl bg-gradient-to-br from-secondary to-accent/40">
          {{ post.image }}
        </div>

        <!-- actions -->
        <div class="flex items-center gap-5 px-4 py-3 border-t border-border">
          <button
            class="flex items-center gap-1.5 text-sm transition-colors"
            :class="liked[post.id] ? 'text-rose-500' : 'text-muted-foreground hover:text-rose-500'"
            @click="toggleLike(post)"
          >
            <HeartIcon class="h-4 w-4" :fill="liked[post.id] ? 'currentColor' : 'none'" />
            {{ post.likes + (liked[post.id] ? 1 : 0) }}
          </button>
          <div class="flex items-center gap-1.5 text-sm text-muted-foreground">
            <MessageCircleIcon class="h-4 w-4" />
            {{ $t('home.title') }}
          </div>
        </div>
      </Card>

      <p v-if="!posts.length" class="text-center text-muted-foreground text-sm py-10">
        {{ $t('common.error') }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { newsApi } from './api'
import { formatRelative } from '../../utils/helpers'
import { hapticFeedback } from '../../services/telegram'
import LoginCard from '../auth/components/LoginCard.vue'
import Card from '../../components/ui/card.vue'
import Badge from '../../components/ui/badge.vue'
import { Heart as HeartIcon, MessageCircle as MessageCircleIcon } from 'lucide-vue-next'

const posts = ref([])
const loading = ref(true)
const liked = reactive({})

async function loadNews() {
  loading.value = true
  try {
    const res = await newsApi.getNews()
    posts.value = res.data || []
  } catch {
    posts.value = []
  } finally {
    loading.value = false
  }
}

function toggleLike(post) {
  liked[post.id] = !liked[post.id]
  hapticFeedback('light')
}

onMounted(loadNews)
</script>
