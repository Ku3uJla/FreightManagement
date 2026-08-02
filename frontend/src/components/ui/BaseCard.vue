<script setup>
/**
 * BaseCard — карточка с опциональным заголовком
 *
 * Slots:
 *   default — тело карточки
 *   header — заголовок (опционально)
 *   actions — кнопки в заголовке (опционально)
 */
</script>

<template>
  <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-[0_4px_20px_rgba(0,0,0,0.06)] transition-all duration-300 hover:shadow-[0_8px_30px_rgba(0,0,0,0.08)]">
    <div v-if="$slots.header || $slots.actions" class="mb-5 flex items-center justify-between gap-4 border-b border-gray-100 pb-4">
      <h2 class="text-lg font-bold text-gray-900">
        <slot name="header" />
      </h2>
      <div v-if="$slots.actions" class="flex shrink-0 items-center gap-2">
        <slot name="actions" />
      </div>
    </div>
    <slot />
  </div>
</template>