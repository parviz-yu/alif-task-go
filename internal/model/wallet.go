package model

// Money is a monetary unit that equals 1⁄100 of the basic monetary unit — diram.
type Money int64

// Wallet ...
type Wallet struct {
	IsIdentified bool   `json:"is_identified"`
	Balance      Money  `json:"balance"`
	UserID       string `json:"-"`
}

type User struct {
	ID   int
	UUID string
}
