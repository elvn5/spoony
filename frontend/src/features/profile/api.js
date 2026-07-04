import api from '../../services/httpClient'

export const telegramApi = {
  getBotInfo: () => api.get('/telegram/bot-info'),
}
