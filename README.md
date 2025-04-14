# Documentation du Projet RealTimeForum

## Vue d'ensemble du projet

Le projet **RealTimeForum** est une application de forum avancée qui supporte les interactions en temps réel via WebSockets, la messagerie privée et des fonctionnalités utilisateur améliorées. Les utilisateurs peuvent s'inscrire, se connecter, créer et commenter des publications, et envoyer des messages privés en temps réel. Le projet utilise **SQLite** pour la gestion de la base de données, **Golang** pour la logique backend, et **JavaScript** pour les interactions frontend. Le frontend est une Single Page Application (SPA) avec des mises à jour d'interface utilisateur dynamiques.

## Table des matières

1. [Vue d'ensemble du projet](#vue-densemble-du-projet)
2. [Architecture technique](#architecture-technique)
3. [Instructions de configuration](#instructions-de-configuration)
4. [Dépendances et prérequis](#dépendances-et-prérequis)
5. [Configuration](#configuration)
6. [Documentation de l'API](#documentation-de-lapi)
7. [Cas d'utilisation courants](#cas-dutilisation-courants)
8. [Guide de dépannage](#guide-de-dépannage)
9. [Considérations de sécurité](#considérations-de-sécurité)
10. [Optimisations de performance](#optimisations-de-performance)
11. [Lignes directrices pour les tests](#lignes-directrices-pour-les-tests)
12. [Processus de déploiement](#processus-de-déploiement)
13. [Procédures de maintenance](#procédures-de-maintenance)
14. [Informations de contact](#informations-de-contact)

## Architecture technique

### Vue d'ensemble de l'architecture

Le projet RealTimeForum est structuré en plusieurs couches pour assurer une séparation claire des préoccupations. Voici une vue d'ensemble de l'architecture :

```mermaid
graph TD
    A[Client (Frontend)] --> B[WebSocket Server]
    A --> C[HTTP Server]
    B --> D[WebSocket Handler]
    C --> E[HTTP Handlers]
    D --> F[WebSocket Service]
    E --> G[Services]
    G --> H[Repositories]
    H --> I[Database (SQLite)]
```

### Modèles de conception et décisions architecturales

- **Modèle MVC (Model-View-Controller)** : Utilisé pour séparer les préoccupations entre la logique métier, l'interface utilisateur et la gestion des données.
- **WebSockets** : Utilisés pour les interactions en temps réel, permettant des mises à jour instantanées sans rechargement de la page.
- **Microservices** : Les différentes fonctionnalités (authentification, gestion des messages, etc.) sont encapsulées dans des services distincts pour une meilleure maintenabilité et évolutivité.

### Documentation des API et interfaces

#### API d'authentification

```http
POST /register
Content-Type: application/json

{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepassword"
}
```

```http
POST /login
Content-Type: application/json

{
    "email": "john@example.com",
    "password": "securepassword"
}
```

#### API de gestion des messages

```http
GET /messages?userId=123
Content-Type: application/json

{
    "messages": [
        {
            "messageId": "uuid",
            "senderId": "uuid",
            "receiverId": "uuid",
            "content": "Hello, World!",
            "sentAt": "2023-10-01T12:00:00Z",
            "is_read": false
        }
    ]
}
```

### Exemples de code avec explications

#### Exemple de gestion des WebSockets

```go
// websocket_handler.go
func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    // Authenticate user
    user, err := authenticateUser(r)
    if err != nil {
        conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
    }

    // Add client to hub
    hub.addClient(conn, user)
    defer hub.removeClient(conn)

    // Handle messages
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            break
        }
        hub.broadcastMessage(message)
    }
}
```

### Options de configuration et leur impact

- **Configuration de la base de données** : Le fichier `config/database.go` gère la connexion à la base de données SQLite. Les options de configuration incluent le chemin du fichier de la base de données et les paramètres de connexion.
- **Configuration du serveur** : Le fichier `config/server.go` définit les paramètres du serveur HTTP, tels que le port d'écoute et les délais d'attente.

## Instructions de configuration

### Prérequis

- **Go** : Version 1.23.3 ou supérieure
- **Node.js** : Version 14 ou supérieure
- **SQLite** : Version 3 ou supérieure

### Installation des dépendances

```sh
go mod tidy
npm install
```

### Configuration de la base de données

```sh
go run cmd/main.go initdb
```

### Lancement du serveur

```sh
go run cmd/main.go
```

## Dépendances et prérequis

### Dépendances Go

```go
require (
    github.com/gofrs/uuid v4.4.0+incompatible
    github.com/mattn/go-sqlite3 v1.14.24
    github.com/gorilla/websocket v1.5.0
)
```

### Dépendances JavaScript

```json
{
    "dependencies": {
        "axios": "^0.21.1",
        "bootstrap": "^5.1.3",
        "jquery": "^3.6.0"
    }
}
```

## Configuration

### Fichier de configuration `config/database.go`

```go
func ConnectDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./forum.db")
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
```

### Fichier de configuration `config/server.go`

```go
func StartServer() {
    http.HandleFunc("/", HomeHandler)
    http.HandleFunc("/ws", WebsocketHandler)
    log.Println("Starting server on :8443")
    if err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil); err != nil {
        log.Fatal(err)
    }
}
```

## Documentation de l'API

### API d'authentification

#### Endpoint `/register`

- **Méthode** : POST
- **Description** : Inscription d'un nouvel utilisateur
- **Corps de la requête** : JSON avec `username`, `email`, et `password`

#### Endpoint `/login`

- **Méthode** : POST
- **Description** : Connexion d'un utilisateur
- **Corps de la requête** : JSON avec `email` et `password`

### API de gestion des messages

#### Endpoint `/messages`

- **Méthode** : GET
- **Description** : Récupération des messages pour un utilisateur
- **Paramètres de requête** : `userId`

## Cas d'utilisation courants

### Inscription d'un nouvel utilisateur

1. L'utilisateur accède à la page d'inscription.
2. L'utilisateur remplit le formulaire avec `username`, `email`, et `password`.
3. Le formulaire est soumis via une requête POST à `/register`.
4. Le serveur traite la requête et enregistre l'utilisateur dans la base de données.

### Connexion d'un utilisateur

1. L'utilisateur accède à la page de connexion.
2. L'utilisateur remplit le formulaire avec `email` et `password`.
3. Le formulaire est soumis via une requête POST à `/login`.
4. Le serveur vérifie les informations d'identification et renvoie un token de session.

## Guide de dépannage

### Problèmes courants et solutions

| **Problème**                  | **Solution**                                                                 |
|--------------------------------|-----------------------------------------------------------------------------|
| **Erreur de connexion à la base de données** | Vérifiez le chemin du fichier de la base de données et les permissions. |
| **Erreur de WebSocket**        | Assurez-vous que le serveur WebSocket est en cours d'exécution.           |
| **Erreur d'authentification**  | Vérifiez les informations d'identification et les tokens de session.       |

## Considérations de sécurité

### Meilleures pratiques de sécurité

- **Utilisation de HTTPS** : Assurez-vous que toutes les communications sont chiffrées.
- **Validation des entrées** : Validez toujours les entrées utilisateur pour éviter les injections SQL et les attaques XSS.
- **Gestion des sessions** : Utilisez des tokens de session sécurisés et expirez-les régulièrement.

## Optimisations de performance

### Conseils de performance

- **Utilisation de la mise en cache** : Cachez les données fréquemment utilisées pour réduire les accès à la base de données.
- **Optimisation des requêtes SQL** : Utilisez des index pour améliorer les performances des requêtes.
- **Gestion des connexions WebSocket** : Limitez le nombre de connexions WebSocket simultanées pour éviter la surcharge du serveur.

## Lignes directrices pour les tests

### Stratégies de test

- **Tests unitaires** : Écrivez des tests unitaires pour chaque fonction et méthode.
- **Tests d'intégration** : Testez les interactions entre les différentes parties du système.
- **Tests de charge** : Simulez des charges élevées pour vérifier la robustesse du système.

### Exemple de test unitaire

```go
func TestUserRegistration(t *testing.T) {
    user := models.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password",
    }
    err := services.RegisterUser(user)
    if err != nil {
        t.Fatalf("Failed to register user: %v", err)
    }
}
```

## Processus de déploiement

### Étapes de déploiement

1. **Préparation** : Assurez-vous que toutes les dépendances sont installées et que la base de données est configurée.
2. **Construction** : Construisez l'application avec `go build`.
3. **Déploiement** : Déployez l'application sur le serveur de production.
4. **Vérification** : Vérifiez que l'application fonctionne correctement en production.

## Procédures de maintenance

### Procédures de mise à jour

1. **Planification** : Planifiez les mises à jour pendant les périodes de faible activité.
2. **Sauvegarde** : Effectuez une sauvegarde complète de la base de données.
3. **Déploiement** : Déployez les nouvelles versions de l'application.
4. **Vérification** : Vérifiez que toutes les fonctionnalités sont opérationnelles après la mise à jour.

### Recommandations de surveillance

- **Surveillance des logs** : Utilisez des outils comme `ELK Stack` pour surveiller les logs.
- **Surveillance des performances** : Utilisez des outils comme `Prometheus` et `Grafana` pour surveiller les performances.

### Stratégies de sauvegarde et de récupération

- **Sauvegarde régulière** : Effectuez des sauvegardes régulières de la base de données.
- **Plan de récupération** : Ayez un plan de récupération en cas de perte de données.

### Liste des tâches de maintenance et des calendriers

| **Tâche**                  | **Fréquence** | **Responsable** |
|----------------------------|---------------|-----------------|
| Sauvegarde de la base de données | Hebdomadaire | Admin Sys       |
| Mise à jour des dépendances | Mensuelle    | Dev Team        |
| Surveillance des logs     | Quotidienne   | Ops Team        |


## Lignes directrices pour les contributions

Pour contribuer au projet, veuillez suivre ces étapes :

1. **Fork** le dépôt.
2. **Créez une branche** pour votre fonctionnalité.
3. **Commitez vos modifications** avec des messages de commit clairs.
4. **Ouvrez une Pull Request** vers la branche `main`.

---

**Note** : Cette documentation est un guide complet pour le projet RealTimeForum. Pour plus d'informations, consultez la [documentation officielle de Go](https://golang.org/doc/) et la [documentation de SQLite](https://www.sqlite.org/docs.html).
