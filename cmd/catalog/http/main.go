package main

import (
	"github.com/open-integration/oi/catalog/services/http"
	"github.com/open-integration/oi/pkg/utils"
)

func main() {
	http.Run(utils.GetEnvOrDefault("PORT", "8080"))
}
