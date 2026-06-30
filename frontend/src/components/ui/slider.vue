<template>
  <SliderRoot
    :class="cn('relative flex w-full touch-none select-none items-center', props.class)"
    v-bind="forwarded"
  >
    <SliderTrack class="relative h-2 w-full grow overflow-hidden rounded-full bg-secondary">
      <SliderRange class="absolute h-full bg-primary" />
    </SliderTrack>
    <SliderThumb
      class="block h-5 w-5 rounded-full border-2 border-primary bg-background ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 shadow-sm shadow-primary/30"
    />
  </SliderRoot>
</template>
<script setup>
import { SliderRange, SliderRoot, SliderThumb, SliderTrack, useForwardPropsEmits } from 'radix-vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  class:        String,
  modelValue:   Array,
  min:          { type: Number, default: 0 },
  max:          { type: Number, default: 100 },
  step:         { type: Number, default: 1 },
  disabled:     Boolean,
})
const emits = defineEmits(['update:modelValue'])
const forwarded = useForwardPropsEmits(props, emits)
</script>
