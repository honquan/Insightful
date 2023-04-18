package dtos

type WsPayload struct {
	UserId    int64      `json:"user_id"`
	Positions []Position `json:"positions"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
