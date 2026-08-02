/**
 * Composable для централизованного управления маппингами статусов
 * (заказы, водители, автомобили, роли пользователей)
 */

// Статусы заказов (из API_ROUTES.md)
export const ORDER_STATUSES = {
  1: { label: 'Новый', class: 'badge-info' },
  2: { label: 'В работе', class: 'badge-warning' },
  3: { label: 'Завершён', class: 'badge-success' },
  4: { label: 'Зарезервировано', class: 'badge-info' },
  5: { label: 'Зарезервировано', class: 'badge-info' },
  '-1': { label: 'Отменён', class: 'badge-danger' },
}

// Статусы водителей и автомобилей (из API_ROUTES.md)
export const RESOURCE_STATUSES = {
  1: { label: 'Активен', class: 'badge-success' },
  2: { label: 'Неактивен', class: 'badge-warning' },
  3: { label: 'Заблокирован', class: 'badge-danger' },
}

// Роли пользователей
export const USER_ROLES = [
  { value: 1, label: 'client' },
  { value: 2, label: 'driver' },
  { value: 3, label: 'admin' },
]

export function useStatus() {
  function getOrderStatusLabel(status) {
    return ORDER_STATUSES[status]?.label || 'Неизвестно'
  }

  function getOrderStatusClass(status) {
    return ORDER_STATUSES[status]?.class || 'badge-info'
  }

  function getResourceStatusLabel(status) {
    return RESOURCE_STATUSES[status]?.label || 'Неизвестно'
  }

  function getResourceStatusClass(status) {
    return RESOURCE_STATUSES[status]?.class || 'badge-info'
  }

  function getRoleLabel(roleString) {
    const role = USER_ROLES.find((r) => r.label === roleString)
    return role ? role.label : roleString
  }

  function getRoleValue(roleString) {
    const role = USER_ROLES.find((r) => r.label === roleString)
    return role ? role.value : ''
  }

  return {
    ORDER_STATUSES,
    RESOURCE_STATUSES,
    USER_ROLES,
    getOrderStatusLabel,
    getOrderStatusClass,
    getResourceStatusLabel,
    getResourceStatusClass,
    getRoleLabel,
    getRoleValue,
  }
}