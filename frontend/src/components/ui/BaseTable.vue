<script setup>
/**
 * BaseTable — таблица с колонками и слотами для ячеек
 *
 * Props:
 *   columns: [{ key, label }] — описание колонок
 *   items: массив данных
 *
 * Slots:
 *   cell-<key> — кастомный рендер ячейки для колонки key
 *   empty — контент при пустом списке
 */
defineProps({
  columns: {
    type: Array,
    required: true,
  },
  items: {
    type: Array,
    default: () => [],
  },
})
</script>

<template>
  <div v-if="items.length === 0" class="py-10 text-center">
    <slot name="empty">
      <h3 class="text-lg font-semibold text-gray-700">Данные отсутствуют</h3>
    </slot>
  </div>
  <div v-else class="overflow-x-auto rounded-2xl border border-gray-100 bg-white shadow-[0_4px_20px_rgba(0,0,0,0.06)]">
    <table class="w-full border-collapse">
      <thead>
        <tr class="border-b border-gray-100 bg-gray-50">
          <th
            v-for="col in columns"
            :key="col.key"
            class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500"
          >
            {{ col.label }}
          </th>
          <th v-if="$slots.actions" class="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-gray-500">
            Действия
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(item, index) in items"
          :key="item.id ?? index"
          class="border-b border-gray-50 transition-colors duration-200 last:border-0 hover:bg-indigo-50/40"
        >
          <td v-for="col in columns" :key="col.key" class="px-5 py-4 text-sm text-gray-700">
            <slot :name="`cell-${col.key}`" :item="item" :value="item[col.key]">
              {{ item[col.key] ?? '—' }}
            </slot>
          </td>
          <td v-if="$slots.actions" class="px-5 py-4">
            <slot name="actions" :item="item" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>