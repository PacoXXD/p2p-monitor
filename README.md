# p2p-monitor

A program that monitor the info  that peers in file share p2p network send back.



	TrackerURL string     `json:"tracker_url"`
	ShareKey  string     `json:"share_key"`
	ChatUrl    string     `json:"chat_url"` // IRC
	IP         string     `json:"ip"`
	Port       string     `json:"port"`
	States     PeerStatus `json:"states"`
  
  
