<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { orderService } from '../api/order'

const router = useRouter()

const loading = ref(false)
const error = ref(null)
const success = ref(false)

const form = ref({
  capacity: '',
  lifting_capacity: '',
  type: '',
  manager_id: '',
  status: 1,
  price: '',
  pickup_address: '',
  delivery_address: '',
  date_start: '',
  date_end: '',
  period_from: '',
  period_to: '',
})

// Преобразование datetime-local в ISO 8601
function toISO(datetimeLocal) {
  if (!datetimeLocal) return null
  // datetime-local возвращает формат "YYYY-MM-DDTHH:mm"
  return new Date(datetimeLocal).toISOString()
}

async function handleSubmit() {
  loading.value = true
  error.value = null
  success.value = false

  try {
    // Формируем payload: только заполненные поля
    const payload = {
      capacity: Number(form.value.capacity),
      lifting_capacity: Number(form.value.lifting_capacity),
    }

    if (form.value.type) payload.type = Number(form.value.type)
    if (form.value.manager_id) payload.manager_id = Number(form.value.manager_id)
    if (form.value.status) payload.status = Number(form.value.status)
    if (form.value.price) payload.price = Number(form.value.price)
    if (form.value.pickup_address) payload.pickup_address = form.value.pickup_address
    if (form.value.delivery_address) payload.delivery_address = form.value.delivery_address
    if (form.value.date_start) payload.date_start = toISO(form.value.date_start)
    if (form.value.date_end) payload.date_end = toISO(form.value.date_end)
    if (form.value.period_from) payload.period_from = toISO(form.value.period_from)
    if (form.value.period_to) payload.period_to = toISO(form.value.period_to)

    const res = await orderService.create(payload)
    success.value = true

    // Перенаправление на страницу созданного заказа
    const createdId = res.data?.id || res.data?.message?.id
    if (createdId) {
      router.push(`/orders/${createdId}`)
    } else {
      router.push('/orders')
    }
  } catch (err) {
    error.value = err.message || 'Ошибка создания заказа'
  } finally {
    loading.value = false
  }
}
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
          <h1 class="text-2xl font-bold text-gray-900 sm:text-3xl">Новый заказ</h1>
          <p class="mt-1 text-sm text-gray-500">Создание заказа на грузоперевозку</p>
        </div>
      </div>
    </div>

    <div v-if="error" class="mb-6 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700">
      {{ error }}
    </div>

    <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)] sm:p-8">
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <!-- Обязательные поля -->
          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="capacity">
              Грузоподъёмность <span class="text-rose-500">*</span>
            </label>
            <input
              id="capacity"
              v-model="form.capacity"
              type="number"
              min="1"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              placeholder="1000"
              required
            />
            <div class="mt-1 text-xs text-gray-500">Минимум 1</div>
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="lifting_capacity">
              Грузоподъёмность (дублирующее) <span class="text-rose-500">*</span>
            </label>
            <input
              id="lifting_capacity"
              v-model="form.lifting_capacity"
              type="number"
              min="1"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              placeholder="500"
              required
            />
            <div class="mt-1 text-xs text-gray-500">Минимум 1</div>
          </div>

          <!-- Опциональные поля -->
          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="type">Тип заказа</label>
            <input
              id="type"
              v-model="form.type"
              type="number"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              placeholder="1"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="status">Статус</label>
            <select
              id="status"
              v-model="form.status"
              class="w-full cursor-pointer rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            >
              <option :value="1">Новый</option>
              <option :value="2">В работе</option>
              <option :value="3">Завершён</option>
              <option :value="-1">Отменён</option>
            </select>
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="manager_id">ID менеджера</label>
            <input
              id="manager_id"
              v-model="form.manager_id"
              type="number"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              placeholder="2"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="price">Цена</label>
            <input
              id="price"
              v-model="form.price"
              type="number"
              step="0.01"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              placeholder="1500.50"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="pickup_address">Адрес забора груза</label>
            <input
              id="pickup_address"
              v-model="form.pickup_address"
              type="text"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              placeholder="г. Москва, ул. Ленина 1"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="delivery_address">Адрес доставки</label>
            <input
              id="delivery_address"
              v-model="form.delivery_address"
              type="text"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
              placeholder="г. Томск, пр. Ленина 40"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="date_start">Дата начала</label>
            <input
              id="date_start"
              v-model="form.date_start"
              type="datetime-local"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="date_end">Дата окончания</label>
            <input
              id="date_end"
              v-model="form.date_end"
              type="datetime-local"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="period_from">Начало периода</label>
            <input
              id="period_from"
              v-model="form.period_from"
              type="datetime-local"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-sm font-semibold text-gray-700" for="period_to">Конец периода</label>
            <input
              id="period_to"
              v-model="form.period_to"
              type="datetime-local"
              class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
            />
          </div>
        </div>

        <div class="mt-6 flex flex-col gap-3 border-t border-gray-100 pt-6 sm:flex-row">
          <button
            type="submit"
            :disabled="loading"
            class="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700 disabled:opacity-60"
          >
            <svg v-if="loading" class="h-4 w-4 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
            <span v-else>Создать заказ</span>
          </button>
          <button
            type="button"
            @click="router.back()"
            class="inline-flex items-center justify-center gap-2 rounded-xl border border-gray-200 bg-white px-5 py-2.5 text-sm font-semibold text-gray-700 transition-all duration-300 hover:border-gray-300 hover:bg-gray-50"
          >
            Отмена
          </button>
        </div>
      </form>
    </div>
  </div>
</template>