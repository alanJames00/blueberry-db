// entry point of the database server
package main

import (
	"fmt"
	"velocitydb/internal/config"
)

func main() {
	// load the configuration
	cfg := config.LoadConfig();

	fmt.Println("port", cfg.ServerPort);
}
