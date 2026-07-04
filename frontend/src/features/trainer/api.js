import api from '../../services/httpClient'

export const trainerApi = {
  getLevels: () => api.get('/levels'),
  getCards: (levelId) => api.get(`/levels/${levelId}/cards`),
  completeLevel: (levelId, stars) => api.post(`/levels/${levelId}/complete`, { stars }),
  getStats: () => api.get('/stats'),
}
