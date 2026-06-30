<template>
  <div class="min-h-screen flex flex-col items-center justify-center px-6 text-center">
    <div class="relative z-10 max-w-sm w-full space-y-8">
      <div class="space-y-4">
        <h1 class="text-4xl font-bold tracking-tight">Welcome</h1>
        <p class="text-muted-foreground text-base leading-relaxed">
          {{ $t('auth.welcome') }}
        </p>
      </div>

      <div class="space-y-3">
        <template v-if="!userStore.isAuthenticated">
          <Button
            variant="gradient"
            size="xl"
            class="w-full"
            :disabled="userStore.loading"
            @click="handleLogin"
          >
            <span v-if="userStore.loading" class="flex items-center gap-2">
              <LoaderIcon class="h-4 w-4 animate-spin" /> {{ $t('common.loading') }}
            </span>
            <span v-else>{{ $t('common.getStarted') }}</span>
          </Button>
          <p v-if="userStore.error" class="text-destructive text-sm">{{ userStore.error }}</p>
        </template>

        <template v-else>
          <Button variant="gradient" size="xl" class="w-full" as="RouterLink" to="/profile">
            {{ $t('nav.profile') }}
          </Button>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import Button from '../components/ui/button.vue'
import { Loader2 as LoaderIcon } from 'lucide-vue-next'

const router = useRouter()
const userStore = useUserStore()

async function handleLogin() {
  try {
    await userStore.loginWithTelegram()
    router.push('/profile')
  } catch {}
}
</script>
