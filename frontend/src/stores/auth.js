import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '../api/auth'

// Auth store: управление состоянием аутентификации.
//
// ВАЖНО: auth-service устанавливает JWT в HttpOnly cookie, поэтому
// фронтенд не может прочитать токен напрямую. Сервер также не возвращает
// ID пользователя в ответе /auth/login (только {"message": "logged"}).
// Поэтому store отслеживает только флаг аутентификации.
// ID пользователя для запросов к order-service вводится пользователем в UI.
export const useAuthStore = defineStore('auth', () => {
  const isAuthenticated = ref(localStorage.getItem('isAuthenticated') === 'true')
  const error = ref(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => isAuthenticated.value)

  // Регистрация нового пользователя
  async function register(payload) {
    loading.value = true
    error.value = null
    try {
      await authService.register(payload)
      return true
    } catch (err) {
      error.value = err.message || 'Ошибка регистрации'
      return false
    } finally {
      loading.value = false
    }
  }

  // Вход в систему
  async function login(payload) {
    loading.value = true
    error.value = null
    try {
      await authService.login(payload)
      isAuthenticated.value = true
      localStorage.setItem('isAuthenticated', 'true')
      return true
    } catch (err) {
      error.value = err.message || 'Ошибка входа'
      return false
    } finally {
      loading.value = false
    }
  }

  // Выход из системы (очищаем локальное состояние; cookie истечёт по времени)
  function logout() {
    isAuthenticated.value = false
    localStorage.removeItem('isAuthenticated')
  }

  function clearError() {
    error.value = null
  }

  return {
    isAuthenticated,
    isLoggedIn,
    error,
    loading,
    register,
    login,
    logout,
    clearError,
  }
})