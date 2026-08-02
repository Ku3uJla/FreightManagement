<script setup>
/**
 * BaseSelect — выпадающий список с label
 *
 * Props:
 *   modelValue: v-model значение
 *   label: подпись
 *   options: [{ value, label }] — массив опций
 *   placeholder: текст пустой опции (пустая строка)
 */
defineProps({
  modelValue: {
    type: [String, Number],
    default: '',
  },
  label: {
    type: String,
    default: '',
  },
  options: {
    type: Array,
    default: () => [],
  },
  placeholder: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div class="mb-5">
    <label v-if="label" class="mb-1.5 block text-sm font-semibold text-gray-700">{{ label }}</label>
    <select
      :value="modelValue"
      @change="emit('update:modelValue', $event.target.value)"
      class="w-full cursor-pointer rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
    >
      <option v-if="placeholder" value="">{{ placeholder }}</option>
      <option v-for="opt in options" :key="opt.value" :value="opt.value">
        {{ opt.label }}
      </option>
    </select>
  </div>
</template>