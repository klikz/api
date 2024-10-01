package models

type Request struct {
	ID        int    `json:"id"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token,omitempty"`
	Name      string `json:"name"`
	ModelID   int    `json:"model_id"`
	ProductID string `json:"product_id"`
	File64    string `json:"file64"`
	Date1     string `json:"date1"`
	Date2     string `json:"date2"`
	C_Time    bool   `json:"c_time"`
	Status_id int    `json:"status_id"`
	Serial    string `json:"serial"`
	PrinterID int    `json:"printerid"`
	Retry     bool   `json:"retry"`
}
