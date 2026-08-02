<script setup>
/**
 * BaseEmptyState — пустое состояние с заголовком, описанием и опциональным action
 *
 * Slots:
 *   default — заголовок
 *   description — описание
 *   action — кнопка действия
 */
</script>

<template>
  <div class="flex flex-col items-center justify-center px-4 py-12 text-center">
    <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-indigo-50 text-3xl">
      📦
    </div>
    <h3 class="text-lg font-semibold text-gray-900"><slot /></h3>
    <p v-if="$slots.description" class="mt-2 max-w-md text-sm text-gray-500">
      <slot name="description" />
    </p>
    <div v-if="$slots.action" class="mt-6">
      <slot name="action" />
    </div>
  </div>
</template>