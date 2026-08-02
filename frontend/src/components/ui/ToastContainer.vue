<script setup>
import { computed } from 'vue'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()

const iconMap = {
  success: '✓',
  error: '✕',
  warning: '⚠',
  info: 'ℹ',
}

const typeClasses = computed(() => ({
  success: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  error: 'bg-rose-50 text-rose-700 border-rose-200',
  warning: 'bg-amber-50 text-amber-700 border-amber-200',
  info: 'bg-indigo-50 text-indigo-700 border-indigo-200',
}))
</script>

<template>
  <Teleport to="body">
    <div class="fixed right-4 top-4 z-[9999] flex max-w-sm flex-col gap-3">
      <TransitionGroup name="toast">
        <div
          v-for="t in toast.toasts"
          :key="t.id"
          class="animate-slide-in-right flex cursor-pointer items-center gap-3 rounded-xl border px-4 py-3 text-sm font-medium shadow-lg"
          :class="typeClasses[t.type]"
          @click="toast.remove(t.id)"
        >
          <span class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-white/60 text-sm font-bold">
            {{ iconMap[t.type] }}
          </span>
          <span class="flex-1">{{ t.message }}</span>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>