<template>
  <input
    :value="modelValue"
    @input="$emit('update:modelValue', $event.target.value)"
    :class="cn(
      'flex h-10 w-full rounded-lg border border-input bg-input/40 px-3 py-2 text-sm text-foreground',
      'placeholder:text-muted-foreground',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:border-ring',
      'disabled:cursor-not-allowed disabled:opacity-50',
      'transition-colors',
      props.class
    )"
    v-bind="$attrs"
  />
</template>
<script setup>
import { cn } from '@/lib/utils'
const props = defineProps({ class: String, modelValue: [String, Number] })
defineEmits(['update:modelValue'])
</script>
