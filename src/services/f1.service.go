package services

import (
	"f1-app/models"
	"fmt"
	"net/http"
	"strings"
)

// getDriversData
// Récupère la liste complète des pilotes depuis les données intégrées.
func getDriversData() []models.Driver {
	return models.DriversData
}

// getConstructorsData
// Récupère la liste complète des écuries depuis les données intégrées.
func getConstructorsData() []models.Constructor {
	return models.ConstructorsData
}

// GetDriverStandingsService
// -------------------------
// Objectif :
//   - Récupérer les pilotes avec filtrage (équipe, nationalité, type) et pagination.
//   - Retourner les données paginées avec les options de filtrage disponibles.
func GetDriverStandingsService(season, teamFilter, nationalityFilter, driverTypeFilter, pageParam, perPageParam string) (*models.PageData, int, error) {

	// Étape 1 : Récupérer tous les pilotes depuis les données intégrées.
	allDrivers := getDriversData()

	// Étape 2 : Appliquer les filtres de recherche.
	filteredDrivers := []models.Driver{}
	for _, driver := range allDrivers {

		if teamFilter != "" && driver.Team != teamFilter {
			continue
		}

		if nationalityFilter != "" && driver.Nationality != nationalityFilter {
			continue
		}

		if driverTypeFilter != "" && driver.DriverType != driverTypeFilter {
			continue
		}
		filteredDrivers = append(filteredDrivers, driver)
	}

	// Étape 3 : Récupérer les valeurs disponibles pour les filtres (équipes, nationalités, types).
	teams := getUniqueTeams(allDrivers)
	nationalities := getUniqueNationalities(allDrivers)
	driverTypes := getUniqueDriverTypes(allDrivers)

	// Étape 4 : Traiter les paramètres de pagination.
	page := 1
	perPage := 10

	if pageParam != "" {
		if n, err := fmt.Sscanf(pageParam, "%d", &page); err != nil || n != 1 || page < 1 {
			page = 1
		}
	}

	if perPageParam != "" {
		if n, err := fmt.Sscanf(perPageParam, "%d", &perPage); err != nil || n != 1 {
			perPage = 10
		}
		if perPage != 10 && perPage != 20 && perPage != 30 {
			perPage = 10
		}
	}

	// Étape 5 : Calculer les indices de pagination.
	totalDrivers := len(filteredDrivers)
	totalPages := (totalDrivers + perPage - 1) / perPage

	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	start := (page - 1) * perPage
	end := start + perPage

	if start > totalDrivers {
		start = totalDrivers
	}
	if end > totalDrivers {
		end = totalDrivers
	}

	paginatedDrivers := filteredDrivers[start:end]

	startIndex := start + 1
	endIndex := end
	if totalDrivers == 0 {
		startIndex = 0
		endIndex = 0
	}

	// Étape 6 : Préparer les données pour le template.
	pageData := &models.PageData{
		Title:       fmt.Sprintf("Pilotes %s", season),
		CurrentPage: "drivers",
		Data: map[string]interface{}{
			"season":            season,
			"drivers":           paginatedDrivers,
			"allDrivers":        allDrivers,
			"teams":             teams,
			"nationalities":     nationalities,
			"driverTypes":       driverTypes,
			"currentPage":       page,
			"perPage":           perPage,
			"totalPages":        totalPages,
			"totalDrivers":      totalDrivers,
			"startIndex":        startIndex,
			"endIndex":          endIndex,
			"teamFilter":        teamFilter,
			"nationalityFilter": nationalityFilter,
			"driverTypeFilter":  driverTypeFilter,
		},
	}

	// Étape 7 : Retourner les données avec le statut HTTP OK.
	return pageData, http.StatusOK, nil
}

// getUniqueTeams
// Extrait la liste unique des équipes à partir de la liste des pilotes.
func getUniqueTeams(drivers []models.Driver) []string {
	teamMap := make(map[string]bool)
	for _, driver := range drivers {
		if driver.Team != "" {
			teamMap[driver.Team] = true
		}
	}
	teams := []string{}
	for team := range teamMap {
		teams = append(teams, team)
	}
	return teams
}

// getUniqueNationalities
// Extrait la liste unique des nationalités à partir de la liste des pilotes.
func getUniqueNationalities(drivers []models.Driver) []string {
	nationalityMap := make(map[string]bool)
	for _, driver := range drivers {
		if driver.Nationality != "" {
			nationalityMap[driver.Nationality] = true
		}
	}
	nationalities := []string{}
	for nationality := range nationalityMap {
		nationalities = append(nationalities, nationality)
	}
	return nationalities
}

// getUniqueDriverTypes
// Extrait la liste unique des types de pilotes à partir de la liste des pilotes.
func getUniqueDriverTypes(drivers []models.Driver) []string {
	typeMap := make(map[string]bool)
	for _, driver := range drivers {
		if driver.DriverType != "" {
			typeMap[driver.DriverType] = true
		}
	}
	types := []string{}
	for driverType := range typeMap {
		types = append(types, driverType)
	}
	return types
}

// GetConstructorStandingsService
// ------------------------------
// Objectif :
//   - Récupérer la liste complète des écuries pour une saison.
//   - Retourner les données formatées pour le template.
func GetConstructorStandingsService(season string) (*models.PageData, int, error) {

	// Étape 1 : Récupérer tous les écuries depuis les données intégrées.
	constructors := getConstructorsData()

	// Étape 2 : Préparer les données pour le template.
	pageData := &models.PageData{
		Title:       fmt.Sprintf("Écuries %s", season),
		CurrentPage: "teams",
		Data: map[string]interface{}{
			"season":       season,
			"constructors": constructors,
		},
	}

	// Étape 3 : Retourner les données avec le statut HTTP OK.
	return pageData, http.StatusOK, nil
}

// SearchService
// ------
// Objectif :
//   - Filtrer les pilotes et écuries selon une requête de recherche.
//   - Recherche insensible à la casse et flexible.
//   - Retourner les pilotes et écuries correspondants.
func SearchService(query string, driversInterface interface{}, constructorsInterface interface{}) ([]models.Driver, []models.Constructor) {
	// Étape 1 : Normaliser la requête de recherche.
	query = strings.ToLower(strings.TrimSpace(query))

	var filteredDrivers []models.Driver
	var filteredConstructors []models.Constructor

	// Étape 2 : Vérifier que la requête n'est pas vide.
	if query == "" {
		return filteredDrivers, filteredConstructors
	}

	// Étape 3 : Convertir les interfaces en listes typées.
	drivers, driversOk := driversInterface.([]models.Driver)
	constructors, constructorsOk := constructorsInterface.([]models.Constructor)

	// Étape 4 : Rechercher les écuries correspondantes.
	var matchedTeamNames []string
	if constructorsOk {
		for _, constructor := range constructors {
			if strings.Contains(strings.ToLower(constructor.Name), query) ||
				strings.Contains(strings.ToLower(constructor.Nationality), query) {
				filteredConstructors = append(filteredConstructors, constructor)
				matchedTeamNames = append(matchedTeamNames, constructor.Name)
			}
		}
	}

	// Étape 5 : Rechercher les pilotes correspondants.
	if driversOk {
		for _, driver := range drivers {

			if strings.Contains(strings.ToLower(driver.GivenName), query) ||
				strings.Contains(strings.ToLower(driver.FamilyName), query) ||
				strings.Contains(strings.ToLower(driver.Code), query) ||
				strings.Contains(driver.PermanentNumber, query) ||
				strings.Contains(strings.ToLower(driver.Nationality), query) {
				filteredDrivers = append(filteredDrivers, driver)
				continue
			}

			for _, teamName := range matchedTeamNames {
				if driver.Team == teamName {
					filteredDrivers = append(filteredDrivers, driver)
					break
				}
			}
		}
	}

	// Étape 6 : Retourner les résultats de la recherche.
	return filteredDrivers, filteredConstructors
}
