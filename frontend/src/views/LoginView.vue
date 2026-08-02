<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import BaseInput from '../components/ui/BaseInput.vue'
import BaseButton from '../components/ui/BaseButton.vue'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToastStore()

const form = ref({
  login: '',
  password: '',
})

async function handleSubmit() {
  const success = await authStore.login(form.value)
  if (success) {
    toast.success('Вы успешно вошли в систему')
    router.push('/')
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-indigo-600 via-indigo-500 to-purple-600 p-4">
    <div class="animate-fade-in w-full max-w-md">
      <div class="mb-8 text-center">
        <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 text-white shadow-lg backdrop-blur">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-8 w-8">
            <path d="M1 3h15v13H1z" />
            <path d="M16 8h4l3 3v5h-7V8z" />
            <circle cx="5.5" cy="18.5" r="2.5" />
            <circle cx="18.5" cy="18.5" r="2.5" />
          </svg>
        </div>
        <h1 class="text-3xl font-bold text-white">Freight Management</h1>
        <p class="mt-1 text-indigo-100">Вход в систему управления грузоперевозками</p>
      </div>

      <div class="rounded-2xl border border-white/20 bg-white p-8 shadow-2xl backdrop-blur">
        <h2 class="mb-6 text-center text-2xl font-bold text-gray-900">Вход в систему</h2>

        <div v-if="authStore.error" class="mb-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
          {{ authStore.error }}
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <BaseInput
            v-model="form.login"
            label="Логин"
            placeholder="Введите логин"
            required
            autocomplete="username"
          />

          <BaseInput
            v-model="form.password"
            label="Пароль"
            type="password"
            placeholder="Введите пароль"
            required
            autocomplete="current-password"
          />

          <BaseButton type="submit" :loading="authStore.loading" size="block">
            Войти
          </BaseButton>
        </form>

        <p class="mt-6 text-center text-sm text-gray-500">
          Нет аккаунта?
          <router-link to="/register" class="font-semibold text-indigo-600 transition-colors hover:text-indigo-700">
            Зарегистрироваться
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>