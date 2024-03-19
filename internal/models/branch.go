package models

// Model branch
type Branch struct {
	Name      string    `json:"name"`
	Manager   string    `json:"manager"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Score     uint8     `json:"score"`
	Employees employees `json:"employees"`
	Financial financial `json:"financial"`
	Menu      menu      `json:"menu"`
}

type employees struct {
	Admins       []employee `json:"admins"`
	Waiters      []employee `json:"waiters"`
	Chefs        []employee `json:"chefs"`
	TotalAdmins  uint8      `json:"total_admins"`
	TotalWaiters uint8      `json:"tatal_waiters"`
	TotalChefs   uint8      `json:"total_chefs"`
}

type employee struct {
	Name  string `json:"name"`
	Years uint8  `years:"name"`
	Sales uint8  `json:"sales"`
}

type financial struct {
	Sales    uint32 `json:"sales"`
	Expenses uint32 `json:"expenses"`
}

type menu struct {
	EntreePlates []string `json:"entree_plates"`
	MainCourse   []string `json:"main_course"`
	Drinks       []string `json:"drinks"`
	Desserts     []string `json:"desserts"`
}
