import axios from 'axios'
import { storage } from './storage'

// Shared axios instance — feature `api.js` files build their endpoint methods
// on top of this. Keep it free of any single feature's business logic.
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

export default api
