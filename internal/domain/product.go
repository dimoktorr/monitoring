package domain

import "time"

type Product struct {
	ID           int
	Name         string
	Description  string
	Price        float64
	InStock      bool
	Category     *Category
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Rating       []*Rating
	Reviews      []*Review
	SKU          string
	Weight       float64
	Dimensions   string
	Manufacturer string
	Barcode      string
}

type Category struct {
	ID          int
	Name        string
	Description string
}

type Rating struct {
	ID        int
	Score     float32
	Comment   string
	User      *User
	ProductID int
	CreatedAt time.Time
}

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Review struct {
	ID        int
	Content   string
	Rating    int
	User      *User
	ProductID int
	CreatedAt time.Time
	UpdatedAt time.Time
}
