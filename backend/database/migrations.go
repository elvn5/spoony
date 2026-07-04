package database

import "log"

func RunMigrations() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			telegram_id BIGINT UNIQUE,
			username VARCHAR(255),
			email VARCHAR(255),
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			avatar_url VARCHAR(500),
			language VARCHAR(10) DEFAULT 'en',
			timezone VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Web guests (users who open Spoony as a normal website, without Telegram).
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS guest_id VARCHAR(64) UNIQUE`,

		// Home feed posts (Facebook-style timeline).
		`CREATE TABLE IF NOT EXISTS news_posts (
			id SERIAL PRIMARY KEY,
			author VARCHAR(255) NOT NULL,
			avatar VARCHAR(16) DEFAULT '🥄',
			title VARCHAR(255) NOT NULL,
			body TEXT NOT NULL,
			image VARCHAR(16) DEFAULT '',
			category VARCHAR(64) DEFAULT 'news',
			likes INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Levels = cities on the England route.
		`CREATE TABLE IF NOT EXISTS levels (
			id SERIAL PRIMARY KEY,
			city VARCHAR(255) NOT NULL,
			title_ru VARCHAR(255) NOT NULL,
			description VARCHAR(500) DEFAULT '',
			emoji VARCHAR(16) DEFAULT '📍',
			order_index INT NOT NULL DEFAULT 0,
			pos_x INT NOT NULL DEFAULT 50,
			pos_y INT NOT NULL DEFAULT 50
		)`,

		// game_type distinguishes the "Find the pair" memory game from other
		// exercise types (e.g. "word_build" — assemble a word from letter blocks).
		`ALTER TABLE levels ADD COLUMN IF NOT EXISTS game_type VARCHAR(32) NOT NULL DEFAULT 'match'`,

		// Vocabulary items — each produces a picture card + a word card.
		`CREATE TABLE IF NOT EXISTS vocab_items (
			id SERIAL PRIMARY KEY,
			level_id INT NOT NULL REFERENCES levels(id) ON DELETE CASCADE,
			word_en VARCHAR(128) NOT NULL,
			word_ru VARCHAR(128) NOT NULL,
			emoji VARCHAR(16) NOT NULL
		)`,

		// Per-user level completion.
		`CREATE TABLE IF NOT EXISTS user_progress (
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			level_id INT NOT NULL REFERENCES levels(id) ON DELETE CASCADE,
			stars INT DEFAULT 0,
			completed BOOLEAN DEFAULT TRUE,
			completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, level_id)
		)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Migration failed: %v\nQuery: %s", err, q)
		}
	}

	seedContent()

	log.Println("Database migrations completed successfully")
}

// seedContent populates the learning content once, on a fresh database.
func seedContent() {
	seedNews()
	seedLevels()
	seedGreetingLevel()
	seedBossLevels()
}

func seedNews() {
	var count int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM news_posts`).Scan(&count); err != nil {
		log.Printf("seedNews: count failed: %v", err)
		return
	}
	if count > 0 {
		return
	}

	type post struct {
		author, avatar, title, body, image, category string
		likes                                        int
	}
	posts := []post{
		{"Spoony", "🥄", "Добро пожаловать в Spoony!", "Привет, друг! Я ложечка Spoony, и я помогу тебе выучить английский язык. Открывай «Тренажёр» и отправляйся в путешествие по Англии! 🇬🇧", "👋", "Новости", 128},
		{"Spoony", "🥄", "Новый город на маршруте — Ливерпуль!", "Мы добавили музыкальный город Ливерпуль. Выучи английские слова про музыку: guitar, drum, piano и другие! 🎸", "🎸", "Обновление", 96},
		{"Учитель Сова", "🦉", "Слово дня: APPLE 🍎", "Apple — это «яблоко». Попробуй сказать вслух: «ЭП-пл». А какое твоё любимое яблоко — красное или зелёное?", "🍎", "Слово дня", 74},
		{"Spoony", "🥄", "Совет: учись каждый день по чуть-чуть", "Лучше заниматься 10 минут каждый день, чем час один раз в неделю. Проходи по одному городу в день — и ты быстро выучишь много слов! ⭐", "📅", "Совет", 61},
		{"Кот Том", "🐱", "Игра «Найди пару» стала ещё веселее", "Теперь, когда ты находишь правильную пару картинка-слово, карточки красиво загораются зелёным! Сможешь пройти город без ошибок? 💚", "🃏", "Игра", 88},
		{"Spoony", "🥄", "Сегодня изучаем животных 🐶", "В городе Оксфорд тебя ждут английские названия животных: cat, dog, rabbit, fox и другие. Удачи, исследователь! 🦊", "🐾", "Слово дня", 53},
	}

	for _, p := range posts {
		_, err := DB.Exec(
			`INSERT INTO news_posts (author, avatar, title, body, image, category, likes)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			p.author, p.avatar, p.title, p.body, p.image, p.category, p.likes,
		)
		if err != nil {
			log.Printf("seedNews: insert failed: %v", err)
		}
	}
	log.Printf("Seeded %d news posts", len(posts))
}

func seedLevels() {
	var count int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM levels`).Scan(&count); err != nil {
		log.Printf("seedLevels: count failed: %v", err)
		return
	}
	if count > 0 {
		return
	}

	type vocab struct{ en, ru, emoji string }
	type lvl struct {
		city, titleRu, desc, emoji string
		posX, posY                 int
		items                      []vocab
	}

	levels := []lvl{
		{"London", "Лондон", "Фрукты", "🏰", 50, 92, []vocab{
			{"apple", "яблоко", "🍎"}, {"banana", "банан", "🍌"}, {"orange", "апельсин", "🍊"},
			{"grape", "виноград", "🍇"}, {"strawberry", "клубника", "🍓"}, {"cherry", "вишня", "🍒"},
		}},
		{"Oxford", "Оксфорд", "Животные", "🎓", 30, 78, []vocab{
			{"cat", "кошка", "🐱"}, {"dog", "собака", "🐶"}, {"rabbit", "кролик", "🐰"},
			{"fox", "лиса", "🦊"}, {"bear", "медведь", "🐻"}, {"lion", "лев", "🦁"},
		}},
		{"Cambridge", "Кембридж", "Еда", "📚", 68, 66, []vocab{
			{"bread", "хлеб", "🍞"}, {"cheese", "сыр", "🧀"}, {"egg", "яйцо", "🥚"},
			{"milk", "молоко", "🥛"}, {"pizza", "пицца", "🍕"}, {"cake", "торт", "🍰"},
		}},
		{"Bristol", "Бристоль", "Море", "🌉", 28, 54, []vocab{
			{"fish", "рыба", "🐟"}, {"whale", "кит", "🐳"}, {"octopus", "осьминог", "🐙"},
			{"crab", "краб", "🦀"}, {"turtle", "черепаха", "🐢"}, {"dolphin", "дельфин", "🐬"},
		}},
		{"Stratford", "Стратфорд", "Предметы", "🎭", 64, 44, []vocab{
			{"book", "книга", "📖"}, {"pencil", "карандаш", "✏️"}, {"scissors", "ножницы", "✂️"},
			{"clock", "часы", "🕐"}, {"key", "ключ", "🔑"}, {"bell", "колокольчик", "🔔"},
		}},
		{"Manchester", "Манчестер", "Спорт", "⚽", 38, 32, []vocab{
			{"ball", "мяч", "⚽"}, {"basketball", "баскетбол", "🏀"}, {"tennis", "теннис", "🎾"},
			{"bicycle", "велосипед", "🚲"}, {"trophy", "кубок", "🏆"}, {"medal", "медаль", "🏅"},
		}},
		{"Liverpool", "Ливерпуль", "Музыка", "🎸", 22, 22, []vocab{
			{"guitar", "гитара", "🎸"}, {"drum", "барабан", "🥁"}, {"piano", "пианино", "🎹"},
			{"microphone", "микрофон", "🎤"}, {"trumpet", "труба", "🎺"}, {"violin", "скрипка", "🎻"},
		}},
		{"York", "Йорк", "Погода", "🏰", 58, 10, []vocab{
			{"sun", "солнце", "☀️"}, {"rain", "дождь", "🌧️"}, {"snow", "снег", "❄️"},
			{"star", "звезда", "⭐"}, {"rainbow", "радуга", "🌈"}, {"cloud", "облако", "☁️"},
		}},
	}

	for i, l := range levels {
		var levelID int
		err := DB.QueryRow(
			`INSERT INTO levels (city, title_ru, description, emoji, order_index, pos_x, pos_y)
			 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
			l.city, l.titleRu, l.desc, l.emoji, i, l.posX, l.posY,
		).Scan(&levelID)
		if err != nil {
			log.Printf("seedLevels: insert level failed: %v", err)
			continue
		}
		for _, v := range l.items {
			if _, err := DB.Exec(
				`INSERT INTO vocab_items (level_id, word_en, word_ru, emoji) VALUES ($1,$2,$3,$4)`,
				levelID, v.en, v.ru, v.emoji,
			); err != nil {
				log.Printf("seedLevels: insert vocab failed: %v", err)
			}
		}
	}
	log.Printf("Seeded %d levels", len(levels))
}

// seedGreetingLevel inserts the "Hello!" greeting level, the route's first
// mini-boss. It runs independently of seedLevels so it also backfills
// databases that were already seeded before this level type existed. Its
// final order_index/position in the route is assigned by seedBossLevels.
func seedGreetingLevel() {
	var count int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM levels WHERE city = 'Hello!'`).Scan(&count); err != nil {
		log.Printf("seedGreetingLevel: count failed: %v", err)
		return
	}
	if count > 0 {
		return
	}

	var levelID int
	err := DB.QueryRow(
		`INSERT INTO levels (city, title_ru, description, emoji, order_index, pos_x, pos_y, game_type)
		 VALUES ($1,$2,$3,$4,0,50,50,'word_build') RETURNING id`,
		"Hello!", "Приветствие и знакомство", "Знакомство", "👋",
	).Scan(&levelID)
	if err != nil {
		log.Printf("seedGreetingLevel: insert level failed: %v", err)
		return
	}

	type vocab struct{ en, ru, emoji string }
	words := []vocab{
		{"hello", "привет", "👋"},
		{"hi", "привет", "🙋"},
		{"bye", "пока", "🙋‍♂️"},
		{"yes", "да", "✅"},
		{"no", "нет", "❌"},
		{"please", "пожалуйста", "🙏"},
		{"thanks", "спасибо", "🙏"},
		{"friend", "друг", "🤝"},
	}
	for _, w := range words {
		if _, err := DB.Exec(
			`INSERT INTO vocab_items (level_id, word_en, word_ru, emoji) VALUES ($1,$2,$3,$4)`,
			levelID, w.en, w.ru, w.emoji,
		); err != nil {
			log.Printf("seedGreetingLevel: insert vocab failed: %v", err)
		}
	}
	log.Printf("Seeded greeting level with %d words", len(words))
}

// seedBossLevels arranges the England route into groups of two "Find the
// pair" cities followed by a "Собери слово" mini-boss that recaps their
// vocabulary, then inserts the three new boss levels the first time this
// runs. The order_index/position updates are safe to repeat on every startup
// (e.g. to fix up databases seeded before this grouping existed).
func seedBossLevels() {
	type routePos struct{ order, x, y int }
	route := map[string]routePos{
		"London":     {0, 50, 92},
		"Oxford":     {1, 30, 78},
		"Hello!":     {2, 46, 72},
		"Cambridge":  {3, 68, 66},
		"Bristol":    {4, 28, 54},
		"Stratford":  {6, 64, 44},
		"Manchester": {7, 38, 32},
		"Liverpool":  {9, 22, 22},
		"York":       {10, 58, 10},
	}
	for city, p := range route {
		if _, err := DB.Exec(
			`UPDATE levels SET order_index = $1, pos_x = $2, pos_y = $3 WHERE city = $4`,
			p.order, p.x, p.y, city,
		); err != nil {
			log.Printf("seedBossLevels: reposition %s failed: %v", city, err)
		}
	}

	type vocab struct{ en, ru, emoji string }
	type boss struct {
		city, titleRu, desc    string
		orderIndex, posX, posY int
		items                  []vocab
	}
	bosses := []boss{
		{"Big Ben", "Биг-Бен", "Повторение: еда и море", 5, 46, 49, []vocab{
			{"bread", "хлеб", "🍞"}, {"cheese", "сыр", "🧀"}, {"egg", "яйцо", "🥚"},
			{"fish", "рыба", "🐟"}, {"whale", "кит", "🐳"}, {"crab", "краб", "🦀"},
		}},
		{"Tower Bridge", "Тауэрский мост", "Повторение: предметы и спорт", 8, 30, 27, []vocab{
			{"book", "книга", "📖"}, {"pencil", "карандаш", "✏️"}, {"key", "ключ", "🔑"},
			{"ball", "мяч", "⚽"}, {"bicycle", "велосипед", "🚲"}, {"trophy", "кубок", "🏆"},
		}},
		{"Stonehenge", "Стоунхендж", "Повторение: музыка и погода", 11, 74, 2, []vocab{
			{"guitar", "гитара", "🎸"}, {"drum", "барабан", "🥁"}, {"piano", "пианино", "🎹"},
			{"sun", "солнце", "☀️"}, {"rain", "дождь", "🌧️"}, {"snow", "снег", "❄️"},
		}},
	}

	for _, b := range bosses {
		var count int
		if err := DB.QueryRow(`SELECT COUNT(*) FROM levels WHERE city = $1`, b.city).Scan(&count); err != nil {
			log.Printf("seedBossLevels: count %s failed: %v", b.city, err)
			continue
		}
		if count > 0 {
			continue
		}

		var levelID int
		err := DB.QueryRow(
			`INSERT INTO levels (city, title_ru, description, emoji, order_index, pos_x, pos_y, game_type)
			 VALUES ($1,$2,$3,'👑',$4,$5,$6,'word_build') RETURNING id`,
			b.city, b.titleRu, b.desc, b.orderIndex, b.posX, b.posY,
		).Scan(&levelID)
		if err != nil {
			log.Printf("seedBossLevels: insert %s failed: %v", b.city, err)
			continue
		}
		for _, v := range b.items {
			if _, err := DB.Exec(
				`INSERT INTO vocab_items (level_id, word_en, word_ru, emoji) VALUES ($1,$2,$3,$4)`,
				levelID, v.en, v.ru, v.emoji,
			); err != nil {
				log.Printf("seedBossLevels: insert vocab for %s failed: %v", b.city, err)
			}
		}
	}
	log.Println("seedBossLevels: route reordered and mini-bosses ensured")
}
