package main

import (
	"github.com/open-integration/oi/catalog/services/github"
	"github.com/open-integration/oi/pkg/utils"
)

func main() {
	github.Run(utils.GetEnvOrDefault("PORT", "8080"))
}
