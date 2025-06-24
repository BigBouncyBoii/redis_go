package server

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func HandleCommand(args []string, db *Database) (string, error) {
	var command string = args[0]
	if len(args) == 0 {
		return "", errors.New("no command provided")
	}
	fmt.Printf("Handling command: %s with args: %v\n", command, args)
	switch command {
	case "ping":
		return handlePingCommand(args)
	case "command":
		return "command", nil
	case "set", "get", "del", "exists":
		return handleStringCommand(args, db)
	case "lpush", "rpush", "lpop", "rpop", "lrange":
		return handleListCommand(args, db)
	case "sadd", "srem", "smembers", "sismember":
		return handleSetCommand(args, db)
	case "hset", "hget", "hdel":
		return handleHashCommand(args, db)
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
	var strcmd StringCommand = StringCommand{
		Command: Command{Action: args[0]},
		Key:    "",
		Value:  "",
	}
	if len(args) < 2 {
		return "", errors.New("not enough arguments for string command")
	}
	strcmd.Key = args[1]
	if len(args) > 2 {
		strcmd.Value = args[2]
	}
	switch strcmd.Action {
	case "set":
		if strcmd.Value == "" {
			return "", errors.New("SET command requires a key and a value")
		}
		db.String[strcmd.Key] = strcmd.Value
		return "ok", nil
	case "get":
		value, exists := db.String[strcmd.Key]
		if !exists {
			return "-1", nil 
		}
		return value, nil
	case "del":
		_, exists := db.String[strcmd.Key]
		if !exists {
			return "0", nil
		}
		delete(db.String, strcmd.Key)
		return "1", nil
	case "exists":
		_, exists := db.String[strcmd.Key]
		if exists {
			return "1", nil
		}
		return "0", nil
	default:
		return "", errors.New("unknown string command: " + strcmd.Action)
	}
}


func handleListCommand(args []string, db *Database) (string, error) {
	var listcmd ListCommand = ListCommand{
		Command: Command{Action: args[0]},
		Key:    "",
		Values: []string{},
		Start:  0,
		End:    -1,
	}
	if len(args) < 2 {
		return "", errors.New("not enough arguments for list command")
	}
	listcmd.Key = args[1]
	switch listcmd.Action {
		case "lpush":
			if len(args) < 3 {
				return "", errors.New("LPUSH command requires a key and at least one value")
			}
			listcmd.Values = args[2:]
			_, exists := db.List[listcmd.Key]; 
			if !exists {
				db.List[listcmd.Key] = &list.List{}
			}
			for _, value := range listcmd.Values {
				db.List[listcmd.Key].PushFront(value)
			}
			return strconv.Itoa(db.List[listcmd.Key].Len()), nil
		case "rpush":
			if len(args) < 3 {
				return "", errors.New("RPUSH command requires a key and at least one value")
			}
			listcmd.Values = args[2:]
			_, exists := db.List[listcmd.Key];
			if !exists {
				db.List[listcmd.Key] = &list.List{}
			}
			for _, value := range listcmd.Values {
				db.List[listcmd.Key].PushBack(value)
			}
			return strconv.Itoa(db.List[listcmd.Key].Len()), nil
		case "lpop":
			list, exists := db.List[listcmd.Key]
			if !exists || list.Len() == 0 {
				return "-1", nil // Return -1 if the list does not exist or is empty
			}
			element := list.Front()
			if element == nil {
				return "-1", nil // Return -1 if the list is empty
			}
			value := element.Value.(string) // Safe type assertion
			list.Remove(element)
			return value, nil
		case "rpop":
			list, exists := db.List[listcmd.Key]
			if !exists || list.Len() == 0 {
				return "-1", nil // Return -1 if the list does not exist or is empty
			}
			element := list.Back()
			if element == nil {
				return "-1", nil // Return -1 if the list is empty
			}
			value := element.Value.(string)
			list.Remove(element)
			return value, nil
		case "lrange":
			if len(args) < 4 {
				return "", errors.New("LRANGE command requires a key, start, and end index")
			}
			listcmd.Start, _ = strconv.Atoi(args[2])
			listcmd.End, _ = strconv.Atoi(args[3])
			list, exists := db.List[listcmd.Key]
			if !exists {
				return "", errors.New("list does not exist: " + listcmd.Key)
			}
			if listcmd.Start < 0 || listcmd.Start > list.Len() || listcmd.End > list.Len() || (listcmd.End < listcmd.Start && listcmd.End != -1) {
				return "", errors.New("invalid range for LRANGE command")
			}
			if listcmd.End == -1 {
				listcmd.End = db.List[listcmd.Key].Len()
			}
			var result []string
			var currentIndex int = 0
			for e := list.Front(); e != nil; e = e.Next() {
				if currentIndex >= listcmd.Start && currentIndex <= listcmd.End {
					value:= e.Value.(string)
					result = append(result, value)
				}
				currentIndex++
			}
			return parseMultiple(result), nil
		default:
			return "", errors.New("unknown list command: " + listcmd.Action)
	}
}


func handleSetCommand(args []string, db *Database) (string, error) {
	var setcmd SetCommand = SetCommand{Command: Command{Action: args[0]}, Key: "", Members: []string{}}
	if len(args) < 2 {
		return "", errors.New("not enough arguments for set command")
	}
	setcmd.Key = args[1]
	switch setcmd.Action {	
		case "sadd":
			if len(args) < 3 {
				return "", errors.New("SADD command requires a key and at least one member")
			}
			setcmd.Members = args[2:]
			_, exists := db.Set[setcmd.Key] 
			if !exists {
				db.Set[setcmd.Key] = make(map[string]struct{})
			}
			for _, member := range setcmd.Members {
				db.Set[setcmd.Key][member] = struct{}{}
			}
			return strconv.Itoa(len(setcmd.Members)), nil
		case "srem":
			if len(args) < 3 {
				return "", errors.New("SREM command requires a key and at least one member")
			}
			setcmd.Members = args[2:]
			_, exists := db.Set[setcmd.Key]
			if !exists {
				return "0", nil 
			}
			removedCount := 0
			for _, member := range setcmd.Members {
				if _, exists := db.Set[setcmd.Key][member]; exists {
					delete(db.Set[setcmd.Key], member)
					removedCount++
				}
			}
			return strconv.Itoa(removedCount), nil
		case "smembers":
			_, exists := db.Set[setcmd.Key]
			if !exists {
				return "-1", nil // Return -1 if the set does not exist
			}
			var members []string
			for member := range db.Set[setcmd.Key] {
				members = append(members, member)
			}
			if len(members) == 0 {
				return "-1", nil // Return -1 if the set is empty
			}
			var result []string
			for _, member := range members {
				result = append(result, member)
			}
			return parseMultiple(result), nil
		case "sismember":
			if len(args) < 3 {
				return "", errors.New("SISMEMBER command requires a key and a member")
			}
			member := args[2]
			_, exists := db.Set[setcmd.Key]
			if !exists {
				return "0", nil // Return 0 if the set does not exist
			}
			if _, exists := db.Set[setcmd.Key][member]; exists {
				return "1", nil // Return 1 if the member exists in the set
			}
			return "0", nil // Return 0 if the member does not exist in the set
		default:
			return "", errors.New("unknown set command: " + setcmd.Action)	
	}
}

func handleHashCommand(args []string, db *Database) (string, error) {
	var hashcmd HashCommand = HashCommand{Command: Command{Action: args[0]}, Key: "", Fields: []string{}, Values: []string{}}
	if len(args) < 2 {
		return "", errors.New("not enough arguments for hash command")
	}
	hashcmd.Key = args[1]
	switch hashcmd.Action {
	case "hset":
		if len(args) < 4 {
			return "", errors.New("HSET command requires a key, field, and value")
		}
		if len(args) % 2 != 0 {
			return "", errors.New("HSET command requires an even number of field-value pairs")
		}
		var fields int
		for i := 2; i < len(args)-1; i += 2 {
			field := args[i]
			value := args[i+1]
			_, exists := db.Hash[hashcmd.Key]
			if !exists {
				db.Hash[hashcmd.Key] = make(map[string]string)
				fields++ 
			}
			db.Hash[hashcmd.Key][field] = value
		}
		return strconv.Itoa(fields), nil
	case "hget":
		if len(args) < 3 {
			return "", errors.New("HGET command requires a key and a field")
		}
		field := args[2]
		_, exists := db.Hash[hashcmd.Key]
		if !exists {
			return "-1", nil // Return -1 if the hash does not exist
		}
		value, exists := db.Hash[hashcmd.Key][field]
		if !exists {
			return "-1", nil // Return -1 if the field does not exist in the hash
		}
		return value, nil
	case "hdel":
		if len(args) < 3 {
			return "", errors.New("HDEL command requires a key and at least one field")
		}
		fields := args[2:]
		_, exists := db.Hash[hashcmd.Key]
		if !exists {
			return "0", nil // Return 0 if the hash does not exist
		}
		removedCount := 0
		for _, field := range fields {
			if _, exists := db.Hash[hashcmd.Key][field]; exists {
				delete(db.Hash[hashcmd.Key], field)
				removedCount++
			}
		}
		return strconv.Itoa(removedCount), nil
	default:
		return "", errors.New("unknown hash command: " + hashcmd.Action)
	}
}

func parseMultiple(args []string) string {
	if len(args) == 0 {
		return "-1"
	}
	var sb strings.Builder
	sb.WriteString("*")
	sb.WriteString(strconv.Itoa(len(args)))
	sb.WriteString("\r\n")
	for _, arg := range args {
		sb.WriteString("$")
		sb.WriteString(strconv.Itoa(len(arg)))
		sb.WriteString("\r\n")
		sb.WriteString(arg)
		sb.WriteString("\r\n")
	}
	return sb.String()
}