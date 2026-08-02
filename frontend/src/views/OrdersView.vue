<script setup>
import { ref, onMounted } from 'vue'
import { orderService } from '../api/order'
import BaseBadge from '../components/ui/BaseBadge.vue'

const loading = ref(true)
const error = ref(null)
const orders = ref([])
const meta = ref({ page: 1, pageSize: 20, total: 0, totalPages: 0 })
const currentPage = ref(1)

// Фильтры
const filters = ref({
  status: '',
  type: '',
  price_min: '',
  price_max: '',
  pickup_address: '',
  delivery_address: '',
})

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

function formatDate(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('ru-RU')
}

async function fetchOrders() {
  loading.value = true
  error.value = null
  try {
    // Очищаем пустые фильтры
    const params = { page: currentPage.value }
    Object.entries(filters.value).forEach(([key, val]) => {
      if (val !== '' && val !== null && val !== undefined) {
        params[key] = val
      }
    })

    const res = await orderService.list(params)
    orders.value = res.data?.data || []
    meta.value = res.data?.meta || { page: 1, pageSize: 20, total: 0, totalPages: 0 }
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить заказы'
    orders.value = []
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  currentPage.value = 1
  fetchOrders()
}

function resetFilters() {
  filters.value = {
    status: '',
    type: '',
    price_min: '',
    price_max: '',
    pickup_address: '',
    delivery_address: '',
  }
  currentPage.value = 1
  fetchOrders()
}

function changePage(page) {
  if (page < 1 || page > meta.value.totalPages) return
  currentPage.value = page
  fetchOrders()
}

onMounted(() => {
  fetchOrders()
})
</script>

<template>
  <div class="animate-fade-in">
    <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 sm:text-3xl">Заказы</h1>
        <p class="mt-1 text-sm text-gray-500">Управление заказами на грузоперевозки</p>
      </div>
      <router-link
        to="/orders/new"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700 hover:shadow-lg"
      >
        + Создать заказ
      </router-link>
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
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Статус</label>
          <select
            v-model="filters.status"
            class="w-full cursor-pointer rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
          >
            <option value="">Все статусы</option>
            <option value="1">Новый</option>
            <option value="2">В работе</option>
            <option value="3">Завершён</option>
            <option value="-1">Отменён</option>
          </select>
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Тип</label>
          <input
            v-model="filters.type"
            type="number"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="Тип заказа"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Цена от</label>
          <input
            v-model="filters.price_min"
            type="number"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="0"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Цена до</label>
          <input
            v-model="filters.price_max"
            type="number"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="∞"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Адрес забора</label>
          <input
            v-model="filters.pickup_address"
            type="text"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="Частичное совпадение"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">Адрес доставки</label>
          <input
            v-model="filters.delivery_address"
            type="text"
            class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            placeholder="Частичное совпадение"
          />
        </div>
      </div>
      <button
        @click="applyFilters"
        class="mt-4 inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700"
      >
        Применить фильтры
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
      <div v-if="orders.length === 0" class="flex flex-col items-center justify-center rounded-2xl border border-gray-100 bg-white py-16 text-center shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
        <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-indigo-50 text-3xl">📦</div>
        <h3 class="text-lg font-semibold text-gray-900">Заказы не найдены</h3>
        <p class="mt-2 text-sm text-gray-500">Попробуйте изменить фильтры или создать новый заказ.</p>
      </div>

      <div v-else>
        <div class="overflow-x-auto rounded-2xl border border-gray-100 bg-white shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
          <table class="w-full border-collapse">
            <thead>
              <tr class="border-b border-gray-100 bg-gray-50">
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">ID</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Адрес забора</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Адрес доставки</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Грузоподъёмность</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Цена</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Статус</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Дата создания</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="order in orders"
                :key="order.id"
                class="border-b border-gray-50 transition-colors duration-200 last:border-0 hover:bg-indigo-50/40"
              >
                <td class="px-5 py-4 text-sm font-semibold text-gray-900">#{{ order.id }}</td>
                <td class="px-5 py-4 text-sm text-gray-700">{{ order.pickup_address || '—' }}</td>
                <td class="px-5 py-4 text-sm text-gray-700">{{ order.delivery_address || '—' }}</td>
                <td class="px-5 py-4 text-sm text-gray-700">{{ order.capacity || '—' }}</td>
                <td class="px-5 py-4 text-sm font-medium text-gray-900">{{ order.price ? `${order.price} ₽` : '—' }}</td>
                <td class="px-5 py-4">
                  <BaseBadge :variant="getStatusVariant(order.status)">
                    {{ getStatusLabel(order.status) }}
                  </BaseBadge>
                </td>
                <td class="px-5 py-4 text-sm text-gray-500">{{ formatDate(order.date_create) }}</td>
                <td class="px-5 py-4">
                  <router-link
                    :to="`/orders/${order.id}`"
                    class="inline-flex items-center gap-1 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600"
                  >
                    Подробнее →
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Пагинация -->
        <div v-if="meta.totalPages > 1" class="mt-6 flex items-center justify-center gap-3">
          <button
            :disabled="currentPage === 1"
            @click="changePage(currentPage - 1)"
            class="inline-flex items-center gap-1 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600 disabled:cursor-not-allowed disabled:opacity-50"
          >
            ← Назад
          </button>
          <span class="rounded-xl bg-indigo-50 px-4 py-2 text-sm font-medium text-indigo-700">
            Страница {{ meta.page }} из {{ meta.totalPages }}
            <span class="ml-1 text-indigo-400">(всего {{ meta.total }})</span>
          </span>
          <button
            :disabled="currentPage === meta.totalPages"
            @click="changePage(currentPage + 1)"
            class="inline-flex items-center gap-1 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Вперёд →
          </button>
        </div>
      </div>
    </div>
  </div>
</template>