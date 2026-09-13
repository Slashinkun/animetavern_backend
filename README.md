# AnimeTavern

Application web de suivi d’animés développée avec React, Go et PostgreSQL.

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

- Créer la base de données de l’application : `CREATE DATABASE nom_de_la_db;`

- Verifier qu’elle a bien été crée : `\l`

- Créer l’utilisateur : `CREATE USER nom_utilisateur WITH PASSWORD 'mdpchoisi'`

- Donner à myuser les permissions sur la base de données : `GRANT ALL PRIVILEGES ON DATABASE nom_de_la_db TO nom_utilisateur;`;

- Quitter PostgreSQL : `\q`

- Se reconnecter à la base de données avec le nom d'utilisateur que vous avez choisi :

Linux :
`psql -U nom_utilisateur -d nom_de_la_db`

Windows :
`psql -U nom_utilisateur -h localhost -p 5432 -d nom_de_la_db`

- Créer les tables de l’application contenu dans le fichier 'tables.sql'

## Serveur 

Installer Go : https://go.dev/doc/install

Cloner le repo : https://github.com/Slashinkun/animetavern_backend

A l’aide d’un terminal, se mettre dans le répertoire du serveur : `cd /server`

Installer les dépendances : `go mod tidy`

Créer le fichier .env dans le répertoire du serveur avec les identifiants que vous avez choisi lors de la création de la base de données :
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=nom_utilisateur
DB_PASSWORD=mdpchoisi
DB_NAME=nom_de_la_db
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
