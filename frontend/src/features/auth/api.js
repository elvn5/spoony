import api from '../../services/httpClient'

export const authApi = {
  telegramLogin: (initData) => api.post('/auth/telegram-login', { init_data: initData }),
  guestLogin: (guestId, name) => api.post('/auth/guest', { guest_id: guestId, name }),
  logout: () => api.post('/auth/logout'),
  getMe: () => api.get('/auth/me'),
  updateProfile: (data) => api.put('/auth/profile', data),
}
