package server

type Database struct {
	String map[string]string
	List   map[string][]string
	Set    map[string]map[string]struct{}
	Hash   map[string]map[string]string
}

type Command struct {
	Action string
}

type PingCommand struct {
	Command
}

type StringCommand struct {
	Command
	Key   string
	Value string
}

type ListCommand struct {
	Command
	Key   string
	Values []string
	Start int 
	End   int
}

func NewDatabase() *Database {
	return &Database{
		String: make(map[string]string),
		List:   make(map[string][]string),
		Set:    make(map[string]map[string]struct{}),
		Hash:   make(map[string]map[string]string),
	}
}
