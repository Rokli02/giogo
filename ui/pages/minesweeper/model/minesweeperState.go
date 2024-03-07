package model

type MinesweeperState uint8

const (
	UNDEFINED MinesweeperState = iota
	WAITING
	START
	RUNNING
	LOSE
	WIN
	END
	LOADING
)

func (s MinesweeperState) ToString() string {
	switch s {
	case UNDEFINED:
		return "UNDEFINED"
	case START:
		return "START"
	case RUNNING:
		return "RUNNING"
	case LOSE:
		return "LOSE"
	case WIN:
		return "WIN"
	case END:
		return "END"
	case LOADING:
		return "LOADING"
	}

	return "UNKNOWN"
}
