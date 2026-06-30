<template>
  <div :class="cn(badgeVariants({ variant }), props.class)">
    <slot />
  </div>
</template>
<script setup>
import { cva } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const badgeVariants = cva(
  'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors',
  {
    variants: {
      variant: {
        default:     'border-transparent bg-primary/20 text-primary',
        secondary:   'border-transparent bg-secondary text-secondary-foreground',
        destructive: 'border-transparent bg-destructive/20 text-destructive',
        outline:     'border-border text-foreground',
        success:     'border-transparent bg-green-500/20 text-green-400',
        warning:     'border-transparent bg-yellow-500/20 text-yellow-400',
      },
    },
    defaultVariants: { variant: 'default' },
  }
)

const props = defineProps({
  variant: { type: String, default: 'default' },
  class:   { type: String, default: '' },
})
</script>
