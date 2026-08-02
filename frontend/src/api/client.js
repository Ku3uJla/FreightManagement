import axios from 'axios'

// Базовый axios-клиент для всех запросов к бэкенду.
// В Docker-окружении frontend-nginx проксирует /api/* → nginx-gateway:80
// (с rewrite, убирающим префикс /api).
// При локальной разработке (vite dev) проксирование настраивается в vite.config.js.
const api = axios.create({
  baseURL: '/api',
  withCredentials: true, // передаём cookie (JWT) между доменами
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Перехватчик ответов: централизованная обработка ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      // Сервер вернул ошибку
      const { status, data } = error.response
      const message = data?.message || data?.error || 'Ошибка сервера'
      return Promise.reject({ status, message, data })
    } else if (error.request) {
      // Запрос отправлен, но ответа нет
      return Promise.reject({ status: 0, message: 'Нет ответа от сервера' })
    }
    // Ошибка при настройке запроса
    return Promise.reject({ status: -1, message: error.message || 'Неизвестная ошибка' })
  }
)

export default api