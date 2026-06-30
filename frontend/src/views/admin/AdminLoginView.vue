<template>
  <div class="min-h-screen flex items-center justify-center bg-background">
    <Card class="p-8 w-full max-w-sm">
      <div class="text-center mb-6">
        <div class="text-5xl mb-3">🌙</div>
        <h1 class="text-xl font-bold">Admin Panel</h1>
        <p class="text-muted-foreground text-sm mt-1">Enter your admin token</p>
      </div>
      <form @submit.prevent="handleLogin" class="space-y-3">
        <Input
          v-model="tokenInput"
          type="password"
          placeholder="Admin token"
          autocomplete="current-password"
        />
        <p v-if="error" class="text-destructive text-xs">{{ error }}</p>
        <Button type="submit" variant="gradient" class="w-full" :disabled="loading">
          <Loader2Icon v-if="loading" class="h-4 w-4 animate-spin" />
          Sign In
        </Button>
      </form>
    </Card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '../../store/admin'
import adminApi from '../../services/adminApi'
import Card from '../../components/ui/card.vue'
import Input from '../../components/ui/input.vue'
import Button from '../../components/ui/button.vue'
import { Loader2 as Loader2Icon } from 'lucide-vue-next'

const router = useRouter()
const adminStore = useAdminStore()

const tokenInput = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  if (!tokenInput.value) return
  loading.value = true
  error.value = ''
  try {
    adminApi.defaults.headers['X-Admin-Token'] = tokenInput.value
    await adminApi.get('/stats')
    adminStore.setToken(tokenInput.value)
    router.push('/admin/dashboard')
  } catch (e) {
    error.value = e.response?.status === 401 ? 'Invalid token' : 'Connection error'
    adminApi.defaults.headers['X-Admin-Token'] = ''
  } finally {
    loading.value = false
  }
}
</script>
