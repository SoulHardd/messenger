# Messenger API Documentation

Документация описывает текущие HTTP-эндпоинты модулей `auth` и `chat`.

Базовый префикс API: `/api/v1`

Для всех защищенных эндпоинтов требуется заголовок:

```http
Authorization: Bearer <access_token>
```

## Auth

### 1. Login

- Метод: `POST`
- Адрес: `/api/v1/auth/login`

Пример запроса:

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "field": "login",
  "value": "johndoe",
  "password": "Password123"
}
```

Альтернативный вариант входа по телефону:

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "field": "phone",
  "value": "+79991234567",
  "password": "Password123"
}
```

Пример успешного ответа: `200 OK`

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "7c1f66f2d62d4d70a4f1..."
}
```

Ошибки:

- `400 Bad Request`:
    - `{"error":"invalid json"}`
    - `{"error": "missing required fields"}`
    - `{"error": "invalid login"}`
    - `{"error": "invalid phone"}`
    - `{"error": "invalid password"}`
      `401 Unauthorized`:
    - `{"error": "incorrect password"}`
- `404 Not Found`:
    - `{"error": "user not found"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

Примечания:

- Поле `field` принимает только значения `login` или `phone`.
- Логин должен соответствовать формату `[A-Za-z0-9]{3,30}`.
- Телефон должен быть в формате `+79991234567`.
- Пароль должен содержать минимум 8 символов и состоять из `[A-Za-z\\d@$!%*#?&]`.

### 2. Register

- Метод: `POST`
- Адрес: `/api/v1/auth/register`

Пример запроса:

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "phone": "+79991234567",
  "login": "johndoe",
  "password": "Password123"
}
```

Пример успешного ответа: `201 Created`

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "7c1f66f2d62d4d70a4f1..."
}
```

Ошибки:

- `400 Bad Request`:
    - `{"error":"invalid json"}`
    - `{"error": "missing required fields"}`
    - `{"error": "invalid phone"}`
    - `{"error": "invalid login"}`
    - `{"error": "invalid password"}`
- `409 Conflict`:
    - `{"error": "phone is already exists"}`
    - `{"error": "login is already exists"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

### 3. Refresh Tokens

- Метод: `POST`
- Адрес: `/api/v1/auth/refresh`

Пример запроса:

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "7c1f66f2d62d4d70a4f1..."
}
```

Пример успешного ответа: `200 OK`

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "new_refresh_token"
}
```

Ошибки:

- `400 Bad Request`:
    - `{"error":"invalid json"}`
    - `{"error": "invalid token"}`
- `404 Not Found`:
    - `{"error": "session not found"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

## Users

Все методы ниже требуют `Authorization: Bearer <access_token>`.

### 1. Get Current User

- Метод: `GET`
- Адрес: `/api/v1/user/me`

Пример запроса:

```http
GET /api/v1/user/me
Authorization: Bearer <access_token>
```

Пример успешного ответа: `200 OK`

```json
{
  "login": "johndoe",
  "phone": "+79991234567",
  "nickname": "John",
  "bio": "Go developer",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `404 Not Found`:
    - `{"error": "user not found"}`
    - `{"error": "profile not found"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

### 2. Update Profile

- Метод: `PATCH`
- Адрес: `/api/v1/user/me`

Пример запроса:

```http
PATCH /api/v1/user/me
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "nickname": "John",
  "bio": "Backend developer",
  "avatar_url": "https://example.com/new-avatar.jpg"
}
```

Пример успешного ответа: `200 OK`

Тело ответа пустое.

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `400 Bad Request`:
    - `{"error":"invalid json"}`
- `404 Not Found`:
    - `{"error": "profile not found"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`


### 3. Get User By Login

- Метод: `GET`
- Адрес: `/api/v1/user/{login}`

Пример запроса:

```http
GET /api/v1/user/johndoe
Authorization: Bearer <access_token>
```

Пример успешного ответа: `200 OK`

```json
{
  "login": "johndoe",
  "nickname": "John",
  "bio": "Go developer",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `404 Not Found`:
    - `{"error": "user not found"}`
    - `{"error": "profile not found"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

### 4. Search Users

- Метод: `GET`
- Адрес: `/api/v1/user/search`

Параметры запроса:

- `q` required, строка поиска
- `limit` optional, по умолчанию `20`

Пример запроса:

```http
GET /api/v1/user/search?q=john&limit=10
Authorization: Bearer <access_token>
```

Пример успешного ответа: `200 OK`

```json
{
  "users": [
    {
      "login": "johndoe",
      "nickname": "John",
      "bio": "Go developer",
      "avatar_url": "https://example.com/avatar.jpg"
    },
    {
      "login": "johnsmith",
      "nickname": "John Smith",
      "bio": "Team lead",
      "avatar_url": ""
    }
  ]
}
```

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`


## Chat

Все методы ниже требуют `Authorization: Bearer <access_token>`.

### 1. Create Private Chat

- Метод: `POST`
- Адрес: `/api/v1/chat/private`

Пример запроса:

```http
POST /api/v1/chat/private
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "user_id": 12
}
```

Пример успешного ответа: `201 Created`

```json
{
  "id": 45
}
```

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `400 Bad Request`:
    - `{"error":"invalid json"}`
    - `{"error": "invalid input data"}`
- `404 Not Found`:
    - `{"error": "user not found"}`
- `409 Conflict`:
    - `{"error": "private chat between these users already exists"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

### 2. Create Group Chat

- Метод: `POST`
- Адрес: `/api/v1/chat/group`

Пример запроса:

```http
POST /api/v1/chat/group
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Backend Team",
  "users": [2, 3, 4]
}
```

Пример успешного ответа: `201 Created`

```json
{
  "id": 46
}
```

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `400 Bad Request`:
    - `{"error":"invalid json"}`
    - `{"error": "invalid input data"}`
- `404 Not Found`:
    - `{"error": "user not found"}`
- `409 Conflict`:
    - `{"error": "participant already exists"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

### 3. Get User Chats

- Метод: `GET`
- Адрес: `/api/v1/chat/chats`

Параметры запроса:

- `limit` optional, по умолчанию `50`
- `cursor` optional, строка пагинации в Base64

Пример запроса:

```http
GET /api/v1/chat/chats?limit=20&cursor=eyJUaW1lIjoiMjAyNi0wMy0yN1QxMDozMDowMFoiLCJJZCI6NDV9
Authorization: Bearer <access_token>
```

Пример успешного ответа: `200 OK`

```json
{
  "chats": [
    {
      "id": 45,
      "type": "private",
      "title": "private",
      "owner_id": 1,
      "last_msg_text": "hello",
      "last_msg_time": "2026-03-27T10:30:00Z",
      "unread_count": 2
    },
    {
      "id": 46,
      "type": "group",
      "title": "Backend Team",
      "owner_id": 1,
      "last_msg_text": "",
      "last_msg_time": "2026-03-27T09:00:00Z",
      "unread_count": 0
    }
  ],
  "next_cursor": "eyJUaW1lIjoiMjAyNi0wMy0yN1QwOTowMDowMFoiLCJJZCI6NDZ9"
}
```

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `400 Bad Request`:
    - `{"error":"invalid query"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

### 4. Get User Chat By ID

- Метод: `GET`
- Адрес: `/api/v1/chat/{chat_id}`

Пример запроса:

```http
GET /api/v1/chat/45
Authorization: Bearer <access_token>
```

Пример успешного ответа: `200 OK`

```json
{
  "id": 45,
  "type": "private",
  "title": "private",
  "owner_id": 1,
  "last_msg_text": "hello",
  "last_msg_time": "2026-03-27T10:30:00Z",
  "unread_count": 2
}
```

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `404 Not Found`:
    - `{"error": "chat not found"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

Примечание:

- При пустом или некорректном `chat_id` контроллер возвращает необработанную ошибку через `ErrorMapper`, поэтому фактически возможен `500 Internal Server Error` вместо ожидаемого `400`.

### 5. Get Chat Participants

- Метод: `GET`
- Адрес: `/api/v1/chat/{chat_id}/participants`

Пример запроса:

```http
GET /api/v1/chat/45/participants
Authorization: Bearer <access_token>
```

Пример успешного ответа: `200 OK`

```json
[
  {
    "chat_id": 45,
    "user_id": 1,
    "role": "admin"
  },
  {
    "chat_id": 45,
    "user_id": 12,
    "role": "member"
  }
]
```

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `400 Bad Request`:
    - `{"error":"invalid query"}`
- `403 Forbidden`:
    - `{"error": "permission denied"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`


### 6. Add Participant

- Метод: `POST`
- Адрес: `/api/v1/chat/participant`

Пример запроса:

```http
POST /api/v1/chat/participant
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "chat_id": 46,
  "user_id": 12,
  "role": "member"
}
```

Пример успешного ответа: `200 OK`

Тело ответа пустое.

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `400 Bad Request`:
    - `{"error":"invalid json"}`
    - `{"error": "invalid role"}`
- `403 Forbidden`:
    - `{"error": "permission denied"}`
- `404 Not Found`:
    - `{"error": "user not found"}`
    - `{"error": "chat or participant not found"}`
    - `{"error": "chat not found"}`
- `409 Conflict`:
    - `{"error": "participant already exists"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

### 7. Update Participant Role

- Метод: `PATCH`
- Адрес: `/api/v1/chat/participant`

Пример запроса:

```http
PATCH /api/v1/chat/participant
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "chat_id": 46,
  "user_id": 12,
  "role": "admin"
}
```

Пример успешного ответа: `200 OK`

Тело ответа пустое.

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `400 Bad Request`:
    - `{"error":"invalid json"}`
    - `{"error": "invalid role"}`
- `403 Forbidden`:
    - `{"error": "permission denied"}`
- `404 Not Found`:
    - `{"error": "participant not found"}`
    - `{"error": "chat not found"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

### 8. Delete Participant

- Метод: `DELETE`
- Адрес: `/api/v1/chat/participant`

Пример запроса:

```http
DELETE /api/v1/chat/participant
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "chat_id": 46,
  "user_id": 12
}
```

Пример успешного ответа: `200 OK`

Тело ответа пустое.

Ошибки:

- `401 Unauthorized`:
    - `{"error":"missing token"}`
    - `{"error":"invalid token"}`
- `400 Bad Request`:
    - `{"error":"invalid json"}`
- `403 Forbidden`:
    - `{"error": "permission denied"}`
- `404 Not Found`:
    - `{"error": "chat or participant not found"}`
    - `{"error": "chat not found"}`
- `500 Internal Server Error`:
    - `{"error": "internal server error"}`

## Notes

- Модуль `message` в документацию не включен, так как реализация еще не завершена.
