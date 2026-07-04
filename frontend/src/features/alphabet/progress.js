import { storage } from '../../services/storage'
import { COMBO_GROUPS } from './data/phonicsWords'

const KEY = 'alphabet_completed_levels'
const TOTAL_LEVELS = 4 + COMBO_GROUPS.length

export function getCompletedAlphabetLevels() {
  return storage.get(KEY) || []
}

export function markAlphabetLevelCompleted(level) {
  const completed = getCompletedAlphabetLevels()
  if (!completed.includes(level)) {
    completed.push(level)
    storage.set(KEY, completed)
  }
}

export function isAlphabetCompleted() {
  return getCompletedAlphabetLevels().includes(TOTAL_LEVELS)
}
