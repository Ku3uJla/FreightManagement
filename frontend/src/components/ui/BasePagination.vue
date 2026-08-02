<script setup>
/**
 * BasePagination — пагинация
 *
 * Props:
 *   page: текущая страница
 *   totalPages: всего страниц
 *   total: всего записей
 */
defineProps({
  page: {
    type: Number,
    default: 1,
  },
  totalPages: {
    type: Number,
    default: 0,
  },
  total: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['change'])
</script>

<template>
  <div v-if="totalPages > 1" class="mt-6 flex items-center justify-center gap-3">
    <button
      :disabled="page === 1"
      @click="emit('change', page - 1)"
      class="inline-flex items-center gap-1 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600 disabled:cursor-not-allowed disabled:opacity-50"
    >
      ← Назад
    </button>
    <span class="rounded-xl bg-indigo-50 px-4 py-2 text-sm font-medium text-indigo-700">
      Страница {{ page }} из {{ totalPages }}
      <span class="ml-1 text-indigo-400">(всего {{ total }})</span>
    </span>
    <button
      :disabled="page === totalPages"
      @click="emit('change', page + 1)"
      class="inline-flex items-center gap-1 rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-all duration-300 hover:border-indigo-300 hover:text-indigo-600 disabled:cursor-not-allowed disabled:opacity-50"
    >
      Вперёд →
    </button>
  </div>
</template>