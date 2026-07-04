# CLAUDE_MEMORY.md

Выжимка проделанной работы по проекту, по дням. Цель — быстро восстановить контекст в новой сессии без перечитывания всей истории git/чата.

Как пользоваться:
- В конце дня (или сессии) добавляй новую запись сверху под заголовком `## YYYY-MM-DD`.
- Кратко: что сделано, что сломано/недоделано, что дальше.
- Не нужно пересказывать код — только суть решений и открытые вопросы.

---

## О проекте

**Spoony** — Telegram Mini App для обучения детей английскому через русский.
Построен на боилерплейте для TMA: Go 1.23 + Gin (backend), Vue 3 + Vite + Tailwind v4 + shadcn-vue (frontend), PostgreSQL 15, всё в одном Docker-контейнере (supervisord: postgres + backend + nginx). CI/CD через GitHub Actions → ghcr.io.

Страницы:
1. **Главная** — лента новостей в стиле Facebook.
2. **Тренажёр слов** — путешествие по городам Англии, в каждом городе игра "найди пару" (карточка-картинка + карточка-слово), правильное совпадение — зелёная подсветка. Города открываются последовательно.
3. **Профиль** — данные пользователя + статистика (пройдено городов, выучено слов, звёзды).

Авторизация: Telegram `initData` внутри Telegram; гостевой вход (`POST /api/auth/guest`) вне Telegram, привязан к `guest_id` в браузере; dev-бypass с demo-пользователем, отключён при `ENV=production`.

Адаптивность: нижний таб-бар на мобильных, левый сайдбар на десктопе (`md:` breakpoint), сетка игры "найди пару" от 3 до 6 колонок.

---

## Правила разработки (постоянные, не устаревают с новыми записями ниже)

**Вся бизнес-логика фич — только в папке `features/`, и в backend, и в frontend.**
Новый код для новой фичи создавать сразу внутри `features/<название>/`, ничего не класть
рядом в старые "технические" папки (`handlers/`, `models/`, `views/`, `store/` и т.п. — они
остаются только для инфраструктуры, общей для всех фич).

Backend (`backend/features/<name>/`, package = `<name>`):
- `handler.go` — HTTP-хендлеры (Gin), `model.go` — структуры данных фичи.
- Общее/инфраструктурное — вне `features/`: `config/`, `database/` (подключение + миграции),
  `middleware/` (только по-настоящему сквозные вещи вроде CORS и JWT-парсинга — `AuthRequired`).
  Если middleware специфичен для одной фичи (пример: `AdminAuth`) — он живёт внутри той фичи, не в `middleware/`.
- `main.go` знает про все фичи и просто их монтирует — сам бизнес-логики не содержит.

Frontend (`frontend/src/features/<name>/`):
- Внутри фичи: её `*View.vue` (страницы), `api.js` (обёртка над HTTP-эндпоинтами фичи),
  `store.js` (если есть Pinia store), и `components/` для компонентов, специфичных только для неё.
- Общее/переиспользуемое — вне `features/`: `components/` (в т.ч. `components/ui/` — примитивы
  shadcn-vue, и `components/games/` — виджеты, используемые несколькими фичами, например
  `WordBuildGame.vue` используется и trainer, и alphabet), `services/` (http-клиент, telegram SDK,
  storage, tts — всё, чем пользуются 2+ фичи), `store/ui.js`, `data/dictionary.js` (используется
  глобальным `WordLookupPopover.vue`), `router/`, `locales/`, `utils/`.
- Если что-то нужно ДВУМ и более фичам — это не бизнес-логика одной фичи, кладётся в общую папку,
  а не дублируется/не импортируется из чужого `features/другая-фича/`.

---

## 2026-07-04 — Рефакторинг на feature-based структуру

Полный перенос существующего кода (backend + frontend) в `features/` по просьбе пользователя,
который заметил отсутствие такой папки при просмотре проекта. Сделано сразу, не только как
правило на будущее.

**Backend** (`backend/features/{auth,news,trainer,telegrambot,admin}/`):
- `auth` — TelegramLogin/GuestLogin/Logout/GetMe/UpdateProfile + `VerifyTelegramInitData` (был
  отдельный пакет `services/`, слит внутрь auth, т.к. используется только там).
- `news` — GetNews (раньше был свален в один файл с trainer-хендлерами, `NewsPost` переименован
  в `Post` внутри пакета news, чтобы не было тавтологии `news.NewsPost`).
- `trainer` — GetLevels/GetLevelCards/CompleteLevel/GetUserStats.
- `telegrambot` — вебхук, регистрация вебхука, /telegram/bot-info, привязка гостя к Telegram.
- `admin` — статистика/список/удаление пользователей + `Auth()` (был `middleware.AdminAuth`,
  переехал внутрь фичи и в admin package, т.к. специфичен только ей).
- Старые `backend/handlers/`, `backend/models/`, `backend/services/` удалены (пустые после переноса).
- `middleware/` и `database/` остались как общая инфраструктура (JWT-парсинг и CORS используются
  всеми фичами; миграции/сид-данные — общая схема для всех таблиц сразу, не разделены по фичам).
- Проверено: `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...` — всё чисто.

**Frontend** (`frontend/src/features/{auth,news,trainer,alphabet,profile,settings,admin}/`):
- `views/` полностью упразднена — все страницы переехали в соответствующие фичи.
- `services/api.js` (один файл на все домены) разделён: общий http-клиент → `services/httpClient.js`,
  а `authApi`/`newsApi`/`trainerApi`/`telegramApi` — по `api.js` внутри своих фич.
- Из `HomeView.vue` (фича news) вынесена карточка логина в `features/auth/components/LoginCard.vue`
  (бизнес-логика входа гостя/telegram не должна жить в news-странице).
- `WordBuildGame.vue` переехал не в фичу, а в `components/games/` — он общий для trainer
  (уровень "Приветствие") и alphabet ("Собери слово").
- `alphabetProgress.js`, `greetingWords.js`, `phonicsWords.js` переехали внутрь `features/alphabet/`
  (использовались только там).
- Проверено: `npm run build` (чисто, все чанки собрались), плюс ручной прогон в браузере по всем
  маршрутам (home/trainer/alphabet × все 10 уровней/profile/settings/admin) — без ошибок в консоли.

**Что не переносил:** `database/migrations.go` остаётся одним файлом на все таблицы (это схема, не
бизнес-логика конкретной фичи) — не разбивал по фичам, чтобы не усложнять миграции.

Изменения на ветке `refactor/feature-based-architecture`, ещё не в PR/main на момент записи.

---

## 2026-07-01 — Текущее состояние (снято на старте сессии)

Репозиторий имеет один коммит `chore: initial commit`, поверх которого лежит большой объём незакоммиченных изменений — по сути, это первая реализация фичи "Spoony" поверх голого TMA-боилерплейта.

**Новые файлы (не в git):**
- `backend/handlers/content.go` (166 строк) — API для контента: новости, уровни/города, вокабуляр, прогресс пользователя.
- `backend/models/content.go` (57 строк) — модели `NewsPost`, `Level`, `VocabItem`, `UserProgress`.
- `frontend/src/views/GameView.vue` (246 строк) — экран игры "найди пару".
- `frontend/src/views/TrainerView.vue` (161 строка) — карта городов/уровней.
- `frontend/src/components/SideNav.vue` (54 строки) — десктопный сайдбар.
- `.github/workflows/ci.yml` — новый CI-воркфлоу.

**Изменённые файлы (существенные):**
- `backend/database/migrations.go` (+166) — схема и seed-данные для нового контента.
- `backend/handlers/auth.go` (+98/-…) — доработка гостевого/demo входа.
- `frontend/src/views/HomeView.vue` (+158) — лента новостей.
- `frontend/src/views/ProfileView.vue` (+67) — статистика пользователя.
- `frontend/src/store/user.js`, `frontend/src/services/api.js`, `frontend/src/router/index.js` — интеграция нового API/страниц.
- Tailwind мигрирован на v4: удалены `postcss.config.js` и `tailwind.config.js`, стили теперь в `frontend/src/styles/tailwind.css`.
- Локали `en.js` / `ru.js` расширены под новые страницы.
- `.github/workflows/production.yml`, `staging.yml` — правки деплой-пайплайнов.

**Открытые вопросы / что дальше:**
- Ничего из этого ещё не закоммичено — стоит проверить и закоммитить логическими порциями.
- Нет данных о том, пройдено ли ручное тестирование новых страниц (Game/Trainer/Home/Profile) в браузере.
- Стоит свериться, что миграции применяются чисто на пустой БД (новая схема добавлена поверх существующей).
