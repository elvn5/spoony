import axios from 'axios'

const adminApi = axios.create({
  baseURL: '/admin/api',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

adminApi.interceptors.request.use(config => {
  const token = sessionStorage.getItem('admin_token')
  if (token) config.headers['X-Admin-Token'] = token
  return config
})

export default adminApi
