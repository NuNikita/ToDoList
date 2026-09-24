# ToDoList

Небольшой API для личного списка задач на Go и PostgreSQL. Можно зарегистрироваться, войти и работать со своими задачами: добавлять, смотреть, менять описание, отмечать выполненными и удалять.

## Запуск

Нужны Go 1.26 и PostgreSQL. Создайте базу данных `ToDoBase` или укажите другое имя в настройках. Таблицы приложение создаст само.

В корне проекта создайте файл `.env` по образцу [`.env.example`](.env.example):

```dotenv
DATABASE_URL=postgres://postgres:your_password@localhost:5432/ToDoBase
JWT_SECRET=your_long_random_secret
```

Замените пароль БД и секрет на свои значения. Затем из корня проекта запустите:

```text
go run .
```

Сервер будет доступен по адресу `http://localhost:8080`. Файл `.env` читается при запуске и не попадает в Git.

## Запросы к API

Для тела запроса используется JSON с заголовком `Content-Type: application/json`. Сначала зарегистрируйтесь через `POST /register`, затем войдите через `POST /login`. В ответ на вход придёт токен:

```json
{"token":"..."}
```

Во всех запросах к задачам передавайте заголовок `Authorization: Bearer <token>`. Токен действует 24 часа. Каждый пользователь видит только свои задачи.

| Метод и адрес | Тело запроса | Ответ |
| --- | --- | --- |
| `POST /register` | `{"login":"demo","password":"example-password"}` | `201 Created`, без тела |
| `POST /login` | `{"login":"demo","password":"example-password"}` | JSON с полем `token` |
| `GET /tasks` | Нет | JSON со списком задач |
| `GET /tasks/{id}` | Нет | JSON с одной задачей |
| `POST /tasks` | `{"title":"Купить молоко","description":"По дороге домой"}` | `201 Created`, текстовое подтверждение |
| `PATCH /tasks/{id}/description` | `{"description":"Новое описание"}` | Текстовое подтверждение |
| `PATCH /tasks/{id}/complete` | Нет | JSON с задачей после изменения |
| `DELETE /tasks/{id}` | Нет | `204 No Content` |

`{id}` — числовой идентификатор задачи. `PATCH /tasks/{id}/complete` переключает состояние: повторный запрос снимет отметку о выполнении. Задача в ответе содержит `id`, `title`, `description`, `completed` и `created_at`.

Тесты можно запустить командой `go test ./...`.
