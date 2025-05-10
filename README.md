# Family Finance Backend

Бэкенд-сервис для управления семейными финансами. Предоставляет API для авторизации пользователей и управления их данными.

## Технологии

- Go 1.21+
- PostgreSQL
- Redis
- JWT для авторизации
- GORM для работы с базой данных

## Структура проекта

```
.
├── config/             # Конфигурация приложения
├── internal/           # Внутренний код приложения
│   ├── db/            # Инициализация и конфигурация базы данных
│   ├── handlers/      # HTTP обработчики
│   ├── middleware/    # Промежуточное ПО
│   ├── models/        # Модели данных
│   ├── repository/    # Слой доступа к данным
│   ├── service/       # Бизнес-логика
│   └── util/          # Вспомогательные функции
└── main.go            # Точка входа в приложение
```

## API Endpoints

### Авторизация

#### Запрос кода для входа
```http
POST /auth/login
Content-Type: application/json

{
    "email": "user@example.com"
}
```
Ответ:
```json
{
    "temp_id": "uuid-временного-идентификатора"
}
```

#### Проверка кода входа
```http
POST /auth/login/verify
Content-Type: application/json

{
    "temp_id": "uuid-временного-идентификатора",
    "code": "123456"
}
```
Ответ:
```json
{
    "token": "jwt-токен"
}
```

#### Запрос кода для регистрации
```http
POST /auth/register
Content-Type: application/json

{
    "email": "newuser@example.com"
}
```
Ответ:
```json
{
    "temp_id": "uuid-временного-идентификатора"
}
```

#### Подтверждение регистрации
```http
POST /auth/register/verify
Content-Type: application/json

{
    "temp_id": "uuid-временного-идентификатора",
    "code": "123456",
    "name": "Иван",
    "surname": "Иванов",
    "nickname": "ivan" // опционально
}
```
Ответ:
```json
{
    "token": "jwt-токен"
}
```

#### Выход из системы
```http
POST /auth/logout
Authorization: Bearer <jwt-токен>
```
Ответ:
```json
{
    "message": "Вы успешно вышли из аккаунта"
}
```

### Пользователь

#### Получение данных пользователя
```http
GET /user/me
Authorization: Bearer <jwt-токен>
```
Ответ:
```json
{
    "id": 1,
    "name": "Иван",
    "surname": "Иванов",
    "nickname": "ivan",
    "email": "user@example.com",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```

#### Обновление данных пользователя
```http
PUT /user/update
Authorization: Bearer <jwt-токен>
Content-Type: application/json

{
    "name": "Новое имя",
    "surname": "Новая фамилия",
    "nickname": "новый_никнейм"
}
```
Ответ:
```json
{
    "id": 1,
    "name": "Новое имя",
    "surname": "Новая фамилия",
    "nickname": "новый_никнейм",
    "email": "user@example.com",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:30:00Z"
}
```

#### Поиск пользователя по email
```http
GET /user/search?email=user@example.com
Authorization: Bearer <jwt-токен>
```
Ответ:
```json
{
    "id": 1,
    "name": "Иван",
    "surname": "Иванов",
    "nickname": "ivan",
    "email": "user@example.com",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
}
```

## Модели данных

### User
```go
type User struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"size:100;not null" json:"name"`
    Surname   string    `gorm:"size:100;not null" json:"surname"`
    Nickname  string    `gorm:"size:100;not null" json:"nickname"`
    Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Конфигурация

Приложение использует переменные окружения для конфигурации. Создайте файл `.env` в корне проекта:

```env
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=family_finance

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## Запуск приложения

1. Установите зависимости:
```bash
go mod download
```

2. Создайте и настройте файл `.env`

3. Запустите приложение:
```bash
go run main.go
```

Приложение будет доступно по адресу `http://localhost:8080`

## Безопасность

- Все пароли и секретные ключи должны храниться в переменных окружения
- JWT токены имеют срок действия 24 часа
- Коды подтверждения действительны 90 секунд
- Все запросы, кроме регистрации и входа, требуют JWT токен
- При выходе из системы токен добавляется в черный список

## Обработка ошибок

Все ошибки возвращаются в формате:
```json
{
    "error": "Краткое описание ошибки",
    "details": "Подробное описание ошибки"
}
```

Коды ответов:
- 200: Успешный запрос
- 400: Неверный запрос
- 401: Ошибка авторизации
- 500: Внутренняя ошибка сервера 