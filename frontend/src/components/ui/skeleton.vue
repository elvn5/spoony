<template>
  <div :class="cn('animate-pulse rounded-md bg-muted', props.class)" v-bind="$attrs" />
</template>
<script setup>
import { cn } from '@/lib/utils'
const props = defineProps({ class: String })
</script>
