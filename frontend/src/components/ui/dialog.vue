<template>
  <DialogRoot v-bind="{ open, defaultOpen, modal }" @update:open="$emit('update:open', $event)">
    <!-- Trigger slot -->
    <DialogTrigger v-if="$slots.trigger" as-child>
      <slot name="trigger" />
    </DialogTrigger>

    <DialogPortal>
      <DialogOverlay
        class="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"
      />
      <DialogContent
        :class="cn(
          'fixed z-50 grid w-full max-w-md gap-4 border border-border bg-popover p-6 shadow-xl',
          'data-[state=open]:animate-in data-[state=closed]:animate-out',
          'data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0',
          'data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom',
          'bottom-0 left-0 right-0 rounded-t-2xl sm:bottom-auto sm:left-1/2 sm:-translate-x-1/2 sm:top-1/2 sm:-translate-y-1/2 sm:rounded-2xl',
          props.class
        )"
      >
        <div class="flex items-center justify-between">
          <DialogTitle v-if="title" class="text-lg font-semibold text-foreground">{{ title }}</DialogTitle>
          <slot name="title" />
          <DialogClose
            class="rounded-md p-1.5 text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors ml-auto"
          >
            <XIcon class="h-4 w-4" />
          </DialogClose>
        </div>
        <DialogDescription v-if="description" class="text-sm text-muted-foreground">{{ description }}</DialogDescription>
        <slot />
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<script setup>
import {
  DialogClose, DialogContent, DialogDescription, DialogOverlay,
  DialogPortal, DialogRoot, DialogTitle, DialogTrigger,
} from 'radix-vue'
import { X as XIcon } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const props = defineProps({
  open:        Boolean,
  defaultOpen: Boolean,
  modal:       { type: Boolean, default: true },
  title:       String,
  description: String,
  class:       String,
})
defineEmits(['update:open'])
</script>
