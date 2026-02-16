package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

var listTemp *template.Template

// getFuncMap
// Retourne un map des fonctions personnalisées disponibles dans les templates.
func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatDuration": func(ms int) string {
			totalSeconds := ms / 1000
			minutes := totalSeconds / 60
			seconds := totalSeconds % 60
			return fmt.Sprintf("%d:%02d", minutes, seconds)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"iterate": func(count int) []int {
			var items []int
			for i := 0; i < count; i++ {
				items = append(items, i)
			}
			return items
		},
	}
}

// Load
// Charge tous les fichiers de templates HTML au démarrage de l'application.
func Load() {
	// Récupérer le répertoire de travail actuel.
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erreur récupération du répertoire de travail : %s", err.Error())
	}

	// Ajuster le chemin selon le répertoire d'exécution.
	if filepath.Base(wd) == "cmd" {
		wd = filepath.Join(wd, "..", "..")
	} else if filepath.Base(wd) == "src" {

		wd = filepath.Join(wd, "..")
	}

	// Construire le pattern pour charger tous les fichiers HTML.
	pattern := filepath.Join(wd, "templates", "*.html")

	// Créer un template avec les fonctions personnalisées et charger tous les fichiers.
	tmpl := template.New("").Funcs(getFuncMap())
	listTemplates, errTemplates := tmpl.ParseGlob(pattern)
	if errTemplates != nil {
		// Arrêter le programme si les templates ne peuvent pas être chargés.
		log.Fatalf("Erreur chargement des templates : %s", errTemplates.Error())
	}
	// Sauvegarder le template global pour l'utiliser partout.
	listTemp = listTemplates
}

// RenderTemplate
// Exécute un template et écrit la réponse HTTP. En cas d'erreur, redirige vers la page d'erreur.
func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	// Étape 1 : Exécuter le template dans un buffer (sans envoyer au client).
	var buffer bytes.Buffer

	errRender := listTemp.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		// Étape 2 : En cas d'erreur, logger l'erreur et rediriger.
		log.Printf("erreur rendu template '%s': %v", name, errRender)

		errorMsg := fmt.Sprintf("Erreur template '%s': %v", name, errRender)
		http.Redirect(w, r, fmt.Sprintf("/error?code=%d&message=%s", http.StatusInternalServerError, url.QueryEscape(errorMsg)), http.StatusSeeOther)
		return
	}

	// Étape 3 : Envoyer le buffer au client en réponse HTTP.
	_, _ = buffer.WriteTo(w)
}
