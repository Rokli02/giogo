package assets

import (
	"fmt"
)

const useEmbededAssets = true

func InitializeAssets() {
	fmt.Println("Initializing Assets")

	initializeConfigs(useEmbededAssets)
	initializeImgs(useEmbededAssets)
}
