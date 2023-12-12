package main

import (
	"fmt"
	"log"
	"net"
)

const BufferLength = 1024

type Server struct {
	listenAddress string
	listener      net.Listener
	couch         chan struct{}
}

func New(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
		couch:         make(chan struct{}),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	defer listener.Close()

	s.listener = listener

	go s.acceptLoop()

	<-s.couch
	close(s.couch)

	return nil
}

func (s *Server) acceptLoop() error {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)

			break
		}

		fmt.Println("new connection to the server: ", connection.RemoteAddr())

		go s.readLoop(connection)
	}

	return nil
}

func (s *Server) readLoop(connection net.Conn) error {
	buffer := make([]byte, BufferLength)

	_, err := connection.Read(buffer)
	if err != nil {
		log.Println("read error: ", err)
	}

	defer connection.Close()

	connection.Write([]byte("HTTP/1.1 200 OK\n" + "\n" + "aboba"))

	return nil
}

func main() {
	server := New(":1337")

	log.Fatal(server.Start())
}
