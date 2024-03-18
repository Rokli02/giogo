package cli

import (
	"fmt"
	"os"
)

const (
	userNameArgFlag = "-u"
	widthArgFlag    = "-w"
	heightArgFlag   = "-h"
	minesArgFlag    = "-m"
)

var longFlagCastMap = map[string]string{
	"--username": userNameArgFlag,
	"--width":    widthArgFlag,
	"--height":   heightArgFlag,
	"--mines":    minesArgFlag,
}

func ProcessAppArgs() map[string]string {
	options := make(map[string]string)
	args := os.Args[1:]
	foundFlag := ""

	for i := 0; i < len(args); i++ {
		if value, hasKey := longFlagCastMap[args[i]]; hasKey {
			args[i] = value
		}

		switch args[i] {
		case userNameArgFlag:
			fallthrough
		case widthArgFlag:
			fallthrough
		case heightArgFlag:
			fallthrough
		case minesArgFlag:
			foundFlag = args[i]
		default:
			if args[i][0] == "-"[0] {
				foundFlag = ""
			}

			if foundFlag == "" {
				break
			}

			value, hasKey := options[foundFlag]
			if !hasKey {
				options[foundFlag] = args[i]

				break
			}

			options[foundFlag] = fmt.Sprintf("%s %s", value, args[i])
		}
	}

	return options
}
