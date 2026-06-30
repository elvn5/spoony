// Telegram Web App SDK wrapper

export const tg = window.Telegram?.WebApp

export function initTelegram() {
  if (!tg) return
  tg.ready()
  tg.expand()
  tg.enableClosingConfirmation()
}

export function getTelegramInitData() {
  return tg?.initData || ''
}

export function getTelegramUser() {
  return tg?.initDataUnsafe?.user || null
}

export function setMainButton(text, callback) {
  if (!tg) return
  tg.MainButton.setText(text)
  tg.MainButton.onClick(callback)
  tg.MainButton.show()
}

export function hideMainButton() {
  tg?.MainButton?.hide()
}

export function setBackButton(callback) {
  if (!tg) return
  tg.BackButton.onClick(callback)
  tg.BackButton.show()
}

export function hideBackButton() {
  tg?.BackButton?.hide()
}

export function showAlert(message) {
  if (tg) {
    tg.showAlert(message)
  } else {
    alert(message)
  }
}

export function showConfirm(message, callback) {
  if (tg) {
    tg.showConfirm(message, callback)
  } else {
    callback(confirm(message))
  }
}

export function hapticFeedback(type = 'light') {
  tg?.HapticFeedback?.impactOccurred(type)
}

export function getThemeParams() {
  return tg?.themeParams || {}
}

export function isTelegramEnvironment() {
  return !!tg?.initData
}
