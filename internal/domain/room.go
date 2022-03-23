package domain

type Room struct {
	UUID    string    `json:"uuid"`
	Game    Game      `json:"game"`
	Players []*Player `json:"players"`
}
