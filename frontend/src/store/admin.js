import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import adminApi from '../services/adminApi'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(sessionStorage.getItem('admin_token') || '')
  const isAuthed = computed(() => token.value !== '')
  const stats = ref(null)

  function setToken(t) {
    token.value = t
    sessionStorage.setItem('admin_token', t)
    adminApi.defaults.headers['X-Admin-Token'] = t
  }

  function logout() {
    token.value = ''
    sessionStorage.removeItem('admin_token')
  }

  async function fetchStats() {
    const res = await adminApi.get('/stats')
    stats.value = res.data
  }

  return { token, isAuthed, stats, setToken, logout, fetchStats }
})
