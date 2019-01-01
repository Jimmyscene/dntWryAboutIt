package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
)

type FileData struct {
	Data string `json:"data"`
}

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

func getDebug(r *http.Request) bool {
	var debug bool
	if mux.Vars(r)["debug"] == "true" {
		debug = true
	} else {
		debug = false
	}
	return debug
}

func getFile(debug bool) string {
	if debug {
		return "pystuff/main_debug.py"
	} else {
		return "pystuff/main.py"
	}
}

func getContents(filepath string) []byte {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return dat
}

func ptyHandler(w http.ResponseWriter, r *http.Request) {
	debug := getDebug(r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %s", err)
		return
	}
	defer conn.Close()
	// You really want to set PYTHONUNBUFFERED=1, otherwise you'll lose 8 hours
	c := exec.Command("python", getFile(debug))
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
			fmt.Printf("Error starting websocket: %s", starterr)
			return
		}
		cPty = NoOpWriter{stdout}
	}
	defer cPty.Close()
	if debug {
		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("Failed to read message: ", err)
					fmt.Println("Terminating Process...")
					err := c.Process.Kill()
					if err != nil {
						fmt.Println("Error terminating process: ", err)
					}
					break
				}
				_, err = cPty.Write(message)
				if err != nil {
					log.Println("Failed to write to webscoket, error:", err)
					continue
				}
			}
		}()
	}

	buf := make([]byte, 128)
	for {
		n, err := cPty.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("Failed to read buffer: "+err.Error()))
				fmt.Printf("Failed to read buffer: %s", err)
				fmt.Println("Terminating Process...")
				err := c.Process.Kill()
				if err != nil {
					fmt.Println("Error terminating process: ", err)
				}
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
	if err := c.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(
						int(status), "Script Exited. Terminating Connection",
					),
				)
				conn.Close()
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	filedata := getContents(getFile(getDebug(r)))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	file := FileData{string(filedata)}
	str, err := json.Marshal(&file)
	if err != nil {
		panic(err)
	}
	w.Write(str)
}

func main() {
	var router = mux.NewRouter()
	router.Path("/pty").Queries("debug", "{debug}").HandlerFunc(ptyHandler)
	router.Path("/file").Queries("debug", "{debug}").HandlerFunc(getFileHandler)
	// http.HandleFunc("/pty", ptyHandler)
	fmt.Println("Serving On localhost:9000")
	fmt.Println(http.ListenAndServe("0.0.0.0:9000", router))
}
