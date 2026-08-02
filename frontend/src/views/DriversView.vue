<script setup>
import { ref, onMounted } from 'vue'
import { resourceService } from '../api/resource'
import BaseBadge from '../components/ui/BaseBadge.vue'
import BaseModal from '../components/ui/BaseModal.vue'
import BaseButton from '../components/ui/BaseButton.vue'

const loading = ref(true)
const error = ref(null)
const drivers = ref([])

// Фильтры
const filters = ref({
  status: '',
  category: '',
})

// Модальное окно создания водителя
const showCreateModal = ref(false)
const creating = ref(false)
const createError = ref('')
const newDriver = ref({ user_id: '' })

// Категории водителя (для просмотра)
const showCategoriesModal = ref(false)
const selectedDriver = ref(null)
const categoriesLoading = ref(false)
const categoriesError = ref('')
const driverCategories = ref([])

// Добавление категории
const showCategoryModal = ref(false)
const addingCategory = ref(false)
const categoryError = ref('')
const newCategory = ref({ DriverID: '', category: '' })

const driverStatusMap = {
  1: { label: 'Активен', variant: 'success' },
  2: { label: 'Неактивен', variant: 'warning' },
  3: { label: 'Заблокирован', variant: 'danger' },
}

function getStatusLabel(status) {
  return driverStatusMap[status]?.label || 'Неизвестно'
}

function getStatusVariant(status) {
  return driverStatusMap[status]?.variant || 'info'
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('ru-RU')
}

async function fetchDrivers() {
  loading.value = true
  error.value = null
  try {
    const params = {}
    if (filters.value.status) params.status = filters.value.status
    if (filters.value.category) params.category = filters.value.category

    const res = await resourceService.listDrivers(params)
    drivers.value = res.data || []
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить водителей'
    drivers.value = []
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  fetchDrivers()
}

function resetFilters() {
  filters.value = { status: '', category: '' }
  fetchDrivers()
}

async function handleCreateDriver() {
  creating.value = true
  createError.value = ''
  try {
    await resourceService.createDriver({ user_id: Number(newDriver.value.user_id) })
    showCreateModal.value = false
    newDriver.value = { user_id: '' }
    fetchDrivers()
  } catch (err) {
    createError.value = err.message || 'Ошибка создания водителя'
  } finally {
    creating.value = false
  }
}

async function viewCategories(driver) {
  selectedDriver.value = driver
  showCategoriesModal.value = true
  categoriesLoading.value = true
  categoriesError.value = ''
  driverCategories.value = []
  try {
    const res = await resourceService.getDriverCategories(driver.id)
    driverCategories.value = res.data || []
  } catch (err) {
    categoriesError.value = err.message || 'Не удалось загрузить категории'
  } finally {
    categoriesLoading.value = false
  }
}

function closeCategories() {
  showCategoriesModal.value = false
  selectedDriver.value = null
  driverCategories.value = []
}

function openCategoryModal(driverId) {
  newCategory.value = { DriverID: driverId, category: '' }
  showCategoryModal.value = true
}

async function handleAddCategory() {
  addingCategory.value = true
  categoryError.value = ''
  try {
    await resourceService.createDriverCategory({
      DriverID: Number(newCategory.value.DriverID),
      category: newCategory.value.category,
    })
    showCategoryModal.value = false
    // Обновляем список категорий, если открыто окно просмотра
    if (selectedDriver.value && selectedDriver.value.id === Number(newCategory.value.DriverID)) {
      viewCategories(selectedDriver.value)
    }
  } catch (err) {
    categoryError.value = err.message || 'Ошибка добавления категории'
  } finally {
    addingCategory.value = false
  }
}

async function handleUpdateStatus(driver) {
  const newStatus = prompt(
    `Новый статус для водителя #${driver.id}:\n1 — Активен\n2 — Неактивен\n3 — Заблокирован`,
    driver.status
  )
  if (newStatus === null) return
  const status = Number(newStatus)
  if (![1, 2, 3].includes(status)) {
    alert('Статус должен быть 1, 2 или 3')
    return
  }
  try {
    await resourceService.updateDriverStatus(driver.id, status)
    driver.status = status
  } catch (err) {
    alert(err.message || 'Ошибка обновления статуса')
  }
}

onMounted(() => {
  fetchDrivers()
})
</script>

<template>
  <div class="animate-fade-in">
    <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 sm:text-3xl">Водители</h1>
        <p class="mt-1 text-sm text-gray-500">Управление профилями водителей</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700 hover:shadow-lg"
      >
        + Добавить водителя
      </button>
    </div>

    <!-- Фильтры -->
    <div class="mb-6 rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
      <div class="mb-4 flex items-center justify-between border-b border-gray-100 pb-4">
        <h2 class="text-lg font-bold text-gray-900">Фильтры</h2>
        <button
          @click="resetFilters"
          class="inline-flex items-center gap-1 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-600 transition-all duration-300 hover:border-rose-200 hover:bg-rose-50 hover:text-rose-600"
        >
          Сбросить
        </button>
      </div>
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Статус</label>
          <select
            v-model="filters.status"
            class="w-full cursor-pointer rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
          >
            <option value="">Все</option>
            <option value="1">Активен</option>
            <option value="2">Неактивен</option>
            <option value="3">Заблокирован</option>
          </select>
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Категория</label>
          <input
            v-model="filters.category"
            type="text"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="Например: C"
          />
        </div>
      </div>
      <button
        @click="applyFilters"
        class="mt-4 inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700"
      >
        Применить
      </button>
    </div>

    <div v-if="error" class="mb-6 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
      {{ error }}
    </div>

    <!-- Skeleton loading -->
    <div v-if="loading" class="space-y-3">
      <div v-for="i in 5" :key="i" class="skeleton h-16 rounded-2xl"></div>
    </div>

    <div v-else>
      <div v-if="drivers.length === 0" class="flex flex-col items-center justify-center rounded-2xl border border-gray-100 bg-white py-16 text-center shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
        <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-indigo-50 text-3xl">👤</div>
        <h3 class="text-lg font-semibold text-gray-900">Водителей не найдено</h3>
        <p class="mt-2 text-sm text-gray-500">Добавьте первого водителя, нажав на кнопку выше.</p>
      </div>

      <div v-else class="overflow-x-auto rounded-2xl border border-gray-100 bg-white shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
        <table class="w-full border-collapse">
          <thead>
            <tr class="border-b border-gray-100 bg-gray-50">
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">ID</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">ID пользователя</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Статус</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Дата создания</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Дата обновления</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="driver in drivers"
              :key="driver.id"
              class="border-b border-gray-50 transition-colors duration-200 last:border-0 hover:bg-indigo-50/40"
            >
              <td class="px-5 py-4 text-sm font-semibold text-gray-900">#{{ driver.id }}</td>
              <td class="px-5 py-4 text-sm text-gray-700">{{ driver.user_id }}</td>
              <td class="px-5 py-4">
                <BaseBadge :variant="getStatusVariant(driver.status)">
                  {{ getStatusLabel(driver.status) }}
                </BaseBadge>
              </td>
              <td class="px-5 py-4 text-sm text-gray-500">{{ formatDate(driver.date_create) }}</td>
              <td class="px-5 py-4 text-sm text-gray-500">{{ formatDate(driver.date_update) }}</td>
              <td class="px-5 py-4">
                <div class="flex gap-2">
                  <button
                    @click="viewCategories(driver)"
                    class="inline-flex items-center gap-1 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600"
                  >
                    Категории
                  </button>
                  <button
                    @click="handleUpdateStatus(driver)"
                    class="inline-flex items-center gap-1 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 transition-all duration-300 hover:border-amber-300 hover:text-amber-600"
                  >
                    Статус
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Модальное окно: создание водителя -->
    <BaseModal v-model="showCreateModal" title="Новый водитель">
      <div v-if="createError" class="mb-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">{{ createError }}</div>
      <form @submit.prevent="handleCreateDriver" class="space-y-4">
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">ID пользователя</label>
          <input
            v-model="newDriver.user_id"
            type="number"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="ID существующего пользователя"
            required
          />
          <div class="mt-1 text-xs text-gray-500">Введите ID пользователя, который станет водителем</div>
        </div>
        <BaseButton type="submit" variant="primary" size="block" :loading="creating">
          Создать
        </BaseButton>
      </form>
    </BaseModal>

    <!-- Модальное окно: категории водителя -->
    <BaseModal v-model="showCategoriesModal" :title="selectedDriver ? `Категории водителя #${selectedDriver.id}` : ''">
      <div v-if="categoriesError" class="mb-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">{{ categoriesError }}</div>
      <div v-if="categoriesLoading" class="flex items-center justify-center gap-3 py-8 text-gray-500">
        <svg class="h-6 w-6 animate-spin text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
      </div>
      <div v-else>
        <div v-if="driverCategories.length === 0" class="py-8 text-center text-sm text-gray-500">
          У этого водителя пока нет категорий.
        </div>
        <div v-else class="space-y-2">
          <div
            v-for="cat in driverCategories"
            :key="cat.id"
            class="flex items-center justify-between rounded-xl bg-gray-50 px-4 py-3 transition-colors hover:bg-indigo-50/60"
          >
            <span class="text-sm font-semibold text-gray-900">#{{ cat.id }}</span>
            <BaseBadge variant="info">{{ cat.category }}</BaseBadge>
            <span class="text-sm text-gray-500">{{ formatDate(cat.date_create) }}</span>
          </div>
        </div>
        <button
          @click="openCategoryModal(selectedDriver.id)"
          class="mt-4 w-full inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700"
        >
          + Добавить категорию
        </button>
      </div>
    </BaseModal>

    <!-- Модальное окно: добавление категории -->
    <BaseModal v-model="showCategoryModal" title="Добавить категорию">
      <div v-if="categoryError" class="mb-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">{{ categoryError }}</div>
      <form @submit.prevent="handleAddCategory" class="space-y-4">
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Категория прав</label>
          <input
            v-model="newCategory.category"
            type="text"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="Например: C, CE, D"
            required
            maxlength="10"
          />
        </div>
        <BaseButton type="submit" variant="primary" size="block" :loading="addingCategory">
          Добавить
        </BaseButton>
      </form>
    </BaseModal>
  </div>
</template>