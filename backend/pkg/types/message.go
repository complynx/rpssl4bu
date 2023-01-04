package types

type Message struct {
	LeftPlayerName    string `json:"left_player_name"`
	RightPlayerName   string `json:"right_player_name"`
	LeftPlayerChoice  Choice `json:"left_player_choice"`
	RightPlayerChoice Choice `json:"right_player_choice"`
	Result            Result `json:"result"`
}
