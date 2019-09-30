package Model

type BankInformation struct {
	UID           int64  `json:"id"`
	AccountNumber int64  `json:"accountNUmber"`
	Amount        int64  `json:"amount"`
	Date          string `json:"date"`
}
