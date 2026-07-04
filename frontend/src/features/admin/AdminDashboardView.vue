<template>
  <div>
    <h1 class="text-xl font-bold mb-6">Dashboard</h1>

    <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 gap-4">
      <Skeleton v-for="i in 6" :key="i" class="h-28 rounded-xl" />
    </div>

    <div v-else class="grid grid-cols-2 md:grid-cols-3 gap-4">
      <Card
        v-for="s in statCards" :key="s.label"
        class="p-5 flex flex-col gap-2"
      >
        <span class="text-3xl">{{ s.icon }}</span>
        <span class="text-3xl font-bold text-foreground">{{ s.value ?? '—' }}</span>
        <span class="text-xs text-muted-foreground uppercase tracking-wider">{{ s.label }}</span>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAdminStore } from './store'
import Card from '../../components/ui/card.vue'
import Skeleton from '../../components/ui/skeleton.vue'

const adminStore = useAdminStore()
const loading = ref(false)

const statCards = computed(() => {
  const s = adminStore.stats
  if (!s) return []
  return [
    { icon: '👤', label: 'Total Users',      value: s.total_users },
    { icon: '🆕', label: 'New Users Today',   value: s.new_users_today },
    { icon: '✈️', label: 'Telegram Users',    value: s.telegram_users },
    { icon: '🌐', label: 'Guest Users',       value: s.guest_users },
    { icon: '✅', label: 'Levels Completed',  value: s.completed_levels },
    { icon: '⭐', label: 'Stars Awarded',     value: s.total_stars },
    { icon: '📰', label: 'Home Feed Posts',   value: s.total_news_posts },
  ]
})

onMounted(async () => {
  loading.value = true
  await adminStore.fetchStats()
  loading.value = false
})
</script>
