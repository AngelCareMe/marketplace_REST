# REST API Маркетплейс

## Обзор

REST API Маркетплейс — это серверное приложение на Go для управления пользователями и их постами в маркетплейсе. Оно предоставляет endpoints для регистрации пользователей, аутентификации и операций CRUD для постов с поддержкой пагинации, сортировки и фильтрации. Приложение использует PostgreSQL для хранения данных, Docker для контейнеризации и фреймворк Gin для обработки HTTP-запросов.

## Возможности

- **Управление пользователями**:
  - Регистрация и вход пользователей с использованием JWT-аутентификации.
  - Получение, обновление и удаление профилей пользователей.
- **Управление постами**:
  - Создание, получение, обновление и удаление постов.
  - Список постов с пагинацией, сортировкой (по `created_at` или `price`) и фильтрацией (по `min_price` и `max_price`).
  - Список постов конкретного пользователя.
  - Обеспечение уникальности постов по `header`, `content` и `author_id`.
- **Безопасность**:
  - Аутентификация на основе JWT для защищённых маршрутов.
  - Проверка прав доступа, чтобы пользователи могли изменять только свои посты или профили.
- **Обработка ошибок**:
  - Корректные HTTP-статусы (200, 201, 400, 401, 403, 404, 500).
  - Подробные сообщения об ошибках для клиентов.
- **Логирование**:
  - Структурированное логирование с использованием Logrus для отладки и мониторинга.

## Технологический стек

- **Бэкенд**: Go (фреймворк Gin)
- **База данных**: PostgreSQL
- **ORM**: Squirrel (построитель SQL-запросов)
- **Аутентификация**: JWT
- **Контейнеризация**: Docker, Docker Compose
- **Тестирование**: Коллекция Postman для тестирования API
- **Логирование**: Logrus
- **Зависимости**:
  - `github.com/Masterminds/squirrel` для построения SQL-запросов
  - `github.com/jackc/pgx/v5` для драйвера PostgreSQL
  - `github.com/sirupsen/logrus` для логирования
  - `github.com/gin-gonic/gin` для маршрутизации HTTP

## Структура проекта

```
marketplace/
├── adapter/
│   ├── post_adapter.go   # Операции с базой данных для постов
│   ├── user_adapter.go   # Операции с базой данных для пользователей
├── entity/
│   ├── post.go           # Определение сущности поста
│   ├── user.go           # Определение сущности пользователя
├── handler/
│   ├── auth_handler.go   # Маршруты аутентификации (регистрация, вход)
│   ├── post_handler.go   # Маршруты для постов
│   ├── user_handler.go   # Маршруты для пользователей
├── service/
│   ├── post_service.go   # Бизнес-логика для постов
│   ├── user_service.go   # Бизнес-логика для пользователей
├── main.go               # Точка входа приложения
├── docker-compose.yml    # Конфигурация Docker Compose
├── Dockerfile            # Конфигурация Docker для приложения Go
├── go.mod                # Зависимости Go-модуля
├── go.sum                # Контрольные суммы Go-модуля
└── README.md             # Этот файл
```

## Установка и настройка

### Требования

- **Docker** и **Docker Compose**
- **Go** (1.21 или новее, если запускаете без Docker)
- **Postman** (опционально, для тестирования API)

### Установка

1. **Клонируйте репозиторий**:
   ```bash
   git clone <repository-url>
   cd marketplace
   ```

2. **Настройте переменные окружения**:
   Создайте файл `.env` в корне проекта со следующими параметрами:
   ```env
   POSTGRES_USER=user
   POSTGRES_PASSWORD=password
   POSTGRES_DB=marketplace
   POSTGRES_HOST=postgres
   POSTGRES_PORT=5432
   JWT_SECRET=your_jwt_secret_key
   ```

3. **Запустите с помощью Docker Compose**:
   ```bash
   docker-compose up --build
   ```
   Это запустит приложение Go (`marketplace:8080`) и базу данных PostgreSQL.

4. **Проверьте запуск**:
   - Убедитесь, что API работает: `curl http://localhost:8080/health`
   - Просмотрите логи: `docker logs marketplace_rest-marketplace-1`

### Схема базы данных

Приложение использует две таблицы:

- **users**:
  ```sql
  CREATE TABLE users (
      id UUID PRIMARY KEY,
      username TEXT UNIQUE NOT NULL,
      hashed_password TEXT NOT NULL,
      created_at TIMESTAMP NOT NULL
  );
  ```

- **posts**:
  ```sql
  CREATE TABLE posts (
      id UUID PRIMARY KEY,
      header TEXT NOT NULL,
      content TEXT NOT NULL,
      image TEXT,
      price DOUBLE PRECISION NOT NULL,
      author_id UUID NOT NULL,
      created_at TIMESTAMP NOT NULL,
      FOREIGN KEY (author_id) REFERENCES users(id),
      CONSTRAINT unique_post UNIQUE (header, content, author_id)
  );
  ```

Для инициализации базы данных выполните следующий SQL в контейнере PostgreSQL:
```bash
docker exec -it marketplace_rest-postgres-1 psql -U user -d marketplace -c "<вышеуказанный SQL>"
```

## Endpoints API

### Аутентификация
- **POST /users/register**: Регистрация нового пользователя.
  - Тело: `{"username": "string", "password": "string"}`
  - Ответ: `200 OK` с данными пользователя и JWT-токеном
- **POST /users/login**: Вход пользователя.
  - Тело: `{"username": "string", "password": "string"}`
  - Ответ: `200 OK` с данными пользователя и JWT-токеном

### Пользователи
- **GET /users/:id**: Получение пользователя по ID (требуется JWT).
  - Ответ: `200 OK` или `404 Not Found`
- **PUT /users/:id**: Обновление пользователя (требуется JWT, право владения).
  - Тело: `{"username": "string", "password": "string"}`
  - Ответ: `200 OK` или `403 Forbidden`
- **DELETE /users/:id**: Удаление пользователя (требуется JWT, право владения).
  - Ответ: `200 OK` или `403 Forbidden`

### Посты
- **POST /posts**: Создание поста (требуется JWT).
  - Тело: `{"header": "string", "content": "string", "image": "string", "price": number}`
  - Ответ: `201 Created` или `400 Bad Request` (например, при дублировании поста)
- **GET /posts/:id**: Получение поста по ID.
  - Ответ: `200 OK` или `404 Not Found`
- **PUT /posts/:id**: Обновление поста (требуется JWT, право владения).
  - Тело: `{"header": "string", "content": "string", "image": "string", "price": number}`
  - Ответ: `200 OK`, `403 Forbidden` или `400 Bad Request`
- **DELETE /posts/:id**: Удаление поста (требуется JWT, право владения).
  - Ответ: `200 OK` или `404 Not Found`
- **GET /posts**: Список всех постов с пагинацией, сортировкой и фильтрацией.
  - Параметры: `page=<int>&pageSize=<int>&sortBy=<created_at|price ASC|DESC>&min_price=<float>&max_price=<float>`
  - Ответ: `200 OK` с постами и общим количеством
- **GET /users/:id/posts**: Список постов по ID пользователя с пагинацией, сортировкой и фильтрацией.
  - Параметры: `page=<int>&pageSize=<int>&sortBy=<created_at|price ASC|DESC>&min_price=<float>&max_price=<float>`
  - Ответ: `200 OK` с постами и общим количеством или `404 Not Found` (пользователь не найден)

## Тестирование

### Использование Postman
1. Импортируйте файл `Marketplace_API_Tests.postman_collection.json` (предоставляется отдельно) в Postman.
2. Настройте окружение с переменными:
   - `base_url`: `http://localhost:8080`
   - `auth_token`: (заполняется после входа)
   - `user_id`: (заполняется после регистрации/входа)
   - `post_id`: (заполняется после создания поста)
3. Выполните коллекцию в порядке:
   - Регистрация пользователя
   - Вход пользователя
   - Получение пользователя
   - Обновление пользователя
   - Создание поста
   - Получение поста
   - Обновление поста
   - Удаление поста
   - Список постов
   - Список постов по пользователю
   - Список постов для несуществующего пользователя
   - Удаление пользователя

### Использование curl
Примеры команд:
```bash
# Регистрация пользователя
curl -X POST http://localhost:8080/users/register -H "Content-Type: application/json" -d '{"username":"testuser3","password":"Test1234!"}'

# Вход
curl -X POST http://localhost:8080/users/login -H "Content-Type: application/json" -d '{"username":"testuser3","password":"Test1234!"}'

# Создание поста
curl -X POST http://localhost:8080/posts -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d '{"header":"Test Post","content":"Content","image":"https://example.com/image.jpg","price":99.99}'

# Список постов
curl -X GET http://localhost:8080/posts?page=1&pageSize=10&sortBy=created_at%20DESC

# Список постов по пользователю
curl -X GET http://localhost:8080/users/302dfa9d-eabb-4a9d-b365-e958d113fbab/posts?page=1&pageSize=10&sortBy=created_at%20DESC -H "Authorization: Bearer <token>"
```

## Отладка

- **Проверка логов**:
  ```bash
  docker logs marketplace_rest-marketplace-1
  ```
- **Доступ к PostgreSQL**:
  ```bash
  docker exec -it marketplace_rest-postgres-1 psql -U user -d marketplace
  ```
  Полезные команды:
  ```sql
  SELECT * FROM users WHERE id = '302dfa9d-eabb-4a9d-b365-e958d113fbab';
  SELECT * FROM posts WHERE author_id = '302dfa9d-eabb-4a9d-b365-e958d113fbab';
  \d posts
  ```

- **Пересборка без кэша** при неприменении изменений:
  ```bash
  docker-compose down
  docker-compose build --no-cache
  docker-compose up
  ```

## Лицензия

Этот проект распространяется под лицензией MIT.
