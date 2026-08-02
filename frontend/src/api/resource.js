import api from './client'

// resource-service: управление водителями и автомобилями
export const resourceService = {
  // ===== Drivers =====

  // POST /drivers — создание профиля водителя
  createDriver(payload) {
    return api.post('/drivers', payload)
  },

  // POST /drivers/categories — добавление категории водителю
  createDriverCategory(payload) {
    return api.post('/drivers/categories', payload)
  },

  // GET /drivers/:id — водитель по ID
  getDriverById(id) {
    return api.get(`/drivers/${id}`)
  },

  // GET /drivers/:id/categories — категории водителя
  getDriverCategories(id) {
    return api.get(`/drivers/${id}/categories`)
  },

  // GET /drivers — список водителей с фильтром
  // params: { status, category }
  listDrivers(params = {}) {
    return api.get('/drivers', { params })
  },

  // PUT /drivers/:id/status — обновление статуса водителя
  // body: { status: int (1-3) }
  updateDriverStatus(id, status) {
    return api.put(`/drivers/${id}/status`, { status })
  },

  // ===== Autos =====

  // POST /autos — создание автомобиля
  createAuto(payload) {
    return api.post('/autos', payload)
  },

  // GET /autos/:id — автомобиль по ID
  getAutoById(id) {
    return api.get(`/autos/${id}`)
  },

  // GET /autos — список автомобилей с фильтром
  // params: { capacity, lifting_capacity, status }
  listAutos(params = {}) {
    return api.get('/autos', { params })
  },

  // PUT /autos/:id/status — обновление статуса автомобиля
  // body: { status: int (1-3) }
  updateAutoStatus(id, status) {
    return api.put(`/autos/${id}/status`, { status })
  },
}

export default resourceService