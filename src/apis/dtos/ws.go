package dtos

type WsPayload struct {
	Event     string   `json:"event"`
	UserId    int64    `json:"user_id"`
	Page      string   `json:"page"`
	Positions Position `json:"positions"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
