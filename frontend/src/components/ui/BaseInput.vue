<script setup>
/**
 * BaseInput — поле ввода с label, подсказкой и required-маркером
 *
 * Props:
 *   modelValue: v-model значение
 *   label: подпись поля
 *   type: text | number | email | password | tel | datetime-local (default: text)
 *   placeholder
 *   required: boolean
 *   hint: подсказка под полем
 *   min, max, step, maxlength — нативные атрибуты
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
  type: {
    type: String,
    default: 'text',
  },
  placeholder: {
    type: String,
    default: '',
  },
  required: {
    type: Boolean,
    default: false,
  },
  hint: {
    type: String,
    default: '',
  },
  min: {
    type: [String, Number],
    default: undefined,
  },
  max: {
    type: [String, Number],
    default: undefined,
  },
  step: {
    type: [String, Number],
    default: undefined,
  },
  maxlength: {
    type: [String, Number],
    default: undefined,
  },
  autocomplete: {
    type: String,
    default: 'off',
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div class="mb-5">
    <label v-if="label" class="mb-1.5 block text-sm font-semibold text-gray-700" :for="$attrs.id">
      {{ label }}
      <span v-if="required" class="text-rose-500">*</span>
    </label>
    <input
      :value="modelValue"
      @input="emit('update:modelValue', $event.target.value)"
      :type="type"
      class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 placeholder-gray-400 transition-all duration-300 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200"
      :placeholder="placeholder"
      :required="required"
      :min="min"
      :max="max"
      :step="step"
      :maxlength="maxlength"
      :autocomplete="autocomplete"
    />
    <div v-if="hint" class="mt-1 text-xs text-gray-500">{{ hint }}</div>
  </div>
</template>