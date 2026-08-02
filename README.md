# FreightManagement

Микросервисная архитектура для управления грузоперевозками.

## Архитектура

```
┌──────────────┐     gRPC :50051     ┌──────────────┐
│  auth-service│                     │  user-service│
│  HTTP :8083  │                     │  HTTP :8080  │
│  gRPC :50053 │                     │  gRPC :50051 │
└──────────────┘                     └──────────────┘
       │                                    │
       │                                    │ RabbitMQ
       │                                    ▼
       │                             ┌──────────────────┐
       │                             │ notification-svc │
       │                             └──────────────────┘
       │                                    ▲
┌──────────────┐     RabbitMQ              │
│ order-service│  order.created/updated ───┤
│  HTTP :8082  │                          │
└──────────────┘                          │
       │                                  │
┌──────────────┐                          │
│resource-svc  │                          │
│  HTTP :8081  │                          │
└──────────────┘                          │
       │                                  │
       ▼                                  ▼
┌──────────────────┐    ┌──────────────────┐
│  nginx (Gateway) │    │    frontend      │
│    HTTP :80      │◄───│  HTTP :3000      │
│  API Reverse Proxy│   │  Vue 3 + Vite    │
└──────────────────┘    └──────────────────┘
```

## Сервисы

| Сервис | Порт | Описание |
|--------|------|----------|
| auth-service | 8083 (HTTP), 50053 (gRPC) | Аутентификация, JWT |
| user-service | 8080 (HTTP), 50051 (gRPC) | Профили пользователей |
| order-service | 8082 (HTTP) | Управление заказами |
| resource-service | 8081 (HTTP) | Водители, автомобили |
| notification-service | — | RabbitMQ consumer (уведомления) |
| nginx (API Gateway) | 80 | Reverse proxy к микросервисам |
| **frontend** | **3000** | **Vue 3 SPA (Composition API, Vite)** |

## Сети Docker

| Сеть | Сервисы | Назначение |
|------|---------|------------|
| `backend` | db, rabbitmq, auth-app, user-app, order-app, resource-app, notification-app, nginx | Внутренняя сеть микросервисов |
| `frontend` | nginx, frontend | Сеть для frontend ↔ API Gateway |

## Запуск через Docker Compose

```bash
# Запуск всех сервисов (включая frontend)
docker-compose up --build

# Frontend будет доступен на http://localhost:3000
# API Gateway (nginx) на http://localhost:80
```

## Локальная разработка frontend

```bash
cd frontend
npm install
npm run dev
```

Vite dev-сервер запустится на `http://localhost:5173` и будет проксировать
запросы `/api/*` к API Gateway на `http://localhost:80`.

## Структура frontend

```
frontend/
├── src/
│   ├── api/           # Axios-клиент и сервисы для каждого микросервиса
│   │   ├── client.js  # Базовый axios-клиент (baseURL: /api)
│   │   ├── auth.js    # auth-service API
│   │   ├── user.js    # user-service API
│   │   ├── order.js   # order-service API
│   │   └── resource.js# resource-service API
│   ├── stores/        # Pinia stores
│   │   └── auth.js    # Состояние аутентификации
│   ├── router/        # Vue Router (history mode)
│   ├── views/         # Страницы приложения
│   ├── assets/        # Глобальные CSS стили
│   ├── App.vue        # Корневой компонент (навигация)
│   └── main.js        # Точка входа
├── Dockerfile         # Многоэтапная сборка (Node build → nginx serve)
├── nginx.conf         # Конфигурация nginx для SPA + API proxy
├── vite.config.js     # Конфигурация Vite (proxy для dev-режима)
└── package.json
```

## Маршрутизация API

Frontend отправляет запросы на `/api/*`, которые проксируются:

- **В Docker:** frontend-nginx → `http://nginx:80/*` (rewrite: убирает `/api`)
- **В dev-режиме:** Vite proxy → `http://localhost:80/*`

| Frontend запрос | → nginx Gateway | → Микросервис |
|-----------------|-----------------|---------------|
| `POST /api/auth/register` | `/auth/register` | auth-service:8083 |
| `POST /api/auth/login` | `/auth/login` | auth-service:8083 |
| `GET /api/user/:id` | `/user/:id` | user-service:8080 |
| `GET /api/orders` | `/orders` | order-service:8082 |
| `POST /api/orders` | `/orders` | order-service:8082 |
| `GET /api/drivers` | `/drivers` | resource-service:8081 |
| `GET /api/autos` | `/autos` | resource-service:8081 |

Полная документация API: [API_ROUTES.md](API_ROUTES.md)