package data

type Payload struct {
	Limit int `json:"limit"`
	Stt   int `json:"stt"`
}

var OutStanding = map[string]string{
	"1": "newProducts",
	"2": "bestSales",
	"3": "popular",
	"4": "highestRate",
}

var OrderBy = map[string]string{
	"latest":       "created_at",
	"best-sales":   "sales",
	"popular":      "popular",
	"highest-rate": "rate",
}
