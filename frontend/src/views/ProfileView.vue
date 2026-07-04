<template>
  <div class="max-w-md md:max-w-2xl mx-auto px-4 pb-24 md:pb-10 pt-5 md:pt-8">
    <h1 class="text-xl font-extrabold mb-5 pt-1">{{ $t('nav.profile') }}</h1>

    <div v-if="userStore.user" class="space-y-4">
      <!-- User card -->
      <Card class="p-4">
        <div class="flex items-center gap-4">
          <Avatar :src="userStore.user.avatar_url" :alt="userStore.user.first_name" size="lg" />
          <div class="flex-1 min-w-0">
            <p class="font-bold text-base">{{ userStore.user.first_name }} {{ userStore.user.last_name }}</p>
            <p class="text-muted-foreground text-sm">@{{ userStore.user.username || 'user' }}</p>
            <p class="text-muted-foreground text-xs mt-0.5">
              {{ $t('profile.memberSince') }} {{ formatDate(userStore.user.created_at) }}
            </p>
          </div>
        </div>
      </Card>

      <!-- Stats -->
      <div>
        <h2 class="text-sm font-semibold text-muted-foreground mb-2 px-1">{{ $t('profile.stats') }}</h2>
        <div class="grid grid-cols-3 gap-3">
          <Card class="p-3 text-center">
            <div class="text-2xl mb-1">🏙️</div>
            <p class="font-extrabold text-lg leading-none">{{ stats.completed_levels }}<span class="text-muted-foreground text-sm">/{{ stats.total_levels }}</span></p>
            <p class="text-[10px] text-muted-foreground mt-1 leading-tight">{{ $t('profile.levelsDone') }}</p>
          </Card>
          <Card class="p-3 text-center">
            <div class="text-2xl mb-1">📖</div>
            <p class="font-extrabold text-lg leading-none">{{ stats.learned_words }}</p>
            <p class="text-[10px] text-muted-foreground mt-1 leading-tight">{{ $t('profile.wordsLearned') }}</p>
          </Card>
          <Card class="p-3 text-center">
            <div class="text-2xl mb-1">⭐</div>
            <p class="font-extrabold text-lg leading-none">{{ stats.total_stars }}</p>
            <p class="text-[10px] text-muted-foreground mt-1 leading-tight">{{ $t('profile.starsEarned') }}</p>
          </Card>
        </div>
      </div>

      <!-- Continue in Telegram (guest accounts on the website only) -->
      <Card v-if="userStore.isGuest && !inTelegram && telegramLink" class="p-4 border-primary/30 bg-primary/5">
        <p class="text-sm font-semibold mb-1">{{ $t('profile.continueInTelegram') }}</p>
        <p class="text-muted-foreground text-xs mb-3">{{ $t('profile.continueInTelegramHint') }}</p>
        <Button as="a" :href="telegramLink" target="_blank" rel="noopener" variant="gradient" class="w-full">
          <SendIcon class="h-4 w-4" /> {{ $t('profile.openTelegramBot') }}
        </Button>
      </Card>

      <!-- Language switch -->
      <Card class="p-4">
        <p class="text-sm font-semibold mb-3">{{ $t('profile.language') }}</p>
        <div class="flex gap-2">
          <Button
            v-for="l in ['ru', 'en']" :key="l"
            :variant="locale === l ? 'default' : 'outline'"
            size="sm"
            class="flex-1"
            @click="setLang(l)"
          >
            {{ $t('languages.' + l) }}
          </Button>
        </div>
      </Card>

      <!-- Logout -->
      <Button
        variant="ghost"
        class="w-full text-destructive hover:text-destructive hover:bg-destructive/10"
        @click="handleLogout"
      >
        <LogOutIcon class="h-4 w-4" /> {{ $t('profile.signOut') }}
      </Button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../store/user'
import { learnApi, telegramApi } from '../services/api'
import { storage } from '../services/storage'
import { isTelegramEnvironment } from '../services/telegram'
import { formatDate } from '../utils/helpers'
import Button from '../components/ui/button.vue'
import Card from '../components/ui/card.vue'
import Avatar from '../components/ui/avatar.vue'
import { LogOut as LogOutIcon, Send as SendIcon } from 'lucide-vue-next'

const router = useRouter()
const userStore = useUserStore()
const { locale } = useI18n()

const stats = reactive({ total_levels: 0, completed_levels: 0, total_stars: 0, learned_words: 0 })
const inTelegram = isTelegramEnvironment()
const telegramLink = ref('')

async function loadStats() {
  try {
    const res = await learnApi.getStats()
    Object.assign(stats, res.data)
  } catch {}
}

async function loadTelegramLink() {
  const guestId = storage.get('guest_id')
  if (!guestId) return
  try {
    const res = await telegramApi.getBotInfo()
    const username = res.data?.username
    if (username) telegramLink.value = `https://t.me/${username}?start=${encodeURIComponent(guestId)}`
  } catch {}
}

function setLang(l) {
  locale.value = l
  localStorage.setItem('lang', l)
}

function handleLogout() {
  userStore.logout()
  router.push('/')
}

onMounted(() => {
  loadStats()
  loadTelegramLink()
})
</script>
