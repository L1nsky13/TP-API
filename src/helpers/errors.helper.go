package helpers

import (
	"net/http"
	"net/url"
	"strconv"
)

// RedirectToError
// ---------------
// Objectif :
//   - Rediriger vers la page /error en passant le code HTTP et le message en paramètres.
//   - Construire une URL sécurisée avec paramètres codés.
//   - Utiliser HTTP status SeeOther (303) pour une redirection POST-to-GET.
func RedirectToError(w http.ResponseWriter, r *http.Request, code int, message string) {

	// Étape 1 : Construire les paramètres de query de manière sécurisée.
	params := url.Values{}
	if code > 0 {
		params.Set("code", strconv.Itoa(code))
	}
	if message != "" {
		params.Set("message", message)
	}

	// Étape 2 : Construire l'URL cible avec paramètres encodés.
	pathTarget := "/error"
	if encodeParams := params.Encode(); encodeParams != "" {
		pathTarget += "?" + encodeParams
	}

	// Étape 3 : Rediriger vers la page d'erreur.
	http.Redirect(w, r, pathTarget, http.StatusSeeOther)
}
