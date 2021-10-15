package main

import (
	gcalendar "github.com/open-integration/oi/catalog/services/google-calendar"
	"github.com/open-integration/oi/pkg/utils"
)

func main() {
	gcalendar.Run(utils.GetEnvOrDefault("PORT", "8080"))
}
