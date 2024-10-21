package velocitydb

import (
	"strings"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":   ping,
	"SET":    set,
	"GET":    get,
	"HSET":   hset,
	"HGET":   hget,
	"CONFIG": config,
	"INFO": info,
}

// PING Command
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

// SET Command
var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	// acquire writers lock and write then Unlock
	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

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
	value, ok := SETs[key]
	SETsMu.RUnlock()

	// check for null value
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
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

	return Value{typ: "bulk", bulk: value}
}

// CONFIG command: Minimal implementation for redis benchmark
func config(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'config' command"}
	}

	if len(args) > 1 && strings.ToUpper(args[0].bulk) == "GET" {
		key := strings.ToLower(args[1].bulk)

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

	// config command is not recognized
	return Value{typ: "error", str: "ERR unsupported CONFIG command"}
}

// INFO command: Minimal implementation
func info(args []Value) Value {
	infoResponse := `# Server
redis_version: velocitydb-0.1
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
	return Value{ typ: "bulk", bulk: infoResponse };
}
