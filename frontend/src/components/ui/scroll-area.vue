<template>
  <ScrollAreaRoot :class="cn('relative overflow-hidden', props.class)" v-bind="$attrs">
    <ScrollAreaViewport class="h-full w-full rounded-[inherit]">
      <slot />
    </ScrollAreaViewport>
    <ScrollAreaScrollbar
      class="flex touch-none select-none transition-colors h-full w-2.5 border-l border-l-transparent p-px"
      orientation="vertical"
    >
      <ScrollAreaThumb class="relative flex-1 rounded-full bg-border" />
    </ScrollAreaScrollbar>
  </ScrollAreaRoot>
</template>
<script setup>
import { ScrollAreaRoot, ScrollAreaScrollbar, ScrollAreaThumb, ScrollAreaViewport } from 'radix-vue'
import { cn } from '@/lib/utils'
const props = defineProps({ class: String })
</script>
