<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import ToastContainer from './components/ui/ToastContainer.vue'

const authStore = useAuthStore()
const route = useRoute()
const mobileMenuOpen = ref(false)

// Показываем навигацию только для аутентифицированных пользователей
const showNav = computed(() => authStore.isLoggedIn && route.meta.requiresAuth !== false)

const navLinks = [
  { to: '/', label: 'Дашборд', exact: true },
  { to: '/orders', label: 'Заказы' },
  { to: '/drivers', label: 'Водители' },
  { to: '/autos', label: 'Автомобили' },
  { to: '/profile', label: 'Профиль' },
]

function handleLogout() {
  authStore.logout()
  mobileMenuOpen.value = false
}
</script>

<template>
  <div class="flex min-h-screen flex-col bg-gray-50">
    <ToastContainer />

    <!-- Navbar -->
    <header v-if="showNav" class="sticky top-0 z-50 border-b border-gray-100 bg-white/80 shadow-sm backdrop-blur-lg">
      <nav class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <router-link to="/" class="flex items-center gap-2.5 text-lg font-bold text-indigo-600 transition-colors hover:text-indigo-700">
          <span class="flex h-9 w-9 items-center justify-center rounded-xl bg-indigo-600 text-white shadow-md shadow-indigo-600/30">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
              <path d="M1 3h15v13H1z" />
              <path d="M16 8h4l3 3v5h-7V8z" />
              <circle cx="5.5" cy="18.5" r="2.5" />
              <circle cx="18.5" cy="18.5" r="2.5" />
            </svg>
          </span>
          Freight Management
        </router-link>

        <!-- Desktop nav -->
        <ul class="hidden items-center gap-1 md:flex">
          <li v-for="link in navLinks" :key="link.to">
            <router-link
              :to="link.to"
              class="rounded-xl px-4 py-2 text-sm font-medium text-gray-600 transition-all duration-300 hover:bg-indigo-50 hover:text-indigo-600"
              :class="{ 'bg-indigo-50 text-indigo-600': link.exact ? $route.path === link.to : $route.path.startsWith(link.to) }"
            >
              {{ link.label }}
            </router-link>
          </li>
          <li class="ml-2">
            <button
              @click="handleLogout"
              class="inline-flex items-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-all duration-300 hover:border-rose-200 hover:bg-rose-50 hover:text-rose-600"
            >
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                <polyline points="16 17 21 12 16 7" />
                <line x1="21" y1="12" x2="9" y2="12" />
              </svg>
              Выйти
            </button>
          </li>
        </ul>

        <!-- Mobile menu button -->
        <button
          @click="mobileMenuOpen = !mobileMenuOpen"
          class="flex h-10 w-10 items-center justify-center rounded-xl text-gray-600 transition-colors hover:bg-gray-100 md:hidden"
        >
          <svg v-if="!mobileMenuOpen" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" class="h-6 w-6">
            <line x1="3" y1="6" x2="21" y2="6" />
            <line x1="3" y1="12" x2="21" y2="12" />
            <line x1="3" y1="18" x2="21" y2="18" />
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" class="h-6 w-6">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </nav>

      <!-- Mobile menu -->
      <Transition name="mobile-menu">
        <div v-if="mobileMenuOpen" class="border-t border-gray-100 bg-white px-4 pb-4 pt-2 md:hidden">
          <ul class="flex flex-col gap-1">
            <li v-for="link in navLinks" :key="link.to">
              <router-link
                :to="link.to"
                @click="mobileMenuOpen = false"
                class="block rounded-xl px-4 py-2.5 text-sm font-medium text-gray-600 transition-colors hover:bg-indigo-50 hover:text-indigo-600"
                :class="{ 'bg-indigo-50 text-indigo-600': link.exact ? $route.path === link.to : $route.path.startsWith(link.to) }"
              >
                {{ link.label }}
              </router-link>
            </li>
            <li class="mt-2 border-t border-gray-100 pt-2">
              <button
                @click="handleLogout"
                class="flex w-full items-center gap-2 rounded-xl px-4 py-2.5 text-sm font-medium text-rose-600 transition-colors hover:bg-rose-50"
              >
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4">
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                  <polyline points="16 17 21 12 16 7" />
                  <line x1="21" y1="12" x2="9" y2="12" />
                </svg>
                Выйти
              </button>
            </li>
          </ul>
        </div>
      </Transition>
    </header>

    <main class="mx-auto w-full max-w-7xl flex-1 px-4 py-8 sm:px-6 lg:px-8">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.mobile-menu-enter-active,
.mobile-menu-leave-active {
  transition: all 0.25s ease;
}

.mobile-menu-enter-from,
.mobile-menu-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>