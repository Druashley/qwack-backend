package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

type Audio struct {
	File string `json:"file"`
}

type User struct {
	Name string `json:"name"`
}

type Message struct {
	User  User  `json:"user"`
	Audio Audio `json:"audio"`
}

var upgrader = websocket.Upgrader{
	//todo: this allows connections from any origin. Prod it should be restricted
	CheckOrigin: func(r *http.Request) bool { return true },
}

func playAudio(file string) error {
	// Full path to ffplay (if necessary)
	// cmd := exec.Command("/usr/bin/ffplay", "-nodisp", "-autoexit", "audio/"+file)
	fmt.Println("Executing:", "ffplay -nodisp -autoexit audio/"+file)
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "audio/"+file)

	// Capture the output and error streams
	cmdOutput, err := cmd.CombinedOutput() // Combines stdout and stderr
	if err != nil {
		log.Printf("Error playing audio: %v\n", err)
		log.Printf("ffplay output: %s\n", cmdOutput)
		return err
	}

	log.Printf("ffplay output: %s\n", cmdOutput)
	return nil
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection established!")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		// log.Println("Received message:", string(message))

		var msg Message

		err = json.Unmarshal(message, &msg)

		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		soundName := msg.Audio.File

		fmt.Println("Playing:", soundName)

		err = playAudio(soundName + ".ogg")
		if err != nil {
			log.Println("Error playing sound:", err)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
