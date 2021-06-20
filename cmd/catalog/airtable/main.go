package main

import (
	"github.com/open-integration/oi/catalog/services/airtable"
	"github.com/open-integration/oi/pkg/utils"
)

func main() {
	airtable.Run(utils.GetEnvOrDefault("PORT", "8080"))
}
