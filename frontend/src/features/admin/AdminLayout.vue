<template>
  <div class="min-h-screen bg-background">
    <!-- Header -->
    <div v-if="adminStore.isAuthed" class="border-b border-border bg-card/50 backdrop-blur-md sticky top-0 z-40">
      <div class="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <span class="text-xl">⚡</span>
          <span class="font-semibold text-sm">Spoony Admin</span>
        </div>
        <div class="flex items-center gap-1">
          <RouterLink
            v-for="tab in tabs" :key="tab.to"
            :to="tab.to"
            class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors"
            :class="route.path.startsWith(tab.to)
              ? 'bg-primary/20 text-primary'
              : 'text-muted-foreground hover:text-foreground hover:bg-accent'"
          >
            {{ tab.icon }} {{ tab.label }}
          </RouterLink>
          <button
            class="ml-3 px-3 py-1.5 rounded-lg text-xs font-medium text-destructive hover:bg-destructive/10 transition-colors"
            @click="handleLogout"
          >
            Exit
          </button>
        </div>
      </div>
    </div>

    <div class="max-w-6xl mx-auto px-4 py-6">
      <RouterView />
    </div>
  </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from './store'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()

const tabs = [
  { to: '/admin/dashboard', icon: '📊', label: 'Dashboard' },
  { to: '/admin/users',     icon: '👤', label: 'Users'     },
  { to: '/admin/content',   icon: '📰', label: 'Content'   },
]

function handleLogout() {
  adminStore.logout()
  router.push('/admin/login')
}
</script>
