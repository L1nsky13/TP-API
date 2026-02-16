package services

import (
	"encoding/json"
	"f1-app/models"
	"fmt"
	"os"
	"path/filepath"
)

const favoritesFileName = "favorites.json"

// GetFavoritesFilePath
// Retourne le chemin absolu vers le fichier des favoris.
func GetFavoritesFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		return favoritesFileName
	}

	if filepath.Base(wd) == "cmd" {
		wd = filepath.Join(wd, "..", "..")
	} else if filepath.Base(wd) == "src" {
		wd = filepath.Join(wd, "..")
	}

	return filepath.Join(wd, favoritesFileName)
}

// LoadFavorites
// Charge les favoris depuis le fichier JSON. Crée le fichier s'il n'existe pas encore.
func LoadFavorites() (*models.Favorites, error) {
	filePath := GetFavoritesFilePath()

	// Vérifier si le fichier existe, sinon le créer.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		favorites := &models.Favorites{
			Drivers:      []string{},
			Constructors: []string{},
		}
		// Créer le fichier avec la structure vide.
		if err := SaveFavorites(favorites); err != nil {
			return nil, err
		}
		return favorites, nil
	}

	// Lire le fichier des favoris.
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture fichier favoris: %w", err)
	}

	var favorites models.Favorites
	if err := json.Unmarshal(data, &favorites); err != nil {
		return nil, fmt.Errorf("erreur décodage JSON favoris: %w", err)
	}

	// Initialiser les slices s'ils sont nil pour éviter les modifications nulles.
	if favorites.Drivers == nil {
		favorites.Drivers = []string{}
	}
	if favorites.Constructors == nil {
		favorites.Constructors = []string{}
	}

	return &favorites, nil
}

// SaveFavorites
// Sauvegarde les favoris dans le fichier JSON.
func SaveFavorites(favorites *models.Favorites) error {
	filePath := GetFavoritesFilePath()

	data, err := json.MarshalIndent(favorites, "", "    ")
	if err != nil {
		return fmt.Errorf("erreur encodage JSON favoris: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("erreur écriture fichier favoris: %w", err)
	}

	return nil
}

// AddDriverToFavorites
// Ajoute un pilote aux favoris s'il n'y est pas déjà.
func AddDriverToFavorites(driverID string) error {
	favorites, err := LoadFavorites()
	if err != nil {
		return err
	}

	// Vérifier si le pilote n'est pas déjà dans les favoris.
	for _, id := range favorites.Drivers {
		if id == driverID {
			return nil
		}
	}

	favorites.Drivers = append(favorites.Drivers, driverID)
	return SaveFavorites(favorites)
}

// RemoveDriverFromFavorites
// Supprime un pilote des favoris.
func RemoveDriverFromFavorites(driverID string) error {
	favorites, err := LoadFavorites()
	if err != nil {
		return err
	}

	// Reconstruire la liste sans le pilote à supprimer.
	newDrivers := []string{}
	for _, id := range favorites.Drivers {
		if id != driverID {
			newDrivers = append(newDrivers, id)
		}
	}

	favorites.Drivers = newDrivers
	return SaveFavorites(favorites)
}

// AddConstructorToFavorites
// Ajoute une écurie aux favoris si elle n'y est pas déjà.
func AddConstructorToFavorites(constructorID string) error {
	favorites, err := LoadFavorites()
	if err != nil {
		return err
	}

	// Vérifier si l'écurie n'est pas déjà dans les favoris.
	for _, id := range favorites.Constructors {
		if id == constructorID {
			return nil
		}
	}

	favorites.Constructors = append(favorites.Constructors, constructorID)
	return SaveFavorites(favorites)
}

// RemoveConstructorFromFavorites
// -----------
// Objectif :
//   - Supprimer une écurie spécifique de la liste des favoris.
//   - Sauvegarder la liste mise à jour.
func RemoveConstructorFromFavorites(constructorID string) error {
	// Étape 1 : Charger les favoris actuels depuis le fichier.
	favorites, err := LoadFavorites()
	if err != nil {
		return err
	}

	// Étape 2 : Reconstruire la liste sans l'écurie à supprimer.
	newConstructors := []string{}
	for _, id := range favorites.Constructors {
		if id != constructorID {
			newConstructors = append(newConstructors, id)
		}
	}

	// Étape 3 : Mettre à jour la liste et sauvegarder.
	favorites.Constructors = newConstructors
	return SaveFavorites(favorites)
}

// IsDriverFavorite
// -----------
// Objectif :
//   - Vérifier si un pilote spécifique est dans les favoris.
//   - Retourner un booléen indiquant la présence du pilote.
func IsDriverFavorite(driverID string) bool {
	// Étape 1 : Charger les favoris depuis le fichier.
	favorites, err := LoadFavorites()
	if err != nil {
		return false
	}

	// Étape 2 : Parcourir la liste des pilotes favoris.
	for _, id := range favorites.Drivers {
		// Étape 3 : Retourner vrai si le pilote est trouvé.
		if id == driverID {
			return true
		}
	}
	// Étape 4 : Retourner faux si le pilote n'existe pas.
	return false
}

// IsConstructorFavorite
// -----------
// Objectif :
//   - Vérifier si une écurie spécifique est dans les favoris.
//   - Retourner un booléen indiquant la présence de l'écurie.
func IsConstructorFavorite(constructorID string) bool {
	// Étape 1 : Charger les favoris depuis le fichier.
	favorites, err := LoadFavorites()
	if err != nil {
		return false
	}

	// Étape 2 : Parcourir la liste des écuries favorites.
	for _, id := range favorites.Constructors {
		// Étape 3 : Retourner vrai si l'écurie est trouvée.
		if id == constructorID {
			return true
		}
	}
	// Étape 4 : Retourner faux si l'écurie n'existe pas.
	return false
}
