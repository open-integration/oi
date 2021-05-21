package main

import (
	"github.com/open-integration/oi/catalog/services/slack"
	"github.com/open-integration/oi/pkg/utils"
)

func main() {
	slack.Run(utils.GetEnvOrDefault("PORT", "8080"))
}
