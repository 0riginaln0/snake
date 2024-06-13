package state

type State byte

const (
	Start State = iota
	Playing
	GameOver
	Win
)
