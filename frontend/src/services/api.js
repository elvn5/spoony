import axios from 'axios'
import { storage } from './storage'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL ? `${import.meta.env.VITE_API_URL}/api` : '/api',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

// Inject auth token
api.interceptors.request.use(config => {
  const token = storage.get('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Handle auth errors
api.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      storage.remove('token')
      storage.remove('user')
      window.location.href = '/'
    }
    return Promise.reject(err)
  }
)

export const authApi = {
  telegramLogin: (initData) => api.post('/auth/telegram-login', { init_data: initData }),
  logout: () => api.post('/auth/logout'),
  getMe: () => api.get('/auth/me'),
  updateProfile: (data) => api.put('/auth/profile', data),
}

export default api
