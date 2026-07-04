# Spoony — Telegram Mini App

**Spoony** is a Telegram Mini App that helps children learn English through Russian.

### Pages
1. **Главная / Home** — a Facebook-style news feed with tips and updates from Spoony.
2. **Тренажёр слов / Word Trainer** — a journey across England where every city is a new exercise. Each exercise is a *"Find the pair"* memory game: tap cards to match a picture with its English word. A correct match triggers a green `box-shadow` glow animation. Cities unlock one after another as you complete them.
3. **Профиль / Profile** — user info plus learning stats (cities completed, words learned, stars earned).

Learning content (news, levels/cities, vocabulary card pairs, and per-user progress) lives in PostgreSQL — see [`backend/database/migrations.go`](backend/database/migrations.go) for the schema and seed data, and [`backend/handlers/content.go`](backend/handlers/content.go) for the API.

### Works inside **and** outside Telegram
- **Inside Telegram** — one-tap login via Telegram `initData`.
- **As a regular website** — a **guest login** (`POST /api/auth/guest`) creates an account tied to a persistent `guest_id` stored in the browser, so a returning visitor keeps their progress. Optional display name. Works in production.
- A separate dev-only bypass signs you in as a **demo kid** when `initData` can't be verified (disabled when `ENV=production`).

### Responsive
The UI adapts to screen size: a bottom tab bar on mobile, and a left **sidebar** with wider content layouts on desktop/laptop (`md:` breakpoint and up). The "Find the pair" grid grows from 3 columns on phones to 6 on large screens.

---

Built with a Go backend, Vue 3 frontend, PostgreSQL, and an admin panel.

## Tech Stack

| Layer     | Technology                                     |
|-----------|------------------------------------------------|
| Backend   | Go 1.23 + Gin                                  |
| Database  | PostgreSQL 15                                  |
| Frontend  | Vue 3 + Vite + Tailwind CSS v4 + shadcn-vue    |
| Container | Docker — single all-in-one container           |
| CI/CD     | GitHub Actions → ghcr.io                       |

All three services (PostgreSQL, Go backend, nginx) run inside a single Docker container managed by supervisord.

---

## Getting Started

### Prerequisites

- Docker & Docker Compose
- A Telegram Bot token — create via [@BotFather](https://t.me/BotFather)
- A [ngrok authtoken](https://dashboard.ngrok.com/get-started/your-authtoken) (free, dev only)

### 1. Configure environment

```bash
cp .env.example .env
```

Minimum required values:

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_BOT_USERNAME=your_bot_username
JWT_SECRET=some-random-32-char-string-here

# Dev only
NGROK_AUTHTOKEN=your_ngrok_authtoken_here
```

### 2. Development

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

This starts:
1. The app container (PostgreSQL + Go backend with hot reload + nginx)
2. An ngrok sidecar that creates a public HTTPS tunnel
3. The backend discovers the tunnel URL and registers the Telegram webhook automatically

Confirm it's working:
```
Discovered ngrok URL: https://xxxx.ngrok-free.app
Telegram webhook registered: https://xxxx.ngrok-free.app/api/webhook/telegram
```

Available locally:
- **App**: http://localhost:80
- **ngrok inspector**: http://localhost:4040

> ngrok URLs change on every restart — webhook re-registers automatically.

### 3. Production

Set your public domain in `.env`:

```env
WEBHOOK_URL=https://your-domain.com
TELEGRAM_MINI_APP_URL=https://your-domain.com
ADMIN_TOKEN=strong-random-token
```

```bash
docker compose up --build -d
```

---

## Project Structure

```
.
├── backend/
│   ├── config/          # Environment config
│   ├── database/        # DB connection & migrations
│   ├── handlers/        # HTTP handlers (auth, webhook, admin)
│   ├── middleware/       # JWT auth, admin token, CORS
│   ├── models/          # Data models (User)
│   ├── services/        # Telegram initData verification
│   └── main.go          # Router & entrypoint
├── frontend/
│   └── src/
│       ├── components/ui/   # shadcn-vue primitives
│       ├── locales/         # i18n (en, ru)
│       ├── services/        # API client, Telegram SDK
│       ├── store/           # Pinia stores (user, ui, admin)
│       ├── views/           # Pages (Home, Profile, Settings)
│       └── views/admin/     # Admin panel (Login, Dashboard, Users)
├── .github/workflows/   # CI/CD (production deploy)
├── Dockerfile           # Multi-stage: go-builder, node-builder, alpine runtime
├── docker-compose.yml   # Production
├── docker-compose.dev.yml  # Dev overlay (ngrok sidecar + hot reload)
├── docker-compose.prod.yml # Server deployment (pulls image from ghcr.io)
├── docker-entrypoint.sh # Initializes PostgreSQL, then starts supervisord
├── nginx.conf           # Prod: serves SPA + proxies /api, /admin/api
├── nginx.dev.conf       # Dev: proxies to Vite dev server
└── supervisord.conf     # Manages postgres + backend + nginx in one container
```

---

## API

### Auth
| Method | Path                       | Auth | Description        |
|--------|----------------------------|------|--------------------|
| POST   | `/api/auth/telegram-login` |      | Login via Telegram |
| POST   | `/api/auth/logout`         |      | Logout             |
| GET    | `/api/auth/me`             | JWT  | Get current user   |
| PUT    | `/api/auth/profile`        | JWT  | Update profile     |

### Telegram Webhook
| Method | Path                    | Description                    |
|--------|-------------------------|--------------------------------|
| POST   | `/api/webhook/telegram` | Receives updates from Telegram |
| GET    | `/api/webhook/info`     | Current webhook status         |

### Admin (requires `X-Admin-Token` header)
| Method | Path                   | Description      |
|--------|------------------------|------------------|
| GET    | `/admin/api/stats`     | User stats       |
| GET    | `/admin/api/users`     | List users       |
| DELETE | `/admin/api/users/:id` | Delete user      |

---

## Environment Variables

| Variable                | Required   | Description                                        |
|-------------------------|------------|----------------------------------------------------|
| `TELEGRAM_BOT_TOKEN`    | Yes        | Bot token from @BotFather                          |
| `TELEGRAM_BOT_USERNAME` | No         | Bot username (without @); powers the "Continue in Telegram" link |
| `TELEGRAM_MINI_APP_URL` | Yes        | Public HTTPS URL of the Mini App frontend          |
| `JWT_SECRET`            | Yes        | Random string, min 32 chars                        |
| `ADMIN_TOKEN`           | Yes (prod) | Token for admin panel (`X-Admin-Token` header)     |
| `WEBHOOK_URL`           | Prod       | Public HTTPS base URL for webhook registration     |
| `NGROK_AUTHTOKEN`       | Dev        | ngrok authtoken for automatic tunnel               |
| `JWT_EXPIRATION_HOURS`  | No         | Token lifetime in hours (default: 720)             |
| `DATABASE_URL`          | No         | Overrides default postgres connection string       |
| `CORS_ALLOWED_ORIGINS`  | No         | Comma-separated allowed origins                    |
| `PORT`                  | No         | Backend port (default: 8080)                       |
| `ENV`                   | No         | `development` or `production`                      |

---

## Git Flow

```
main  ─────────────────────────────────────────► release/x.x.x
  │                                                     │
  ├── feat/my-feature                                   │
  ├── fix/some-bug            (merge via PR)            │
  └── chore/update-deps                                 │
                                              tag + deploy prod
```

### Branches

| Branch      | Purpose                                                 |
|-------------|---------------------------------------------------------|
| `main`      | Main branch. Gated by CI; merge here doesn't deploy anywhere |
| `feat/*`    | New features                                           |
| `fix/*`     | Bug fixes                                              |
| `chore/*`   | Dependency updates, config changes                     |
| `release/*` | Release branch. Push → auto-deploy to **production**   |

### Rules

1. Direct commits to `main` and `release/*` are **forbidden** — use Pull Requests.
2. Always branch from up-to-date `main`:
   ```bash
   git checkout main && git pull
   git checkout -b feat/your-feature
   ```
3. Delete branches after merge.
4. For a release, create `release/x.x.x` from `main` — CI extracts the version from the branch name.

---

## Deploy

### Production

**Trigger:** push to `release/x.x.x` (or manual `workflow_dispatch`).

**What happens:**
1. CI extracts the version from the branch name.
2. Builds and pushes the image to `ghcr.io` with tags `latest` and `x.x.x`.
3. Connects to the server via SSH, generates `.env` from GitHub Secrets, and runs
   `docker compose --project-name spoony_prod up -d` in `/opt/spoony`.

### Environment diagram

```
Developer → PR → main → [CI]
                  │
                  └── release/1.0.0 → [Deploy] → Production (ghcr.io:latest)
```

### Required GitHub Secrets

Repo-level (shared, same for every environment) or environment-level (`production`) —
values that only apply to production go on the `production` environment.

| Secret                  | Description                                        |
|-------------------------|-----------------------------------------------------|
| `SERVER_HOST`           | Server IP or hostname                              |
| `SERVER_USER`           | SSH user                                            |
| `SERVER_SSH_KEY`        | Private SSH key (matches a key in the server's `~/.ssh/authorized_keys`) |
| `TELEGRAM_BOT_TOKEN`    | Bot token from @BotFather                          |
| `TELEGRAM_BOT_USERNAME` | Bot username (without @)                           |
| `JWT_SECRET`            | Random string, min 32 chars                        |
| `ADMIN_TOKEN`           | Token for the admin panel (`X-Admin-Token` header) |
| `MINI_APP_URL`          | Public HTTPS URL of the Mini App frontend          |
| `WEBHOOK_URL`           | Public HTTPS base URL for Telegram webhook registration |
| `CORS_ALLOWED_ORIGINS`  | Comma-separated allowed origins                    |
| `EXPOSE_PORT`           | Host port to bind (optional, defaults to `80`)     |

### Manual deploy (emergency)

```bash
docker pull ghcr.io/<org>/spoony:latest
docker compose --project-name spoony_prod up -d
```
