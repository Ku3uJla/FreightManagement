<script setup>
import { ref, onMounted } from 'vue'
import { resourceService } from '../api/resource'
import BaseBadge from '../components/ui/BaseBadge.vue'
import BaseModal from '../components/ui/BaseModal.vue'
import BaseButton from '../components/ui/BaseButton.vue'

const loading = ref(true)
const error = ref(null)
const autos = ref([])

// Фильтры
const filters = ref({
  capacity: '',
  lifting_capacity: '',
  status: '',
})

// Модальное окно создания автомобиля
const showCreateModal = ref(false)
const creating = ref(false)
const createError = ref('')
const newAuto = ref({
  status: 1,
  capacity: '',
  lifting_capacity: '',
  number: '',
  required_category: '',
})

const autoStatusMap = {
  1: { label: 'Активен', variant: 'success' },
  2: { label: 'Неактивен', variant: 'warning' },
  3: { label: 'Заблокирован', variant: 'danger' },
}

function getStatusLabel(status) {
  return autoStatusMap[status]?.label || 'Неизвестно'
}

function getStatusVariant(status) {
  return autoStatusMap[status]?.variant || 'info'
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('ru-RU')
}

async function fetchAutos() {
  loading.value = true
  error.value = null
  try {
    const params = {}
    if (filters.value.capacity) params.capacity = filters.value.capacity
    if (filters.value.lifting_capacity) params.lifting_capacity = filters.value.lifting_capacity
    if (filters.value.status) params.status = filters.value.status

    const res = await resourceService.listAutos(params)
    autos.value = res.data || []
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить автомобили'
    autos.value = []
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  fetchAutos()
}

function resetFilters() {
  filters.value = { capacity: '', lifting_capacity: '', status: '' }
  fetchAutos()
}

async function handleCreateAuto() {
  creating.value = true
  createError.value = ''
  try {
    const payload = {
      status: Number(newAuto.value.status),
      capacity: Number(newAuto.value.capacity),
      lifting_capacity: Number(newAuto.value.lifting_capacity),
      number: newAuto.value.number,
      required_category: newAuto.value.required_category,
    }
    await resourceService.createAuto(payload)
    showCreateModal.value = false
    newAuto.value = {
      status: 1,
      capacity: '',
      lifting_capacity: '',
      number: '',
      required_category: '',
    }
    fetchAutos()
  } catch (err) {
    createError.value = err.message || 'Ошибка создания автомобиля'
  } finally {
    creating.value = false
  }
}

async function handleUpdateStatus(auto) {
  const newStatus = prompt(
    `Новый статус для автомобиля #${auto.id} (${auto.number}):\n1 — Активен\n2 — Неактивен\n3 — Заблокирован`,
    auto.status
  )
  if (newStatus === null) return
  const status = Number(newStatus)
  if (![1, 2, 3].includes(status)) {
    alert('Статус должен быть 1, 2 или 3')
    return
  }
  try {
    await resourceService.updateAutoStatus(auto.id, status)
    auto.status = status
  } catch (err) {
    alert(err.message || 'Ошибка обновления статуса')
  }
}

onMounted(() => {
  fetchAutos()
})
</script>

<template>
  <div class="animate-fade-in">
    <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 sm:text-3xl">Автомобили</h1>
        <p class="mt-1 text-sm text-gray-500">Управление автопарком</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700 hover:shadow-lg"
      >
        + Добавить автомобиль
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
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Грузоподъёмность</label>
          <input
            v-model="filters.capacity"
            type="number"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="Точное значение"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Грузоподъёмность (дублирующее)</label>
          <input
            v-model="filters.lifting_capacity"
            type="number"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="Точное значение"
          />
        </div>
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
      <div v-if="autos.length === 0" class="flex flex-col items-center justify-center rounded-2xl border border-gray-100 bg-white py-16 text-center shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
        <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-indigo-50 text-3xl">🚛</div>
        <h3 class="text-lg font-semibold text-gray-900">Автомобилей не найдено</h3>
        <p class="mt-2 text-sm text-gray-500">Добавьте первый автомобиль, нажав на кнопку выше.</p>
      </div>

      <div v-else class="overflow-x-auto rounded-2xl border border-gray-100 bg-white shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
        <table class="w-full border-collapse">
          <thead>
            <tr class="border-b border-gray-100 bg-gray-50">
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">ID</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Номер</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Грузоподъёмность</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Грузоподъёмность (дубл.)</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Требуемая категория</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Статус</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Дата создания</th>
              <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="auto in autos"
              :key="auto.id"
              class="border-b border-gray-50 transition-colors duration-200 last:border-0 hover:bg-indigo-50/40"
            >
              <td class="px-5 py-4 text-sm font-semibold text-gray-900">#{{ auto.id }}</td>
              <td class="px-5 py-4 text-sm font-bold text-gray-900">{{ auto.number }}</td>
              <td class="px-5 py-4 text-sm text-gray-700">{{ auto.capacity }}</td>
              <td class="px-5 py-4 text-sm text-gray-700">{{ auto.lifting_capacity }}</td>
              <td class="px-5 py-4">
                <BaseBadge variant="info">{{ auto.required_category }}</BaseBadge>
              </td>
              <td class="px-5 py-4">
                <BaseBadge :variant="getStatusVariant(auto.status)">
                  {{ getStatusLabel(auto.status) }}
                </BaseBadge>
              </td>
              <td class="px-5 py-4 text-sm text-gray-500">{{ formatDate(auto.date_create) }}</td>
              <td class="px-5 py-4">
                <button
                  @click="handleUpdateStatus(auto)"
                  class="inline-flex items-center gap-1 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 transition-all duration-300 hover:border-amber-300 hover:text-amber-600"
                >
                  Статус
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Модальное окно: создание автомобиля -->
    <BaseModal v-model="showCreateModal" title="Новый автомобиль">
      <div v-if="createError" class="mb-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">{{ createError }}</div>
      <form @submit.prevent="handleCreateAuto" class="space-y-4">
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Гос. номер <span class="text-rose-500">*</span></label>
          <input
            v-model="newAuto.number"
            type="text"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="А123БВ"
            required
            maxlength="10"
          />
          <div class="mt-1 text-xs text-gray-500">До 10 символов</div>
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Грузоподъёмность <span class="text-rose-500">*</span></label>
          <input
            v-model="newAuto.capacity"
            type="number"
            min="1"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="1000"
            required
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Грузоподъёмность (дублирующее) <span class="text-rose-500">*</span></label>
          <input
            v-model="newAuto.lifting_capacity"
            type="number"
            min="1"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="500"
            required
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Требуемая категория прав <span class="text-rose-500">*</span></label>
          <input
            v-model="newAuto.required_category"
            type="text"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="C"
            required
            maxlength="10"
          />
          <div class="mt-1 text-xs text-gray-500">До 10 символов</div>
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Статус</label>
          <select
            v-model="newAuto.status"
            class="w-full cursor-pointer rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
          >
            <option :value="1">Активен</option>
            <option :value="2">Неактивен</option>
            <option :value="3">Заблокирован</option>
          </select>
        </div>
        <BaseButton type="submit" variant="primary" size="block" :loading="creating">
          Создать
        </BaseButton>
      </form>
    </BaseModal>
  </div>
</template>