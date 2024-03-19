package utils

import (
	"fmt"
	"giogo/assets"
	"math/rand"
)

func GenerateRandomName() string {
	firstnameIndex := rand.Int31n(int32(len(assets.NameConfig.Firstnames)))
	lastnameIndex := rand.Int31n(int32(len(assets.NameConfig.Lastnames)))

	return fmt.Sprintf("%s %s", assets.NameConfig.Firstnames[firstnameIndex], assets.NameConfig.Lastnames[lastnameIndex])
}
