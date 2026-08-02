import api from './client'

// user-service: управление профилем пользователя
export const userService = {
  // GET /user/:id — профиль пользователя по ID
  getById(id) {
    return api.get(`/user/${id}`)
  },

  // PUT /user/:id/role — обновление роли пользователя
  // body: int (новая роль)
  updateRole(id, role) {
    return api.put(`/user/${id}/role`, role, {
      headers: { 'Content-Type': 'application/json' },
    })
  },
}

export default userService