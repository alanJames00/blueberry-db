// entry point of the database server
package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"velocitydb/internal/config"
	"velocitydb/pkg/velocitydb"
)

func main() {
	// load the configuration
	cfg := config.LoadConfig();

	fmt.Println("server listening on port", cfg.ServerPort);

	// setup aof
	aof, err := velocitydb.NewAof(cfg.AofFilePath);
	if err != nil {
		fmt.Println(err);
		return;
	}
	defer aof.Close();

	// reload previous commands from aof file
	aof.Read(func(value velocitydb.Value) {
		command := strings.ToUpper(value.GetArray()[0].GetBulk());
		args := value.GetArray()[1:];

		handler, ok := velocitydb.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args);
	})
	
	
	// listen on the port
	ln, err := net.Listen("tcp", cfg.ServerPort);
	if err != nil {
		log.Fatalf("error starting server: %v", err);
	}

	// accept and handle connections
	conn, err := ln.Accept();
	if err != nil {
		fmt.Printf("error accepting connection: %v\n", err);
	}
	defer conn.Close();

	// handler loop
	for {

		resp := velocitydb.NewResp(conn);
		value, err := resp.Read();
		if err != nil {
			fmt.Println(err);
			return;
		}
		
		if value.GetType() != "array" {
			fmt.Printf("Invalid Request, expected array\n");
			continue;
		} 
		
		if len(value.GetArray()) == 0 {
			fmt.Printf("Invalid Request, expected array length > 0\n");
			continue;
		}

		command := strings.ToUpper(value.GetArray()[0].GetBulk());
		args := value.GetArray()[1:];

		// debug
		fmt.Println("command", command, args);

		writer := velocitydb.NewWriter(conn);

		handler, ok := velocitydb.Handlers[command];
		if !ok {
			fmt.Printf("Invalid Command: %v\n", command);
			writer.Write(*velocitydb.NewValue("string", "", 0, "", nil));
			continue;
		}

		// write to aof for set and hset commands
		if command == "SET" || command == "HSET" {
			aof.Write(value);
		}

		result := handler(args);
		writer.Write(result);
	}

}
