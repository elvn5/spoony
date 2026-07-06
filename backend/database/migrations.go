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

		// Grammar cards for "theory" levels — kept short so theory stays a
		// small share of the course next to the practice games.
		`CREATE TABLE IF NOT EXISTS theory_slides (
			id SERIAL PRIMARY KEY,
			level_id INT NOT NULL REFERENCES levels(id) ON DELETE CASCADE,
			order_index INT NOT NULL DEFAULT 0,
			title_ru VARCHAR(255) NOT NULL,
			body_ru TEXT NOT NULL,
			example_en VARCHAR(255) DEFAULT '',
			example_ru VARCHAR(255) DEFAULT ''
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

		// Per-user "First Steps" (alphabet) completion. Unlike the trainer's
		// levels, the alphabet's 10 levels (4 base + 6 letter-combo groups)
		// are defined in frontend data, not a DB table, so level_id here is
		// just the 1-10 index the frontend already uses — no FK.
		`CREATE TABLE IF NOT EXISTS alphabet_progress (
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			level_id INT NOT NULL,
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
	seedBookUnits()
	repositionRoute()
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
// vocabulary (each boss's words are a subset of the two levels right before
// it, except "Hello!", which opens the whole route and so has nothing to
// recap — same as London opening the first "Find the pair" group). It then
// inserts any boss levels that don't exist yet. The order_index/position
// updates are safe to repeat on every startup (e.g. to fix up databases
// seeded before this grouping existed).
func seedBossLevels() {
	type routePos struct{ order, x, y int }
	route := map[string]routePos{
		"Hello!":     {0, 50, 99},
		"London":     {1, 50, 92},
		"Oxford":     {2, 30, 78},
		"Cambridge":  {4, 68, 66},
		"Bristol":    {5, 28, 54},
		"Stratford":  {7, 64, 44},
		"Manchester": {8, 38, 32},
		"Liverpool":  {10, 22, 22},
		"York":       {11, 58, 10},
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
		{"Windsor Castle", "Виндзорский замок", "Повторение: фрукты и животные", 3, 46, 72, []vocab{
			{"apple", "яблоко", "🍎"}, {"banana", "банан", "🍌"}, {"orange", "апельсин", "🍊"},
			{"cat", "кошка", "🐱"}, {"dog", "собака", "🐶"}, {"fox", "лиса", "🦊"},
		}},
		{"Big Ben", "Биг-Бен", "Повторение: еда и море", 6, 46, 49, []vocab{
			{"bread", "хлеб", "🍞"}, {"cheese", "сыр", "🧀"}, {"egg", "яйцо", "🥚"},
			{"fish", "рыба", "🐟"}, {"whale", "кит", "🐳"}, {"crab", "краб", "🦀"},
		}},
		{"Tower Bridge", "Тауэрский мост", "Повторение: предметы и спорт", 9, 30, 27, []vocab{
			{"book", "книга", "📖"}, {"pencil", "карандаш", "✏️"}, {"key", "ключ", "🔑"},
			{"ball", "мяч", "⚽"}, {"bicycle", "велосипед", "🚲"}, {"trophy", "кубок", "🏆"},
		}},
		{"Stonehenge", "Стоунхендж", "Повторение: музыка и погода", 12, 74, 2, []vocab{
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

// seedBookUnits seeds seven units based on "New English File Beginner"
// (Oxford). Each unit is a compact theory level (grammar cards, ~10% of the
// content), a couple of "Find the pair" vocabulary levels, and a
// sentence-build boss that recaps the unit: the user assembles English
// sentences out of the words they just learned. Inserts are idempotent
// (skipped when the level's city name already exists).
func seedBookUnits() {
	type slide struct{ title, body, exEn, exRu string }
	type vocab struct{ en, ru, emoji string }
	type lvl struct {
		city, titleRu, desc, emoji, gameType string
		items                                []vocab
		slides                               []slide
	}

	units := []lvl{
		// ---- Unit 1: Hello! (verb be, numbers, countries) ----
		{"Verb BE", "Теория: глагол be", "Юнит 1", "📖", "theory", nil, []slide{
			{"Глагол be — я и ты", "В английском «есть» не пропускается: I am (я есть), you are (ты есть). Кратко: I'm, you're.", "I'm Molly. You're late.", "Я Молли. Ты опоздал."},
			{"Он, она, оно", "he is (he's) — он, she is (she's) — она, it is (it's) — оно. Мы и они: we are (we're), they are (they're).", "She's from Italy. We're American.", "Она из Италии. Мы американцы."},
			{"Вопрос и отрицание", "Вопрос — be вперёд: Are you...? Отрицание — not после be: I'm not.", "Are you Tom? — No, I'm not.", "Ты Том? — Нет."},
		}},
		{"Numbers", "Числа 0–10", "Юнит 1", "🔢", "match", []vocab{
			{"one", "один", "1️⃣"}, {"two", "два", "2️⃣"}, {"three", "три", "3️⃣"},
			{"five", "пять", "5️⃣"}, {"seven", "семь", "7️⃣"}, {"ten", "десять", "🔟"},
		}, nil},
		{"Countries", "Страны", "Юнит 1", "🌍", "match", []vocab{
			{"Russia", "Россия", "🇷🇺"}, {"England", "Англия", "🇬🇧"}, {"the USA", "США", "🇺🇸"},
			{"Italy", "Италия", "🇮🇹"}, {"Spain", "Испания", "🇪🇸"}, {"Japan", "Япония", "🇯🇵"},
		}, nil},
		{"Nice to meet you", "Собери предложение: знакомство", "Босс юнита 1", "👑", "sentence_build", []vocab{
			{"I am from Russia", "Я из России", "🇷🇺"},
			{"You are late", "Ты опоздал", "⏰"},
			{"Nice to meet you", "Приятно познакомиться", "🤝"},
			{"Where are you from?", "Откуда ты?", "🌍"},
			{"We are American", "Мы американцы", "🇺🇸"},
		}, nil},

		// ---- Unit 2: Things & family (plurals, possessives, colours) ----
		{"A or An?", "Теория: артикли и множественное число", "Юнит 2", "📖", "theory", nil, []slide{
			{"a или an?", "a — перед согласным звуком (a pen, a key), an — перед гласным (an umbrella, an ID card).", "It's a key. It's an umbrella.", "Это ключ. Это зонт."},
			{"Множественное число", "Обычно + s: pen → pens. Запомни: man → men, woman → women, child → children, person → people.", "one child, two children", "один ребёнок, двое детей"},
			{"Чей это?", "Принадлежность — 's: Anna's car (машина Анны). Мой — my, твой — your, его — his, её — her.", "It's my sister's bag.", "Это сумка моей сестры."},
		}},
		{"Small Things", "Мелкие вещи", "Юнит 2", "🎒", "match", []vocab{
			{"pen", "ручка", "🖊️"}, {"bag", "сумка", "👜"}, {"umbrella", "зонт", "☂️"},
			{"watch", "часы", "⌚"}, {"glasses", "очки", "👓"}, {"wallet", "кошелёк", "👛"},
		}, nil},
		{"Family", "Семья", "Юнит 2", "👨‍👩‍👧", "match", []vocab{
			{"mother", "мама", "👩"}, {"father", "папа", "👨"}, {"sister", "сестра", "👧"},
			{"brother", "брат", "👦"}, {"grandmother", "бабушка", "👵"}, {"grandfather", "дедушка", "👴"},
		}, nil},
		{"Colours", "Цвета", "Юнит 2", "🎨", "match", []vocab{
			{"red", "красный", "🔴"}, {"blue", "синий", "🔵"}, {"green", "зелёный", "🟢"},
			{"yellow", "жёлтый", "🟡"}, {"black", "чёрный", "⚫"}, {"white", "белый", "⚪"},
		}, nil},
		{"Family Album", "Собери предложение: семья и вещи", "Босс юнита 2", "👑", "sentence_build", []vocab{
			{"This is my family", "Это моя семья", "👨‍👩‍👧‍👦"},
			{"My car is red", "Моя машина красная", "🚗"},
			{"It is an umbrella", "Это зонт", "☂️"},
			{"My brother is tall", "Мой брат высокий", "🧍"},
			{"The keys are old", "Ключи старые", "🔑"},
		}, nil},

		// ---- Unit 3: My day (present simple, food, jobs) ----
		{"Present Simple", "Теория: настоящее время", "Юнит 3", "📖", "theory", nil, []slide{
			{"Я делаю каждый день", "I / you / we / they + глагол без изменений: I drink coffee. We watch TV.", "I watch TV in the evening.", "Я смотрю телевизор вечером."},
			{"Он и она + s", "После he / she / it к глаголу добавляется -s: he works, she speaks.", "He speaks English at work.", "Он говорит по-английски на работе."},
			{"Вопросы и отрицания", "Вопрос: Do you...? / Does she...? Отрицание: don't / doesn't + глагол.", "Do you like tea? — No, I don't.", "Ты любишь чай? — Нет."},
		}},
		{"Food & Drink", "Еда и напитки", "Юнит 3", "☕", "match", []vocab{
			{"coffee", "кофе", "☕"}, {"tea", "чай", "🍵"}, {"water", "вода", "💧"},
			{"juice", "сок", "🧃"}, {"sandwich", "сэндвич", "🥪"}, {"salad", "салат", "🥗"},
		}, nil},
		{"Jobs", "Профессии", "Юнит 3", "💼", "match", []vocab{
			{"doctor", "врач", "🧑‍⚕️"}, {"teacher", "учитель", "🧑‍🏫"}, {"policeman", "полицейский", "👮"},
			{"student", "студент", "🧑‍🎓"}, {"factory worker", "рабочий", "👷"}, {"waiter", "официант", "🤵"},
		}, nil},
		{"My Day", "Собери предложение: мой день", "Босс юнита 3", "👑", "sentence_build", []vocab{
			{"I drink coffee every morning", "Я пью кофе каждое утро", "☕"},
			{"She works in a hospital", "Она работает в больнице", "🏥"},
			{"Do you like tea?", "Ты любишь чай?", "🍵"},
			{"He speaks English at work", "Он говорит по-английски на работе", "💬"},
			{"We eat fast food", "Мы едим фастфуд", "🍔"},
		}, nil},

		// ---- Unit 4: Habits (adverbs of frequency, can/can't) ----
		{"Always or Never?", "Теория: наречия частоты и can", "Юнит 4", "📖", "theory", nil, []slide{
			{"Как часто?", "always — всегда, usually — обычно, sometimes — иногда, never — никогда. Ставятся перед глаголом.", "I always get up at seven.", "Я всегда встаю в семь."},
			{"can — могу", "can + глагол — умение или разрешение. Отрицание: can't.", "You can't park here.", "Здесь нельзя парковаться."},
			{"Вопрос с can", "Can вперёд: Can you swim? Короткий ответ: Yes, I can / No, I can't.", "Can you swim? — Yes, I can.", "Ты умеешь плавать? — Да."},
		}},
		{"Daily Routine", "Распорядок дня", "Юнит 4", "⏰", "match", []vocab{
			{"get up", "вставать", "⏰"}, {"have breakfast", "завтракать", "🥐"}, {"go to work", "ехать на работу", "🚌"},
			{"have lunch", "обедать", "🍽️"}, {"watch TV", "смотреть телевизор", "📺"}, {"go to bed", "ложиться спать", "🛏️"},
		}, nil},
		{"Action Verbs", "Глаголы действия", "Юнит 4", "🏃", "match", []vocab{
			{"swim", "плавать", "🏊"}, {"dance", "танцевать", "💃"}, {"sing", "петь", "🎤"},
			{"run", "бегать", "🏃"}, {"drive", "водить машину", "🚗"}, {"play tennis", "играть в теннис", "🎾"},
		}, nil},
		{"Habits", "Собери предложение: привычки", "Босс юнита 4", "👑", "sentence_build", []vocab{
			{"I always get up at seven", "Я всегда встаю в семь", "⏰"},
			{"You can't park here", "Здесь нельзя парковаться", "🚫"},
			{"Can you swim?", "Ты умеешь плавать?", "🏊"},
			{"She never watches TV", "Она никогда не смотрит телевизор", "📺"},
			{"We sometimes play tennis", "Мы иногда играем в теннис", "🎾"},
		}, nil},

		// ---- Unit 5: The past (past simple) ----
		{"Past Simple", "Теория: прошедшее время", "Юнит 5", "📖", "theory", nil, []slide{
			{"was и were", "Прошлое глагола be: I / he / she / it was, you / we / they were.", "They were famous.", "Они были знамениты."},
			{"Правильные глаголы + ed", "Прошедшее время — добавь -ed: work → worked, change → changed.", "It changed my life.", "Это изменило мою жизнь."},
			{"Неправильные глаголы", "Учим наизусть: go → went, have → had, get → got, buy → bought.", "We went to Rome.", "Мы поехали в Рим."},
		}},
		{"Travel Verbs", "Глаголы путешествий", "Юнит 5", "🧳", "match", []vocab{
			{"buy a ticket", "купить билет", "🎫"}, {"stay at a hotel", "жить в отеле", "🏨"}, {"meet a friend", "встретить друга", "🤝"},
			{"take photos", "фотографировать", "📸"}, {"lose your keys", "потерять ключи", "🗝️"}, {"rent a car", "арендовать машину", "🚙"},
		}, nil},
		{"Yesterday", "Собери предложение: прошлое", "Босс юнита 5", "👑", "sentence_build", []vocab{
			{"I was at home yesterday", "Вчера я был дома", "🏠"},
			{"They were in Rome last week", "Они были в Риме на прошлой неделе", "🇮🇹"},
			{"We bought some food", "Мы купили еды", "🛒"},
			{"She got up early", "Она встала рано", "🌅"},
			{"He went on holiday in July", "Он уехал в отпуск в июле", "🏖️"},
		}, nil},

		// ---- Unit 6: Places (there is / there are) ----
		{"There is / are", "Теория: оборот there is/are", "Юнит 6", "📖", "theory", nil, []slide{
			{"There is — есть (один)", "There is (there's) + один предмет: There's a hotel near here — рядом есть отель.", "There is a shop in the village.", "В деревне есть магазин."},
			{"There are — есть (много)", "There are + несколько: There are two banks. Вопрос: Is there...? / Are there...?", "There are a lot of cafés.", "Там много кафе."},
			{"Было: there was / were", "Прошлое: there was (один), there were (много).", "There was a castle here.", "Здесь был замок."},
		}},
		{"In a Hotel", "В отеле", "Юнит 6", "🏨", "match", []vocab{
			{"bed", "кровать", "🛏️"}, {"shower", "душ", "🚿"}, {"bathroom", "ванная", "🛁"},
			{"lamp", "лампа", "💡"}, {"lift", "лифт", "🛗"}, {"reception", "ресепшен", "🛎️"},
		}, nil},
		{"Places", "Места в городе", "Юнит 6", "🏙️", "match", []vocab{
			{"bank", "банк", "🏦"}, {"hospital", "больница", "🏥"}, {"museum", "музей", "🖼️"},
			{"supermarket", "супермаркет", "🛒"}, {"church", "церковь", "⛪"}, {"beach", "пляж", "🏖️"},
		}, nil},
		{"My Town", "Собери предложение: город", "Босс юнита 6", "👑", "sentence_build", []vocab{
			{"There is a hotel near here", "Рядом есть отель", "🏨"},
			{"There are two banks", "Там два банка", "🏦"},
			{"There was a castle here", "Здесь был замок", "🏰"},
			{"Is there a park near here?", "Рядом есть парк?", "🌳"},
			{"The museum is old", "Музей старый", "🖼️"},
		}, nil},

		// ---- Unit 7: Plans (like + -ing, going to, weather) ----
		{"Going to", "Теория: планы на будущее", "Юнит 7", "📖", "theory", nil, []slide{
			{"like + -ing", "Люблю что-то делать: like + глагол с -ing: I like reading. She likes dancing.", "She likes dancing.", "Она любит танцевать."},
			{"be going to — собираюсь", "План на будущее: be going to + глагол: I'm going to travel.", "We're going to visit Rome.", "Мы собираемся посетить Рим."},
			{"Предсказания", "be going to — и для предсказаний по очевидным признакам.", "Look at the clouds! It's going to rain.", "Посмотри на тучи! Сейчас пойдёт дождь."},
		}},
		{"Weather", "Погода", "Юнит 7", "🌦️", "match", []vocab{
			{"sunny", "солнечно", "🌞"}, {"windy", "ветрено", "💨"}, {"hot", "жарко", "🥵"},
			{"cold", "холодно", "🥶"}, {"snowy", "снежно", "🌨️"}, {"foggy", "туманно", "🌫️"},
		}, nil},
		{"Activities", "Активности", "Юнит 7", "⛺", "match", []vocab{
			{"shopping", "шопинг", "🛍️"}, {"camping", "кемпинг", "⛺"}, {"fishing", "рыбалка", "🎣"},
			{"cycling", "велоспорт", "🚴"}, {"painting", "рисование", "🎨"}, {"skiing", "лыжи", "⛷️"},
		}, nil},
		{"Future Plans", "Собери предложение: планы", "Босс юнита 7", "👑", "sentence_build", []vocab{
			{"I like reading in bed", "Я люблю читать в кровати", "📖"},
			{"We are going to travel to Italy", "Мы собираемся поехать в Италию", "✈️"},
			{"It is going to rain", "Скоро пойдёт дождь", "🌧️"},
			{"She likes dancing", "Она любит танцевать", "💃"},
			{"What are you going to do?", "Что ты собираешься делать?", "❓"},
		}, nil},
	}

	// Book units continue the route after the original 13 levels (order 0-12).
	const startOrder = 13

	for i, u := range units {
		var count int
		if err := DB.QueryRow(`SELECT COUNT(*) FROM levels WHERE city = $1`, u.city).Scan(&count); err != nil {
			log.Printf("seedBookUnits: count %s failed: %v", u.city, err)
			continue
		}
		if count > 0 {
			continue
		}

		var levelID int
		err := DB.QueryRow(
			`INSERT INTO levels (city, title_ru, description, emoji, order_index, pos_x, pos_y, game_type)
			 VALUES ($1,$2,$3,$4,$5,50,50,$6) RETURNING id`,
			u.city, u.titleRu, u.desc, u.emoji, startOrder+i, u.gameType,
		).Scan(&levelID)
		if err != nil {
			log.Printf("seedBookUnits: insert %s failed: %v", u.city, err)
			continue
		}

		for _, v := range u.items {
			if _, err := DB.Exec(
				`INSERT INTO vocab_items (level_id, word_en, word_ru, emoji) VALUES ($1,$2,$3,$4)`,
				levelID, v.en, v.ru, v.emoji,
			); err != nil {
				log.Printf("seedBookUnits: insert vocab for %s failed: %v", u.city, err)
			}
		}
		for j, s := range u.slides {
			if _, err := DB.Exec(
				`INSERT INTO theory_slides (level_id, order_index, title_ru, body_ru, example_en, example_ru)
				 VALUES ($1,$2,$3,$4,$5,$6)`,
				levelID, j, s.title, s.body, s.exEn, s.exRu,
			); err != nil {
				log.Printf("seedBookUnits: insert slide for %s failed: %v", u.city, err)
			}
		}
	}
	log.Printf("seedBookUnits: %d book unit levels ensured", len(units))
}

// repositionRoute lays every level out along the winding map path: y runs
// from the bottom (97%) to the top (2%) in route order and x zigzags. Runs
// on every startup so the map stays evenly spaced as levels are added.
func repositionRoute() {
	rows, err := DB.Query(`SELECT id FROM levels ORDER BY order_index ASC, id ASC`)
	if err != nil {
		log.Printf("repositionRoute: query failed: %v", err)
		return
	}
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	rows.Close()

	if len(ids) < 2 {
		return
	}

	xs := []int{50, 26, 64, 32, 72, 38, 20, 58}
	n := len(ids)
	for i, id := range ids {
		y := 97 - (95*i)/(n-1)
		x := xs[i%len(xs)]
		if _, err := DB.Exec(`UPDATE levels SET pos_x = $1, pos_y = $2 WHERE id = $3`, x, y, id); err != nil {
			log.Printf("repositionRoute: update %d failed: %v", id, err)
		}
	}
}
