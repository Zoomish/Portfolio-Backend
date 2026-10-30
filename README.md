# Portfolio backend

Go-сервис для сайта-визитки. На билде Vercel читает последний успешный снимок GitHub и LinkedIn из Postgres. Redis нет. Telegram-бот в этот процесс не входит.

## Эндпоинты

- `GET /api/github/repos`
- `GET /api/linkedin/all`
- `GET /api/linkedin/posts`
- `GET /api/linkedin/profile`
- `GET /api/linkedin/contact`
- `GET /api/app` — процесс жив, без проверки базы
- `GET /api/ready` — база отвечает

Если снимка ещё нет, контентные ручки отвечают `503`. Билд фронта в этом случае оставляет предыдущий `site-content.json`.

## Запуск

```bash
go run ./cmd/api
```

```bash
make test
make build
```

Порт берётся из `PORT`, по умолчанию `3000`. Миграции применяются при старте.

## Переменные

Обязательные: `GITHUB_TOKEN`, `RAPIDAPI_KEY`, `RAPIDAPI_HOST`, `POSTGRES_HOST`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DATABASE`.

`POSTGRES_PORT` по умолчанию `5432`. В `APP_ENV=prod` (или `NODE_ENV=prod`) соединение с Postgres идёт по SSL без проверки сертификата. Для локальной базы без SSL: `POSTGRES_SSL=disable` или `APP_ENV=dev`.

Необязательные:

- `LINKEDIN_USERNAME` — по умолчанию `zoomish`
- `CONTENT_REFRESH_INTERVAL` — по умолчанию `6h`
- `DB_KEEPALIVE_INTERVAL` — по умолчанию `15m`
- `DB_KEEPALIVE_ENABLED` — `false` выключает keepalive
- `TELEGRAM_TOKEN` и `ADMIN_TELEGRAM_ID` — алерт после повторяющихся ошибок refresh и keepalive
- `ALERT_AFTER_FAILURES` — по умолчанию `3`

Секреты только из окружения или локального `.env`. `.env` не коммитится.

## Деплой

Render собирает этот сервис из `Dockerfile`, не через `npm start`. URL сервиса не меняйте.

Секреты те же, что локально: `GITHUB_TOKEN`, `RAPIDAPI_KEY`, `RAPIDAPI_HOST`, `POSTGRES_*`, `NODE_ENV=prod`. `REDIS_URL` не используется.
