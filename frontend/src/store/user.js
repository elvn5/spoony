import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../services/api'
import { storage } from '../services/storage'
import { getTelegramInitData } from '../services/telegram'

export const useUserStore = defineStore('user', () => {
  const user = ref(storage.get('user'))
  const token = ref(storage.get('token'))
  const loading = ref(false)
  const error = ref(null)

  const isAuthenticated = computed(() => !!token.value)

  async function loginWithTelegram() {
    loading.value = true
    error.value = null
    try {
      const initData = getTelegramInitData()
      const res = await authApi.telegramLogin(initData || 'dev_mode')
      token.value = res.data.token
      user.value = res.data.user
      storage.set('token', token.value)
      storage.set('user', user.value)
    } catch (e) {
      error.value = e.response?.data?.error || 'Login failed'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchMe() {
    try {
      const res = await authApi.getMe()
      user.value = res.data
      storage.set('user', user.value)
    } catch {}
  }

  function logout() {
    authApi.logout().catch(() => {})
    token.value = null
    user.value = null
    storage.clear()
  }

  return {
    user, token, loading, error,
    isAuthenticated,
    loginWithTelegram, fetchMe, logout,
  }
})
