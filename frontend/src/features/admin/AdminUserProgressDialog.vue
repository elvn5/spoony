<template>
  <Dialog :open="open" @update:open="$emit('update:open', $event)" title="User progress" :description="subtitle">
    <div v-if="loading" class="py-8 text-center text-muted-foreground text-sm">Loading…</div>
    <div v-else-if="!detail" class="py-8 text-center text-destructive text-sm">Failed to load user</div>
    <div v-else class="space-y-4 min-w-0">
      <div class="rounded-xl border border-border bg-secondary/40 p-3 text-sm">
        <p class="font-medium">{{ detail.user.first_name }} {{ detail.user.last_name }}</p>
        <p class="text-xs text-muted-foreground mt-0.5">
          {{ detail.user.guest_id ? '🌐 Guest' : '✈️ Telegram · TG ' + detail.user.telegram_id }}
          <span v-if="detail.user.username"> · @{{ detail.user.username }}</span>
        </p>
      </div>

      <div class="max-h-96 overflow-y-auto space-y-2 pr-1">
        <div
          v-for="lvl in detail.progress" :key="lvl.id"
          class="flex items-center gap-3 rounded-xl border border-border p-3"
        >
          <span class="text-lg shrink-0">{{ lvl.game_type === 'word_build' ? '👑' : '🃏' }}</span>
          <div class="min-w-0 flex-1">
            <p class="text-sm font-medium truncate">{{ lvl.city }}</p>
            <p class="text-xs text-muted-foreground truncate">{{ lvl.title_ru }}</p>
          </div>
          <Select
            class="!w-20 h-9 text-xs"
            :value="edits[lvl.id].stars"
            @change="e => onStarsChange(lvl.id, e.target.value)"
          >
            <option v-for="n in [0,1,2,3]" :key="n" :value="n">{{ n }} ⭐</option>
          </Select>
          <Switch
            :checked="edits[lvl.id].completed"
            @update:checked="v => onToggleCompleted(lvl.id, v)"
          />
          <Button
            size="sm" variant="ghost" class="text-xs h-8 text-destructive hover:text-destructive shrink-0"
            :disabled="!lvl.completed && !lvl.stars"
            @click="resetLevel(lvl.id)"
          >
            Reset
          </Button>
        </div>
      </div>

      <div>
        <p class="text-xs font-semibold text-muted-foreground mb-2 px-1">First Steps (alphabet)</p>
        <div class="max-h-72 overflow-y-auto space-y-2 pr-1">
          <div
            v-for="lvl in alphabetLevels" :key="lvl.id"
            class="flex items-center gap-3 rounded-xl border border-border p-3"
          >
            <span class="text-lg shrink-0">🔤</span>
            <p class="text-sm font-medium truncate min-w-0 flex-1">{{ lvl.label }}</p>
            <Switch
              :checked="alphaEdits[lvl.id]"
              @update:checked="v => onToggleAlphabet(lvl.id, v)"
            />
          </div>
        </div>
      </div>
    </div>
  </Dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import adminApi from './api'
import { COMBO_GROUPS } from '../alphabet/data/phonicsWords'
import Dialog from '../../components/ui/dialog.vue'
import Select from '../../components/ui/select.vue'
import Switch from '../../components/ui/switch.vue'
import Button from '../../components/ui/button.vue'

const alphabetLevels = [
  { id: 1, label: 'Letters A–M' },
  { id: 2, label: 'Letters N–Z' },
  { id: 3, label: 'Match pairs' },
  { id: 4, label: 'Build words' },
  ...COMBO_GROUPS.map((g, i) => ({ id: 5 + i, label: g.combos.join(', ') })),
]

const props = defineProps({
  open: { type: Boolean, default: false },
  userId: { type: [Number, String], default: null },
})
defineEmits(['update:open'])

const loading = ref(false)
const detail = ref(null)
const edits = reactive({})
const alphaEdits = reactive({})

const subtitle = computed(() => detail.value ? `#${detail.value.user.id} — tap a row to edit stars or mark complete` : '')

async function load() {
  if (!props.userId) return
  loading.value = true
  detail.value = null
  try {
    const { data } = await adminApi.get(`/users/${props.userId}`)
    detail.value = data
    for (const lvl of data.progress) {
      edits[lvl.id] = { stars: lvl.stars, completed: lvl.completed }
    }
    for (const lvl of alphabetLevels) {
      alphaEdits[lvl.id] = (data.alphabet_progress || []).includes(lvl.id)
    }
  } catch {
    detail.value = null
  } finally {
    loading.value = false
  }
}

function onToggleAlphabet(levelId, value) {
  alphaEdits[levelId] = value
  if (value) {
    adminApi.put(`/users/${props.userId}/alphabet-progress/${levelId}`).catch(() => {})
  } else {
    adminApi.delete(`/users/${props.userId}/alphabet-progress/${levelId}`).catch(() => {})
  }
}

async function save(levelId) {
  const e = edits[levelId]
  await adminApi.put(`/users/${props.userId}/progress/${levelId}`, {
    stars: Number(e.stars),
    completed: e.completed,
  })
  const lvl = detail.value.progress.find(l => l.id === levelId)
  if (lvl) {
    lvl.stars = Number(e.stars)
    lvl.completed = e.completed
  }
}

function onStarsChange(levelId, value) {
  edits[levelId].stars = Number(value)
  if (Number(value) > 0) edits[levelId].completed = true
  save(levelId)
}

function onToggleCompleted(levelId, value) {
  edits[levelId].completed = value
  save(levelId)
}

async function resetLevel(levelId) {
  await adminApi.delete(`/users/${props.userId}/progress/${levelId}`)
  edits[levelId] = { stars: 0, completed: false }
  const lvl = detail.value.progress.find(l => l.id === levelId)
  if (lvl) {
    lvl.stars = 0
    lvl.completed = false
  }
}

watch(() => [props.open, props.userId], ([isOpen]) => {
  if (isOpen) load()
})
</script>
