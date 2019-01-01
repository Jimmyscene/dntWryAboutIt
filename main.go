package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
	"github.com/kr/pty"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ptyHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %s", err)
		return
	}
	defer conn.Close()
	c := exec.Command("python3", "pystuff/main.py")
	cPty, err := pty.Start(c)
	if err != nil {
		fmt.Printf("Error starting websocket: %s", err)
		return
	}
	defer cPty.Close()
	var done bool
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := cPty.Read(buf)
			if err != nil {
				if err == io.EOF {
					conn.WriteMessage(websocket.TextMessage, []byte("\nScript Exited. Terminating Connection\n"))
					break
				} else {
					conn.WriteMessage(websocket.TextMessage, []byte("Failed to read buffer: "+err.Error()))
					fmt.Printf("Failed to read buffer: %s", err)
					break
				}
			}
			err = conn.WriteMessage(websocket.TextMessage, buf[0:n])
			if err != nil {
				fmt.Printf("Failed writing to ws: %s", err)
				continue
			}
		}
		done = true
	}()

	for !done {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message: ", err)
			// Just keep going boys
			break
		}
		_, err = cPty.Write(message)

		if err != nil {
			log.Println("Failed to write to webscoket, error:", err)
			continue
		}
	}

}

func test() (err error) {
	// Create arbitrary command.
	c := exec.Command("python3", "main.py")
	c.Start()
	err = c.Wait()
	return
}

func main() {
	http.HandleFunc("/pty", ptyHandler)
	fmt.Println("Serving On localhost:9000")
	fmt.Println(http.ListenAndServe("localhost:9000", nil))
}
