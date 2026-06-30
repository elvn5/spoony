<template>
  <AvatarRoot :class="cn('relative flex shrink-0 overflow-hidden rounded-full', sizeClass, props.class)">
    <AvatarImage
      v-if="src"
      :src="src"
      :alt="alt"
      class="aspect-square h-full w-full object-cover"
    />
    <AvatarFallback
      :class="cn('flex h-full w-full items-center justify-center rounded-full bg-primary/20 text-primary font-semibold', fallbackClass)"
    >
      {{ fallback || alt?.[0]?.toUpperCase() || '?' }}
    </AvatarFallback>
  </AvatarRoot>
</template>
<script setup>
import { AvatarFallback, AvatarImage, AvatarRoot } from 'radix-vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  src:           String,
  alt:           String,
  fallback:      String,
  size:          { type: String, default: 'md' },
  fallbackClass: String,
  class:         String,
})

const sizeClass = {
  sm:  'h-8 w-8 text-xs',
  md:  'h-10 w-10 text-sm',
  lg:  'h-16 w-16 text-xl',
  xl:  'h-20 w-20 text-2xl',
}[props.size] ?? 'h-10 w-10 text-sm'
</script>
