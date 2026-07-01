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

    <!-- Login prompt for guests -->
    <Card v-if="!userStore.isAuthenticated" class="p-5 mb-5 border-primary/30 bg-primary/5">
      <!-- Inside Telegram: one-tap login -->
      <template v-if="inTelegram">
        <p class="text-sm text-foreground/90 mb-3">{{ $t('auth.loginToLearn') }}</p>
        <Button variant="gradient" class="w-full" :disabled="userStore.loading" @click="handleTelegram">
          <LoaderIcon v-if="userStore.loading" class="h-4 w-4 animate-spin" />
          <span v-else>{{ $t('common.getStarted') }}</span>
        </Button>
      </template>

      <!-- On the web: continue as a guest with an optional name -->
      <template v-else>
        <p class="text-sm text-foreground/90 mb-3">{{ $t('auth.webHint') }}</p>
        <form class="flex flex-col gap-3 sm:flex-row" @submit.prevent="handleGuest">
          <Input v-model="name" :placeholder="$t('auth.namePlaceholder')" maxlength="24" class="flex-1" />
          <Button type="submit" variant="gradient" :disabled="userStore.loading" class="shrink-0">
            <LoaderIcon v-if="userStore.loading" class="h-4 w-4 animate-spin" />
            <span v-else>{{ $t('auth.continueAsGuest') }}</span>
          </Button>
        </form>
      </template>

      <p v-if="userStore.error" class="text-destructive text-xs mt-2">{{ userStore.error }}</p>
    </Card>

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
import { learnApi } from '../services/api'
import { useUserStore } from '../store/user'
import { formatRelative } from '../utils/helpers'
import { hapticFeedback, isTelegramEnvironment } from '../services/telegram'
import Card from '../components/ui/card.vue'
import Badge from '../components/ui/badge.vue'
import Button from '../components/ui/button.vue'
import Input from '../components/ui/input.vue'
import { Heart as HeartIcon, MessageCircle as MessageCircleIcon, Loader2 as LoaderIcon } from 'lucide-vue-next'

const userStore = useUserStore()
const posts = ref([])
const loading = ref(true)
const liked = reactive({})
const name = ref('')
const inTelegram = isTelegramEnvironment()

async function loadNews() {
  loading.value = true
  try {
    const res = await learnApi.getNews()
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

async function handleTelegram() {
  try {
    await userStore.loginWithTelegram()
  } catch {}
}

async function handleGuest() {
  try {
    await userStore.loginAsGuest(name.value.trim())
  } catch {}
}

onMounted(loadNews)
</script>
