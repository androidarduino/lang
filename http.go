package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lang/board"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BoardConfig struct {
	Roles []string          `json:"roles"`
	Meta  map[string]string `json:"meta"`
}

// create a board room and return the room number with specified roles
// JSON checkIn(roles); // returns: "开房成功！房间号为 123456789，可以邀请你的朋友到 vrcats.com/ww/123456789 开始游戏" + ws url
// a websocket conne
func checkIn(w http.ResponseWriter, req *http.Request) {
	rolesJson := req.FormValue("roles")
	n, nickName := req.FormValue("seat"), req.FormValue("nick")
	seatNumber, _ := strconv.Atoi(n)

	log.Println("input json", rolesJson)

	var config BoardConfig
	json.Unmarshal([]byte(rolesJson), &config)

	log.Println("unmarshalled json", config)

	//把输入值变成数组，创建一个局，把每一个元素加入其中
	board := new(board.Board)
	globalBoardId = globalBoardId + 1
	board.Id = globalBoardId
	boards[board.Id] = board
	//TODO：转换输入数据为json对象，从中取出roles数组
	//roles := parseRoles(roleJson)
	//TODO: 从输入里获取本局配置
	//meta := map[string]string {"女巫自救": "不能"}
	//roles := []string {"预言家","女巫","猎人","白痴","村民","村民","村民","村民","狼人","狼人","狼人","狼人",}
	roles := config.Roles
	meta := config.Meta
	board.New(board.Id, roles, meta)

	sitDown := board.ViewCard(seatNumber, nickName)
	board.Seats[seatNumber].Label("房主")
	//Process sitDown, extract seat number and assign
	responses := strings.Split(sitDown, "\n")
	actualSeat := responses[0]
	actualRole := responses[1]

	content, err := ioutil.ReadFile("html/ops.html")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reading host html...")

	message := string(content)
	message = strings.ReplaceAll(message, "1000001", strconv.Itoa(board.Id))
	message = strings.ReplaceAll(message, "991", actualSeat)
	message = strings.ReplaceAll(message, "isHost=false", "isHost=true")
	message = strings.ReplaceAll(message, "未知身份", actualRole)
	message = strings.ReplaceAll(message, "vrcats", nickName)
	message = strings.ReplaceAll(message, "1234", strconv.Itoa(board.SeatsCount))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//message := fmt.Sprintf("开房成功！ 房间号为%d，可以邀请你的好友到 abc.com/%d 开始游戏\n", board.Id, board.Id)
	fmt.Fprintf(w, message)
}

// sit down in a room numbered board, preffered seat number number, with nick name nick
// returns: "为您选择了3号座位，目前12人中已有7人入座。您的身份是：<b>狼人<b>"
func sitDown(w http.ResponseWriter, req *http.Request) {
	b, n, k := req.FormValue("board"), req.FormValue("number"), req.FormValue("nick")
	boardId, _ := strconv.Atoi(b)
	board, ok := boards[boardId]
	if !ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<h2>404: 房间号没找到，请确认房间号后重试。</h2>")
		return
	}
	s, _ := strconv.Atoi(n)
	log.Println(s, k)
	sitDown := board.ViewCard(s, k)
	responses := strings.Split(sitDown, "\n")
	log.Println(responses)
	actualSeat := responses[0]
	actualRole := responses[1]

	content, err := ioutil.ReadFile("html/ops.html")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reading host html...")

	message := string(content)
	message = strings.ReplaceAll(message, "1000001", b)
	message = strings.ReplaceAll(message, "991", actualSeat)
	message = strings.ReplaceAll(message, "未知身份", actualRole)
	message = strings.ReplaceAll(message, "vrcats", k)
	message = strings.ReplaceAll(message, "1234", strconv.Itoa(board.SeatsCount))

	f3, _ := strconv.Atoi(responses[3])
	f2, _ := strconv.Atoi(responses[2])

	instruction := board.ToHtml()
	log.Println(instruction)
	msg := &Message{BoardId: board.Id, Body: instruction}
	hub.host <- msg

	if f3 == f2-1 {
		time.Sleep(100 * time.Millisecond)
		instruction := "房间已满，请房主开始游戏。"
		msg := &Message{BoardId: board.Id, Body: instruction}
		hub.host <- msg
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, message)
}

func operate(w http.ResponseWriter, req *http.Request) {
	n, number, action := req.FormValue("board"), req.FormValue("number"), req.FormValue("action")
	num1, num2, num3, skill, card := req.FormValue("num1"), req.FormValue("num2"), req.FormValue("num3"), req.FormValue("skill"), req.FormValue("card")
	log.Println(fmt.Sprintf("request form values: %v", req.Form))

	boardId, _ := strconv.Atoi(n)
	board, ok := boards[boardId]
	if !ok {
		fmt.Fprintf(w, "房间号没找到，请确认房间号后重试。")
		return
	}
	oldState := board.State
	log.Println(fmt.Sprintf("Calling TakeAction with parameters board: %s seat number %s action %s num1 %s num2 %s num3 %s skill %s card %s", n, number, action, num1, num2, num3, skill, card))
	fmt.Fprintf(w, board.TakeAction(number, action, num1, num2, num3, skill, card))
	// If state changes, send instruction to the board's host
	log.Println(fmt.Sprintf("State changed, old: %s new: %s, equals: \n", oldState, board.State, oldState == board.State))
	if oldState == board.State {
		log.Println("State not changed, skipping host notification...")
	} else {
		log.Println("State changed, sending notification to hosts...")
		instruction := "所有人请闭眼，5，4，3，2，1。" + board.SM[board.State][1]
		msg := &Message{BoardId: board.Id, Body: instruction}
		hub.host <- msg
		//如果黑商被反噬，幸运儿操作的轮次却不存在幸运儿，此时要假装幸运儿还在，进行延时操作
		_, ok := board.GetMeta("被反噬")
		if ok && board.SM[board.State][0] == "幸运儿" {
			time.Sleep(100 * time.Millisecond)
			luckyId := 991323 //Magic seat number，可以不操作
			url := fmt.Sprintf("operate?board=%d&number=%d&action=不操作&num1=0&num2=0&num3=0&skill=&card=", board.Id, luckyId)
			msg := &Message{BoardId: board.Id, Body: "fakeOps('" + url + "');"}
			hub.host <- msg
		}
	}
}

var boards map[int]*board.Board
var hub *Hub
var globalBoardId int

func main() {
	fmt.Println("Creating board pool...")
	boards = map[int]*board.Board{}
	globalBoardId = 1000000
	hub = newHub()
	go hub.run()
	log.Println("Websocket hub started...")

	// web service api entries
	http.HandleFunc("/checkIn", checkIn)
	http.HandleFunc("/sitDown", sitDown)
	http.HandleFunc("/operate", operate)
	http.Handle("/", http.FileServer(http.Dir("./html")))
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
