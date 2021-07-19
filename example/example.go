package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"time"
)

type ConnectMessage struct {
	RoomId   string `json:"roomId"`
	UserName string `json:"userName"`
	Apikey   string `json:"apikey"`
	Action   string `json:"action"`
}

type RecievedMessage struct {
	Type      string `json:"type"`
	From      User   `json:"from"`
	Timestamp int    `json:"timestamp"`
	RoomId    string `json:"roomid"`
	Message   string `json:"message"`
	Msgid     int    `json:"msgid"`
}

type SendMessage struct {
	Message   string `json:"message"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}

type User struct {
	Name string `json:"name"`
}

type Settings struct {
	StreamerName string
	BotName      string
	Apikey       string
	Commands     map[string]string
}

var SETTINGS = Settings{}

const SETTINGS_FILE = "settings.json"
const ACTION_CONNECT = "connect"
const ACTION_MESSAGE = "message"

func main() {
	fmt.Println("starting up")

	b, err := ioutil.ReadFile(SETTINGS_FILE)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &SETTINGS)
	if err != nil {
		panic(err)
	}

	messagesToSend := make(chan SendMessage)
	messagesRecieved := make(chan RecievedMessage)

	origin := "https://pomf.tv"
	url := "wss://pomf.tv/websocket/"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}

	firstMsg := ConnectMessage{SETTINGS.StreamerName, SETTINGS.BotName, SETTINGS.Apikey, ACTION_CONNECT}
	a, _ := json.Marshal(firstMsg)
	fmt.Println("sending:")
	fmt.Println(string(a))
	ws.Write(append(a, '\n'))

	go sendMessages(ws, messagesToSend)
	go recieveMessages(ws, messagesRecieved)

	go handleMessages(messagesToSend, messagesRecieved)

	<-make(chan int) // wait forever
}

func sendMessages(ws *websocket.Conn, messagesToSend <-chan SendMessage) {
	for msg := range messagesToSend {
		a, _ := json.Marshal(msg)
		fmt.Println("sending:")
		fmt.Println(string(a))
		ws.Write(a)
	}
}

func recieveMessages(ws *websocket.Conn, messagesRecieved chan<- RecievedMessage) {
	for {
		frame := make([]byte, 1024)
		n := 0
		n, err := ws.Read(frame)
		if err != nil {
			panic(err)
		}
		if n >= 1023 {
			// message was too big for our buffer, just ignore it
			fmt.Println("Skipping too large message")
			continue
		}
		msg := RecievedMessage{}
		err = json.Unmarshal(frame[:n], &msg)
		if err != nil {
			fmt.Println("Unable to unmarshal message:")
			fmt.Println(err)
			fmt.Println(frame[:n])
			fmt.Println(string(frame[:n]))
			continue
		}
		fmt.Println("recieved:")
		fmt.Println(msg)
		messagesRecieved <- msg
	}
}

func handleMessages(messagesToSend chan<- SendMessage, messagesRecieved <-chan RecievedMessage) {
	for msg := range messagesRecieved {
		if msg.RoomId == SETTINGS.StreamerName && msg.From.Name != SETTINGS.BotName && msg.Type == ACTION_MESSAGE {
			if cmd, exists := SETTINGS.Commands[msg.Message]; exists {
				messagesToSend <- SendMessage{cmd, ACTION_MESSAGE, time.Now().Unix()}
			}
		}
	}
}
