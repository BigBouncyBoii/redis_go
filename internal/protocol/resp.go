package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

func ParseCommand(input []byte) ([]string, error){
	if len(input) == 0 || input[0] != '*'{
		return nil, errors.New("invalid RESP format: missing array marker")
	}
	reader := bufio.NewReader(bytes.NewReader(input))
	args, err := readBulkString(reader)
	if err != nil {
		return nil, errors.New("failed to read bulk string: " + err.Error())
	}
	if len(args) == 0 {
		return nil, errors.New("no command found")
	}
	return args, nil
}

func readBulkString(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, errors.New("failed to read command line: " + err.Error())
	}
	count, err := strconv.Atoi(line[1:len(line)-2])
	if err != nil{
		return nil, errors.New("invalid command count: " + err.Error())
	}
	args := make([]string, count)
	for i := 0; i < count; i++ {
		bulkHeader, err := reader.ReadString('\n')
		if err != nil {
			return nil, errors.New("failed to read bulk string header for argument " + strconv.Itoa(i+1) + ": " + err.Error())
		}
		if len(bulkHeader) < 4 || bulkHeader[0] != '$' {
			return nil, errors.New("invalid bulk string header for argument " + strconv.Itoa(i+1) + ": " + bulkHeader)
		}
		length, err := strconv.Atoi(bulkHeader[1:len(bulkHeader)-2])
        if err != nil {
            return nil, errors.New("invalid bulk string length for argument " + strconv.Itoa(i+1) + ": " + err.Error())
        }
		bulkData := make([]byte, length)
		_, err = io.ReadFull(reader, bulkData)
		if err != nil {
			return nil, errors.New("failed to read bulk string data for argument " + strconv.Itoa(i+1) + ": " + err.Error())
		}
		args[i] = strings.ToLower(string(bulkData))
		reader.Discard(2) // Discard the trailing \r\n
	}
	return args, nil
}

func isInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func ParseRespone(response string) []byte {
	if response == "pong"{
		return []byte("+PONG\r\n")
	} else if response == "ok" {
		return []byte("+OK\r\n")
	} else if response == "-1" {	
		return []byte("$-1\r\n")
	} else if isInteger(response) {
		return []byte(":" + response + "\r\n")
	} else if strings.HasPrefix(response, "*") {
		return []byte(response)
	} else {
		return []byte("$" + strconv.Itoa(len(response)) + "\r\n" + response + "\r\n")
	}
}