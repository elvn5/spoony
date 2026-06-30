<template>
  <!-- Native select wrapped for consistent dark styling -->
  <select
    :class="cn(
      'flex h-10 w-full rounded-lg border border-input bg-input/40 px-3 py-2 text-sm text-foreground',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:border-ring',
      'disabled:cursor-not-allowed disabled:opacity-50 transition-colors appearance-none',
      'bg-[url(\'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJNNCA2TDggMTBMMTIgNiIgc3Ryb2tlPSIjOTU4ZGFkIiBzdHJva2Utd2lkdGg9IjEuNSIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIi8+PC9zdmc+\')] bg-no-repeat bg-[right_0.75rem_center]',
      props.class
    )"
    v-bind="$attrs"
  >
    <slot />
  </select>
</template>
<script setup>
import { cn } from '@/lib/utils'
const props = defineProps({ class: String })
</script>
