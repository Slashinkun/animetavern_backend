# Projet

Un site qui permet de traquer les animés qu'on a vu, qu'on veut voir.
Mais aussi, de noter les animés et ecrire des critiques qui seront vu par d'autre utilisateur. Egalement, l'utilisateur peut personnaliser sa page utilisateur en mettant une photo de profil et une bannière . De plus, il pourra mettre ses animes favoris.


## API
On utilise l'API Jikan : 
https://jikan.moe/

Elle permet de récuperer les données sur les animés (nom,image,genre,personnages,realisateur...)

On utilisera une API REST interne pour les actions de l'utilisateur.

## Fonctionnalités
- Creer un compte
- Personnaliser sa page (photo de profil et bannière) (si temps possible)
- Ajouter un anime à sa liste et changer son statut
- Noter et ecrire une critique sur un animé (si temps possible)
- Rechercher un animé depuis la barre de recherche
- Mettre en favoris un animé

## Cas d'utilisations
- Johnny a un compte sur le site, il se connecte et cherche un anime qu'il veut regarder. Il en trouve un,il lit les avis dessus et decide de le mettre dans sa liste a regarder.
- Michel a fini de regarder un anime et il decide d'ecrire une critique et mettre une note. (si temps possible)
- Diana veut personnaliser son compte et cherche une photo de profil et une  sur Internet
Elle trouve les images et decide de les rajouter sur son profil   (si temps possible)
-Frank a fini un animé et l'a bien aimé. Il décide le rajouter à ses favoris.

## BDD
4 tables :

Utilisateur : idUser int,pseudo string,email string,password string,photo_profil string ,banniere string

Anime : idAnime int , title string ,image string ,note int

UserAnime : id int , idAnime int ,idUser int ,status string, episodeViewed int, favorite boolean

Critique : idUser int, id anime , note int , contenu string, dateEcriture date    (a voir avec le temps disponible)

## Mise a jour des données

Les données sont mis a jour via l'API Rest de l'application
L'api de Jikan (API externe) est appellé quand on veut recuperer des informations à propos d'un animé


## Description du serveur 

On utilise un API REST (CRUD) afin de permet à l'utilisateur de mettre à jour sa liste d'animé, ecrire des critiques (a voir avec le temps disponible ), consulter une page d'un animé et s'inscrire

On aura une page de login, une page liste de l'utilisateur, une page pour les infos de l'animé et une page pour ecrire un critique (a voir avec le temps disponible).


## Description des requetes 

GET site/anime/id
Récupere la page de l'animé

POST site/login
Se connecter au compte

POST site/register
S'inscrire sur le site

POST site/anime/id/reviews/  (à voir avec le temps disponible)
Envoie la critique à la BDD

GET site/user/id
Récupere la page de l'utilisateur

GET site/animes?search=motclé
Recherche les animés avec le mot-clé



## Schéma global

Un site pour traquer sa progression dans le visionnage d'animé.

2 API :

1 externe : Jikan
1 interne : api rest de l'application

Une BDD avec plusieurs tables qui sera utilisé par l'api rest de l'application
Les actions fait sur la BDD se fait à l'aide de l'api interne.
# animetavern_backend
