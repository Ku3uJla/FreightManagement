/**
 * Утилиты для работы с датами без внешних зависимостей
 */

const MONTHS = [
  'января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
  'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря',
]

const MONTHS_SHORT = [
  'янв', 'фев', 'мар', 'апр', 'май', 'июн',
  'июл', 'авг', 'сен', 'окт', 'ноя', 'дек',
]

/**
 * Форматирование даты по шаблону
 * Поддерживаемые токены: dd, MM, yyyy, HH, mm
 * @param {string|Date} date - дата
 * @param {string} pattern - шаблон (например 'dd.MM.yyyy' или 'dd.MM.yyyy HH:mm')
 * @returns {string}
 */
export function format(date, pattern = 'dd.MM.yyyy') {
  const d = typeof date === 'string' ? new Date(date) : date
  if (isNaN(d.getTime())) return '—'

  const day = String(d.getDate()).padStart(2, '0')
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const year = d.getFullYear()
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')

  return pattern
    .replace('dd', day)
    .replace('MM', month)
    .replace('yyyy', year)
    .replace('HH', hours)
    .replace('mm', minutes)
}

/**
 * Относительное время на русском ("2 дня назад", "только что")
 * @param {string|Date} date - дата
 * @returns {string}
 */
export function formatDistanceToNow(date) {
  const d = typeof date === 'string' ? new Date(date) : date
  if (isNaN(d.getTime())) return '—'

  const now = new Date()
  const diffMs = now - d
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)
  const diffMonth = Math.floor(diffDay / 30)
  const diffYear = Math.floor(diffDay / 365)

  if (diffSec < 60) return 'только что'
  if (diffMin < 60) return `${diffMin} ${plural(diffMin, 'минуту', 'минуты', 'минут')} назад`
  if (diffHour < 24) return `${diffHour} ${plural(diffHour, 'час', 'часа', 'часов')} назад`
  if (diffDay < 30) return `${diffDay} ${plural(diffDay, 'день', 'дня', 'дней')} назад`
  if (diffMonth < 12) return `${diffMonth} ${plural(diffMonth, 'месяц', 'месяца', 'месяцев')} назад`
  return `${diffYear} ${plural(diffYear, 'год', 'года', 'лет')} назад`
}

/**
 * Русская плюрализация
 * @param {number} n - число
 * @param {string} one - форма для 1
 * @param {string} few - форма для 2-4
 * @param {string} many - форма для 5+
 * @returns {string}
 */
function plural(n, one, few, many) {
  const mod10 = n % 10
  const mod100 = n % 100
  if (mod10 === 1 && mod100 !== 11) return one
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) return few
  return many
}