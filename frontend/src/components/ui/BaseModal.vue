<script setup>
/**
 * BaseModal — модальное окно с overlay
 *
 * Props:
 *   modelValue: v-model (boolean) — открыто/закрыто
 *   title: заголовок окна
 *
 * Emits:
 *   update:modelValue — закрытие
 */
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])

function close() {
  emit('update:modelValue', false)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-[1000] flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
        @click.self="close"
      >
        <div class="animate-scale-in w-full max-w-lg max-h-[90vh] overflow-y-auto rounded-2xl border border-gray-100 bg-white p-6 shadow-2xl">
          <div class="mb-5 flex items-center justify-between border-b border-gray-100 pb-4">
            <h2 class="text-lg font-bold text-gray-900">{{ title }}</h2>
            <button
              @click="close"
              class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
            >
              ✕
            </button>
          </div>
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>