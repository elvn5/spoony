import { createI18n } from 'vue-i18n'
import en from './locales/en'
import ru from './locales/ru'

export const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('lang') || 'ru',
  fallbackLocale: 'en',
  messages: { en, ru },
})
