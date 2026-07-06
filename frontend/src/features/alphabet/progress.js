import { reactive } from 'vue'
import { storage } from '../../services/storage'
import { alphabetApi } from './api'
import { COMBO_GROUPS } from './data/phonicsWords'

const LEGACY_KEY = 'alphabet_completed_levels'
const TOTAL_LEVELS = 4 + COMBO_GROUPS.length

// Shared reactive state: any component reading state.completed re-renders
// automatically once loadAlphabetProgress() resolves, wherever it was
// triggered from — no need to await it at every call site.
//
// Deliberately not cached forever: an admin can unlock/reset a level for a
// stuck kid at any time, and that needs to show up the moment the app next
// checks (route guard, nav re-render) without requiring a full reload —
// same as how the trainer's own levels are re-fetched on every visit.
// loadPromise only dedupes calls that overlap in-flight.
const state = reactive({ completed: [] })
let loadPromise = null

export function loadAlphabetProgress() {
  if (loadPromise) return loadPromise

  loadPromise = (async () => {
    try {
      const res = await alphabetApi.getProgress()
      state.completed = res.data || []

      // One-time migration: older clients kept progress in localStorage
      // only, with nothing on the server. Replay it once — once the server
      // has any entries this condition is false forever after, so it's
      // safe to re-check on every load.
      const legacy = storage.get(LEGACY_KEY)
      if (legacy?.length && state.completed.length === 0) {
        for (const level of legacy) {
          if (!state.completed.includes(level)) {
            await alphabetApi.completeLevel(level)
            state.completed.push(level)
          }
        }
      }
    } catch {
      state.completed = []
    } finally {
      loadPromise = null
    }
  })()

  return loadPromise
}

export function getCompletedAlphabetLevels() {
  return state.completed
}

export function markAlphabetLevelCompleted(level) {
  if (!state.completed.includes(level)) {
    state.completed.push(level)
    alphabetApi.completeLevel(level).catch(() => {})
  }
}

export function isAlphabetCompleted() {
  return state.completed.includes(TOTAL_LEVELS)
}
