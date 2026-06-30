<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold">Users</h1>
      <span class="text-muted-foreground text-sm">{{ total }} total</span>
    </div>

    <!-- Search -->
    <div class="relative mb-4">
      <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
      <Input v-model="search" placeholder="Search by name, username or Telegram ID…" class="pl-9"
        @input="onSearch" />
    </div>

    <!-- Table -->
    <Card class="overflow-hidden">
      <div v-if="loading" class="p-6 text-center text-muted-foreground text-sm">Loading…</div>
      <div v-else-if="!users.length" class="p-6 text-center text-muted-foreground text-sm">No users found</div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-border text-left">
              <th class="px-4 py-3 text-xs text-muted-foreground font-medium uppercase tracking-wider">ID</th>
              <th class="px-4 py-3 text-xs text-muted-foreground font-medium uppercase tracking-wider">User</th>
              <th class="px-4 py-3 text-xs text-muted-foreground font-medium uppercase tracking-wider">Joined</th>
              <th class="px-4 py-3 text-xs text-muted-foreground font-medium uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id" class="border-b border-border/50 hover:bg-accent/30 transition-colors">
              <td class="px-4 py-3 text-muted-foreground font-mono text-xs">{{ u.id }}</td>
              <td class="px-4 py-3">
                <p class="font-medium">{{ u.first_name }} {{ u.last_name }}</p>
                <p class="text-xs text-muted-foreground">{{ u.username ? '@' + u.username : '' }} · TG {{ u.telegram_id }}</p>
              </td>
              <td class="px-4 py-3 text-muted-foreground text-xs">{{ u.created_at.slice(0, 10) }}</td>
              <td class="px-4 py-3">
                <Button size="sm" variant="ghost"
                  class="text-xs h-7 text-destructive hover:text-destructive"
                  @click="deleteUser(u)">Delete</Button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </Card>

    <!-- Pagination -->
    <div v-if="total > perPage" class="flex items-center justify-center gap-3 mt-4">
      <Button variant="ghost" size="sm" :disabled="page === 1" @click="page--; load()">← Prev</Button>
      <span class="text-muted-foreground text-sm">Page {{ page }} / {{ Math.ceil(total / perPage) }}</span>
      <Button variant="ghost" size="sm" :disabled="page * perPage >= total" @click="page++; load()">Next →</Button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import adminApi from '../../services/adminApi'
import Card from '../../components/ui/card.vue'
import Input from '../../components/ui/input.vue'
import Button from '../../components/ui/button.vue'
import { Search as SearchIcon } from 'lucide-vue-next'

const users  = ref([])
const total  = ref(0)
const page   = ref(1)
const search = ref('')
const loading = ref(false)
const perPage = 20

let searchTimer = null
function onSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { page.value = 1; load() }, 400)
}

async function load() {
  loading.value = true
  try {
    const { data } = await adminApi.get('/users', { params: { page: page.value, search: search.value } })
    users.value = data.users || []
    total.value = data.total
  } finally {
    loading.value = false
  }
}

async function deleteUser(u) {
  if (!confirm(`Delete ${u.first_name} and all their data? This cannot be undone.`)) return
  await adminApi.delete(`/users/${u.id}`)
  users.value = users.value.filter(x => x.id !== u.id)
}

onMounted(load)
</script>
