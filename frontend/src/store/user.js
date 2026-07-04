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
  // A guest account has no linked Telegram identity yet.
  const isGuest = computed(() => isAuthenticated.value && !user.value?.telegram_id)

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

  async function loginAsGuest(name) {
    loading.value = true
    error.value = null
    try {
      let guestId = storage.get('guest_id')
      if (!guestId) {
        guestId = (crypto?.randomUUID?.() || 'g_' + Math.random().toString(36).slice(2) + Date.now().toString(36))
        storage.set('guest_id', guestId)
      }
      const res = await authApi.guestLogin(guestId, name)
      token.value = res.data.token
      user.value = res.data.user
      if (res.data.guest_id) storage.set('guest_id', res.data.guest_id)
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
    // Keep guest_id & lang so a returning guest recovers their account/progress.
    storage.remove('token')
    storage.remove('user')
  }

  return {
    user, token, loading, error,
    isAuthenticated, isGuest,
    loginWithTelegram, loginAsGuest, fetchMe, logout,
  }
})
