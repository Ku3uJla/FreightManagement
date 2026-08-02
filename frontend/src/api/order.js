import api from './client'

// order-service: управление заказами
export const orderService = {
  // POST /orders — создание заказа
  create(payload) {
    return api.post('/orders', payload)
  },

  // GET /orders — список заказов с фильтрацией и пагинацией
  // params: OrderFilter + page
  list(params = {}) {
    return api.get('/orders', { params })
  },

  // GET /orders/:id — заказ по ID
  getById(id) {
    return api.get(`/orders/${id}`)
  },

  // PUT /orders/:id — обновление полей заказа
  update(id, payload) {
    return api.put(`/orders/${id}`, payload)
  },

  // PATCH /orders/:id/status — обновление статуса заказа
  // body: int (новый статус)
  updateStatus(id, status) {
    return api.patch(`/orders/${id}/status`, status, {
      headers: { 'Content-Type': 'application/json' },
    })
  },

  // POST /orders/:id/manager — назначение менеджера на заказ
  assignManager(id) {
    return api.post(`/orders/${id}/manager`)
  },

  // GET /users/:user_id/orders — заказы пользователя с пагинацией
  getOrdersByUser(userId, page = 1) {
    return api.get(`/users/${userId}/orders`, { params: { page } })
  },

  // GET /drivers/:driver_id/orders — заказы водителя с пагинацией
  getOrdersByDriver(driverId, page = 1) {
    return api.get(`/drivers/${driverId}/orders`, { params: { page } })
  },
}

export default orderService