import api from '../../services/httpClient'

export const newsApi = {
  getNews: () => api.get('/news'),
}
