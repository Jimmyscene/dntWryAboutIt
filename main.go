package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
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

type NoOpWriter struct {
	io.ReadCloser
}

func (c NoOpWriter) Write([]byte) (int, error) {
	return 0, nil
}
func ptyHandler(w http.ResponseWriter, r *http.Request) {
	var debug bool
	if mux.Vars(r)["debug"] == "true" {
		debug = true
	} else {
		debug = false
	}
	fmt.Println("DEBUG IS: ", debug)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %s", err)
		return
	}
	defer conn.Close()
	// You really want to set PYTHONUNBUFFERED=1, otherwise you'll lose 8 hours
	c := exec.Command("python3", "pystuff/main.py")
	var cPty io.ReadWriteCloser
	if debug {
		cPty, err = pty.Start(c)
		if err != nil {
			fmt.Printf("Error starting websocket: %s", err)
			return
		}
	} else {
		stdout, err := c.StdoutPipe()
		if err != nil {
			fmt.Printf("Error starting websocket: %s", err)
			return
		}
		starterr := c.Start()
		if starterr != nil {
			fmt.Printf("Error starting websocket: %s", err)
			return
		}
		cPty = NoOpWriter{stdout}
	}
	// defer cPty.Close()
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
			// For some reason, the StdoutPipe doesn't have \r\n, but only \n, which breaks the xterm render
			// TODO: Should probably check to see if byte already has the \r before the \n, but adding it doesn't appear to break anything
			data := bytes.Replace(buf[0:n], []byte{10}, []byte{13, 10}, -1)
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Printf("Failed writing to ws: %s", err)
				continue
			}
		}
		done = true
	}()
	if debug {
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
	} else {
		c.Wait()
	}

}

func main() {
	var router = mux.NewRouter()
	router.Path("/pty").Queries("debug", "{debug}").HandlerFunc(ptyHandler)
	// http.HandleFunc("/pty", ptyHandler)
	fmt.Println("Serving On localhost:9000")
	fmt.Println(http.ListenAndServe("localhost:9000", router))
}
