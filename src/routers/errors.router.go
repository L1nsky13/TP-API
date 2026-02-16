package routers

import (
	"f1-app/controllers"
	"net/http"
)

// errorRouter
// -----------
// Objectif :
//   - Enregistrer la route de gestion des erreurs.
//   - Afficher la page d'erreur avec codes d'erreur personnalisés.
func errorRouter(router *http.ServeMux) {
	// Étape 1 : Enregistrer la route /error pour afficher les pages d'erreur.
	router.HandleFunc("/error", controllers.ErrorDisplay)
}
