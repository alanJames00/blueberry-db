package velocitydb

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
}

// PING Command
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}
