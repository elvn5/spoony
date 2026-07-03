import { storage } from './storage'

const KEY = 'alphabet_completed_levels'

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
