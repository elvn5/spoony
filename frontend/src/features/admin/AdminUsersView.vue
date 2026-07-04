<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold">Users</h1>
      <span class="text-muted-foreground text-sm">{{ total }} total</span>
    </div>

    <!-- Search + filters -->
    <div class="flex flex-col sm:flex-row gap-3 mb-4">
      <div class="relative flex-1">
        <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
        <Input v-model="search" placeholder="Search by name, username or Telegram ID…" class="pl-9"
          @input="onSearch" />
      </div>
      <Select :value="type" class="sm:w-44" @change="e => { type = e.target.value; page = 1; load() }">
        <option value="all">All users</option>
        <option value="telegram">✈️ Telegram only</option>
        <option value="guest">🌐 Guest only</option>
      </Select>
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
              <th class="px-4 py-3 text-xs text-muted-foreground font-medium uppercase tracking-wider">Type</th>
              <th class="px-4 py-3 text-xs text-muted-foreground font-medium uppercase tracking-wider">Joined</th>
              <th class="px-4 py-3 text-xs text-muted-foreground font-medium uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id" class="border-b border-border/50 hover:bg-accent/30 transition-colors">
              <td class="px-4 py-3 text-muted-foreground font-mono text-xs">{{ u.id }}</td>
              <td class="px-4 py-3">
                <p class="font-medium">{{ u.first_name }} {{ u.last_name }}</p>
                <p class="text-xs text-muted-foreground">
                  {{ u.username ? '@' + u.username : '' }}
                  <span v-if="u.telegram_id">· TG {{ u.telegram_id }}</span>
                </p>
              </td>
              <td class="px-4 py-3">
                <Badge :variant="u.guest_id ? 'secondary' : 'default'" class="text-[10px]">
                  {{ u.guest_id ? '🌐 Guest' : '✈️ Telegram' }}
                </Badge>
              </td>
              <td class="px-4 py-3 text-muted-foreground text-xs">{{ u.created_at.slice(0, 10) }}</td>
              <td class="px-4 py-3 flex gap-1">
                <Button size="sm" variant="ghost" class="text-xs h-7" @click="openProgress(u)">Progress</Button>
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

    <AdminUserProgressDialog v-model:open="progressOpen" :user-id="selectedUserId" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import adminApi from './api'
import Card from '../../components/ui/card.vue'
import Input from '../../components/ui/input.vue'
import Select from '../../components/ui/select.vue'
import Badge from '../../components/ui/badge.vue'
import Button from '../../components/ui/button.vue'
import AdminUserProgressDialog from './AdminUserProgressDialog.vue'
import { Search as SearchIcon } from 'lucide-vue-next'

const users  = ref([])
const total  = ref(0)
const page   = ref(1)
const search = ref('')
const type   = ref('all')
const loading = ref(false)
const perPage = 20

const progressOpen = ref(false)
const selectedUserId = ref(null)

let searchTimer = null
function onSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { page.value = 1; load() }, 400)
}

async function load() {
  loading.value = true
  try {
    const { data } = await adminApi.get('/users', {
      params: { page: page.value, search: search.value, type: type.value },
    })
    users.value = data.users || []
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function openProgress(u) {
  selectedUserId.value = u.id
  progressOpen.value = true
}

async function deleteUser(u) {
  if (!confirm(`Delete ${u.first_name} and all their data? This cannot be undone.`)) return
  await adminApi.delete(`/users/${u.id}`)
  users.value = users.value.filter(x => x.id !== u.id)
}

onMounted(load)
</script>
