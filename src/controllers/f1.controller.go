package controllers

import (
	"f1-app/helpers"
	"f1-app/models"
	"f1-app/services"
	"f1-app/templates"
	"fmt"
	"net/http"
)

// DriversHandler
// ---------------
// Objectif :
//   - Afficher la liste des pilotes F1 avec système de filtrage et pagination.
//   - Récupérer les paramètres de filtrage (équipe, nationalité, type) et pagination depuis l'URL.
//   - En cas de succès : rendre le template "drivers" avec les données.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func DriversHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est GET.
	if r.Method != http.MethodGet {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Récupérer les paramètres depuis l'URL.
	season := r.URL.Query().Get("season")
	if season == "" {
		season = "2025"
	}
	teamFilter := r.URL.Query().Get("team")
	nationalityFilter := r.URL.Query().Get("nationality")
	driverTypeFilter := r.URL.Query().Get("driverType")
	pageParam := r.URL.Query().Get("page")
	perPageParam := r.URL.Query().Get("perPage")

	// Étape 3 : Appeler services.GetDriverStandingsService avec les filtres.
	data, status, err := services.GetDriverStandingsService(season, teamFilter, nationalityFilter, driverTypeFilter, pageParam, perPageParam)

	// Étape 4 : Vérifier si status != http.StatusOK ou err != nil.
	// Si erreur → helpers.RedirectToError(...) + fmt.Println(err) + return.
	if status != http.StatusOK || err != nil {
		fmt.Println("Erreur lors de la récupération des pilotes:", err)
		helpers.RedirectToError(w, r, status, "Impossible de récupérer les pilotes")
		return
	}

	// Étape 5 : Si succès, rendre le template "drivers" avec les données.
	templates.RenderTemplate(w, r, "drivers", data)
}

// SearchHandler
// -------------
// Objectif :
//   - Gérer la recherche globale dans les pilotes ET les écuries.
//   - Récupérer toutes les données puis filtrer selon la query de recherche.
//   - En cas de succès : rendre le template "search" avec les résultats.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func SearchHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est GET.
	if r.Method != http.MethodGet {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Récupérer la query de recherche depuis l'URL.
	query := r.URL.Query().Get("q")

	// Étape 3 : Récupérer TOUTES les données des pilotes.
	driversData, statusDrivers, errDrivers := services.GetDriverStandingsService("2025", "", "", "", "", "")
	if statusDrivers != http.StatusOK || errDrivers != nil {
		helpers.RedirectToError(w, r, statusDrivers, "Erreur lors de la recherche")
		return
	}

	// Étape 4 : Récupérer TOUTES les données des écuries.
	teamsData, statusTeams, errTeams := services.GetConstructorStandingsService("2025")
	if statusTeams != http.StatusOK || errTeams != nil {
		helpers.RedirectToError(w, r, statusTeams, "Erreur lors de la recherche")
		return
	}

	// Étape 5 : Filtrer les résultats selon la query de recherche.
	filteredDrivers, filteredTeams := services.SearchService(query, driversData.Data["allDrivers"], teamsData.Data["constructors"])

	// Étape 6 : Préparer les données pour le template.
	data := &models.PageData{
		Title:       "Search Results",
		CurrentPage: "search",
		Data: map[string]interface{}{
			"query":        query,
			"drivers":      filteredDrivers,
			"constructors": filteredTeams,
		},
	}

	// Étape 7 : Rendre le template "search" avec les résultats.
	templates.RenderTemplate(w, r, "search", data)
}

// TeamsHandler
// ------------
// Objectif :
//   - Afficher la liste de toutes les écuries F1 pour une saison donnée.
//   - Récupérer la saison depuis l'URL (par défaut 2025).
//   - En cas de succès : rendre le template "teams" avec les données.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func TeamsHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est GET.
	if r.Method != http.MethodGet {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Récupérer la saison depuis l'URL (par défaut 2025).
	season := r.URL.Query().Get("season")
	if season == "" {
		season = "2025"
	}

	// Étape 3 : Appeler services.GetConstructorStandingsService.
	data, status, err := services.GetConstructorStandingsService(season)

	// Étape 4 : Vérifier le statut et l'erreur.
	// Si erreur → helpers.RedirectToError(...) + fmt.Println(err) + return.
	if status != http.StatusOK || err != nil {
		fmt.Println("Erreur lors de la récupération des écuries:", err)
		helpers.RedirectToError(w, r, status, "Impossible de récupérer les écuries")
		return
	}

	// Étape 5 : Si succès, rendre le template "teams" avec les données.
	templates.RenderTemplate(w, r, "teams", data)
}

// IndexHandler
// ------------
// Objectif :
//   - Afficher la page d'accueil avec un aperçu des pilotes et écuries.
//   - Vérifier que l'URL est exactement "/".
//   - Récupérer les données des pilotes et écuries pour la saison 2025.
//   - En cas de succès : rendre le template "index" avec les données.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que l'URL est exactement "/".
	if r.URL.Path != "/" {
		helpers.RedirectToError(w, r, http.StatusNotFound, "Page non trouvée")
		return
	}

	// Étape 2 : Vérifier que la méthode HTTP est GET.
	if r.Method != http.MethodGet {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 3 : Définir la saison actuelle.
	season := "2025"

	// Étape 4 : Récupérer les données des pilotes.
	driversData, statusDrivers, errDrivers := services.GetDriverStandingsService(season, "", "", "", "", "")
	if statusDrivers != http.StatusOK || errDrivers != nil {
		helpers.RedirectToError(w, r, statusDrivers, "Impossible de charger la page d'accueil")
		return
	}

	// Étape 5 : Récupérer les données des écuries.
	teamsData, statusTeams, errTeams := services.GetConstructorStandingsService(season)
	if statusTeams != http.StatusOK || errTeams != nil {
		helpers.RedirectToError(w, r, statusTeams, "Impossible de charger la page d'accueil")
		return
	}

	// Étape 6 : Préparer les données pour le template.
	pageData := models.PageData{
		Title: fmt.Sprintf("Saison %s - Accueil", season),
		Data: map[string]interface{}{
			"season":       season,
			"drivers":      driversData.Data["drivers"],
			"constructors": teamsData.Data["constructors"],
		},
	}

	// Étape 7 : Rendre le template "index" avec les données.
	templates.RenderTemplate(w, r, "index", pageData)
}

// TeamDetailHandler
// -----------------
// Objectif :
//   - Afficher la page de détails d'une écurie spécifique.
//   - Extraire l'ID de l'écurie depuis l'URL.
//   - Récupérer les informations de l'écurie et de ses pilotes.
//   - Vérifier si l'écurie est dans les favoris.
//   - En cas de succès : rendre le template "teams-detail" avec les données.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func TeamDetailHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est GET.
	if r.Method != http.MethodGet {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Extraire l'ID de l'écurie depuis l'URL.
	constructorID := r.URL.Path[len("/teams/"):]
	if constructorID == "" {
		helpers.RedirectToError(w, r, http.StatusNotFound, "Écurie non trouvée")
		return
	}

	// Étape 3 : Récupérer toutes les écuries.
	teamsData, status, err := services.GetConstructorStandingsService("2025")
	if status != http.StatusOK || err != nil {
		helpers.RedirectToError(w, r, status, "Impossible de récupérer les écuries")
		return
	}

	// Étape 4 : Convertir les données et rechercher l'écurie demandée.
	constructors, ok := teamsData.Data["constructors"].([]models.Constructor)
	if !ok {
		helpers.RedirectToError(w, r, http.StatusInternalServerError, "Erreur de données")
		return
	}

	var team *models.Constructor
	for _, c := range constructors {
		if c.ConstructorID == constructorID {
			team = &c
			break
		}
	}

	if team == nil {
		helpers.RedirectToError(w, r, http.StatusNotFound, "Écurie non trouvée")
		return
	}

	// Étape 5 : Récupérer tous les pilotes.
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

	// Étape 6 : Filtrer les pilotes de cette équipe.
	var teamDrivers []models.Driver
	for _, d := range allDrivers {
		if d.Team == team.Name ||
			(d.Team == "Haas" && team.Name == "Haas F1 Team") ||
			(d.Team == "Red Bull" && team.Name == "Red Bull Racing") {
			teamDrivers = append(teamDrivers, d)
		}
	}

	// Étape 7 : Vérifier si l'écurie est dans les favoris.
	isFavorite := services.IsConstructorFavorite(constructorID)

	// Étape 8 : Préparer les données pour le template.
	data := map[string]interface{}{
		"Team":       team,
		"Drivers":    teamDrivers,
		"isFavorite": isFavorite,
	}

	// Étape 9 : Rendre le template "teams-detail" avec les données.
	templates.RenderTemplate(w, r, "teams-detail", data)
}

// DriverDetailHandler
// -------------------
// Objectif :
//   - Afficher la page de détails d'un pilote spécifique.
//   - Extraire l'ID du pilote depuis l'URL.
//   - Récupérer les informations du pilote et de son équipe.
//   - Vérifier si le pilote est dans les favoris.
//   - En cas de succès : rendre le template "drivers-detail" avec les données.
//   - En cas d'erreur : rediriger vers une page d'erreur.
func DriverDetailHandler(w http.ResponseWriter, r *http.Request) {

	// Étape 1 : Vérifier que la méthode HTTP est GET.
	if r.Method != http.MethodGet {
		helpers.RedirectToError(w, r, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Étape 2 : Extraire l'ID du pilote depuis l'URL.
	driverPathPrefix := "/drivers/"
	if len(r.URL.Path) <= len(driverPathPrefix) {
		helpers.RedirectToError(w, r, http.StatusBadRequest, "ID du pilote manquant")
		return
	}
	driverID := r.URL.Path[len(driverPathPrefix):]

	// Étape 3 : Récupérer tous les pilotes.
	driversData, statusDrivers, errDrivers := services.GetDriverStandingsService("2025", "", "", "", "", "")
	if statusDrivers != http.StatusOK || errDrivers != nil {
		helpers.RedirectToError(w, r, statusDrivers, "Impossible de récupérer les pilotes")
		return
	}

	// Étape 4 : Convertir les données et rechercher le pilote demandé.
	allDrivers, ok := driversData.Data["allDrivers"].([]models.Driver)
	if !ok {
		helpers.RedirectToError(w, r, http.StatusInternalServerError, "Erreur de données des pilotes")
		return
	}

	var driver *models.Driver
	for _, d := range allDrivers {
		if d.DriverID == driverID {
			driver = &d
			break
		}
	}

	if driver == nil {
		helpers.RedirectToError(w, r, http.StatusNotFound, "Pilote non trouvé")
		return
	}

	// Étape 5 : Récupérer toutes les écuries.
	teamsData, statusTeams, errTeams := services.GetConstructorStandingsService("2025")
	if statusTeams != http.StatusOK || errTeams != nil {
		helpers.RedirectToError(w, r, statusTeams, "Impossible de récupérer les écuries")
		return
	}

	constructors, ok := teamsData.Data["constructors"].([]models.Constructor)
	if !ok {
		helpers.RedirectToError(w, r, http.StatusInternalServerError, "Erreur de données des écuries")
		return
	}

	// Étape 6 : Trouver l'équipe du pilote.
	var team *models.Constructor
	for _, c := range constructors {
		if c.Name == driver.Team ||
			(driver.Team == "Haas" && c.Name == "Haas F1 Team") ||
			(driver.Team == "Red Bull" && c.Name == "Red Bull Racing") {
			team = &c
			break
		}
	}

	// Étape 7 : Vérifier si le pilote est dans les favoris.
	isFavorite := services.IsDriverFavorite(driverID)

	// Étape 8 : Préparer les données pour le template.
	data := map[string]interface{}{
		"Driver":     driver,
		"Team":       team,
		"season":     "2025",
		"isFavorite": isFavorite,
	}

	// Étape 9 : Rendre le template "drivers-detail" avec les données.
	templates.RenderTemplate(w, r, "drivers-detail", data)
}
