# AnimeTavern

Application web de suivi d’animés développée avec React, Go et PostgreSQL.

Lien vers le site : [Site](https://animetavern-frontend.onrender.com/) (très instable)

## Backend

https://github.com/Slashinkun/animetavern_backend VOUS ÊTES ICI

## Frontend

https://github.com/Slashinkun/animetavern_frontend

# Installation

Pour déployer l'application localement, veuillez suivre les instructions suivantes :

## Base de données 

- Installer PostgreSQL sur votre machine

  Linux : suivre les instructions d’installation selon la distribution en tant que `sudo` https://www.postgresql.org/download/linux/

  Windows : https://www.postgresql.org/download/windows/

- Demarrer postgreSQL en tant que postgres : `sudo -u postgres psql` (Linux) ou `psql -U postgres -h localhost -p 5432` (Windows)

- Créer la base de données de l’application : `CREATE DATABASE animetavern_db;`

- Verifier qu’elle a bien été crée : `\l`

- Se connecter à animetavern_db : `\c animetavern_db`
- Créer l’utilisateur myuser : `CREATE USER myuser WITH PASSWORD 'mypassword'`;
- Donner à myuser les permissions sur la base de données : `ALTER DATABASE animetavern_db OWNER TO myuser`;
- Créer les tables de l’application :

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

```sql
CREATE TABLE anime (
    id INTEGER PRIMARY KEY,
    title VARCHAR(50),
    image VARCHAR(100),
    episodes INTEGER DEFAULT 0
);
```

```sql
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
```

```sql
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
```

## Serveur 

Installer Go : https://go.dev/doc/install

Cloner le repo : https://github.com/Slashinkun/animetavern_backend

A l’aide d’un terminal, se mettre dans le répertoire du serveur : `cd /server`

Installer les dépendances : `go mod tidy`

Créer le fichier .env dans le répertoire du serveur avec ces paramètres :

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=animetavern_db
DB_SSLMODE=disable
```

Vérifier que le serveur démarre : `go run main.go`

## Client

Installer NodeJS : https://nodejs.org/fr

Cloner le repo : https://github.com/Slashinkun/animetavern_frontend

A l’aide d’un terminal, se mettre dans le répertoire du client

Installer les dépendences : `npm install`

Vérifier que le client marche : `npm run dev`

## Démarrer l’application :

Démarrer le serveur et le client dans 2 terminaux séparés
