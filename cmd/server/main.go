// entry point of the database server
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"velocitydb/internal/config"
	"velocitydb/internal/logger"
	"velocitydb/pkg/velocitydb"
)

func main() {
	// load the configuration
	cfg := config.LoadConfig()

	// setup logging
	logger.InitLogger(cfg.LogLevel)

	// setup aof
	aof, err := velocitydb.NewAof(cfg.AofFilePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading aof file: %s", err.Error()))
		return
	}
	defer aof.Close()

	// reload previous commands from aof file
	logger.Info(fmt.Sprintf("restoring previous database state from: %s", cfg.AofFilePath))

	aof.Read(func(value velocitydb.Value) {
		command := strings.ToUpper(value.GetArray()[0].GetBulk())
		args := value.GetArray()[1:]

		handler, ok := velocitydb.Handlers[command]
		if !ok {
			logger.Debug(fmt.Sprintf("Invalid command: %s", command))
			return
		}

		handler(args)
	})

	logger.Info(fmt.Sprintf("previous database state restored successfully"))

	// listen on the port
	ln, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		logger.Error("error starting server. err: " + err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	logger.Info("velocitydb started on port" + cfg.ServerPort)

	// accept and handle connections in loop

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Error("error accepting clients. err: " + err.Error())
			continue
		}

		go handleConnection(conn, aof);

	}

}

// goroutine to handle individual connection
func handleConnection(conn net.Conn, aof *velocitydb.Aof) {
	defer conn.Close()

	for {
		resp := velocitydb.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			logger.Error(fmt.Sprintf("error reading command: %s", err.Error()))
			return
		}

		if value.GetType() != "array" {
			logger.Error("Invalid Request, expected array")
			continue
		}

		if len(value.GetArray()) == 0 {
			logger.Error("Invalid Request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.GetArray()[0].GetBulk())
		args := value.GetArray()[1:]

		// debug command

		writer := velocitydb.NewWriter(conn)
	
		// for QUIT command, gracefully close the connection with client: early closing
		if command == "QUIT" {
			// debug
			logger.Debug("command executed: QUIT")

			// Send OK response before closing the connection
			writer.Write(*velocitydb.NewValue("string", "OK", 0, "", nil))	
			conn.Close();
			return;
		}

		handler, ok := velocitydb.Handlers[command]
		if !ok {
			logger.Error("Invalid Command: " + command)
			writer.Write(*velocitydb.NewValue("string", "", 0, "", nil))
			continue
		}
	
		// write to aof for set and hset commands
		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
