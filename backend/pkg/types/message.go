package types

type Message struct {
	P1     string `json:"left_player_name"`
	P2     string `json:"right_player_name"`
	Result Result `json:"result"`
}
