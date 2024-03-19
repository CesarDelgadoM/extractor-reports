package models

// Model restaurant data
type Restaurant struct {
	Name        string `json:"name"`
	Founder     string `json:"founder"`
	Location    string `json:"location"`
	Country     string `json:"country"`
	Fundation   string `json:"fundation"`
	Headquarter string `json:"headquarter"`
}
