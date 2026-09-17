# AnimeTavern - Backend

Backend de l'application web **AnimeTavern**, une application de suivi d'animés développée avec **Go** et **PostgreSQL**.

Ce projet a été réalisé dans le cadre de l'UE **PC3R** à **Sorbonne Université**.

## Projet

- **Backend** : https://github.com/Slashinkun/animetavern_backend
- **Frontend** : https://github.com/Slashinkun/animetavern_frontend

# Installation

Pour déployer le backend localement, veuillez suivre les instructions suivantes :

## Base de données 

- Installer [PostgreSQL](https://www.postgresql.org/download/) sur votre machine

- Démarrer postgreSQL :
  - Linux : `sudo -u postgres psql`
  - Windows :  `psql -U postgres -h localhost -p 5432` 

- Créer la base de données de l’application : `CREATE DATABASE nom_de_la_db;`

- Verifier qu’elle a bien été crée : `\l`

- Créer l’utilisateur : `CREATE USER nom_utilisateur WITH PASSWORD 'mdpchoisi'`

- Donner les permissions à l'utilisateur : `GRANT ALL PRIVILEGES ON DATABASE nom_de_la_db TO nom_utilisateur;`

- Quitter PostgreSQL : `\q`

- Se connecter à la base de données avec l'utilisateur créé :

  - Linux : `psql -U nom_utilisateur -d nom_de_la_db`

  - Windows : `psql -U nom_utilisateur -h localhost -p 5432 -d nom_de_la_db`

- Créer les tables à partir du fichier 'tables.sql'

## Serveur 

- Installer [Go] (https://go.dev/doc/install)

- Cloner le repositoire : https://github.com/Slashinkun/animetavern_backend

  ```git clone https://github.com/Slashinkun/animetavern_backend```

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


- Démarrer le serveur : `go run main.go`


