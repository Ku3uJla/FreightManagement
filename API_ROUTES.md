# API Routes — Freight Management

Полная документация всех API эндпоинтов (HTTP/REST и gRPC) микросервисов проекта Freight Management.

---

## Краткое Оглавление

| Сервис | HTTP Порт | gRPC Порт | Описание |
|--------|-----------|-----------|----------|
| [auth-service](#auth-service) | `8083` | `50053` | Аутентификация и авторизация (регистрация, вход, JWT) |
| [user-service](#user-service) | `8080` | `50051` | Управление пользователями (профиль, роли) |
| [order-service](#order-service) | `8082` | — | Управление заказами (создание, обновление, фильтрация) |
| [resource-service](#resource-service) | `8081` | — | Управление ресурсами (водители, автомобили) |
| [notification-service](#notification-service) | — | — | Сервис уведомлений (RabbitMQ consumer, без HTTP/gRPC) |

### Дополнительные порты инфраструктуры

| Сервис | Порт | Описание |
|--------|------|----------|
| PostgreSQL (db) | `5432` | База данных `freightmanagementdb` |
| RabbitMQ | `5672` | Очередь сообщений (exchange: `app.events`) |

---

## auth-service

Сервис аутентификации. Отвечает за регистрацию пользователей, вход в систему и выдачу JWT-токенов.
Также предоставляет gRPC-методы для получения логина и email пользователя по ID.

- **HTTP порт:** `8083`
- **gRPC порт:** `50053`
- **Framework:** Gin (HTTP), gRPC
- **Зависимости:** PostgreSQL, user-service (gRPC клиент на `:50051`), RabbitMQ

### HTTP REST Эндпоинты

#### Health Check

| Метод | Path | Хэндлер | Назначение | Аутентификация | Входные данные | Успешный ответ |
|-------|------|---------|------------|----------------|----------------|----------------|
| `GET` | `/` | inline func | Проверка работоспособности сервиса | Public | — | `200` `{"message": "auth-service"}` |

#### Auth — Аутентификация

| Метод | Path | Хэндлер | Назначение | Аутентификация | Входные данные | Успешный ответ |
|-------|------|---------|------------|----------------|----------------|----------------|
| `POST` | `/auth/register` | `AuthController.SignUp` | Регистрация нового пользователя (создание учётной записи auth + профиля user через gRPC) | Public | [CreateAuthRequest](#createauthrequest) | `200` `{"message": "Signed UP"}` |
| `POST` | `/auth/login` | `AuthController.Login` | Вход в систему по логину и паролю, выдача JWT в cookie | Public | [Auth (Login)](#auth-login) | `200` `{"message": "logged"}` + Set-Cookie `JWT` |

##### CreateAuthRequest

```json
{
  "full_name": "string (required)",
  "login": "string (required)",
  "email": "string (required, email format)",
  "phone": "string (required)",
  "password": "string (required)"
}
```

##### Auth (Login)

```json
{
  "login": "string (required)",
  "password": "string (required)"
}
```

> **Примечание:** При успешном логине устанавливается cookie `JWT` (HttpOnly, срок 3600с / 1 час) с JWT-токеном, содержащим `ID` пользователя и `exp` (24 часа).

### gRPC Эндпоинты

**Сервис:** `AuthService`  
**Пакет:** `auth`  
**Порт:** `50053`

| Method | Request | Response | Хэндлер | Назначение |
|--------|---------|----------|---------|------------|
| `GetEmail` | [GetUserRequest](#getuserrequest) | [GetEmailResponse](#getemailresponse) | `AuthGrpcController.GetEmail` | Получение email пользователя по ID |
| `GetLogin` | [GetUserRequest](#getuserrequest) | [GetLoginResponse](#getloginresponse) | `AuthGrpcController.GetLogin` | Получение логина пользователя по ID |

##### GetUserRequest

```protobuf
message GetUserRequest {
  int64 id = 1;
}
```

##### GetEmailResponse

```protobuf
message GetEmailResponse {
  string email = 1;
}
```

##### GetLoginResponse

```protobuf
message GetLoginResponse {
  string Login = 1;
}
```

### gRPC Клиент

auth-service выступает в роли gRPC **клиента** для `user-service` (`UserService` на порту `50051`):
- Вызывает `CreateUser` при регистрации нового пользователя
- Вызывает `GetUser` для получения данных пользователя

---

## user-service

Сервис управления пользователями. Хранит профили пользователей (роль, телефон, email, ФИО).
Создаёт пользователей через gRPC (вызывается auth-service) и публикует событие `user.created` в RabbitMQ.

- **HTTP порт:** `8080`
- **gRPC порт:** `50051`
- **Framework:** Gin (HTTP), gRPC
- **Зависимости:** PostgreSQL, RabbitMQ

### HTTP REST Эндпоинты

#### User — Управление профилем

| Метод | Path | Хэндлер | Назначение | Аутентификация | Входные данные | Успешный ответ |
|-------|------|---------|------------|----------------|----------------|----------------|
| `GET` | `/user/:id` | `UserController.GetByID` | Получение профиля пользователя по ID | Public | URL param: `id` (int) | `200` `{"message": <User>}` |
| `PUT` | `/user/:id/role` | `UserController.UpdateRole` | Обновление роли пользователя | Public | URL param: `id` (int) + Body: `int` (role) | `200` (без тела) |

##### User (модель ответа)

```json
{
  "id": 1,
  "role": "client",
  "phone": "+1234567890",
  "email": "user@example.com",
  "name": "John Doe",
  "dateCreate": "2026-01-01T00:00:00Z",
  "dateUpdate": "2026-01-01T00:00:00Z"
}
```

##### UpdateRole (Request Body)

```json
2
```

> Тело запроса — просто число (int), представляющее новую роль пользователя.

### gRPC Эндпоинты

**Сервис:** `UserService`  
**Пакет:** `user`  
**Порт:** `50051`

| Method | Request | Response | Хэндлер | Назначение |
|--------|---------|----------|---------|------------|
| `CreateUser` | [CreateUserRequest](#createuserrequest) | [CreateUserResponse](#createuserresponse) | `UserGrpcController.CreateUser` | Создание профиля пользователя (вызывается auth-service при регистрации) |
| `GetUser` | [GetUserRequest](#getuserrequest-1) | [GetUserResponse](#getuserresponse-1) | `UserGrpcController.GetUser` | Получение данных пользователя по ID |

##### CreateUserRequest

```protobuf
message CreateUserRequest {
  int64 id = 1;
  string email = 2;
  string full_name = 3;
  string phone = 4;
}
```

##### CreateUserResponse

```protobuf
message CreateUserResponse {
  bool success = 1;
}
```

##### GetUserRequest

```protobuf
message GetUserRequest {
  int64 id = 1;
}
```

##### GetUserResponse

```protobuf
message GetUserResponse {
  string role = 1;
  string phone = 2;
  string full_name = 3;
  string email = 4;
}
```

### RabbitMQ События (Publisher)

user-service публикует следующие события в exchange `app.events`:

| Routing Key | Payload | Назначение |
|-------------|---------|------------|
| `user.created` | [UserCreatedPayload](#usercreatedpayload) | Регистрация нового пользователя |

##### UserCreatedPayload

```json
{
  "user_id": "string",
  "email": "string",
  "name": "string",
  "role": "client | driver | admin",
  "created_at": "2026-01-01T00:00:00Z"
}
```

---

## order-service

Сервис управления заказами. Создание, обновление, фильтрация заказов, назначение менеджеров.
Публикует события в RabbitMQ при создании, обновлении и отмене заказов.

- **HTTP порт:** `8082`
- **gRPC порт:** — (порт `50052` проброшен в docker-compose, но gRPC-сервер не запущен в коде)
- **Framework:** Gin (HTTP)
- **Зависимости:** PostgreSQL, RabbitMQ

### Аутентификация

order-service имеет middleware (`middleware.Middleawre`), которое проверяет заголовки:

| Заголовок | Тип | Описание |
|-----------|-----|----------|
| `X-User-ID` | string | ID пользователя (обязательно) |
| `X-User-Role` | string | Роль пользователя |
| `Authorization` | string | Токен авторизации |

> **Важно:** В текущей версии middleware определено, но **не применено** в `routes.go`. Тем не менее, хэндлеры `CreateOrder` и `AssignManager` используют `c.MustGet("userID").(int)`, что требует применения middleware для корректной работы.

### HTTP REST Эндпоинты

#### Orders — Управление заказами

| Метод | Path | Хэндлер | Назначение | Аутентификация | Входные данные | Успешный ответ |
|-------|------|---------|------------|----------------|----------------|----------------|
| `POST` | `/orders` | `OrderHandler.CreateOrder` | Создание нового заказа | Auth required (`X-User-ID`) | [CreateOrderRequest](#createorderrequest) | `201` `<Order>` |
| `GET` | `/orders/:id` | `OrderHandler.GetOrderByID` | Получение заказа по ID | Public | URL param: `id` (int) | `200` `<Order>` |
| `PUT` | `/orders/:id` | `OrderHandler.UpdateOrder` | Обновление полей заказа | Public | URL param: `id` (int) + [UpdateOrderRequest](#updateorderrequest) | `200` `{"message": "order updated"}` |
| `PATCH` | `/orders/:id/status` | `OrderHandler.UpdateStatus` | Обновление статуса заказа | Public | URL param: `id` (int) + Body: `int` (status) | `200` `{"message": "status updated"}` |
| `POST` | `/orders/:id/manager` | `OrderHandler.AssignManager` | Назначение менеджера на заказ | Auth required (`X-User-ID`) | URL param: `id` (int) | `200` `{"message": "manager assigned"}` |
| `GET` | `/orders` | `OrderHandler.GetOrders` | Получение списка заказов с фильтрацией и пагинацией | Public | Query: [OrderFilter](#orderfilter) + `page` (int, default 1) | `200` `{"data": [<Order>], "meta": {...}}` |

#### Orders by User/Driver — Заказы по пользователю/водителю

| Метод | Path | Хэндлер | Назначение | Аутентификация | Входные данные | Успешный ответ |
|-------|------|---------|------------|----------------|----------------|----------------|
| `GET` | `/users/:user_id/orders` | `OrderHandler.GetOrdersByUser` | Получение заказов пользователя с пагинацией | Public | URL param: `user_id` (int) + Query: `page` (int, default 1) | `200` `{"data": [<Order>], "meta": {...}}` |
| `GET` | `/drivers/:driver_id/orders` | `OrderHandler.GetOrdersByDriver` | Получение заказов водителя с пагинацией | Public | URL param: `driver_id` (int) + Query: `page` (int, default 1) | `200` `{"data": [<Order>], "meta": {...}}` |

##### CreateOrderRequest

```json
{
  "capacity": 1000,
  "lifting_capacity": 500,
  "type": 1,
  "manager_id": 2,
  "status": 1,
  "price": 1500.50,
  "pickup_address": "г. Москва, ул. Ленина 1",
  "delivery_address": "г. Томск, пр. Ленина 40",
  "date_start": "2026-01-15T09:00:00Z",
  "date_end": "2026-01-20T18:00:00Z",
  "period_from": "2026-01-15T09:00:00Z",
  "period_to": "2026-01-15T12:00:00Z"
}
```

| Поле | Тип | Обязательно | Описание |
|------|-----|-------------|----------|
| `capacity` | int | ✅ Да (min=1) | Грузоподъёмность |
| `lifting_capacity` | int | ✅ Да (min=1) | Грузоподъёмность (дублирующее поле) |
| `type` | *int | Нет | Тип заказа |
| `manager_id` | *int | Нет | ID менеджера |
| `status` | *int | Нет | Статус (по умолчанию 1) |
| `price` | *float64 | Нет | Цена |
| `pickup_address` | *string | Нет | Адрес забора груза |
| `delivery_address` | *string | Нет | Адрес доставки |
| `date_start` | *time.Time | Нет | Дата начала |
| `date_end` | *time.Time | Нет | Дата окончания |
| `period_from` | *time.Time | Нет | Начало периода |
| `period_to` | *time.Time | Нет | Конец периода |

##### UpdateOrderRequest

```json
{
  "status": 2,
  "manager_id": 3,
  "price": 2000.00,
  "pickup_address": "новый адрес",
  "delivery_address": "новый адрес",
  "date_start": "2026-01-16T09:00:00Z",
  "date_end": "2026-01-21T18:00:00Z",
  "period_from": "2026-01-16T09:00:00Z",
  "period_to": "2026-01-16T12:00:00Z"
}
```

> Все поля опциональны (указываются только те, которые нужно обновить).

##### OrderFilter (Query параметры)

| Параметр | Тип | Описание |
|----------|-----|----------|
| `user_id` | int | Фильтр по ID пользователя |
| `status` | int | Фильтр по статусу |
| `manager_id` | int | Фильтр по ID менеджера |
| `type` | int | Фильтр по типу заказа |
| `capacity_min` | int | Минимальная грузоподъёмность |
| `capacity_max` | int | Максимальная грузоподъёмность |
| `lifting_capacity_min` | int | Минимальная грузоподъёмность |
| `lifting_capacity_max` | int | Максимальная грузоподъёмность |
| `price_min` | float | Минимальная цена |
| `price_max` | float | Максимальная цена |
| `date_create_from` | time | Дата создания (от) |
| `date_create_to` | time | Дата создания (до) |
| `date_start_from` | time | Дата начала (от) |
| `date_start_to` | time | Дата начала (до) |
| `date_end_from` | time | Дата окончания (от) |
| `date_end_to` | time | Дата окончания (до) |
| `period_from_from` | time | Начало периода (от) |
| `period_from_to` | time | Начало периода (до) |
| `period_to_from` | time | Конец периода (от) |
| `period_to_to` | time | Конец периода (до) |
| `pickup_address` | string | Адрес забора (частичное совпадение) |
| `delivery_address` | string | Адрес доставки (частичное совпадение) |
| `driver_id` | int | Фильтр по ID водителя (через driver_auto) |
| `auto_id` | int | Фильтр по ID автомобиля (через driver_auto) |
| `limit` | int | Лимит записей (default 20) |
| `offset` | int | Смещение (default 0) |
| `page` | int | Номер страницы (default 1, используется вместо limit/offset) |

##### Order (модель ответа)

```json
{
  "id": 1,
  "user_id": 5,
  "capacity": 1000,
  "lifting_capacity": 500,
  "type": 1,
  "manager_id": 2,
  "status": 1,
  "price": 1500.50,
  "pickup_address": "г. Москва, ул. Ленина 1",
  "delivery_address": "г. Томск, пр. Ленина 40",
  "date_start": "2026-01-15T09:00:00Z",
  "date_end": "2026-01-20T18:00:00Z",
  "period_from": "2026-01-15T09:00:00Z",
  "period_to": "2026-01-15T12:00:00Z",
  "date_create": "2026-01-10T12:00:00Z",
  "date_update": "2026-01-10T12:30:00Z"
}
```

##### Response с пагинацией (meta)

```json
{
  "data": [<Order>, ...],
  "meta": {
    "page": 1,
    "pageSize": 20,
    "total": 100,
    "totalPages": 5
  }
}
```

### RabbitMQ События (Publisher)

order-service публикует следующие события в exchange `app.events`:

| Routing Key | Payload | Когда вызывается |
|-------------|---------|------------------|
| `order.created` | [OrderPayload](#orderpayload) | При создании заказа (`POST /orders`) |
| `order.updated` | [OrderPayload](#orderpayload) | При обновлении статуса (`PATCH /orders/:id/status`, status ≠ -1) |
| `order.canceled` | [OrderPayload](#orderpayload) | При отмене заказа (`PATCH /orders/:id/status`, status = -1) |

##### OrderPayload

```json
{
  "order_id": "string",
  "user_id": "string",
  "user_email": "string (omitempty)",
  "cargo_type": "string (omitempty)",
  "weight": 0.0,
  "price": 1500.50,
  "status": "created | in_progress | completed | canceled",
  "timestamp": "2026-01-10T12:00:00Z"
}
```

---

## resource-service

Сервис управления ресурсами. Управляет водителями (создание, статусы, категории) и автомобилями (создание, статусы, фильтрация).

- **HTTP порт:** `8081`
- **gRPC порт:** — (не используется)
- **Framework:** Gin (HTTP)
- **Зависимости:** PostgreSQL

### HTTP REST Эндпоинты

#### Drivers — Управление водителями

| Метод | Path | Хэндлер | Назначение | Аутентификация | Входные данные | Успешный ответ |
|-------|------|---------|------------|----------------|----------------|----------------|
| `POST` | `/drivers` | `DriverHandler.CreateDriver` | Создание профиля водителя | Public | [Driver](#driver-model) | `201` `{"message": "driver created", "id": <int>}` |
| `POST` | `/drivers/categories` | `DriverHandler.CreateDriverCategory` | Добавление категории водителю | Public | [DriverCategory](#drivercategory-model) | `201` `{"message": "driver category added"}` |
| `GET` | `/drivers/:id` | `DriverHandler.GetDriverByID` | Получение водителя по ID | Public | URL param: `id` (int) | `200` `<Driver>` |
| `GET` | `/drivers/:id/categories` | `DriverHandler.GetDriverCategories` | Получение категорий водителя | Public | URL param: `id` (int) | `200` `[<DriverCategory>]` |
| `GET` | `/drivers` | `DriverHandler.GetDriversByFilter` | Получение списка водителей с фильтром | Public | Query: `status` (int), `category` (string) | `200` `[<Driver>]` |
| `PUT` | `/drivers/:id/status` | `DriverHandler.UpdateDriverStatus` | Обновление статуса водителя | Public | URL param: `id` (int) + Body: `{"status": <int 1-3>}` | `200` `{"message": "driver status updated"}` |

#### Autos — Управление автомобилями

| Метод | Path | Хэндлер | Назначение | Аутентификация | Входные данные | Успешный ответ |
|-------|------|---------|------------|----------------|----------------|----------------|
| `POST` | `/autos` | `AutoHandler.CreateAuto` | Создание автомобиля | Public | [Auto](#auto-model) | `201` `{"message": "auto created", "id": <int>}` |
| `GET` | `/autos/:id` | `AutoHandler.GetAutoByID` | Получение автомобиля по ID | Public | URL param: `id` (int) | `200` `<Auto>` |
| `GET` | `/autos` | `AutoHandler.GetAutosWithFilter` | Получение списка автомобилей с фильтром | Public | Query: `capacity` (int), `lifting_capacity` (int), `status` (int) | `200` `[<Auto>]` |
| `PUT` | `/autos/:id/status` | `AutoHandler.UpdateAutoStatus` | Обновление статуса автомобиля | Public | URL param: `id` (int) + Body: `{"status": <int 1-3>}` | `200` `{"message": "status updated"}` |

##### Driver (модель)

```json
{
  "id": 1,
  "user_id": 5,
  "status": 1,
  "date_create": "2026-01-01T00:00:00Z",
  "date_update": "2026-01-01T00:00:00Z"
}
```

##### DriverCategory (модель)

```json
{
  "id": 1,
  "DriverID": 1,
  "category": "C",
  "date_create": "2026-01-01T00:00:00Z",
  "date_update": "2026-01-01T00:00:00Z"
}
```

> **Примечание:** Поле `DriverID` не имеет JSON-тега, поэтому сериализуется с именем `DriverID`.

##### Auto (модель)

```json
{
  "id": 1,
  "status": 1,
  "capacity": 1000,
  "lifting_capacity": 500,
  "number": "А123БВ",
  "required_category": "C",
  "date_create": "2026-01-01T00:00:00Z",
  "date_update": "2026-01-01T00:00:00Z"
}
```

| Поле | Тип | Обязательно | Описание |
|------|-----|-------------|----------|
| `status` | int | Нет | Статус (1-3) |
| `capacity` | int | ✅ Да | Грузоподъёмность |
| `lifting_capacity` | int | ✅ Да | Грузоподъёмность |
| `number` | string | ✅ Да | Гос. номер (до 10 символов) |
| `required_category` | string | ✅ Да | Требуемая категория прав (до 10 символов) |

##### UpdateStatus (Request Body для водителей и авто)

```json
{
  "status": 2
}
```

> Статус должен быть в диапазоне 1-3.

---

## notification-service

Сервис уведомлений. Не имеет HTTP или gRPC эндпоинтов. Работает как RabbitMQ consumer — слушает события от других сервисов и отправляет уведомления через провайдер (MockNotificationProvider).

- **HTTP порт:** — (нет HTTP-сервера)
- **gRPC порт:** — (нет gRPC-сервера)
- **Зависимости:** RabbitMQ

### RabbitMQ Consumer

**Exchange:** `app.events` (topic)  
**Queue:** `notifications_queue`  
**Auto-ack:** `true`

#### Слушаемые события (Routing Keys)

| Routing Key | Payload | Действие |
|-------------|---------|----------|
| `user.created` | [UserCreatedPayload](#usercreatedpayload-1) | Отправка приветственного email: `"Добро пожаловать, {name}!"` |
| `order.created` | [OrderPayload](#orderpayload-1) | Отправка email: `"Ваш заказ #{order_id} успешно создан на сумму {amount}$"` |
| `order.updated` | [OrderPayload](#orderpayload-1) | Отправка email: `"Статус вашего заказа #{order_id} изменен на: {status}"` |
| `driver.order` | [DriverOrderPayload](#driverorderpayload) | Отправка email: `"Вы привязаны к заказу #{order_id}! {message}"` |

##### UserCreatedPayload

```json
{
  "user_id": "string",
  "email": "string",
  "name": "string",
  "created_at": "2026-01-01T00:00:00Z"
}
```

##### OrderPayload

```json
{
  "order_id": "string",
  "user_id": "string",
  "user_email": "string",
  "driver_id": "string (omitempty)",
  "status": "string",
  "amount": 1500.50,
  "updated_at": "2026-01-01T00:00:00Z"
}
```

##### DriverOrderPayload

```json
{
  "driver_id": "string",
  "driver_email": "string",
  "order_id": "string",
  "message": "string",
  "assigned_at": "2026-01-01T00:00:00Z"
}
```

### Провайдер уведомлений

Используется `MockNotificationProvider` — заглушка, которая логирует отправку:

```
[PROVIDER SEND] ──> Кому: {email} | Сообщение: {message}
```

---

## Архитектура взаимодействия сервисов

```
┌──────────────┐     gRPC :50051     ┌──────────────┐
│              │ ──────────────────> │              │
│  auth-service│                     │  user-service│
│  HTTP :8083  │                     │  HTTP :8080  │
│  gRPC :50053 │                     │  gRPC :50051 │
└──────────────┘                     └──────────────┘
       │                                    │
       │                                    │ RabbitMQ
       │                                    │ user.created
       │                                    ▼
       │                             ┌──────────────────┐
       │                             │                  │
       │                             │ notification-svc │
       │                             │ (RabbitMQ consumer│
       │                             │  no HTTP/gRPC)   │
       │                             └──────────────────┘
       │                                    ▲
       │                                    │
┌──────────────┐     RabbitMQ              │
│              │  order.created ───────────┤
│ order-service│  order.updated ───────────┤
│  HTTP :8082  │  order.canceled ──────────┘
└──────────────┘
       │
       │
┌──────────────┐
│              │
│resource-svc  │
│  HTTP :8081  │
└──────────────┘
```

### gRPC взаимодействия

| Клиент | Сервер | Сервис | Методы |
|--------|--------|--------|--------|
| auth-service | user-service:50051 | `UserService` | `CreateUser`, `GetUser` |

### RabbitMQ события

| Publisher | Routing Key | Consumer |
|-----------|-------------|----------|
| user-service | `user.created` | notification-service |
| order-service | `order.created` | notification-service |
| order-service | `order.updated` | notification-service |
| order-service | `order.canceled` | notification-service |

---

## Статусы и коды ошибок

### Общие коды ответов

| Код | Описание |
|-----|----------|
| `200` | OK — успешный запрос |
| `201` | Created — ресурс успешно создан |
| `400` | Bad Request — неверный формат данных |
| `401` | Unauthorized — отсутствует или неверный токен |
| `403` | Forbidden — доступ запрещён |
| `404` | Not Found — ресурс не найден |
| `500` | Internal Server Error — внутренняя ошибка сервера |

### Статусы заказов

| Значение | Описание |
|----------|----------|
| `1` | Новый (по умолчанию) |
| `2` | В работе |
| `3` | Завершён |
| `4` | (зарезервировано) |
| `5` | (зарезервировано) |
| `-1` | Отменён |

### Статусы водителей и автомобилей

| Значение | Описание |
|----------|----------|
| `1` | Активен |
| `2` | Неактивен |
| `3` | Заблокирован |