# 🏎️ API - Formula 1 2025

## 📋 Présentation du Projet

### Thème
Le thème sera la Formule 1 (site en version Anglaise) avec l'API d'Ergast F1 API **(https://api.jolpi.ca/ergast/)** tout en permettant aux utilisateurs d'explorer les pilotes et leurs écuries de la saison 2025 et de pouvoir les mettre en favori à travers un site dynamique en **HTML/CSS/GO**.

### Fonctionnalités
- 🏁 **Affichage des Pilotes** : Liste complète des pilotes F1 avec détails individuels
- 🚗 **Affichage des Écuries** : Liste des constructeurs et leurs informations
- 🔍 **Recherche Globale** : Recherche unifiée dans les pilotes et les écuries
- ❤️ **Système de Favoris** : Ajouter/supprimer des favoris (persistance locale)
- 📊 **Filtrage Avancé** : Par équipe, nationalité, type de pilote (titulaire, test, réserve)
- 📄 **Pagination** : Navigation efficace à travers les données
- 🎵 **Ambiance F1** : Son au changement de page (Max Verstappen) + Une musique par page

---

## 🚀 Installation et Lancement

### Prérequis
- **go 1.25.3**
- **Navigateur Web** (Chrome, Firefox, Safari, Edge)

### Installation

1. **Lancer l'application**
```bash
cd src/cmd
go run main.go
```

Le serveur démarre sur : **http://localhost:8080**

2. **Structure du projet**
```
.
├── src/
│   ├── cmd/
│   │   └── main.go                     # Point d'entrée de l'application
│   ├── controllers/                    
│   │       ├── errors.controller.go    # Gestion des pages d'erreur (404, 500, etc.)
│   │       ├── f1.controller.go        # Handlers pour pilotes, équipes, recherche
│   │       └── favorites.controller.go # Handlers pour ajouter/retirer favoris
│   ├── helpers/                        
│   │       └── errors.helper.go        # Fonctions d'aide pour redirection erreurs
│   ├── models/
│   │       ├── data.model.go           # Modèle pour les détails des pilotes et des écuries                        
│   │       ├── errors.model.go         # Modèle pour gestion d'erreurs
│   │       └── f1.model.go             # Modèles Driver, Constructor, PageData
│   ├── routers/
│   │       ├── errors.router.go        # Routes pour pages d'erreur
│   │       ├── f1.router.go            # Routes pour pilotes, équipes, favoris
│   │       └── main.router.go          # Routeur principal + fichiers statiques
│   ├── services/
│   │       ├── f1.services.go          # Filtrage, pagination, recherche
│   │       └── favorites.service.go    # Gestion des favoris (CRUD)
│   ├── templates/
│   │       └── templates.go            # Rendu des templates HTML
│   └── go.mod                          # Dépendances Go
├── templates/                          
│       ├── about.html                  # Page À Propos avec FAQ projet
│       ├── drivers-detail.html         # Détail d'un pilote spécifique
│       ├── drivers.html                # Liste des pilotes avec filtres
│       ├── error.html                  # Page d'erreur générique
│       ├── favorites.html              # Liste des favoris utilisateur
│       ├── index.html                  # Accueil du site
│       ├── search.html                 # Résultats de recherche globale
│       ├── teams-detail.html           # Détail d'une écurie spécifique
│       └── teams.html                  # Liste des écuries
├── assets/
│       ├── *.css                       # Feuilles de style (header, drivers, teams, etc.)
│       ├── *.js                        # Scripts clients (audio persistence)
│       ├── *.mp3                       # Fichiers audio (F1 themes)
│       ├── *.ttf                       # Polices Formula 1 officielles
│       └── formula1-logo.webp          # Logo et images F1
├── favorites.json                      # Favoris stockés (JSON)
└── README.md                           # Documentation
```

---

## 🛣️ Routes et Endpoints

### Routes Frontend (HTML)

| Route | Méthode | Description |
|-------|---------|-------------|
| `/` | GET | Page d'accueil avec présentation F1 |
| `/drivers` | GET | Liste complète des pilotes avec filtres |
| `/drivers/:id` | GET | Détails d'un pilote spécifique |
| `/teams` | GET | Liste de toutes les écuries |
| `/teams/:id` | GET | Détails d'une écrie spécifique |
| `/search` | GET | Page de résultats de recherche globale |
| `/favorites` | GET | Liste des favoris de l'utilisateur |
| `/about` | GET | Page À Propos avec FAQ projet |

### Routes d'Actions (API Interne)

| Route | Méthode | Description |
|--------|---------|-------------|
| `/add-favorite` | POST | Ajouter un pilote/écurie aux favoris |
| `/remove-favorite` | POST | Retirer un pilote/écurie des favoris |

### Ressources Statiques

| Type | Endpoint | Description |
|------|----------|-------------|
| CSS | `/static/*.css` | Feuilles de style |
| JS | `/static/*.js` | Fichiers JavaScript |
| MP3 | `/static/*.mp3` | Fichiers audio |
| Fonts | `/static/*.ttf` | Polices Formula 1 |
| Images | `/static/*.webp` | Images et logos |

---

## 📡 API Externe - Ergast F1 API

### Endpoints Exploités

L'application utilise les données de l'**Ergast F1 API** (https://api.jolpi.ca/ergast/)

#### 1. **Pilotes (Drivers)**
```
GET https://api.jolpi.ca/ergast/f1/2025/drivers/?format=json
```
- **Paramètres** : season (année), driverId (optionnel)
- **Données** : Tous les pilotes de la saison avec infos personnelles
- **Filtres appliqués** : Équipe, nationalité, type (titulaire/test/réserve)

##### Format de Réponse avant modification (Pilote titulaire)
```json
{
    "MRData": {
        "xmlns": "",
        "series": "f1",
        "url": "https://api.jolpi.ca/ergast/f1/2025/drivers/",
        "limit": "30",
        "offset": "0",
        "total": "30",
        "DriverTable": {
            "season": "2025",
            "Drivers": [
                {
                    "driverId": "albon",
                    "permanentNumber": "23",
                    "code": "ALB",
                    "url": "http://en.wikipedia.org/wiki/Alexander_Albon",
                    "givenName": "Alexander",
                    "familyName": "Albon",
                    "dateOfBirth": "1996-03-23",
                    "nationality": "Thai"
                }
            "Suite..."
            ]
        }
    }
}            
```

##### Rajout des autres paramètres

- **image** : Changement de nom (url -> image) avec une photo du pilote (de la tête au pied)
- **driverType** : Poste du pilote (Titulaire, réserve et essai)

##### Format de Réponse après modification (Pilote titulaire)
```json
{
    "MRData": {
        "xmlns": "",
        "series": "f1",
        "url": "https://api.jolpi.ca/ergast/f1/2025/drivers/",
        "limit": "30",
        "offset": "0",
        "total": "30",
        "DriverTable": {
            "season": "2025",
            "Drivers": [
                {
                    "driverId": "albon",
                    "permanentNumber": "23",
                    "code": "ALB",
                    "image": "https://media.formula1.com/image/upload/c_fill,w_720/q_auto/v1740000000/common/f1/2025/williams/alealb01/2025williamsalealb01right.webp",
                    "givenName": "Alexander",
                    "familyName": "Albon",
                    "dateOfBirth": "1996-03-23",
                    "nationality": "Thai",
                    "team": "Williams",
                    "driverType" : "Race Driver"
                },
            "Suite..."
            ]
        }
    }
}            
```

##### Format de Réponse avant modification (Pilote d'essai)
```json
{
    "MRData": {
        "xmlns": "",
        "series": "f1",
        "url": "https://api.jolpi.ca/ergast/f1/2025/drivers/",
        "limit": "30",
        "offset": "0",
        "total": "30",
        "DriverTable": {
            "season": "2025",
            "Drivers": [
            "Suite..."
                {
                    "driverId": "paul_aron",
                    "givenName": "Paul",
                    "familyName": "Aron"
                }
            "Suite..."
            ]
        }
    }
}            
```

##### Rajout des autres paramètres

- **permanentNumber** : Numéro de sa monoplace
- **code** : 3 premières lettres de son nom de famille
- **image** : Photo de profil du pilote sur le jeu F1 Manager 2024
- **dateOfBirth** : Date de naissance
- **nationality** : Sa nationalité
- **team** : Son écurie
- **driverType** : Poste du pilote (Titulaire, réserve et essai)

##### Format de Réponse après modification (Pilote d'essai)
```json
{
    "MRData": {
        "xmlns": "",
        "series": "f1",
        "url": "https://api.jolpi.ca/ergast/f1/2025/drivers/",
        "limit": "30",
        "offset": "0",
        "total": "30",
        "DriverTable": {
            "season": "2025",
            "Drivers": [
            "Suite..."
                {
                    "driverId": "paul_aron",
                    "permanentNumber": "61",
                    "code": "ARO",
                    "image": "https://image-service.zaonce.net/eyJidWNrZXQiOiJmcm9udGllci1jbXMiLCJrZXkiOiJmMW1hbmFnZXIvMjAyNC9kcml2ZXJzL2hlYWRzaG90cy9mMi9hcm8ucG5nIiwiZWRpdHMiOnsicmVzaXplIjp7IndpZHRoIjo1MDB9fX0=",
                    "givenName": "Paul",
                    "familyName": "Aron",
                    "dateOfBirth": "2004-02-04",
                    "nationality": "Estonian",
                    "team": "Alpine",
                    "driverType" : "Test Driver"
                }
            "Suite..."
            ]
        }
    }
}            
```

---

#### 2. **Constructeurs (Teams)**
```
GET /api/f1/2025/constructors.json
```
- **Données** : Toutes les écuries et leurs informations
- **Format** : Nom, nationalité, historique

##### Format de Réponse avant modification (Écurie)

```json
{
    "MRData": {
        "xmlns": "",
        "series": "f1",
        "url": "https://api.jolpi.ca/ergast/f1/2025/constructors/",
        "limit": "30",
        "offset": "0",
        "total": "10",
        "ConstructorTable": {
            "season": "2025",
            "Constructors": [
                {
                    "constructorId": "alpine",
                    "url": "https://en.wikipedia.org/wiki/Alpine_F1_Team",
                    "name": "Alpine F1 Team",
                    "nationality": "French"
                }
            "Suite..."
            ]
        }
    }
}  
```
##### Rajout des autres paramètres

- **icon** : Logo de l'écurie
- **image** : Changement de nom (url -> image) et de la photo de la monoplace (2025)
- **teamColor** : Couleur qui correspond le plus à l'écurie

---

##### Format de Réponse avant modification (Écurie)

```json
{
    "MRData": {
        "xmlns": "",
        "series": "f1",
        "url": "https://api.jolpi.ca/ergast/f1/2025/constructors/",
        "limit": "30",
        "offset": "0",
        "total": "10",
        "ConstructorTable": {
            "season": "2025",
            "Constructors": [
                {
                    "constructorId": "alpine",
                    "icon": "https://media.formula1.com/image/upload/c_fit,h_64/q_auto/v1740000000/common/f1/2025/alpine/2025alpinelogowhite.webp",
                    "image": "https://media.formula1.com/image/upload/c_lfill,w_3392/q_auto/v1740000000/common/f1/2025/alpine/2025alpinecarright.webp",
                    "name": "Alpine F1 Team",
                    "nationality": "French",
                    "teamColor": "#005081"
                }
            "Suite..."
            ]
        }
    }
}  
```
---

## 🎨 Fonctionnalités Avancées

### Système de Filtrage
```
GET /drivers?season=2025&team=ferrari&nationality=Italian&driverType=RACE_DRIVER
```

Paramètres supportés :
- `season` : Année (défaut: 2025)
- `team` : Nom de l'équipe
- `nationality` : Nationalité
- `driverType` : RACE_DRIVER, TEST_DRIVER, RESERVE_DRIVER
- `page` : Numéro de page
- `perPage` : Éléments par page

### Recherche Globale
```
GET /search?q=verstappen
```
Recherche simultanée dans :
- Noms de pilotes (givenName + surname)
- Noms d'écuries
- Nationalités

### Favoris (localStorage)
Structure de stockage :
```javascript
{
  "favorites": ["verstappen", "leclerc", "ferrari"]
}
```

### Audio Immersif
- Persistance du lecteur F1 (position, état lecture)

---

## 📝 Scripts Clients

### `audio-persistence.js`
- Sauvegarde l'état du lecteur audio
- Restaure la position de lecture
- Sync localStorage au déchargement

---

## 🐛 Gestion des Erreurs

### Pages d'Erreur Dédiées
- **301 Moved Permanently** :  Redirection permanente vers un autre URL
- **400 Bad Request** : Requête invalide ou mal formulée
- **404 Not Found** : Ressource demandée introuvable
- **500 Internal Server Error** : Erreur interne du serveur

### Fallback
- Redirection automatique vers page erreur
- Affichage message utilisateur

---

## 📄 Licence

Ce projet a été réalisé par Théodore NAJMAN, B1 - Informatique