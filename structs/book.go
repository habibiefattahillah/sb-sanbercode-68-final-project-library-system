package structs

type Book struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Stock      int    `json:"stock"`
	CategoryID string `json:"category_id"`
}
