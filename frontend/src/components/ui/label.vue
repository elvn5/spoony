<template>
  <label
    :class="cn('text-xs font-medium text-muted-foreground leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70', props.class)"
    v-bind="$attrs"
  >
    <slot />
  </label>
</template>
<script setup>
import { cn } from '@/lib/utils'
const props = defineProps({ class: String })
</script>
