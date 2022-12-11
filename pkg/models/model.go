package models

type PeerStatus string

const (
	Online  PeerStatus = "online"
	Offline PeerStatus = "offline"
)

type Peer struct {
	TrackerURL string     `json:"tracker_url"`
	ShareKey  string     `json:"share_key"`
	ChatUrl    string     `json:"chat_url"` // IRC
	IP         string     `json:"ip"`
	Port       string     `json:"port"`
	States     PeerStatus `json:"states"`
}
