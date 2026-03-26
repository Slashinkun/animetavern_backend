# Projet

Un site qui permet de traquer les animés qu'on a vu, qu'on veut voir.
Mais aussi, de noter les animés et ecrire des critiques qui seront vu par d'autre utilisateur. Egalement, l'utilisateur peut personnaliser sa page utilisateur en mettant une photo de profil et une bannière . De plus, il pourra mettre ses animes favoris.


## API
On utilise l'API Jikan : 
https://jikan.moe/

Elle permet de récuperer les données sur les animés (nom,image,genre,personnages,realisateur...)

On utilisera une API REST interne si besoin pour les actions de l'utilisateur.

## Fonctionnalités
- Creer un compte
- Personnaliser sa page (photo de profil et bannière) (si temps et ressources possible)
- Ajouter un anime à sa liste et changer son statut de visionnages, épisodes vus (en cours de visionnage ,prevoit de voir, terminé, arreté)
- Noter et ecrire une critique sur un animé
- Rechercher un animé depuis la barre de recherche
- Mettre en favoris un animé

## Cas d'utilisations
- Johnny a un compte sur le site, il se connecte et cherche un anime qu'il veut regarder. Il en trouve un,il lit les avis dessus et decide de le mettre dans sa liste a regarder.
- Michel a fini de regarder un anime et il decide d'ecrire une critique et mettre une note. 
- Diana veut personnaliser son compte et cherche une photo de profil et une  sur Internet
Elle trouve les images et decide de les rajouter sur son profil   (si temps possible)
- Frank a fini un animé et l'a bien aimé. Il décide le rajouter à ses favoris.

## BDD
4 tables :

Utilisateur : idUser int,pseudo string,email string,password  (string,photo_profil string ,banniere string)

Anime : idAnime int , title string , note int (note calculé depuis la BDD ? ou pris de l'API externe)

UserAnime : id int , idAnime int ,idUser int ,status string, episodeViewed int, favorite boolean  

Critique : idUser int, id anime , note int , contenu string, dateEcriture date    

## Mise a jour des données

Les données sont mis a jour via l'API Rest de l'application 
L'api de Jikan (API externe) est appellé quand on veut recuperer des informations à propos d'un animé


## Description du serveur 

On utilise un API REST afin de permet à l'utilisateur de mettre à jour sa liste d'animé, ecrire des critiques (a voir avec le temps disponible ), consulter une page d'un animé et s'inscrire.

On aura une page de login, une page perso de l'utilisateur, une page pour les infos de l'animé , une page pour ecrire un critique et une page pour modifier sa page perso (si temps et ressource possible).


## Description du client
- Login -> POST /login
- Register -> POST /register
- AnimePage -> GET anime/id
- UserPage -> GET user/id
- Reviews -> GET /anime/id/reviews


## Description des requetes 
Récupere la page de l'anime
GET url_du_site/anime/id

Reponse:

HTTP 1.1 200 OK

info de l'animé (JSON)

Se connecter au compte

POST url_du_site/login

Reponse : 

cookie de session (JWT probablement)

S'inscrire sur le site

POST url_du_site/register

pseudo : 

mdp : 


Envoie la critique à la BDD

POST url_du_site/anime/id/reviews/  

note : 

contenu :

Récupere la page de l'utilisateur

GET url_du_site/user/id

info de l'utilisateur

cookie de session (important)


Recherche les animés avec le mot-clé

GET url_du_site/search

Reponse :

le json avec les infos des animes qui corresponde



## Schéma global

Un site pour traquer sa progression dans le visionnage d'animé.

2 API :

1 externe : Jikan

1 interne : api rest de l'application

Une BDD avec plusieurs tables qui sera utilisé par l'api rest de l'application

Les actions fait sur la BDD se fait à l'aide de l'api interne.

