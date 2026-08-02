<script setup>
import { ref } from 'vue'
import { userService } from '../api/user'
import BaseBadge from '../components/ui/BaseBadge.vue'
import BaseButton from '../components/ui/BaseButton.vue'

const userId = ref('')
const loading = ref(false)
const error = ref(null)
const user = ref(null)

const updating = ref(false)
const updateMessage = ref('')
const updateError = ref('')
const newRole = ref('')

// Маппинг ролей (на основе API_ROUTES.md)
const roleOptions = [
  { value: 1, label: 'client' },
  { value: 2, label: 'driver' },
  { value: 3, label: 'admin' },
]

async function fetchUser() {
  if (!userId.value) return
  loading.value = true
  error.value = null
  user.value = null
  try {
    const res = await userService.getById(userId.value)
    // API возвращает { message: <User> }
    user.value = res.data?.message || res.data
    if (user.value) {
      // Определяем числовое значение роли для селекта
      const roleMatch = roleOptions.find((r) => r.label === user.value.role)
      newRole.value = roleMatch ? roleMatch.value : ''
    }
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить профиль'
  } finally {
    loading.value = false
  }
}

async function handleUpdateRole() {
  if (!user.value || !newRole.value) return
  updating.value = true
  updateMessage.value = ''
  updateError.value = ''
  try {
    await userService.updateRole(user.value.id, Number(newRole.value))
    const role = roleOptions.find((r) => r.value === Number(newRole.value))
    if (role) {
      user.value.role = role.label
    }
    updateMessage.value = 'Роль успешно обновлена'
  } catch (err) {
    updateError.value = err.message || 'Ошибка обновления роли'
  } finally {
    updating.value = false
  }
}

function formatDateTime(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('ru-RU')
}
</script>

<template>
  <div class="animate-fade-in">
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900 sm:text-3xl">Профиль пользователя</h1>
      <p class="mt-1 text-sm text-gray-500">Просмотр и управление пользователями</p>
    </div>

    <!-- Поиск пользователя по ID -->
    <div class="mb-6 rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
      <div class="mb-4 border-b border-gray-100 pb-4">
        <h2 class="text-lg font-bold text-gray-900">Поиск пользователя</h2>
      </div>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-end">
        <div class="flex-1">
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">ID пользователя</label>
          <input
            v-model="userId"
            type="number"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="Введите ID пользователя"
            @keyup.enter="fetchUser"
          />
        </div>
        <BaseButton @click="fetchUser" :loading="loading" :disabled="!userId">
          Найти
        </BaseButton>
      </div>
    </div>

    <div v-if="error" class="mb-6 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
      {{ error }}
    </div>

    <div v-if="user">
      <div v-if="updateMessage" class="mb-4 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-medium text-emerald-700">
        {{ updateMessage }}
      </div>
      <div v-if="updateError" class="mb-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
        {{ updateError }}
      </div>

      <!-- Информация о пользователе -->
      <div class="mb-6 rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
        <div class="mb-5 flex items-center justify-between border-b border-gray-100 pb-4">
          <h2 class="text-lg font-bold text-gray-900">Информация</h2>
          <BaseBadge variant="info">{{ user.role }}</BaseBadge>
        </div>
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <div class="space-y-3">
            <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
              <span class="text-sm text-gray-500">ID</span>
              <span class="text-sm font-semibold text-gray-900">{{ user.id }}</span>
            </div>
            <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
              <span class="text-sm text-gray-500">ФИО</span>
              <span class="text-sm font-semibold text-gray-900">{{ user.name }}</span>
            </div>
            <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
              <span class="text-sm text-gray-500">Email</span>
              <span class="text-sm font-semibold text-gray-900">{{ user.email }}</span>
            </div>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
              <span class="text-sm text-gray-500">Телефон</span>
              <span class="text-sm font-semibold text-gray-900">{{ user.phone }}</span>
            </div>
            <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
              <span class="text-sm text-gray-500">Дата создания</span>
              <span class="text-sm font-semibold text-gray-900">{{ formatDateTime(user.dateCreate) }}</span>
            </div>
            <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
              <span class="text-sm text-gray-500">Дата обновления</span>
              <span class="text-sm font-semibold text-gray-900">{{ formatDateTime(user.dateUpdate) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Обновление роли -->
      <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
        <div class="mb-5 border-b border-gray-100 pb-4">
          <h2 class="text-lg font-bold text-gray-900">Управление ролью</h2>
        </div>
        <div class="flex flex-col gap-4 sm:flex-row sm:items-end">
          <div class="flex-1">
            <label class="mb-1.5 block text-sm font-semibold text-gray-700">Роль</label>
            <select
              v-model="newRole"
              class="w-full cursor-pointer rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            >
              <option v-for="role in roleOptions" :key="role.value" :value="role.value">
                {{ role.label }}
              </option>
            </select>
          </div>
          <BaseButton @click="handleUpdateRole" :loading="updating">
            Обновить роль
          </BaseButton>
        </div>
      </div>
    </div>

    <div v-else-if="!loading && !error" class="flex flex-col items-center justify-center rounded-2xl border border-gray-100 bg-white py-16 text-center shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
      <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-indigo-50 text-3xl">👤</div>
      <h3 class="text-lg font-semibold text-gray-900">Введите ID пользователя</h3>
      <p class="mt-2 text-sm text-gray-500">Введите ID в поле выше и нажмите «Найти», чтобы просмотреть профиль.</p>
    </div>
  </div>
</template>