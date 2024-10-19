package velocitydb

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET": set,
	"GET": get,
}

// PING Command
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

// SET Command
var SETs = map[string]string{};
var SETsMu = sync.RWMutex{};

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command" }
	}

	key := args[0].bulk;
	value := args[1].bulk;

	// acquire writers lock and write then Unlock
	SETsMu.Lock();
	SETs[key] = value;
	SETsMu.Unlock();

	return Value{typ: "string", str: "OK"};
}

// GET Command
func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"};
	}

	key := args[0].bulk;

	// acquire readers' lock, read and then unlock
	SETsMu.RLock();
	value, ok := SETs[key];
	SETsMu.RUnlock();

	// check for null value
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value};
}

// HSET command
var HSETs = map[string]map[string]string{};
var HSETsMu = sync.RWMutex{};

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"};
	}

	hash := args[0].bulk;
	key := args[1].bulk;
	value := args[2].bulk;

	// acquire writer's lock, set and unlock
	HSETsMu.Lock();
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{};
	}
	HSETs[hash][key] = value;
	HSETsMu.Unlock();

	return Value{typ: "string", str: "OK"};
}

// HGET command 
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{ typ: "error", str: "ERR wrong number of arguments for 'hget' command"};
	}

	hash := args[0].bulk;
	key := args[1].bulk;

	// acquire readers' lock, read and unlock
	HSETsMu.RLock();
	value, ok := HSETs[hash][key];
	HSETsMu.RUnlock();

	// null check
	if !ok {
		return Value{typ: "null"};
	}

	return Value{ typ: "bulk", str: value };
}
