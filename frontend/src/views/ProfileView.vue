<template>
  <div class="max-w-md mx-auto px-4 pb-24 pt-4">
    <h1 class="text-xl font-bold mb-5 pt-1">{{ $t('nav.profile') }}</h1>

    <div v-if="userStore.user" class="space-y-4">
      <!-- User card -->
      <Card class="p-4">
        <div class="flex items-center gap-4">
          <Avatar :src="userStore.user.avatar_url" :alt="userStore.user.first_name" size="lg" />
          <div class="flex-1 min-w-0">
            <p class="font-semibold text-base">{{ userStore.user.first_name }} {{ userStore.user.last_name }}</p>
            <p class="text-muted-foreground text-sm">@{{ userStore.user.username || 'user' }}</p>
            <p class="text-muted-foreground text-xs mt-0.5">
              {{ $t('profile.memberSince') }} {{ formatDate(userStore.user.created_at) }}
            </p>
            <p class="text-muted-foreground text-xs mt-0.5">
              TG ID: {{ userStore.user.telegram_id }}
            </p>
          </div>
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
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import { formatDate } from '../utils/helpers'
import Button from '../components/ui/button.vue'
import Card from '../components/ui/card.vue'
import Avatar from '../components/ui/avatar.vue'
import { LogOut as LogOutIcon } from 'lucide-vue-next'

const router = useRouter()
const userStore = useUserStore()

function handleLogout() {
  userStore.logout()
  router.push('/')
}
</script>
