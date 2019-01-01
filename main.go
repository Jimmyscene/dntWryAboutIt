package main

import (
	"log"
	"net/http"
	"os/exec"
	"io"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
type NoOpWrite struct {
	io.ReadCloser
}
func (n NoOpWrite) Write([]byte) (int,  error) {
	return 0, nil
}

func ptyHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("Starting new pty")
	debug := false
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %s", err)
		return
	}
	defer conn.Close()
	c := exec.Command("python3", "pystuff/main.py")
	var cPty io.ReadWriteCloser
	if debug {
		cPty, err = pty.Start(c)
	} else {
		ccPty, pipeerr := c.StdoutPipe()
		err = c.Start()
		if pipeerr != nil {
			fmt.Printf("Error getting stdoutPipe: %s\n", pipeerr.Error())
		}
		cPty = NoOpWrite{ccPty}
	}
	if err != nil {
		fmt.Printf("Error starting websocket: %s", err)
		return
	}
	defer cPty.Close()
	go func() {
		fmt.Println("In the goloop")
		buf := make([]byte, 128)
		for {
			fmt.Println("Reading stuff")
			n, err := cPty.Read(buf)
			fmt.Println(n)
			fmt.Println(string(buf))
			if err != nil {
				fmt.Println("Error: " + err.Error())
				if err == io.EOF {
					conn.WriteMessage(websocket.TextMessage, []byte("Script Exited. Terminating Connection"))
					break
				} else {
					conn.WriteMessage(websocket.TextMessage, []byte("Failed to read buffer: "+ err.Error()))
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
	}()
	if debug {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Failed to read message: %s", err)
				// Just keep going boys
				continue
			}
			_, err = cPty.Write(message)

			if err != nil {
				log.Println("write:", err)
				continue
			}
		}
	} else {
		for {}
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
	http.ListenAndServe("localhost:9000", nil)
	if err := test(); err != nil {
		log.Fatal(err)
	}
}
