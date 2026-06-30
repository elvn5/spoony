<template>
  <component
    :is="as || 'button'"
    :class="cn(buttonVariants({ variant, size }), props.class)"
    v-bind="$attrs"
  >
    <slot />
  </component>
</template>

<script setup>
import { cva } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 active:scale-[0.98]',
  {
    variants: {
      variant: {
        default:     'bg-primary text-primary-foreground hover:bg-primary/90 shadow-sm shadow-primary/30',
        secondary:   'bg-secondary text-secondary-foreground hover:bg-secondary/80',
        ghost:       'hover:bg-accent hover:text-accent-foreground',
        destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/90',
        outline:     'border border-input bg-transparent hover:bg-accent hover:text-accent-foreground',
        link:        'text-primary underline-offset-4 hover:underline p-0 h-auto',
        gradient:    'bg-gradient-to-r from-violet-600 to-indigo-600 text-white hover:from-violet-500 hover:to-indigo-500 shadow-lg shadow-violet-900/40',
      },
      size: {
        default: 'h-10 px-4 py-2',
        sm:      'h-8 rounded-md px-3 text-xs',
        lg:      'h-12 rounded-xl px-8 text-base',
        xl:      'h-14 rounded-xl px-8 text-base',
        icon:    'h-10 w-10',
        'icon-sm': 'h-8 w-8',
      },
    },
    defaultVariants: { variant: 'default', size: 'default' },
  }
)

const props = defineProps({
  variant: { type: String, default: 'default' },
  size:    { type: String, default: 'default' },
  as:      { type: [String, Object], default: 'button' },
  class:   { type: String, default: '' },
})
</script>
