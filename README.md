# AnimeTavern - Backend

Backend de l'application web **AnimeTavern**, une application de suivi d'animés développée avec **Go** et **PostgreSQL**.

Ce projet a été réalisé dans le cadre de l'UE **PC3R** à **Sorbonne Université**.

## Fonctionnalités

- Gestion des animés
- Suivi de progression
- Gestion des utilisateurs
- Authentification

## Pré-requis

- [PostgreSQL](https://www.postgresql.org/download/)
- [Go](https://go.dev/doc/install)

## Installation

Pour déployer le backend localement, veuillez suivre les instructions suivantes :

### Base de données 


- Se connecter à PostgreSQL :
  - Linux : `sudo -u postgres psql`
  - Windows :  `psql -U postgres -h localhost -p 5432` 

- Créer la base de données de l’application : `CREATE DATABASE nom_de_la_db;`

- Vérifier qu’elle a bien été créée : `\l`

- Créer l’utilisateur dédié à l'application : `CREATE USER nom_utilisateur WITH PASSWORD 'mdpchoisi';`

- Donnez-lui les permissions sur la base de données : `GRANT ALL PRIVILEGES ON DATABASE nom_de_la_db TO nom_utilisateur;`

- Quitter PostgreSQL : `\q`

- Se connecter à la base de données avec l'utilisateur créé :

  - Linux : `psql -U nom_utilisateur -d nom_de_la_db`

  - Windows : `psql -U nom_utilisateur -h localhost -p 5432 -d nom_de_la_db`

- Créer les tables à partir du fichier `tables.sql`

   `psql -U nom_utilisateur -d nom_de_la_db -f tables.sql`

## Serveur 

- Cloner le dépôt : https://github.com/Slashinkun/animetavern_backend

  ```bash
  git clone https://github.com/Slashinkun/animetavern_backend
  ```

- Se placer dans le répertoire du projet : `cd animetavern_backend`

- Installer les dépendances : `go mod tidy`

- Créer un fichier .env à la racine du projet avec les identifiants de la base de données :
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=nom_utilisateur
DB_PASSWORD=mdpchoisi
DB_NAME=nom_de_la_db
DB_SSLMODE=disable
```
- Démarrer le serveur : `go run .`

## Frontend

Pour interagir avec le backend, il est recommandé d'installer également le frontend de l'application.

Le dépôt du frontend est disponible ici : https://github.com/Slashinkun/animetavern_frontend

