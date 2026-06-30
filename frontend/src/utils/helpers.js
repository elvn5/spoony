import { format, formatDistanceToNow, parseISO } from 'date-fns'

export function formatDate(dateStr) {
  if (!dateStr) return ''
  try {
    return format(parseISO(dateStr), 'MMM d, yyyy')
  } catch {
    return dateStr
  }
}

export function formatRelative(dateStr) {
  if (!dateStr) return ''
  try {
    return formatDistanceToNow(parseISO(dateStr), { addSuffix: true })
  } catch {
    return dateStr
  }
}

export const MOODS = [
  { value: 'peaceful',   label: 'Peaceful',   color: 'bg-blue-500/20 text-blue-300',   emoji: '😌' },
  { value: 'excited',    label: 'Excited',    color: 'bg-yellow-500/20 text-yellow-300', emoji: '😄' },
  { value: 'anxious',    label: 'Anxious',    color: 'bg-orange-500/20 text-orange-300', emoji: '😰' },
  { value: 'sad',        label: 'Sad',        color: 'bg-indigo-500/20 text-indigo-300', emoji: '😢' },
  { value: 'confused',   label: 'Confused',   color: 'bg-purple-500/20 text-purple-300', emoji: '😕' },
  { value: 'happy',      label: 'Happy',      color: 'bg-green-500/20 text-green-300',  emoji: '😊' },
  { value: 'scared',     label: 'Scared',     color: 'bg-red-500/20 text-red-300',      emoji: '😨' },
  { value: 'nostalgic',  label: 'Nostalgic',  color: 'bg-pink-500/20 text-pink-300',    emoji: '🥹' },
  { value: 'mysterious', label: 'Mysterious', color: 'bg-violet-500/20 text-violet-300', emoji: '🔮' },
]

export function getMoodStyle(mood) {
  return MOODS.find(m => m.value === mood) || { color: 'bg-white/10 text-white/70', emoji: '💤', label: mood }
}

export function lucidityStars(level) {
  return Math.round((level / 10) * 5)
}

export function truncate(str, len = 120) {
  if (!str || str.length <= len) return str
  return str.slice(0, len).trim() + '…'
}
