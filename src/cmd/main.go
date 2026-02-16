package main

import (
	"f1-app/routers"
	"f1-app/templates"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Chargement des templates au démarrage (fail fast si besoin dans Load()).
	templates.Load()

	// Construction du routeur principal (toutes les routes sont enregistrées dedans)
	mux := routers.MainRouter()

	// Démarrage du serveur HTTP
	addr := "localhost:8080"
	fmt.Printf("Serveur prêt sur http://%s\n", addr)
	// IMPORTANT : on passe bien "mux" à ListenAndServe pour utiliser NOTRE routeur,
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		// En cas d'erreur au lancement, on log et on sort.
		log.Fatalf("Erreur lancement serveur : %s\n", err.Error())
	}
}
