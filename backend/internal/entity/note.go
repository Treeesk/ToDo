package entity

type Note struct {
	ID      int    `json:"id"`
	User_id int    `json:"user_id"`
	Text    string `json:"text"`
}
