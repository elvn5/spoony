<template>
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
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '../store'
import { isTelegramEnvironment } from '../../../services/telegram'
import Card from '../../../components/ui/card.vue'
import Button from '../../../components/ui/button.vue'
import Input from '../../../components/ui/input.vue'
import { Loader2 as LoaderIcon } from 'lucide-vue-next'

const userStore = useUserStore()
const name = ref('')
const inTelegram = isTelegramEnvironment()

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
</script>
