package domain

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Class string `json:"class_name"`
}
