---
name: frontend-design
description: Используй этот навык при разработке UI/UX фронтенда на Vue 3 + Vite + Tailwind CSS — современный дизайн, карточки, адаптивная сетка, анимации, тени. Переписывает вьюхи с красивым дизайном, сохраняя маршруты и API-интеграцию.
---

# Frontend Design (Vue 3 + Vite + Tailwind CSS)

## Стек проекта

- Vue 3 (Composition API, `<script setup>`)
- Vite
- Tailwind CSS (утилитные классы, без UI-библиотек)
- Vue Router (все существующие маршруты сохранять)
- Pinia (сторы `auth`, `toast`)
- API-клиенты через `src/api/*.js`, ходят через nginx (reverse proxy `/api`)

## Принципы дизайна

1. **Карточки** — скругления `rounded-2xl`, тени `shadow-lg` / `shadow-xl`, лёгкие границы (`border border-gray-100` или `border-gray-200`).
2. **Адаптивная сетка** — `grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6` для карточек; `flex flex-col md:flex-row` для хедеров/панелей.
3. **Анимации** — плавные переходы `transition-all duration-300`, hover-эффекты `hover:-translate-y-1 hover:shadow-2xl`, появление элементов с `animate-fade-in` (ключевые кадры в `main.css`).
4. **Тени** — `shadow-lg`, `shadow-xl`, `hover:shadow-2xl`; мягкие тени `shadow-[0_4px_20px_rgba(0,0,0,0.06)]`.
5. **Цветовая палитра** — основной цвет `indigo-600`, акценты `emerald-500` (успех), `amber-500` (предупреждение), `rose-500` (ошибка). Фон страниц — `bg-gray-50`, поверхности — `bg-white`.
6. **Типографика** — заголовки `text-2xl font-bold text-gray-900`, подзаголовки `text-gray-500`. Отступы `space-y-4`, `space-y-6`.
7. **Кнопки** — `rounded-xl px-5 py-2.5 font-medium transition-all`, primary `bg-indigo-600 text-white hover:bg-indigo-700`, secondary `bg-white text-gray-700 border border-gray-200 hover:bg-gray-50`.
8. **Загрузка** — скелетоны или `animate-pulse`; спины — `animate-spin`.
9. **Бейджи** — `inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold` с цветами статусов.
10. **Пустые состояния** — центрированная иконка (emoji или inline SVG), заголовок, поясняющий текст, кнопка действия.

## Требования

- Сохранять **все маршруты** из `src/router/index.js` и все `name` маршрутов.
- Сохранять **все API-вызовы** (те же endpoints, методы, параметры) — только менять дизайн, не бэкенд.
- Все UI-компоненты (`BaseButton`, `BaseCard`, `BaseBadge`, `BaseInput`, `BaseSelect`, `BaseModal`, `BaseTable`, `BasePagination`, `BaseEmptyState`, `BaseLoading`, `ToastContainer`) перевести на Tailwind-классы.
- Tailwind подключён через `@tailwindcss/vite` плагин в `vite.config.js` и директивы `@import "tailwindcss"` в `src/assets/main.css`.
- После переписывания пересобрать контейнер: `docker compose up -d --build frontend`.

## Структура проекта

- `frontend/src/views/` — страницы (Login, Register, Dashboard, Orders, OrderDetail, OrderCreate, Drivers, Autos, Profile, NotFound)
- `frontend/src/router/index.js` — маршруты
- `frontend/src/api/` — клиенты API (client, auth, user, order, resource)
- `frontend/src/stores/` — Pinia-сторы (auth, toast)
- `frontend/src/components/ui/` — переиспользуемые UI-компоненты

## Проверка качества

- Дизайн адаптивный (mobile-first): проверять на `sm`, `md`, `lg` брейкпоинтах.
- Анимации не ломают layout (использовать `transform` / `opacity`, не `width/height` где возможно).
- Все формы валидируются и используют `@submit.prevent`.
- Семантическая разметка: `<header>`, `<main>`, `<section>`, `<footer>` где уместно.