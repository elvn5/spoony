<template>
  <Teleport to="body">
    <div
      v-if="entry"
      class="fixed z-[70] pointer-events-none px-3 py-2 rounded-xl bg-popover/95 border border-border shadow-xl text-center max-w-[220px]"
      :style="{ top: top + 'px', left: left + 'px', transform: 'translate(-50%, -100%)' }"
    >
      <p class="text-sm font-bold leading-tight whitespace-nowrap">
        {{ word }} <span class="font-normal text-muted-foreground">/{{ entry.ipa }}/</span>
      </p>
      <p class="text-primary font-semibold text-sm leading-tight mt-0.5">{{ entry.ru }}</p>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { DICTIONARY } from '../data/dictionary'

const word = ref('')
const entry = ref(null)
const top = ref(0)
const left = ref(0)
let debounceTimer = null

function normalize(text) {
  return text.trim().toLowerCase().replace(/[^a-z']/g, '')
}

function hide() {
  entry.value = null
}

function updateFromSelection() {
  const selection = window.getSelection()
  const raw = selection?.toString().trim() || ''
  if (!selection || selection.isCollapsed || !raw) {
    hide()
    return
  }

  const found = DICTIONARY[normalize(raw)]
  if (!found) {
    hide()
    return
  }

  const rect = selection.getRangeAt(0).getBoundingClientRect()
  if (!rect.width && !rect.height) {
    hide()
    return
  }

  top.value = Math.max(rect.top - 10, 10)
  left.value = Math.min(Math.max(rect.left + rect.width / 2, 110), window.innerWidth - 110)
  word.value = raw
  entry.value = found
}

function onSelectionChange() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(updateFromSelection, 150)
}

onMounted(() => {
  document.addEventListener('selectionchange', onSelectionChange)
  window.addEventListener('scroll', hide, true)
})
onUnmounted(() => {
  document.removeEventListener('selectionchange', onSelectionChange)
  window.removeEventListener('scroll', hide, true)
  clearTimeout(debounceTimer)
})
</script>
