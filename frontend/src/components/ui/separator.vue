<template>
  <SeparatorRoot
    :class="cn('shrink-0 bg-border', orientation === 'horizontal' ? 'h-px w-full' : 'h-full w-px', props.class)"
    v-bind="{ orientation, decorative }"
  />
</template>
<script setup>
import { Separator as SeparatorRoot } from 'radix-vue'
import { cn } from '@/lib/utils'
const props = defineProps({
  orientation: { type: String, default: 'horizontal' },
  decorative:  { type: Boolean, default: true },
  class:       String,
})
</script>
