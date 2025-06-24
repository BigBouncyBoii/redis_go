package server

import (
	"container/list"
)

type Database struct {
	String map[string]string
	List   map[string]*list.List
	Set    map[string]map[string]struct{}
	Hash   map[string]map[string]string
	ZSet   map[string]map[string]float64
}


type Command struct {
	Action string
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

type SetCommand struct {
	Command
	Key    string
	Members []string
}

type HashCommand struct {
	Command
	Key    string
	Fields []string
	Values []string
}

func NewDatabase() *Database {
	return &Database{
		String: make(map[string]string),
		List:   make(map[string]*list.List),
		Set:    make(map[string]map[string]struct{}),
		Hash:   make(map[string]map[string]string),
		ZSet:   make(map[string]map[string]float64),
	}
}
