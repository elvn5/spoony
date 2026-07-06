import api from '../../services/httpClient'

export const alphabetApi = {
  getProgress: () => api.get('/alphabet-progress'),
  completeLevel: (levelId) => api.post(`/alphabet-progress/${levelId}/complete`),
}
