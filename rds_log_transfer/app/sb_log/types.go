package sb_log

type JsonCharacterClass struct {
	ClassType int `json:"class_type"`
	Level     int `json:"level"`
	Exp       int `json:"exp"`
}

type JsonRequestHeaderLocation struct {
	SessionId string `json:"session_id"`
	MapId     string `json:"map_id"`
	PosX      string `json:"pos_x"`
	PosY      string `json:"pos_y"`
	PosZ      string `json:"pos_z"`
}
