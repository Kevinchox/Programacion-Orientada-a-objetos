package products

// type Product struct {
// 	ID          string    `json:"id"`
// 	Name        string    `json:"name"`
// 	Description string    `json:"description"`
// 	Price       float64   `json:"price"`
// 	Stock       int       `json:"stock"`
// 	Category    string    `json:"category"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// }

type ProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}
