package database

import "log"

func InitTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT NOT NULL,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS anime (
			id SERIAL PRIMARY KEY,
			title VARCHAR(50),
			image VARCHAR(100),
			episodes INTEGER DEFAULT 0
		);`,

		`CREATE TABLE IF NOT EXISTS user_anime (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			anime_id INTEGER NOT NULL,
			status VARCHAR(50) DEFAULT 'PLANNING',
			note INTEGER,
			favorite BOOLEAN DEFAULT FALSE,
			viewed_episodes INTEGER DEFAULT 0,
			UNIQUE (user_id, anime_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (anime_id) REFERENCES anime(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			anime_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			rating INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (user_id, anime_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (anime_id) REFERENCES anime(id) ON DELETE CASCADE
		);`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatal("Erreur création table:", err)
		}
	}

	log.Println("Tables créées / vérifiées")
}
