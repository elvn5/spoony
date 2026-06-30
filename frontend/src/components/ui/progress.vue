<template>
  <ProgressRoot
    :class="cn('relative h-2 w-full overflow-hidden rounded-full bg-secondary', props.class)"
    v-bind="{ modelValue, max }"
  >
    <ProgressIndicator
      class="h-full w-full flex-1 bg-primary transition-all"
      :style="{ transform: `translateX(-${100 - ((modelValue ?? 0) / (max ?? 100)) * 100}%)` }"
    />
  </ProgressRoot>
</template>
<script setup>
import { ProgressIndicator, ProgressRoot } from 'radix-vue'
import { cn } from '@/lib/utils'
const props = defineProps({
  modelValue: Number,
  max:        { type: Number, default: 100 },
  class:      String,
})
</script>
