package main

import "sync"

func ping(args []Value) Value {
  if len(args) == 0 {
    return Value{typ: "string", str: "PONG"}
  }
  return Value{typ: "string", str: args[0].bulk}
}

var Handlers = map[string]func([]Value) Value {
  "PING": ping,
  "SET": set,
  "GET": get,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set (args []Value) Value {
  if len(args) != 2 {
    return Value{typ: "error", str: "ERR wrong number"}
  }

  key := args[0].bulk
  value := args[1].bulk

  SETsMu.Lock()
  SETs[key] = value
  SETsMu.Unlock()

  return Value{typ: "string", str: "OK"}
}

func get (args []Value) Value {
  if len (args) != 1 {
    return Value {typ: "error", str: "ERR wrong number"}
  }

  key := args[0].bulk

  SETsMu.RLock()
  value, ok := SETs[key]
  SETsMu.RUnlock()

  if !ok {
    return  Value{typ: "error", str: "key not found"}
  }

  return Value{typ: "string", str: value}
}
