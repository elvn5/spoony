<template>
  <Card class="p-4 flex flex-col items-center text-center gap-1">
    <span class="text-3xl select-none">{{ icon }}</span>
    <span class="text-2xl font-bold">{{ value }}</span>
    <span class="text-xs text-muted-foreground">{{ label }}</span>
  </Card>
</template>
<script setup>
import Card from './ui/card.vue'
defineProps({ icon: String, label: String, value: [String, Number] })
</script>
