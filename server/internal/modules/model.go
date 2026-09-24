package modules

type Module struct {
	ID           string `json:"id"`
	UniversityID string `json:"university_id"`
	Code         string `json:"code"`
	Title        string `json:"title"`
}
