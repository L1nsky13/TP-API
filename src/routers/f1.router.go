package routers

import (
	"f1-app/controllers"
	"net/http"
)

// f1Router
// -----------
// Objectif :
//   - Enregistrer toutes les routes F1 (pilotes, écuries, recherche, favoris).
//   - Configurer les handlers pour les pages principales de l'application.
func f1Router(router *http.ServeMux) {
	// Étape 1 : Enregistrer la route racine vers l'index.
	router.HandleFunc("/", controllers.IndexHandler)

	// Étape 2 : Enregistrer les routes de navigation principales.
	router.HandleFunc("/drivers", controllers.DriversHandler)
	router.HandleFunc("/teams", controllers.TeamsHandler)
	router.HandleFunc("/search", controllers.SearchHandler)

	// Étape 3 : Enregistrer les routes de détail avec paramètres dynamiques.
	router.HandleFunc("/teams/", controllers.TeamDetailHandler)
	router.HandleFunc("/drivers/", controllers.DriverDetailHandler)

	// Étape 4 : Enregistrer les routes de gestion des favoris.
	router.HandleFunc("/favorites", controllers.FavoritesHandler)
	router.HandleFunc("/add-favorite", controllers.AddFavoriteHandler)
	router.HandleFunc("/remove-favorite", controllers.RemoveFavoriteHandler)

	// Étape 5 : Enregistrer la route supplémentaire.
	router.HandleFunc("/about", controllers.AboutHandler)
}
