package main

import (
	"fmt"
	"net/http"
	"lang/board"
	"strconv"
	"log"
//	"encoding/json"
)

// create a board room and return the room number with specified roles
// JSON checkIn(roles); // returns: "开房成功！房间号为 123456789，可以邀请你的朋友到 vrcats.com/ww/123456789 开始游戏" + ws url
// a websocket conne
func checkIn(w http.ResponseWriter, req *http.Request) {
	rolesJson := req.FormValue("roles")
	//var roles 
	//json.Unmarshal(rolesJson, roles)

	fmt.Fprintf(w, rolesJson)
	//把输入值变成数组，创建一个局，把每一个元素加入其中
	board := new(board.Board)
	boardId = boardId + 1
	//TODO：转换输入数据为json对象，从中取出roles数组
	//roles := parseRoles(roleJson)
	//TODO: 从输入里获取本局配置
	meta := map[string]string {"女巫自救": "不能"}
	roles := []string {"预言家","女巫","猎人","白痴","村民","村民","村民","村民","狼人","狼人","狼人","狼人",}
	board.New(boardId, roles, meta)
	message := fmt.Sprintf("开房成功！ 房间号为%n，可以邀请你的好友到 abc.com/%n 开始游戏\n", board.Id, board.Id)
	fmt.Fprintf(w, message)
}

// sit down in a room numbered board, preffered seat number number, with nick name nick
 // returns: "为您选择了3号座位，目前12人中已有7人入座。您的身份是：<b>狼人<b>"
func sitDown(w http.ResponseWriter, req *http.Request) {
	b, n, k := req.FormValue("board"), req.FormValue("number"), req.FormValue("nick")
	boardId, _ := strconv.Atoi(b)
	board := boards[boardId]
	s, _ := strconv.Atoi(n)
	message := board.ViewCard(s, k)
	fmt.Fprintf(w, message)
}

func operate(w http.ResponseWriter, req *http.Request) {
	n, number, action := req.FormValue("board"), req.FormValue("number"), req.FormValue("action")
	num1, num2, num3, skill, card := req.FormValue("num1"), req.FormValue("num2"), req.FormValue("num3"), req.FormValue("skill"), req.FormValue("card")

	boardId, _ = strconv.Atoi(n)
	board := boards[boardId]
	fmt.Fprintf(w, board.TakeAction(number, action, num1, num2, num3, skill, card))
	// Send instruction to the board's host
	instruction := board.SM[board.State][1]
	msg := &Message{BoardId: board.Id, Body: instruction}
	hub.host <- msg
}

var boards map[int]*board.Board
var boardId int
var hub *Hub

func main() {
	fmt.Println("Creating board pool...")
	boards = map[int]*board.Board {}
	boardId = 1000000
	hub = newHub()
	go hub.run()
	log.Println("Websocket hub started...")

	// web service api entries
	http.HandleFunc("/checkIn", checkIn)
	http.HandleFunc("/sitDown", sitDown)
	http.HandleFunc("/operate", operate)
	// web socket service entry
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Println("Handlers intialized, starting server at port 80...")

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
