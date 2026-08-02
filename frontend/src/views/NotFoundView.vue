<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="animate-fade-in flex min-h-[60vh] flex-col items-center justify-center px-4 text-center">
    <p class="bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-8xl font-extrabold text-transparent">404</p>
    <h1 class="mt-4 text-2xl font-bold text-gray-900">Страница не найдена</h1>
    <p class="mt-2 max-w-md text-sm text-gray-500">
      Запрашиваемая страница не существует или была перемещена.
    </p>
    <button
      @click="router.push('/')"
      class="mt-8 inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-6 py-3 text-sm font-semibold text-white shadow-md shadow-indigo-600/20 transition-all duration-300 hover:bg-indigo-700 hover:shadow-lg"
    >
      ← Вернуться на главную
    </button>
  </div>
</template>