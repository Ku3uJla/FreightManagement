<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { orderService } from '../api/order'
import BaseBadge from '../components/ui/BaseBadge.vue'

const route = useRoute()
const router = useRouter()
const orderId = route.params.id

const loading = ref(true)
const error = ref(null)
const order = ref(null)

const updating = ref(false)
const updateMessage = ref('')
const updateError = ref('')
const newStatus = ref('')

const orderStatusMap = {
  1: { label: 'Новый', variant: 'info' },
  2: { label: 'В работе', variant: 'warning' },
  3: { label: 'Завершён', variant: 'success' },
  '-1': { label: 'Отменён', variant: 'danger' },
}

function getStatusLabel(status) {
  return orderStatusMap[status]?.label || 'Неизвестно'
}

function getStatusVariant(status) {
  return orderStatusMap[status]?.variant || 'info'
}

function formatDateTime(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('ru-RU')
}

async function fetchOrder() {
  loading.value = true
  error.value = null
  try {
    const res = await orderService.getById(orderId)
    // API возвращает заказ напрямую или в { message: <Order> }
    order.value = res.data?.id ? res.data : res.data?.message
    if (order.value) {
      newStatus.value = order.value.status
    }
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить заказ'
  } finally {
    loading.value = false
  }
}

async function handleUpdateStatus() {
  if (newStatus.value === '' || newStatus.value === null) return
  updating.value = true
  updateMessage.value = ''
  updateError.value = ''
  try {
    await orderService.updateStatus(orderId, Number(newStatus.value))
    updateMessage.value = 'Статус успешно обновлён'
    if (order.value) {
      order.value.status = Number(newStatus.value)
    }
  } catch (err) {
    updateError.value = err.message || 'Ошибка обновления статуса'
  } finally {
    updating.value = false
  }
}

async function handleAssignManager() {
  updating.value = true
  updateMessage.value = ''
  updateError.value = ''
  try {
    await orderService.assignManager(orderId)
    updateMessage.value = 'Менеджер назначен'
  } catch (err) {
    updateError.value = err.message || 'Ошибка назначения менеджера'
  } finally {
    updating.value = false
  }
}

onMounted(() => {
  fetchOrder()
})
</script>

<template>
  <div class="animate-fade-in">
    <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-3">
        <button
          @click="router.back()"
          class="flex h-10 w-10 items-center justify-center rounded-xl border border-gray-200 bg-white text-gray-600 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600"
        >
          ←
        </button>
        <div>
          <h1 class="text-2xl font-bold text-gray-900 sm:text-3xl">Заказ #{{ orderId }}</h1>
          <p class="mt-1 text-sm text-gray-500">Детальная информация о заказе</p>
        </div>
      </div>
    </div>

    <!-- Skeleton loading -->
    <div v-if="loading" class="space-y-4">
      <div class="skeleton h-64 rounded-2xl"></div>
      <div class="skeleton h-40 rounded-2xl"></div>
    </div>

    <div v-else>
      <div v-if="error" class="mb-6 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
        {{ error }}
      </div>

      <div v-else-if="order">
        <div v-if="updateMessage" class="mb-4 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-medium text-emerald-700">
          {{ updateMessage }}
        </div>
        <div v-if="updateError" class="mb-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
          {{ updateError }}
        </div>

        <!-- Информация о заказе -->
        <div class="mb-6 rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
          <div class="mb-5 flex items-center justify-between border-b border-gray-100 pb-4">
            <h2 class="text-lg font-bold text-gray-900">Информация о заказе</h2>
            <BaseBadge :variant="getStatusVariant(order.status)">
              {{ getStatusLabel(order.status) }}
            </BaseBadge>
          </div>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-3">
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">ID</span>
                <span class="text-sm font-semibold text-gray-900">{{ order.id }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">ID пользователя</span>
                <span class="text-sm font-semibold text-gray-900">{{ order.user_id }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Грузоподъёмность</span>
                <span class="text-sm font-semibold text-gray-900">{{ order.capacity }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Тип</span>
                <span class="text-sm font-semibold text-gray-900">{{ order.type || '—' }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Цена</span>
                <span class="text-sm font-semibold text-indigo-600">{{ order.price ? `${order.price} ₽` : '—' }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Менеджер</span>
                <span class="text-sm font-semibold text-gray-900">{{ order.manager_id || '—' }}</span>
              </div>
            </div>
            <div class="space-y-3">
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Адрес забора</span>
                <span class="text-sm font-semibold text-gray-900">{{ order.pickup_address || '—' }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Адрес доставки</span>
                <span class="text-sm font-semibold text-gray-900">{{ order.delivery_address || '—' }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Дата начала</span>
                <span class="text-sm font-semibold text-gray-900">{{ formatDateTime(order.date_start) }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Дата окончания</span>
                <span class="text-sm font-semibold text-gray-900">{{ formatDateTime(order.date_end) }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Период с</span>
                <span class="text-sm font-semibold text-gray-900">{{ formatDateTime(order.period_from) }}</span>
              </div>
              <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
                <span class="text-sm text-gray-500">Период по</span>
                <span class="text-sm font-semibold text-gray-900">{{ formatDateTime(order.period_to) }}</span>
              </div>
            </div>
          </div>
          <div class="mt-4 grid grid-cols-1 gap-3 border-t border-gray-100 pt-4 sm:grid-cols-2">
            <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
              <span class="text-sm text-gray-500">Дата создания</span>
              <span class="text-sm font-semibold text-gray-900">{{ formatDateTime(order.date_create) }}</span>
            </div>
            <div class="flex justify-between rounded-xl bg-gray-50 px-4 py-2.5">
              <span class="text-sm text-gray-500">Дата обновления</span>
              <span class="text-sm font-semibold text-gray-900">{{ formatDateTime(order.date_update) }}</span>
            </div>
          </div>
        </div>

        <!-- Управление статусом -->
        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
          <div class="mb-5 border-b border-gray-100 pb-4">
            <h2 class="text-lg font-bold text-gray-900">Управление статусом</h2>
          </div>
          <div class="flex flex-col gap-4 sm:flex-row sm:items-end">
            <div class="flex-1">
              <label class="mb-1.5 block text-sm font-semibold text-gray-700">Новый статус</label>
              <select
                v-model="newStatus"
                class="w-full cursor-pointer rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              >
                <option value="1">Новый</option>
                <option value="2">В работе</option>
                <option value="3">Завершён</option>
                <option value="-1">Отменён</option>
              </select>
            </div>
            <button
              @click="handleUpdateStatus"
              :disabled="updating"
              class="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700 disabled:opacity-60"
            >
              <svg v-if="updating" class="h-4 w-4 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              <span v-else>Обновить статус</span>
            </button>
            <button
              @click="handleAssignManager"
              :disabled="updating"
              class="inline-flex items-center justify-center gap-2 rounded-xl border border-gray-200 bg-white px-5 py-2.5 text-sm font-semibold text-gray-700 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600 disabled:opacity-60"
            >
              Назначить меня менеджером
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>