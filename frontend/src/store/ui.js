import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUIStore = defineStore('ui', () => {
  const notifications = ref([])
  let notifId = 0

  function notify(message, type = 'info', duration = 3000) {
    const id = ++notifId
    notifications.value.push({ id, message, type })
    setTimeout(() => dismiss(id), duration)
  }

  function dismiss(id) {
    notifications.value = notifications.value.filter(n => n.id !== id)
  }

  function success(msg) { notify(msg, 'success') }
  function error(msg) { notify(msg, 'error', 5000) }
  function info(msg) { notify(msg, 'info') }

  return {
    notifications,
    notify, dismiss, success, error, info,
  }
})
