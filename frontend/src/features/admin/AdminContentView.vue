<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold">Home Feed Content</h1>
      <Button size="sm" variant="gradient" @click="openCreate">+ New Post</Button>
    </div>

    <Card class="overflow-hidden">
      <div v-if="loading" class="p-6 text-center text-muted-foreground text-sm">Loading…</div>
      <div v-else-if="!posts.length" class="p-6 text-center text-muted-foreground text-sm">No posts yet</div>
      <div v-else class="divide-y divide-border">
        <div v-for="p in posts" :key="p.id" class="flex items-start gap-3 p-4">
          <span class="text-2xl shrink-0">{{ p.avatar }}</span>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2 flex-wrap">
              <p class="font-semibold text-sm">{{ p.title }}</p>
              <Badge variant="secondary" class="text-[10px]">{{ p.category }}</Badge>
            </div>
            <p class="text-xs text-muted-foreground mt-0.5">{{ p.author }} · {{ p.created_at.slice(0, 10) }} · ❤️ {{ p.likes }}</p>
            <p class="text-sm text-muted-foreground mt-1.5 line-clamp-2">{{ p.body }}</p>
          </div>
          <div class="flex gap-1 shrink-0">
            <Button size="sm" variant="ghost" class="text-xs h-7" @click="openEdit(p)">Edit</Button>
            <Button size="sm" variant="ghost" class="text-xs h-7 text-destructive hover:text-destructive"
              @click="deletePost(p)">Delete</Button>
          </div>
        </div>
      </div>
    </Card>

    <Dialog v-model:open="formOpen" :title="editing ? 'Edit post' : 'New post'">
      <form @submit.prevent="save" class="space-y-3">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <Label class="mb-1 block">Author</Label>
            <Input v-model="form.author" required />
          </div>
          <div>
            <Label class="mb-1 block">Avatar emoji</Label>
            <Input v-model="form.avatar" placeholder="🥄" />
          </div>
        </div>
        <div>
          <Label class="mb-1 block">Title</Label>
          <Input v-model="form.title" required />
        </div>
        <div>
          <Label class="mb-1 block">Body</Label>
          <Textarea v-model="form.body" rows="4" required />
        </div>
        <div class="grid grid-cols-3 gap-3">
          <div>
            <Label class="mb-1 block">Image emoji</Label>
            <Input v-model="form.image" placeholder="🎉" />
          </div>
          <div>
            <Label class="mb-1 block">Category</Label>
            <Input v-model="form.category" placeholder="news" />
          </div>
          <div>
            <Label class="mb-1 block">Likes</Label>
            <Input v-model.number="form.likes" type="number" min="0" />
          </div>
        </div>
        <p v-if="error" class="text-destructive text-xs">{{ error }}</p>
        <div class="flex justify-end gap-2 pt-2">
          <Button type="button" variant="ghost" @click="formOpen = false">Cancel</Button>
          <Button type="submit" variant="gradient" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</Button>
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import adminApi from './api'
import Card from '../../components/ui/card.vue'
import Dialog from '../../components/ui/dialog.vue'
import Input from '../../components/ui/input.vue'
import Textarea from '../../components/ui/textarea.vue'
import Label from '../../components/ui/label.vue'
import Badge from '../../components/ui/badge.vue'
import Button from '../../components/ui/button.vue'

const posts = ref([])
const loading = ref(false)
const formOpen = ref(false)
const saving = ref(false)
const error = ref('')
const editing = ref(null)

function emptyForm() {
  return { author: '', avatar: '🥄', title: '', body: '', image: '', category: 'news', likes: 0 }
}
const form = reactive(emptyForm())

async function load() {
  loading.value = true
  try {
    const { data } = await adminApi.get('/news')
    posts.value = data || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  Object.assign(form, emptyForm())
  error.value = ''
  formOpen.value = true
}

function openEdit(p) {
  editing.value = p
  Object.assign(form, {
    author: p.author, avatar: p.avatar, title: p.title, body: p.body,
    image: p.image, category: p.category, likes: p.likes,
  })
  error.value = ''
  formOpen.value = true
}

async function save() {
  saving.value = true
  error.value = ''
  try {
    if (editing.value) {
      await adminApi.put(`/news/${editing.value.id}`, form)
    } else {
      await adminApi.post('/news', form)
    }
    formOpen.value = false
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || 'Failed to save post'
  } finally {
    saving.value = false
  }
}

async function deletePost(p) {
  if (!confirm(`Delete "${p.title}"?`)) return
  await adminApi.delete(`/news/${p.id}`)
  posts.value = posts.value.filter(x => x.id !== p.id)
}

onMounted(load)
</script>
