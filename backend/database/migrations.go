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
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Migration failed: %v\nQuery: %s", err, q)
		}
	}

	log.Println("Database migrations completed successfully")
}
