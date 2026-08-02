<script setup>
/**
 * BaseLoading — индикатор загрузки (спиннер)
 *
 * Props:
 *   text: текст рядом со спиннером (опционально)
 */
defineProps({
  text: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <div class="flex items-center justify-center gap-3 py-12 text-gray-500">
    <svg
      class="h-6 w-6 animate-spin text-indigo-600"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
    </svg>
    <span v-if="text" class="text-sm font-medium">{{ text }}</span>
  </div>
</template>