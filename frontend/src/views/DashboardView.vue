<script setup>
import { ref, onMounted } from 'vue'
import { orderService } from '../api/order'
import { resourceService } from '../api/resource'
import BaseBadge from '../components/ui/BaseBadge.vue'

const loading = ref(true)
const error = ref(null)
const stats = ref({
  orders: 0,
  drivers: 0,
  autos: 0,
  activeOrders: 0,
})

const recentOrders = ref([])

onMounted(async () => {
  try {
    // Загружаем данные параллельно
    const [ordersRes, driversRes, autosRes] = await Promise.all([
      orderService.list({ page: 1 }),
      resourceService.listDrivers(),
      resourceService.listAutos(),
    ])

    const ordersData = ordersRes.data?.data || []
    stats.value.orders = ordersRes.data?.meta?.total || ordersData.length
    stats.value.activeOrders = ordersData.filter(
      (o) => o.status === 1 || o.status === 2
    ).length
    stats.value.drivers = (driversRes.data || []).length
    stats.value.autos = (autosRes.data || []).length
    recentOrders.value = ordersData.slice(0, 5)
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить данные дашборда'
  } finally {
    loading.value = false
  }
})

// Маппинг статусов заказов
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

const statCards = [
  { key: 'orders', label: 'Всего заказов', icon: '📋', color: 'from-indigo-500 to-indigo-600' },
  { key: 'activeOrders', label: 'Активные заказы', icon: '🚚', color: 'from-emerald-500 to-emerald-600' },
  { key: 'drivers', label: 'Водители', icon: '👤', color: 'from-amber-500 to-amber-600' },
  { key: 'autos', label: 'Автомобили', icon: '🚛', color: 'from-purple-500 to-purple-600' },
]
</script>

<template>
  <div class="animate-fade-in">
    <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 sm:text-3xl">Дашборд</h1>
        <p class="mt-1 text-sm text-gray-500">Обзор системы управления грузоперевозками</p>
      </div>
      <router-link
        to="/orders/new"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700 hover:shadow-lg"
      >
        + Создать заказ
      </router-link>
    </div>

    <!-- Skeleton loading -->
    <div v-if="loading" class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
      <div v-for="i in 4" :key="i" class="skeleton h-32 rounded-2xl"></div>
    </div>

    <div v-else>
      <div v-if="error" class="mb-6 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
        {{ error }}
      </div>

      <!-- Статистика -->
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <div
          v-for="(card, index) in statCards"
          :key="card.key"
          class="animate-fade-in group rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)] transition-all duration-300 hover:-translate-y-1 hover:shadow-xl"
          :class="`animate-delay-${(index + 1) * 100}`"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wider text-gray-500">{{ card.label }}</p>
              <p class="mt-2 text-3xl font-bold text-gray-900">{{ stats[card.key] }}</p>
            </div>
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br text-xl text-white shadow-lg" :class="card.color">
              {{ card.icon }}
            </div>
          </div>
        </div>
      </div>

      <!-- Последние заказы -->
      <div class="mt-8 rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
        <div class="mb-5 flex items-center justify-between border-b border-gray-100 pb-4">
          <h2 class="text-lg font-bold text-gray-900">Последние заказы</h2>
          <router-link
            to="/orders"
            class="inline-flex items-center gap-1 rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-xs font-medium text-gray-700 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600"
          >
            Все заказы →
          </router-link>
        </div>

        <div v-if="recentOrders.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
          <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-indigo-50 text-3xl">📦</div>
          <h3 class="text-lg font-semibold text-gray-900">Заказов пока нет</h3>
          <p class="mt-2 text-sm text-gray-500">Создайте первый заказ, чтобы увидеть его здесь.</p>
          <router-link
            to="/orders/new"
            class="mt-6 inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700"
          >
            Создать заказ
          </router-link>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full border-collapse">
            <thead>
              <tr class="border-b border-gray-100 bg-gray-50">
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">ID</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Адрес забора</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Адрес доставки</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Цена</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Статус</th>
                <th class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">Дата создания</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="order in recentOrders"
                :key="order.id"
                class="border-b border-gray-50 transition-colors duration-200 last:border-0 hover:bg-indigo-50/40"
              >
                <td class="px-5 py-4 text-sm">
                  <router-link :to="`/orders/${order.id}`" class="font-semibold text-indigo-600 transition-colors hover:text-indigo-700">
                    #{{ order.id }}
                  </router-link>
                </td>
                <td class="px-5 py-4 text-sm text-gray-700">{{ order.pickup_address || '—' }}</td>
                <td class="px-5 py-4 text-sm text-gray-700">{{ order.delivery_address || '—' }}</td>
                <td class="px-5 py-4 text-sm font-medium text-gray-900">{{ order.price ? `${order.price} ₽` : '—' }}</td>
                <td class="px-5 py-4">
                  <BaseBadge :variant="getStatusVariant(order.status)">
                    {{ getStatusLabel(order.status) }}
                  </BaseBadge>
                </td>
                <td class="px-5 py-4 text-sm text-gray-500">{{ formatDate(order.date_create) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>