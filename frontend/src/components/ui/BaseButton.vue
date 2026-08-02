<script setup>
import { computed } from 'vue'

/**
 * BaseButton — универсальная кнопка с вариантами и размерами
 *
 * Props:
 *   variant: primary | danger | success | outline (default: primary)
 *   size: sm | md | block (default: md)
 *   loading: boolean — показывает спиннер и блокирует кнопку
 *   disabled: boolean
 *   type: button | submit (default: button)
 */
const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (v) => ['primary', 'danger', 'success', 'outline'].includes(v),
  },
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md', 'block'].includes(v),
  },
  loading: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  type: {
    type: String,
    default: 'button',
  },
})

const isDisabled = computed(() => props.disabled || props.loading)

const variantClasses = computed(() => {
  const map = {
    primary:
      'bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-300 shadow-md shadow-indigo-600/20',
    danger:
      'bg-rose-500 text-white hover:bg-rose-600 focus:ring-rose-300 shadow-md shadow-rose-500/20',
    success:
      'bg-emerald-500 text-white hover:bg-emerald-600 focus:ring-emerald-300 shadow-md shadow-emerald-500/20',
    outline:
      'bg-white text-gray-700 border border-gray-200 hover:bg-gray-50 hover:border-gray-300 focus:ring-gray-200',
  }
  return map[props.variant] || map.primary
})

const sizeClasses = computed(() => {
  if (props.size === 'sm') return 'px-3 py-1.5 text-xs rounded-lg'
  if (props.size === 'block') return 'w-full px-5 py-3 text-sm rounded-xl'
  return 'px-5 py-2.5 text-sm rounded-xl'
})
</script>

<template>
  <button
    :type="type"
    class="inline-flex items-center justify-center gap-2 font-semibold transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-1 disabled:opacity-60 disabled:cursor-not-allowed active:scale-[0.98]"
    :class="[variantClasses, sizeClasses]"
    :disabled="isDisabled"
  >
    <svg
      v-if="loading"
      class="h-4 w-4 animate-spin"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
    </svg>
    <slot v-if="!loading" />
  </button>
</template>