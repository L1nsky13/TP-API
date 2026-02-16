package models

// PageData
// Structure pour passer les données aux templates HTML.
type PageData struct {
	Title       string
	CurrentPage string
	Data        map[string]interface{}
	Error       string
}

// Driver
// Structure représentant un pilote F1 avec ses données personnelles et sa écurie.
type Driver struct {
	DriverID        string `json:"driverId"`
	PermanentNumber string `json:"permanentNumber"`
	Code            string `json:"code"`
	Image           string `json:"image"`
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	DateOfBirth     string `json:"dateOfBirth"`
	Nationality     string `json:"nationality"`
	Team            string `json:"team"`
	DriverType      string `json:"driverType"`
}

// Constructor
// Structure représentant une écurie F1 avec ses informations visuelles et son identifiant.
type Constructor struct {
	ConstructorID string `json:"constructorId"`
	Icon          string `json:"icon"`
	Image         string `json:"image"`
	Name          string `json:"name"`
	Nationality   string `json:"nationality"`
	TeamColor     string `json:"teamColor"`
}

// MRData
// Structure contenant les métadonnées de la réponse API F1.
type MRData struct {
	Xmlns  string `json:"xmlns"`
	Series string `json:"series"`
	URL    string `json:"url"`
	Limit  string `json:"limit"`
	Offset string `json:"offset"`
	Total  string `json:"total"`
}

// F1Response
// Structure enveloppe pour les réponses API F1 contenant les métadonnées.
type F1Response struct {
	MRData MRData `json:"MRData"`
}

// Favorites
// Structure pour stocker les pilotes et écuries favoris de l'utilisateur.
type Favorites struct {
	Drivers      []string `json:"drivers"`
	Constructors []string `json:"constructors"`
}
