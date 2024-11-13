// implements basic password auth and manages auth state of each client connection
package blueberrydb 

import (
	"net"
	"blueberrydb/internal/logger"
)

// map to keep track of client auth state
var clientAuthState = map[net.Conn]bool{}

// AUTH command
func Auth(args []Value, conn net.Conn, password string) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'AUTH' command"}
	}

	// compare password
	providedPassword := args[0].bulk
	if password == providedPassword {
		// set the auth state
		clientAuthState[conn] = true
		logger.Info("client authentication successfully")
		return Value{typ: "string", str: "OK"}
	}

	logger.Info("client authentication failed")
	return Value{typ: "error", str: "ERR invalid password"}
}

// check auth state of a connection
func CheckAuth(conn net.Conn) bool {
	state, exists := clientAuthState[conn]
	if !exists {
		return false
	}

	return state
}
