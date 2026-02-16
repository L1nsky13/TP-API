package routers

import (
	"net/http"
	"os"
	"path/filepath"
)

// MainRouter
// ----------
// Objectif :
//   - Initialiser et configurer le routeur principal de l'application.
//   - Enregistrer toutes les routes métier (erreurs, F1).
//   - Configurer le serveur de fichiers statiques pour CSS, JS, images et audio.
//   - Retourner le routeur configuré prêt à être utilisé.
func MainRouter() *http.ServeMux {

	// Étape 1 : Créer le routeur principal avec http.ServeMux.
	mainRouter := http.NewServeMux()

	// Étape 2 : Enregistrer les routes de gestion des erreurs.
	errorRouter(mainRouter)

	// Étape 3 : Enregistrer les routes métier de Formule 1.
	f1Router(mainRouter)

	// Étape 4 : Déterminer le chemin du répertoire des assets.
	wd, _ := os.Getwd()
	assetsPath := filepath.Join(wd, "..", "..")

	if filepath.Base(wd) == "src" {
		assetsPath = filepath.Join(wd, "..")
	}
	assetsPath = filepath.Join(assetsPath, "assets")

	// Étape 5 : Créer le serveur de fichiers statiques.
	fileServer := http.FileServer(http.Dir(assetsPath))

	// Étape 6 : Enregistrer la route /static/ pour servir les fichiers statiques.
	mainRouter.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Étape 7 : Retourner le routeur configuré.
	return mainRouter
}
