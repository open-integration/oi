package main

import (
	"github.com/open-integration/oi/catalog/services/trello"
	"github.com/open-integration/oi/pkg/utils"
)

func main() {
	trello.Run(utils.GetEnvOrDefault("PORT", "8080"))
}
