# ==============================================================================
# Freight Management — Frontend Dockerfile
# ==============================================================================
# Многоэтапная сборка:
#   Stage 1 (build): Node.js + Vite собирают статические файлы Vue SPA
#   Stage 2 (serve): nginx:alpine отдаёт собранные файлы и проксирует /api
# ==============================================================================

# ---- Stage 1: Build Vue SPA ----
FROM node:20-alpine AS build

WORKDIR /app

# Копируем package-файлы и устанавливаем зависимости
COPY package.json package-lock.json* ./
RUN npm install

# Копируем исходный код и собираем проект
COPY . .
RUN npm run build

# ---- Stage 2: Serve with nginx ----
FROM nginx:alpine AS serve

# Копируем собственную конфигурацию nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем собранные статические файлы из stage build
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]