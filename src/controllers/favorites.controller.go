package controllers

import (
	"f1-app/helpers"
	"f1-app/models"
	"f1-app/services"
	"f1-app/templates"
	"net/http"
	"strings"
)

// FavoritesHandler
// ----------------
// Objectif :
//   - Afficher la page des favoris avec les pilotes et écuries favoris.
//   - Charger les favoris depuis le fichier JSON.
//   - Récupérer toutes les données puis filtrer pour n'afficher que les favoris.
//   - En cas de succès : rendre le template "favorites" avec les données.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func FavoritesHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est GET.
	if r.Method != http.MethodGet {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Charger les favoris depuis le fichier JSON.
	favorites, err := services.LoadFavorites()
	if err != nil {
		helpers.RedirectToError(w, r, http.StatusInternalServerError, "Impossible de charger les favoris")
		return
	}

	// Étape 3 : Récupérer tous les pilotes.
	driversData, statusDrivers, errDrivers := services.GetDriverStandingsService("2025", "", "", "", "", "")
	if statusDrivers != http.StatusOK || errDrivers != nil {
		helpers.RedirectToError(w, r, statusDrivers, "Impossible de récupérer les pilotes")
		return
	}

	allDrivers, ok := driversData.Data["allDrivers"].([]models.Driver)
	if !ok {
		helpers.RedirectToError(w, r, http.StatusInternalServerError, "Erreur de données des pilotes")
		return
	}

	// Étape 4 : Récupérer toutes les écuries.
	teamsData, statusTeams, errTeams := services.GetConstructorStandingsService("2025")
	if statusTeams != http.StatusOK || errTeams != nil {
		helpers.RedirectToError(w, r, statusTeams, "Impossible de récupérer les écuries")
		return
	}

	allConstructors, ok := teamsData.Data["constructors"].([]models.Constructor)
	if !ok {
		helpers.RedirectToError(w, r, http.StatusInternalServerError, "Erreur de données des écuries")
		return
	}

	// Étape 5 : Construire une map des pilotes par ID pour des recherches O(1).
	driverByID := make(map[string]models.Driver, len(allDrivers))
	for _, driver := range allDrivers {
		driverByID[driver.DriverID] = driver
	}

	// Étape 6 : Filtrer les pilotes favoris.
	favoriteDrivers := make([]models.Driver, 0, len(favorites.Drivers))
	for _, driverID := range favorites.Drivers {
		if driver, exists := driverByID[driverID]; exists {
			favoriteDrivers = append(favoriteDrivers, driver)
		}
	}

	// Étape 7 : Construire une map des écuries par ID pour des recherches O(1).
	constructorByID := make(map[string]models.Constructor, len(allConstructors))
	for _, constructor := range allConstructors {
		constructorByID[constructor.ConstructorID] = constructor
	}

	// Étape 8 : Filtrer les écuries favorites.
	favoriteConstructors := make([]models.Constructor, 0, len(favorites.Constructors))
	for _, constructorID := range favorites.Constructors {
		if constructor, exists := constructorByID[constructorID]; exists {
			favoriteConstructors = append(favoriteConstructors, constructor)
		}
	}

	// Étape 9 : Préparer les données pour le template.
	data := &models.PageData{
		Title:       "My Favorites",
		CurrentPage: "favorites",
		Data: map[string]interface{}{
			"drivers":      favoriteDrivers,
			"constructors": favoriteConstructors,
			"season":       "2025",
		},
	}

	// Étape 10 : Rendre le template "favorites" avec les données.
	templates.RenderTemplate(w, r, "favorites", data)
}

// AddFavoriteHandler
// ------------------
// Objectif :
//   - Ajouter un élément (pilote ou écurie) aux favoris.
//   - Récupérer les paramètres type, id et returnUrl depuis le formulaire.
//   - Appeler le service approprié selon le type (driver ou constructor).
//   - Rediriger vers la page d'origine après l'ajout.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est POST.
	if r.Method != http.MethodPost {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Récupérer les paramètres du formulaire.
	itemType := r.FormValue("type")
	itemID := r.FormValue("id")
	returnURL := r.FormValue("returnUrl")

	if itemID == "" || itemType == "" {
		helpers.RedirectToError(w, r, http.StatusBadRequest, "Paramètres manquants")
		return
	}

	// Étape 3 : Ajouter aux favoris selon le type (driver ou constructor).
	var err error
	switch itemType {
	case "driver":
		err = services.AddDriverToFavorites(itemID)
	case "constructor":
		err = services.AddConstructorToFavorites(itemID)
	default:
		helpers.RedirectToError(w, r, http.StatusBadRequest, "Type invalide")
		return
	}

	// Étape 4 : Vérifier s'il y a eu une erreur lors de l'ajout.
	if err != nil {
		helpers.RedirectToError(w, r, http.StatusInternalServerError, "Impossible d'ajouter aux favoris")
		return
	}

	// Étape 5 : Rediriger vers la page d'origine.
	if returnURL == "" {
		returnURL = "/"
	}
	http.Redirect(w, r, returnURL, http.StatusSeeOther)
}

// RemoveFavoriteHandler
// ---------------------
// Objectif :
//   - Supprimer un élément (pilote ou écurie) des favoris.
//   - Récupérer les paramètres type, id et returnUrl depuis le formulaire.
//   - Appeler le service approprié selon le type (driver ou constructor).
//   - Rediriger vers la page d'origine après la suppression.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est POST.
	if r.Method != http.MethodPost {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Récupérer les paramètres du formulaire.
	itemType := r.FormValue("type")
	itemID := r.FormValue("id")
	returnURL := r.FormValue("returnUrl")

	if itemID == "" || itemType == "" {
		helpers.RedirectToError(w, r, http.StatusBadRequest, "Paramètres manquants")
		return
	}

	// Étape 3 : Supprimer des favoris selon le type (driver ou constructor).
	var err error
	switch itemType {
	case "driver":
		err = services.RemoveDriverFromFavorites(itemID)
	case "constructor":
		err = services.RemoveConstructorFromFavorites(itemID)
	default:
		helpers.RedirectToError(w, r, http.StatusBadRequest, "Type invalide")
		return
	}

	// Étape 4 : Vérifier s'il y a eu une erreur lors de la suppression.
	if err != nil {
		helpers.RedirectToError(w, r, http.StatusInternalServerError, "Impossible de supprimer des favoris")
		return
	}

	// Étape 5 : Rediriger vers la page d'origine.
	if returnURL == "" {
		returnURL = "/"
	}
	http.Redirect(w, r, returnURL, http.StatusSeeOther)
}

// AboutHandler
// ------------
// Objectif :
//   - Afficher la page About de l'application.
//   - Vérifier que l'URL est exactement "/about" ou "/about/".
//   - En cas de succès : rendre le template "about" avec les données.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func AboutHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est GET.
	if r.Method != http.MethodGet {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Extraire le path après /about et vérifier qu'il est valide.
	path := strings.TrimPrefix(r.URL.Path, "/about")
	if path != "" && path != "/" {
		helpers.RedirectToError(w, r, http.StatusNotFound, "Page non trouvée")
		return
	}

	// Étape 3 : Préparer les données pour le template.
	data := &models.PageData{
		Title:       "About",
		CurrentPage: "about",
		Data:        map[string]interface{}{},
	}

	// Étape 4 : Rendre le template "about" avec les données.
	templates.RenderTemplate(w, r, "about", data)
}
