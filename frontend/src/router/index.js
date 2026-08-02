import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/DashboardView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/RegisterView.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/orders',
    name: 'orders',
    component: () => import('../views/OrdersView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/orders/:id',
    name: 'order-detail',
    component: () => import('../views/OrderDetailView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/orders/new',
    name: 'order-create',
    component: () => import('../views/OrderCreateView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/drivers',
    name: 'drivers',
    component: () => import('../views/DriversView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/autos',
    name: 'autos',
    component: () => import('../views/AutosView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('../views/ProfileView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('../views/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Навигационный guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next({ name: 'login' })
  } else if (to.meta.guestOnly && authStore.isLoggedIn) {
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router