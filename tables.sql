CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE anime (
    id INTEGER PRIMARY KEY,
    title VARCHAR(255),
    image TEXT,
    episodes INTEGER DEFAULT 0
);    


CREATE TABLE user_anime (
    id SERIAL PRIMARY KEY,

    user_id INTEGER NOT NULL,
    anime_id INTEGER NOT NULL,

    status VARCHAR(50) DEFAULT 'PLANNING',
    note INTEGER,
    favorite BOOLEAN DEFAULT FALSE,
    viewed_episodes INTEGER DEFAULT 0,

    CONSTRAINT unique_user_anime UNIQUE (user_id, anime_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (anime_id) REFERENCES anime(id) ON DELETE CASCADE
);


CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    anime_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    rating INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (anime_id) REFERENCES anime(id) ON DELETE CASCADE,
    UNIQUE (user_id, anime_id)
);
