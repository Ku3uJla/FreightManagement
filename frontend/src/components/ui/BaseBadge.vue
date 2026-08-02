<script setup>
import { computed } from 'vue'

/**
 * BaseBadge — цветной бейдж для статусов
 *
 * Props:
 *   variant: info | success | warning | danger (default: info)
 */
const props = defineProps({
  variant: {
    type: String,
    default: 'info',
    validator: (v) => ['info', 'success', 'warning', 'danger'].includes(v),
  },
})

const variantClasses = computed(() => {
  const map = {
    info: 'bg-indigo-50 text-indigo-700 ring-indigo-200',
    success: 'bg-emerald-50 text-emerald-700 ring-emerald-200',
    warning: 'bg-amber-50 text-amber-700 ring-amber-200',
    danger: 'bg-rose-50 text-rose-700 ring-rose-200',
  }
  return map[props.variant] || map.info
})
</script>

<template>
  <span
    class="inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold ring-1 ring-inset"
    :class="variantClasses"
  >
    <slot />
  </span>
</template>