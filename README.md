## Financial System API

REST API приложение для управления простой финансовой системой с ролями **client**, **manager**, **admin**.  
Реализованы счета, вклады, зарплатный проект, аудит действий и логическая отмена действий.

---

## Стек

- Go `1.25.3`
- PostgreSQL `18-bookworm`
- `github.com/gorilla/mux` — Роут
- `github.com/swaggo/swag` + `http-swagger` — Swagger UI

---

## Быстрый старт

### 1. Локальный запуск (без Docker)

1. Установить Go `1.25.3` и PostgreSQL.
2. Создать БД и пользователя под значения из `.env`.
3. Применить миграции из папки `migrations`.
4. Установить зависимости:

```bash
go mod download
```

1. Запустить приложение:

```bash
go run ./cmd/app
```

API будет доступен на `http://localhost:8080`.  
Swagger UI: `http://localhost:8080/swagger/index.html`.

### 2. Запуск через Docker Compose

Требуется Docker и Docker Compose.

```bash
docker-compose up --build
```

- Сервис `db`: PostgreSQL `18-bookworm`, данные сохраняются в volume `postgres_data`.
- Сервис `app`: Go‑приложение, собирается из `Dockerfile`, стартует после успешного healthcheck БД.

Переменные окружения берутся из `.env`. Внутри контейнера приложение подключается к хосту `db`.

---

## Конфигурация

Используется файл `.env` (см. пример в репозитории):

- **HTTP**
  - `HTTP_PORT` — порт HTTP‑сервера (по умолчанию `8080`)
  - `HTTP_TIMEOUT`, `HTTP_IDLE_TIMEOUT` — тайм-ауты
- **Database**
  - `DB_HOST` — хост БД (`localhost` локально, `db` в Docker)
  - `DB_PORT` — порт БД (`5432`)
  - `DB_USER`, `DB_PASSWORD`, `DB_NAME` — доступ к PostgreSQL
- **JWT**
  - `JWT_SECRET` — секрет для подписания токенов
  - `JWT_EXPIRE` — время жизни токена (например, `24h`)

---

## Роли и функционал

### Auth (без JWT)

- `POST /auth/sign-up` — регистрация клиента (создаёт пользователя с ролью `client`, `is_active=false`, требует подтверждения менеджером).
- `POST /auth/sign-in` — вход по `email`/`password`, возвращает JWT `Bearer <token>` (только для активных клиентов).

### Client (роль `client`, требуется Bearer JWT)

- **Банки и предприятия**
  - `GET /client/banks` — список банков.
  - `GET /client/enterprises` — список предприятий.
- **Счета**
  - `POST /client/accounts` — открыть счёт в банке.
  - `DELETE /client/accounts/{id}` — закрыть счёт (баланс должен быть 0).
  - `POST /client/accounts/transfer` — перевод со счёта на счёт или вклад того же пользователя.
  - `GET /client/accounts/history?account_id=` — история операций по счёту.
- **Вклады**
  - `POST /client/deposits` — открыть вклад.
  - `DELETE /client/deposits/{id}` — закрыть вклад (баланс 0).
  - `POST /client/deposits/transfer` — перевод с вклада на счёт или вклад.
  - `POST /client/deposits/{id}/accumulate` — пополнение вклада со счёта.
- **Зарплатный проект**
  - `POST /client/salary-project/apply` — подать заявку на зарплатный проект от предприятия.
  - `POST /client/salary-project/receive` — получить зарплату по одобренной заявке на выбранный счёт/вклад.

### Manager (роль `manager`, требуется Bearer JWT)

- `POST /manager/users/{id}/approve` — подтвердить регистрацию клиента (`is_active=true`).
- `GET /manager/users/{id}/history` — объединённая история операций по всем счетам пользователя.
- `POST /manager/accounts/{id}/block` / `/unblock` — заблокировать/разблокировать счёт.
- `POST /manager/deposits/{id}/block` / `/unblock` — заблокировать/разблокировать вклад.
- `GET /manager/enterprises` — предприятия со списками сотрудников.
- `POST /manager/enterprises/{id}/employees` — добавить клиента как сотрудника предприятия.
- `DELETE /manager/enterprises/{enterprise_id}/employees/{user_id}` — удалить сотрудника (его pending‑заявки по предприятию отклоняются).
- `POST /manager/salary-project/applications/{id}/approve` — одобрить заявку на зарплатный проект (при недостатке средств предприятия вернёт ошибку).

### Admin (роль `admin`, требуется Bearer JWT)

- `GET /admin/logs` — все логи действий (таблица `action_logs`, сортировка по дате).
- `POST /admin/logs/{id}/undo` — логический откат конкретного действия по записи лога (поддерживаются действия клиентов и менеджеров; для `auth_sign_`* недоступно).
- `POST /admin/logs/undo-all` — откат всех отменяемых действий (по убыванию времени, пропуская уже отменённые и `auth_sign_`*).

Полный контракт и схемы запросов/ответов можно посмотреть в Swagger UI (`/swagger/index.html`).

---

## Undo‑механизм

Каждое важное действие клиента/менеджера логируется в таблицу `action_logs` с типом действия и деталями.  
Администратор может:

- откатить **одно** действие по id лога (`POST /admin/logs/{id}/undo`),
- откатить **все** возможные действия (`POST /admin/logs/undo-all`).

