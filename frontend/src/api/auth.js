import api from './client'

// auth-service: регистрация и вход (JWT устанавливается в HttpOnly cookie)
export const authService = {
  // POST /auth/register
  register(payload) {
    // payload: { full_name, login, email, phone, password }
    return api.post('/auth/register', payload)
  },

  // POST /auth/login
  login(payload) {
    // payload: { login, password }
    return api.post('/auth/login', payload)
  },
}

export default authService