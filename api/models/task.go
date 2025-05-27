package models

type Dialogue struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Trans string `json:"trans"`
}

type Task struct {
	ID        int        `json:"id"`
	Creator   string     `json:"creator"`
	Filename  string     `json:"filename"`
	Dialogues []Dialogue `json:"dialogues"`
	Status    string     `json:"status"` // e.g. open, closed
	Changes   []Change   `json:"changes"`
}

type Change struct {
	ID         int    `json:"id"`
	TaskID     int    `json:"task_id"`
	DialogueID int    `json:"dialogue_id"`
	User       string `json:"user"`
	NewTrans   string `json:"new_trans"`
	Status     string `json:"status"` // pending, approved, rejected
}
