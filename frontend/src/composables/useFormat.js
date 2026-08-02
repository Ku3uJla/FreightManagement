import { format, formatDistanceToNow } from '../utils/date'

/**
 * Composable для форматирования дат и чисел
 */
export function useFormat() {
  function formatDate(dateStr) {
    if (!dateStr) return '—'
    return format(dateStr, 'dd.MM.yyyy')
  }

  function formatDateTime(dateStr) {
    if (!dateStr) return '—'
    return format(dateStr, 'dd.MM.yyyy HH:mm')
  }

  function formatRelative(dateStr) {
    if (!dateStr) return '—'
    return formatDistanceToNow(dateStr)
  }

  function formatPrice(value) {
    if (value === null || value === undefined || value === '') return '—'
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
      minimumFractionDigits: 2,
    }).format(value)
  }

  function formatNumber(value) {
    if (value === null || value === undefined || value === '') return '—'
    return new Intl.NumberFormat('ru-RU').format(value)
  }

  return {
    formatDate,
    formatDateTime,
    formatRelative,
    formatPrice,
    formatNumber,
  }
}