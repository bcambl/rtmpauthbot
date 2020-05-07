package models

// OnEvent is used to marshal the POST body from on_ events
// on_publish
// on_publish_done
// on_play
// on_play_done
type OnEvent struct {
	ClientID int    `json:"clientid"`
	Call     string `json:"call"`
	App      string `json:"app"`
	Name     string `json:"name"`
	TcURL    string `json:"tcurl"`
	Addr     string `json:"addr"`
}
