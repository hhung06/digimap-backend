// @title           Digimap Backend API
// @version         1.0
// @description     Digimap backend REST API
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"github.com/hhung06/digimap-backend/cmd"
)

func main() {
	cmd.Execute()
}
