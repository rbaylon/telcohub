package models

type Category struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"` // e.g., hex or Tailwind class
}
