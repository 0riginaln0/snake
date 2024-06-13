package difficulty

type Level byte

const (
	Easy Level = iota
	StillEz
	Medium
	Hard
)

func (d Level) String() string {
	var difficulty string
	switch d {
	case Easy:
		difficulty = "easy"
	case StillEz:
		difficulty = "still ez"
	case Medium:
		difficulty = "medium"
	case Hard:
		difficulty = "hard"
	}
	return difficulty
}

// Sets difficulty level to "Easy"
func (d *Level) Reset() {
	*d = Easy
}
