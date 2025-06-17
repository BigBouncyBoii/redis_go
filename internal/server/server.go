package server

import (
	"fmt"
	"net"
)

type Server struct {
	listenAddr string;
	listen  net.Listener
	quitch chan struct{}
	msgch chan Message
}

type Message struct {
	Conn net.Conn
	From string
	Payload []byte
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch: make(chan struct{}),
		msgch: make(chan Message, 100), 
	}
}

func (s *Server) Start() error{
	listen, err := net.Listen("tcp", s.listenAddr)
	fmt.Printf("Starting server on %s\n", s.listenAddr)
	if err != nil{
		return err
	}
	defer listen.Close()
	s.listen = listen

	go s.acceptLoop()
	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *Server) acceptLoop(){
	for{
		conn, err := s.listen.Accept()
		if err != nil {
			fmt.Println("Accept error", err)
			continue
		}
		fmt.Println("New connection from", conn.RemoteAddr())
		go s.readLoop(conn)
	}	
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 1024)
	for{
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Read error", err)
			continue
		}
		msg := buff[:n]
		s.msgch <- Message{
			Conn: conn,
			From: conn.RemoteAddr().String(), 
			Payload: msg,
		}
	}
}

func (s* Server) GetMsgChannel() chan Message {
	return s.msgch
}

func (s *Server) Stop() {
	close(s.quitch)
}

func (s *Server) WriteResponse(conn net.Conn, response []byte) error {
	_, err := conn.Write(response)
	return err
}
