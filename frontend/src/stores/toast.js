import { defineStore } from 'pinia'
import { ref } from 'vue'

let toastId = 0

/**
 * Toast store — глобальные уведомления (замена alert())
 * Поддерживает типы: success, error, info, warning
 */
export const useToastStore = defineStore('toast', () => {
  const toasts = ref([])

  function show(message, type = 'info', duration = 4000) {
    const id = ++toastId
    toasts.value.push({ id, message, type })
    if (duration > 0) {
      setTimeout(() => remove(id), duration)
    }
    return id
  }

  function success(message, duration) {
    return show(message, 'success', duration)
  }

  function error(message, duration) {
    return show(message, 'error', duration ?? 6000)
  }

  function info(message, duration) {
    return show(message, 'info', duration)
  }

  function warning(message, duration) {
    return show(message, 'warning', duration)
  }

  function remove(id) {
    const index = toasts.value.findIndex((t) => t.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }

  function clear() {
    toasts.value = []
  }

  return {
    toasts,
    show,
    success,
    error,
    info,
    warning,
    remove,
    clear,
  }
})