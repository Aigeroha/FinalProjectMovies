package models


type Customer struct {
	ID           int    `json:"customer_id"`
	Nickname     string `json:"nickname"`
	PasswordHash string `json:"-"` 
	Phone        string `json:"phone"`
}


type CustomerWallet struct {
	AccountID  int `json:"account_id"`
	CustomerID int `json:"customer_id"`
	Balance    int `json:"balance"` 
}


type RegisterInput struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}


type LoginInput struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}


type CustomerProfileResponse struct {
	ID       int    `json:"customer_id"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Balance  int    `json:"balance"`
}