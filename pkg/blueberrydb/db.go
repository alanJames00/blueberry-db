package blueberrydb

import (
	"blueberrydb/internal/logger"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Handlers = map[string]func([]Value) Value{
	"PING":   ping,
	"SET":    set,
	"DEL":    deleteKeys,
	"GET":    get,
	"HSET":   hset,
	"HGET":   hget,
	"CONFIG": config,
	"INFO":   info,
	"EXPIRE": expire,
}

// PING Command
func ping(args []Value) Value {
	if len(args) == 0 {
		// debug
		logger.Debug("command recieved: PING")

		return Value{typ: "string", str: "PONG"}

	}

	// debug
	logger.Debug(fmt.Sprintf("command executed: PING %s", args[0].bulk))

	return Value{typ: "string", str: args[0].bulk}
}

// SET Command

type SetValStruct struct {
	value     string
	expiresAt int64 // 0 -> no expiration
}

var SETs = map[string]SetValStruct{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	// acquire writers lock and write then Unlock
	SETsMu.Lock()
	SETs[key] = SetValStruct{
		value:     value,
		expiresAt: 0, // no expiration by default
	}
	SETsMu.Unlock()

	// debug
	logger.Debug(fmt.Sprintf("command executed: SET %s %s", args[0].bulk, args[1].bulk))

	return Value{typ: "string", str: "OK"}
}

// GET Command
func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	// acquire readers' lock, read and then unlock
	SETsMu.RLock()
	entry, ok := SETs[key]
	SETsMu.RUnlock()

	// check for null value
	if !ok {
		return Value{typ: "null"}
	}

	// check if key is expired
	if entry.expiresAt > 0 && time.Now().Unix() > entry.expiresAt {
		// delete the expired key and return null
		SETsMu.Lock()
		delete(SETs, key)
		SETsMu.Unlock()

		return Value{typ: "null"}
	}

	// debug
	logger.Debug(fmt.Sprintf("command executed: GET %s", key))

	return Value{typ: "bulk", bulk: entry.value}
}

// DELETE command: delete multiple keys and values
func deleteKeys(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'del' command"}
	}

	key := args[0].bulk

	// acquire lock, delete the key value pair and unlock
	SETsMu.Lock()
	delete(SETs, key)
	SETsMu.Unlock()

	return Value{typ: "string", str: "1"}
}

// EXPIRE command: set an expire on a key
func expire(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for  'expire' command"}
	}

	key := args[0].bulk
	secondString := args[1].bulk

	seconds, err := strconv.Atoi(secondString)
	if err != nil || seconds <= 0 {
		return Value{typ: "error", str: "Invalid seconds argument"}
	}

	// get the current time in unix timestamp and add the seconds to get the expiresAt
	currentTime := time.Now().Unix()
	expiresAt := currentTime + int64(seconds)

	// acquire write lock and set the expire on key
	SETsMu.Lock()
	entry, ok := SETs[key]
	if ok {
		entry.expiresAt = expiresAt
		SETs[key] = entry
	}
	SETsMu.Unlock()

	if !ok {
		return Value{typ: "null"} // key does not exists
	}

	logger.Debug(fmt.Sprintf("command executed: EXPIRE %s %s", key, secondString))

	return Value{typ: "string", str: "1"}
}

// HSET command
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	// acquire writer's lock, set and unlock
	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	// debug
	logger.Debug(fmt.Sprintf("command executed: HSET %s %s %s", hash, key, value))

	return Value{typ: "string", str: "OK"}
}

// HGET command
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	// acquire readers' lock, read and unlock
	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	// null check
	if !ok {
		return Value{typ: "null"}
	}

	// debug
	logger.Debug(fmt.Sprintf("command executed: HGET %s %s", hash, key))

	return Value{typ: "bulk", bulk: value}
}

// CONFIG command: Minimal implementation for redis benchmark
func config(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'config' command"}
	}

	if len(args) > 1 && strings.ToUpper(args[0].bulk) == "GET" {
		key := strings.ToLower(args[1].bulk)

		// debug
		logger.Debug(fmt.Sprintf("command executed: CONFIG GET %s", key))

		// simulate basic config responses
		switch key {
		case "maxmemory":
			return Value{typ: "array", str: "", bulk: "", array: []Value{
				{typ: "bulk", bulk: "maxmemory"},
				{typ: "bulk", bulk: "0"}, // TODO: placeholder value
			}}
		case "timeout":
			return Value{typ: "array", str: "", bulk: "", array: []Value{
				{typ: "bulk", bulk: "timeout"},
				{typ: "bulk", bulk: "0"}, // Placeholder value
			}}
		case "save":
			return Value{typ: "array", str: "", bulk: "", array: []Value{
				{typ: "bulk", bulk: "save"},
				{typ: "bulk", bulk: "3600 1 300 100 60 10000"},
			}}
		default:
			// Return empty array for unrecognized config keys
			return Value{typ: "array", array: []Value{}}
		}
	}

	// debug
	logger.Error("commmand errored: unsupported CONFIG command")

	// config command is not recognized
	return Value{typ: "error", str: "ERR unsupported CONFIG command"}
}

// INFO command: Minimal implementation
func info(args []Value) Value {
	infoResponse := `# Server
redis_version: blueberrydb-0.1
uptime_in_seconds: 12345
uptime_in_days: 0
# Clients
connected_clients: 1
# Memory
used_memory: 2048
# Persistence
rdb_last_save_time: 0
# Stats
total_connections_received: 1
total_commands_processed: 1
# CPU
used_cpu_sys: 0.00
used_cpu_user: 0.00
# Keyspace
db0:keys=1,expires=0,avg_ttl=0
`

	// debug
	logger.Debug("commmand executed: INFO")

	return Value{typ: "bulk", bulk: infoResponse}
}
