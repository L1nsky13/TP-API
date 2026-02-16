package controllers

import (
	"f1-app/models"
	"f1-app/templates"
	"net/http"
)

// ErrorDisplay
// -----------
// Objectif :
//   - Extraire le code et le message d'erreur des paramètres de requête.
//   - Afficher la page d'erreur avec les informations fournies.
func ErrorDisplay(w http.ResponseWriter, r *http.Request) {
	// Étape 1 : Récupérer les paramètres de query (code et message).
	data := models.Error{
		Code:    r.FormValue("code"),
		Message: r.FormValue("message"),
	}

	// Étape 2 : Déléguer le rendu au helper de templates (avec gestion d'erreur).
	templates.RenderTemplate(w, r, "error", data)
}
