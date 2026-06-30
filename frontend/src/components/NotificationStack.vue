<template>
  <Teleport to="body">
    <div class="fixed top-safe inset-x-0 z-[60] flex flex-col items-center gap-2 pt-4 px-4 pointer-events-none">
      <TransitionGroup name="notif">
        <div
          v-for="n in uiStore.notifications"
          :key="n.id"
          @click="uiStore.dismiss(n.id)"
          :class="[
            'pointer-events-auto flex items-center gap-3 max-w-sm w-full px-4 py-3 rounded-xl shadow-xl border text-sm font-medium cursor-pointer transition-all',
            n.type === 'success' ? 'bg-emerald-950/95 border-emerald-800/60 text-emerald-200' :
            n.type === 'error'   ? 'bg-rose-950/95 border-rose-800/60 text-rose-200' :
                                   'bg-popover/95 border-border text-foreground'
          ]"
        >
          <component :is="iconFor(n.type)" class="h-4 w-4 shrink-0" />
          <span class="flex-1">{{ n.message }}</span>
          <XIcon class="h-3.5 w-3.5 opacity-50 shrink-0" />
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { useUIStore } from '../store/ui'
import { CheckCircle2, XCircle, Info, X as XIcon } from 'lucide-vue-next'

const uiStore = useUIStore()

function iconFor(type) {
  return { success: CheckCircle2, error: XCircle, info: Info }[type] || Info
}
</script>

<style scoped>
.notif-enter-active { transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1); }
.notif-leave-active { transition: all 0.2s ease-in; }
.notif-enter-from   { opacity: 0; transform: translateY(-12px) scale(0.95); }
.notif-leave-to     { opacity: 0; transform: translateY(-8px) scale(0.95); }
.notif-move         { transition: transform 0.2s ease; }
</style>
