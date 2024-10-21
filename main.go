package main

import (
	"database/sql" // Used for connecting to a database
	"log"          // Used for logging messages

	_ "github.com/lib/pq"                         // Import driver for postgres (assuming DB driver is postgres)
	"github.com/nibir1/banking_system/api"        // Import banking system API package
	db "github.com/nibir1/banking_system/db/sqlc" // Import banking system sqlc package for database access
	"github.com/nibir1/banking_system/util"       // Import banking system utility package
)

func main() {
	// Load configuration from the current directory (".") using the util package
	config, err := util.LoadConfig(".")
	if err != nil {
		// Log a fatal error message if configuration cannot be loaded
		log.Fatal("cannot load configuration:", err)
	}

	// Open a database connection using the configured driver and source
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		// Log a fatal error message if database connection fails
		log.Fatal("Cannot connect to db:", err)
	}

	// Create a new sqlc store using the established database connection
	store := db.NewStore(conn)

	// Create a new API server instance using the configuration and store
	server, err := api.NewServer(config, store)
	if err != nil {
		// Log a fatal error message if server creation fails
		log.Fatal("cannot create server:", err)
	}

	// Start the server listening on the configured server address
	err = server.Start(config.ServerAddress)
	if err != nil {
		// Log a fatal error message if server fails to start
		log.Fatal("cannot start server", err)
	}
}
