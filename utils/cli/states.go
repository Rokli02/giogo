package cli

import (
	"fmt"
	"giogo/utils"
)

var (
	State    map[string]string = make(map[string]string)
	Username string
	Width    uint16 = 8
	Height   uint16 = 12
	Mines    uint16 = 12
)

func InitializeState() {
	fmt.Println("Initializing State variables")
	State = ProcessAppArgs()

	if value, hasKey := State[userNameArgFlag]; hasKey {
		Username = value
	} else {
		Username = utils.GenerateRandomName()
	}

	if value, hasKey := State[widthArgFlag]; hasKey {
		var width int
		fmt.Sscanf(value, "%d", &width)

		read, err := fmt.Sscanf(value, "%d", &width)
		if read == 1 && err == nil && width > 0 {
			Width = uint16(width)
		}
	}

	if value, hasKey := State[heightArgFlag]; hasKey {
		var height int
		fmt.Sscanf(value, "%d", &height)

		read, err := fmt.Sscanf(value, "%d", &height)
		if read == 1 && err == nil && height > 0 {
			Height = uint16(height)
		}
	}

	if value, hasKey := State[minesArgFlag]; hasKey {
		var mines int
		fmt.Sscanf(value, "%d", &mines)

		read, err := fmt.Sscanf(value, "%d", &mines)
		if read == 1 && err == nil && mines > 0 {
			Mines = uint16(mines)
		}
	}
}
