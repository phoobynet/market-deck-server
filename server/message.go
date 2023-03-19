package server

type Message struct {
	Event EmittableEvents `json:"event"`
	Data  interface{}     `json:"data"`
}
