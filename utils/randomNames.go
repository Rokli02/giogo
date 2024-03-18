package utils

import (
	"fmt"
	"math/rand"
)

var firstnames = []string{
	"Rigó",
	"Erős",
	"Vajonőr",
	"Cikáyn",
	"Tompoló",
	"Mikorka",
	"Buzi",
	"Galambos",
	"Techno",
	"Dundi",
	"Kokkeró",
	"ZéTéNy",
}

var lastnames = []string{
	"Kálmán",
	"Jancsi",
	"Galambfos",
	"Pepe",
	"Gregory",
	"Wise Mystical Tree",
	"Lacatusu",
}

func GenerateRandomName() string {
	firstnameIndex := rand.Int31n(int32(len(firstnames)))
	lastnameIndex := rand.Int31n(int32(len(lastnames)))

	return fmt.Sprintf("%s %s", firstnames[firstnameIndex], lastnames[lastnameIndex])
}
