package server

import (
	"errors"
	"fmt"
)

func HandleCommand(args []string, db *Database) (string, error) {
	var command string = args[0]
	fmt.Printf("\n%s", command)
	if len(args) == 0 {
		return "", errors.New("no command provided")
	}
	switch command {
	case "ping":
		return handlePingCommand(args)
	case "set", "get", "del", "exists":
		return handleStringCommand(args, db)
	case "lpush", "rpush", "lpop", "rpop":
		// Handle list commands
	case "sadd", "srem", "smembers":
		// Handle set commands
	case "hset", "hget", "hdel":
		// Handle hash commands
	case "flushall", "flushdb":
		// Handle flush commands
	case "ttl", "expire":
		// Handle TTL and expire commands
	default:
		return "", errors.New("unknown command: " + command)
	}
	return "", nil
}

func handlePingCommand(args []string) (string, error) {
	length := len(args)
	switch length {
	case 1:
		return "pong", nil
	case 2:
		return args[1], nil
	default:
		return "", errors.New("invalid number of arguments for PING command")
	}
}

func handleStringCommand(args []string, db *Database) (string, error) {
	if len(args) < 2 {
		return "", errors.New("not enough arguments for string command")
	}
	command := args[0]
	key := args[1]
	switch command {
	case "set":
		if len(args) != 3 {
			return "", errors.New("SET command requires exactly 2 arguments: key and value")
		}
		value := args[2]
		db.String[key] = value
		return "ok", nil
	case "get":
		value, exists := db.String[key]
		if !exists {
			return "-1", nil 
		}
		return value, nil
	case "del":
		if len(args) != 2 {
			return "", errors.New("DEL command requires exactly 1 argument: key")
		}
		_, exists := db.String[key]
		if !exists {
			return "0", nil
		}
		delete(db.String, key)
		return "1", nil
	case "exists":
		if len(args) != 2 {
			return "", errors.New("EXISTS command requires exactly 1 argument: key")
		}
		_, exists := db.String[key]
		if exists {
			return "1", nil
		}
		return "0", nil
	default:
		return "", errors.New("unknown string command: " + command)
	}
}