<template>
  <div class="max-w-md mx-auto px-4 pb-24 pt-4">
    <h1 class="text-xl font-bold mb-5 pt-1">{{ $t('nav.settings') }}</h1>

    <div class="space-y-4">
      <Card class="p-4">
        <p class="text-xs text-muted-foreground uppercase tracking-wider mb-3">
          {{ $t('settings.language') }}
        </p>
        <Select v-model="selectedLang" @change="changeLang">
          <option value="en">{{ $t('languages.en') }}</option>
          <option value="ru">{{ $t('languages.ru') }}</option>
        </Select>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Card from '../components/ui/card.vue'
import Select from '../components/ui/select.vue'

const { locale } = useI18n()
const selectedLang = ref(locale.value)

function changeLang() {
  locale.value = selectedLang.value
  localStorage.setItem('lang', selectedLang.value)
}
</script>
