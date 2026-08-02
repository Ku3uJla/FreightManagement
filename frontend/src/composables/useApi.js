import { ref } from 'vue'
import { useToastStore } from '../stores/toast'

/**
 * Composable для стандартизации API-запросов
 * Управляет состояниями loading, error, data
 * Автоматически показывает toast при ошибке
 *
 * @param {Function} apiFn - функция API-сервиса
 * @param {Object} options - { immediate, defaultData, showErrorToast }
 * @returns {{ data, loading, error, execute, reset }}
 */
export function useApi(apiFn, options = {}) {
  const {
    immediate = false,
    defaultData = null,
    showErrorToast = true,
  } = options

  const data = ref(defaultData)
  const loading = ref(false)
  const error = ref(null)
  const toast = useToastStore()

  async function execute(...args) {
    loading.value = true
    error.value = null
    try {
      const result = await apiFn(...args)
      data.value = result.data
      return result
    } catch (err) {
      error.value = err.message || 'Ошибка запроса'
      if (showErrorToast) {
        toast.error(error.value)
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  function reset() {
    data.value = defaultData
    loading.value = false
    error.value = null
  }

  if (immediate) {
    execute()
  }

  return {
    data,
    loading,
    error,
    execute,
    reset,
  }
}