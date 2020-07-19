package board

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sort"
	"math/rand"
	"time"
)

type Board struct{
	Id int
	WSUrl string
	Seats map[int]*Player
	SeatsCount int
	ActivePlayer *Player
	State string
	SM map[string][]string
	meta map[string]string
	log, report []string
}

func (b *Board) Log(l string) {
	// TODO: remove next line when debug is done
	fmt.Println(l)
	b.log = append(b.log, l)
}

func (b *Board) Report(l string) {
	b.report = append(b.report, l)
}

// Get player by seat number
func (b *Board) Player(n int) *Player {
	return b.Seats[n]
}

// Create a new board
func (b *Board) New(id int, roles []string, meta map[string]string) Board {
	log.Println(id, roles, meta)
	b.Seats = make(map[int]*Player)
	b.Id = id
	b.State = "setup"
	b.report = append(b.report, "昨夜信息：\n")
	//TODO 处理盗贼的牌，数量-2
	b.SeatsCount = len(roles)
	b.meta = meta
	fmt.Printf("New board created: %v\n", b)
	b.AddRoles(roles)
	b.SM = map[string][]string {
		"setup": {"房主", "生成房间并洗牌，如有盗贼生成盗贼选择，请大家入座查看身份", "allSeated","begin"},
		"begin": {"房主", "所有玩家已经入座完毕并查看了身份，请点击确认开始游戏", "done","10"},
		"10": {"混血儿", "混血儿请睁眼，选择一个号码作为爸爸","done","20"},
		"20": {"野孩子", "野孩子请睁眼，请选择一个号码作为榜样，榜样死亡你将变成狼人", "done","22"},
		"22": {"企鹅", "企鹅请睁眼，请选择一位玩家冷冻他的技能", "done","25"},
		"25": {"迷妹", "迷妹请睁眼，请选择你的偶像", "done","30"},
		"30": {"盗贼", "盗贼请睁眼查看状态，从两张牌中选择一张作为自己的身份，如果有狼牌必须选择狼牌", "done","40"},
		"40": {"丘比特", "丘比特请睁眼，请选择两位玩家成为情侣。如果这两位玩家一好人一狼，你们三个会成为第三方", "done","43"},
		"43": {"所有人", "所有人请睁眼，并点击查看结果，确认你有没有被连为情侣", "allViewed","45"},
		"45": {"情侣", "情侣请睁眼互认，但不要交流身份。如果你们一好人一狼，你们三个会成为第三方", "done","50"},
		"50": {"机械狼", "机械狼请睁眼，本轮是否要选择一位玩家学习他的技能，学习后你的身份将与他一样", "done","55"},
		"55": {"魔术师", "魔术师请睁眼，请选择是否要交换场上两位玩家今晚的身份", "done","60"},
		"60": {"两姐妹", "两姐妹请睁眼互认。你们投票时必须投同一位玩家，请相互确认后点击确认", "done","70"},
		"70": {"三兄弟", "三兄弟请睁眼互认。你们投票时必须投统一为玩家，请相互确认后点击确认", "done","90"},
		"90": {"黑商", "黑商请睁眼，选择一位玩家并赋予他一项神技，若选中狼人，则他不会得到技能而你会被反噬而死", "done","100"},
		"100": {"所有人", "所有人请睁眼病点击查看结果，确认你有没有被黑商赋予技能", "allViewed","110"},
		"110": {"幸运儿", "幸运儿请睁眼，今天你是否要使用你得到的技能，如果使用请选择技能和目标", "done","120"},
		"120": {"狼人", "所有狼人和小女孩请睁眼，请用手语商量战术，并选择一个目标进行狙杀", "done","130"},
		"130": {"狼兄", "狼兄狼弟请睁眼，互认身份之后请狼兄点击确认继续", "done","140"},
		"140": {"狼美人", "狼美人请睁眼，请选择魅惑一位玩家", "done","150"},
		"150": {"恶魔", "恶魔请睁眼，请选择一位玩家查验他的身份是不是神", "done","160"},
		"160": {"预言家", "预言家请睁眼，请选择一位玩家查验他的身份是好是坏", "done","170"},
		"170": {"女巫", "女巫请睁眼，请点击查看结果查看哪位玩家死亡，然后选择是否使用解药和毒药中的一瓶", "done","180"},
		"180": {"猎人", "猎人请睁眼，请点击查看结果查看开枪状态", "done","190"},
		"190": {"守卫", "守卫请睁眼，请选择要守卫的玩家", "done","200"},
		"200": {"狼枪", "狼枪请睁眼，请点击查看结果查看开枪状态", "done","205"},
		"205": {"机械狼", "机械狼请睁眼，请选择是否要使用学到的技能", "done","208"},
		"208": {"种枪", "种狼请睁眼，请选择你是否要感染今天被杀的玩家", "done","210"},
		"210": {"通灵师", "通灵师请睁眼，请选择一个玩家查验他的具体身份", "done","220"},
		"220": {"狐狸", "狐狸请睁眼，请选择一位玩家查验他和他身边两位活着的玩家有没有狼人", "done","230"},
		"230": {"乌鸦", "乌鸦请睁眼，请选择是否要诽谤一位玩家。受到诽谤后投票放逐时此玩家会多出一票", "done","250"},
		"250": {"禁言长老", "禁言长老请睁眼，请选择是否要禁止一位玩家明天的发言", "done","270"},
		"270": {"名媛", "名媛请睁眼，请选择要与哪位玩家共渡春宵。此玩家今晚中刀中毒不会死亡，但你若死亡他也将殉葬。若连睡两晚，对方将死亡。", "done","280"},
		"280": {"迷妹", "迷妹请睁眼，是否要对你的偶像使用粉或黑的技能，粉可以减少一票，黑可以增加一票", "done","290"},
		"290": {"潜行者", "潜行者请睁眼，请选择是否要暗杀你白天投票的玩家", "done","end"},
		"end": {"房主", "", "done", "EXIT"},
		"EXIT": {"房主", "请所有人整理表情，想竞选警长的玩家请举手。5，4，3，2，1，所有人睁眼。选举结束后请房主查看昨夜信息。", "done", ""},
	}
	//开牌前洗牌
	b.shuffle()
	return *b
}

func (b* Board) inOperatorGroup(seatNumber int) bool {
	player := b.Player(seatNumber)
	role := b.SM[b.State][0]
	//log.Println(fmt.Sprintf("正在检查：%v 是否属于 %s 团队，结果%b。", player, role, player.HasLabel(role)))

	return player.HasLabel(role)
}

func (board* Board) TakeAction(seatNumber string, action string, n1 string, n2 string, n3 string, skill string, card string) string {
	seat, err := strconv.Atoi(seatNumber)
	if (err != nil) {
		return "座位号错误"
	}
	num1, err := strconv.Atoi(n1)
	if (err != nil) {
		return fmt.Sprintf("操作数错误%s ", n1)
	}
	num2, err := strconv.Atoi(n2)
	if (err != nil) {
		return fmt.Sprintf("操作数错误%s ", n2)
	}
	num3, err := strconv.Atoi(n3)
	if (err != nil) {
		return fmt.Sprintf("操作数错误%d ", n3)
	}

	if num1 > board.SeatsCount || num2 > board.SeatsCount || num3 > board.SeatsCount || seat > board.SeatsCount {
		return "参数错误。"
	}
	if action != "房主开局" && (board.State == "setup" || board.State == "begin") {
		return "尚未开局，还不能进行操作。"
	}
	//检查是否有权限做这个事情，发起action的人身份必须与当前轮次操作人身份相符
	if (!board.inOperatorGroup(seat)) {
		operator := board.SM[board.State][0]
		return fmt.Sprintf("当前是 %s 操作的轮次，您没有操作权限。", operator)
	}
	//TODO: 检查是否被冻，如果是任何一个狼人被冻，则所有狼人都不可操作
	board.ActivePlayer = board.Player(seat)
	message := ""
	switch(action) {
		case "房主开局":
			message = board.startGame()
		case "所有人入座":
			//用seat代替preferred seat，用action代替用户的nick name
			//TODO: 分析一下从http.go调用好还是action调用好
			message = board.ViewCard(seat, action)
		case "房主关闭":
			message = board.endGame()
		case "房主重新发牌":
			message = board.shuffle()
		case "房主查看昨夜结果":
			message = board.lastNightResult()
		case "不操作":
			message = board.skip()
		case "猎人状态":
			message = board.confirm()
		case "女巫查看":
			message = board.witchResult()
		case "全体查看情侣":
			message = board.checkValentine()
		case "全体查看幸运儿":
			message = board.checkLucky()
		case "讨论完毕":
			message = board.skip()
		case "狼人杀":
			message = board.slaughter(num1)
		case "女巫毒":
			message = board.poison(num1)
		case "女巫救":
			message = board.heal(num1)
		case "预言家验":
			message = board.examine(num1)
		case "狼美人连":
			message = board.charm(num1)
		case "丘比特连":
			message = board.qubit(num1, num2)
		case "守卫守":
			message = board.guard(num1)
		case "黑商给":
			message = board.endow(num1, skill)
		case "幸运儿验":
			message = board.examine(num1)
		case "幸运儿毒":
			message = board.poison(num1)
		case "狐狸验":
			message = board.foxExamine(num1)
		case "企鹅冻":
			message = board.freeze(num1)
		case "乌鸦诽谤":
			message = board.defamation(num1)
		case "盗贼选":
			message = board.burgular(card)
		case "魔术师交换":
			message = board.swap(num1, num2)
		case "名媛睡":
			message = board.sleep(num1)
		case "潜行者暗杀":
			message = board.assasin(num1)
		case "禁言长老禁言":
			message = board.forbidden(num1)
		case "通灵师验":
			message = board.psychicExamine(num1)
		case "机械狼学":
			message = board.learn(num1)
		case "机械狼验":
			message = board.psychicExamine(num1)
		case "机械狼毒":
			message = board.poison(num1)
		case "机械狼守":
			message = board.guardPoison(num1)
		case "种狼感染":
			message = board.infect()
		case "混血儿混":
			message = board.mix(num1)
		case "野孩子混":
			message = board.father(num1)
		case "迷妹迷":
			message = board.fan(num1)
		case "迷妹粉":
			message = board.idol()
		case "迷妹黑":
			message = board.hate()
		case "恶魔验":
			message = board.deamonExamine(num1)
	}
	log.Println("返回信息：" + message)
	return message
}

// Add roles to the board, TODO: consider burgler case
func (b *Board) AddRoles(roles []string) {
	for i, role := range roles {
		player := new(Player)
		player.New(role)
		b.Seats[i+1] = player
		player.Seat = i+1
	}
}

func (b *Board) TakenSeatsCount() int {
	takenSeats := 0
	for k,_ := range b.Seats {
		if b.Seats[k].HasLabel("已入座") {
			takenSeats ++
		}
	}
	return takenSeats
}

func (b *Board) TakeSeat(seatNumber int, nickName string) (int,int,int) {
	totalSeats := len(b.Seats)
	takenSeats := b.TakenSeatsCount()
	if takenSeats == totalSeats {
		return -1, totalSeats, takenSeats
	}
	takenSeats = takenSeats + 1
	if !b.Seats[seatNumber].HasLabel("已入座") {
		b.Seats[seatNumber].Label("已入座")
		b.Seats[seatNumber].Seat = seatNumber
		b.Seats[seatNumber].Nick = nickName
		return seatNumber, totalSeats, takenSeats
	}
	for k,_ := range b.Seats {
		if !b.Seats[k].HasLabel("已入座") {
			b.Seats[k].Label("已入座")
			return k, totalSeats, takenSeats
		}
	}
	return -1, totalSeats, takenSeats
}

func (b *Board) hasRole(role string) bool {
	for _, player := range b.Seats {
		if player.Role == role {
			return true
		}
	}
	return false
}

func (b *Board) nextStep() {
	//TODO： 有丘比特或黑商的板子，不能跳过43和100两个状态
	//跳过所有不存在的角色
	if b.State == "EXIT" {
		return
	}
	s := b.State
	nextState := b.SM[b.State][3]
	b.State = nextState
	for  b.State!="EXIT" && !b.hasRole(b.SM[b.State][0]) {
		nextState := b.SM[b.State][3]
		b.State = nextState
	}
	b.Log(fmt.Sprintf("正在从%s状态向%s状态转换......", s, b.State))
}

func (b *Board) skip() string {
	b.nextStep()
	return "进入下一轮次"
}

func (b *Board) allSeated() bool {
	return b.TakenSeatsCount() == b.SeatsCount
}

func (b *Board) allViewed() bool {
	for _, player := range b.Seats {
		if !player.HasLabel(b.State+"已查看") {
			return false
		}
	}
	return true
}

func (b *Board) checkSkill(skill string) (int, bool, string) {
	left, ok := b.ActivePlayer.Skills[skill]
	if !(ok) {
		return  -1, false, fmt.Sprintf("您没有%s的技能，请不要乱点。", skill)
	}
	if (left > 0) {
		b.ActivePlayer.Skills[skill] = left - 1
		return left - 1, true, ""
	} else {
		return left, false, ""
	}
}

func (b *Board) startGame() string {
	//开始游戏
	b.nextStep()
	return "游戏开始"
}
func (b *Board) ViewCard(prefferedSeat int, nickName string) string {
	//调用TakeSeat
	if prefferedSeat < 1 || prefferedSeat > len(b.Seats) {
		prefferedSeat = 1
	}
	seat, total_seats, taken_seats := b.TakeSeat(prefferedSeat, nickName)
	if (total_seats == taken_seats) {
		return fmt.Sprintf("%d\n%s\n%d\n%d", 0, "房间已满", total_seats, taken_seats)
	}
	role := b.Seats[seat].Role
	b.Log(fmt.Sprintf("玩家 %s 占据了%d号座位，他的身份是%s, %d个座位中已经有%d个座位有人。",nickName, seat, role, total_seats, taken_seats))
	b.Println()
	return fmt.Sprintf("%d\n%s\n%d\n%d", seat, role, total_seats, taken_seats)
}
func (b *Board) endGame() string {
	//发出自杀信号，要求管理器移除这个板子
	return "TODO"
}
func (b *Board) shuffle() string {
	//保持座位不变，重新安排玩家的身份和能力，只有在开局以前可以调用
	rand.Seed(time.Now().UnixNano())
	for i := len(b.Seats); i > 0; i-- { // Fisher–Yates shuffle
	    j := rand.Intn(i) + 1
	    log.Println("swapping %d and %d", i, j)
	    b.Seats[i], b.Seats[j] = b.Seats[j], b.Seats[i]
	}
	for idx, player := range b.Seats {
		log.Println(idx)
		log.Println(player)
		player.Seat = idx
	}
	b.Println()
	return "洗牌成功。"
}

func (b *Board) Println() {
	log.Println("Here is the board ==================")
	for _, p := range b.Seats {
		log.Println(*p)
	}
	log.Println(" ==================")

}

func (b *Board) lastNightResult() string {
	//TODO: 防止重复计算, set a flag when it is already checked
	if b.State == "EXIT" {
		//return strings.Join(b.report, "\n")
	}

	log.Println("正在计算昨夜结果....")
	//TODO: 计算昨夜死讯，需要考虑的情况非常多
	//{"被刀","被毒","被救","被连","被守","被睡","被潜","情侣","被守毒","被感染"}
	b.Println()


	death := []*Player {}
	for _, p := range b.Seats {
		//计算每个玩家是否死亡，或半死
		if p.HasLabel("被睡") {
			continue
		}
		if p.HasLabel("被刀") {
			if p.HasLabel("被救") {
				if p.HasLabel("被守") {
					death = append(death, p)
				} else {
					//healed
				}
			} else {
				death = append(death, p)
			}
		}
		if p.HasLabel("被毒") {
			if p.HasLabel("被守毒") {
				//guarded
			} else {
				death = append(death, p)
			}
		}
		if p.HasLabel("被潜") || p.HasLabel("被掏空") || p.HasLabel("被反噬") {
			death = append(death, p)
		}
	}
	log.Println(fmt.Sprintf("found death: %v", death))
	//处理被连被睡情侣多死的情况，要循环四次，因为可能有连环死
	for i := 0; i < 4; i++ {
		for _, p := range death {
			if p.HasLabel("狼美人") {
				for _, k := range b.Seats {
					if k.HasLabel("被连") {
						death = append(death, k)
					}
				}
			}
			if p.HasLabel("名媛") {
				for _, k := range b.Seats {
					if k.HasLabel("被睡") {
						death = append(death, k)
					}
				}			
			}
			if p.HasLabel("情侣") {
				for _, k := range b.Seats {
					if k.HasLabel("被睡") {
						death = append(death, k)
					}
				}	
			}
		}
	}
	log.Println(fmt.Sprintf("found death afterward: %v", death))
	dn := map[int]bool {}

	for _, p := range death {
		dn[p.Seat] = true
	}
	for n, _ := range dn {
		b.report = append(b.report, fmt.Sprintf("昨晚%d号玩家死亡。", n))
	}

	if len(dn) == 0 {
		b.report = append(b.report, fmt.Sprintf("昨晚是平安夜。"))
	}

	sort.Strings(b.report)
	b.nextStep()
	return strings.Join(b.report, "\n")
}
//用作猎人，狼枪进行确认开枪状态
func (b *Board) confirm() string {
	seat := b.ActivePlayer.Seat
	if (!b.inOperatorGroup(seat)) {
		operator := b.SM[b.State][0]
		return fmt.Sprintf("当前是 %s 操作的轮次，您没有操作权限。", operator)
	}
	b.nextStep()
	if (b.ActivePlayer.HasLabel("被毒")||b.ActivePlayer.HasLabel("被冻")) {
		return "你昨晚被毒或被冻了，明天无法开枪。"
	} else {
		return "你昨晚没有被毒，明天可以正常开枪。"
	}
}
//用于全体玩家查看情侣信息，所有玩家查看完成后才可以进入下一步
func (b *Board) checkValentine() string {
	b.ActivePlayer.Label(b.State+"已查看")
	if b.allViewed() {
		b.nextStep()
	}
	if b.ActivePlayer.HasLabel("情侣") {
		otherPlayer := b.ActivePlayer
		for _, player := range b.Seats {
			if player.HasLabel("情侣") && player != b.ActivePlayer {
				otherPlayer = player
			}
		}
		return fmt.Sprintf("您和%d号玩家'%s'被连为情侣，你们将同生共死。\n如果是一好人一狼人，你们与丘比特将组成第三方阵营。", otherPlayer.Seat, otherPlayer.Nick)
	} else {
		return "很遗憾，您没有被连为情侣。"
	}
}
//用于全体玩家查看幸运儿信息，所有玩家查看完成后才可以进入下一步
func (b *Board) checkLucky() string {
	b.ActivePlayer.Label(b.State+"已查看")
	if b.allViewed() {
		b.nextStep()
	}
	if (b.ActivePlayer.HasLabel("幸运儿")) {
		return fmt.Sprintf("恭喜您，您被选作幸运儿，您现在拥有如下技能：%v 。", b.ActivePlayer.Skills)
	} else {
		return "很遗憾，您没有被黑商选中，下次请继续努力！"
	}
}
//用于女巫查看昨夜杀人信息
func (b *Board) witchResult() string {
	player := b.ActivePlayer
	n, ok := player.Skills["女巫救"]
	if (!ok || n < 1) {
		return "您没有解药，看不到刀人结果。"
	} else {
		last, ok := b.meta["昨晚中刀"]
		if (!ok) {
			return "昨晚狼人空刀。"
		}
		return fmt.Sprintf("昨晚%s号玩家被刀，如果要救请选择菜单中的 女巫救 。", last)
	}
	//不需要修改状态机，只有女巫选择救或毒时才修改状态机
}
func (b *Board) slaughter(num1 int) string {
	if num1 == 0 {
		b.nextStep()
		return "狼人空刀成功。"
	}
	player := b.ActivePlayer
	target := b.Seats[num1]
	if (target.HasLabel("恶魔")||target.HasLabel("狼枪")||target.HasLabel("狼美人")) {
		return fmt.Sprintf("%d号玩家是不可以自刀的角色，请重新选择。", target.Seat)
	}
	n, ok := player.Skills["狼人刀"]
	if (!ok || n < 1) {
		return "您没有刀，或者没有刀人的技能。"
	} else {
		b.meta["昨晚中刀"] = strconv.Itoa(target.Seat)
		target.Label("被刀")
		b.Log(fmt.Sprintf("%d号狼人刀了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
	}
	return "狼人刀人成功。"
}
func (b *Board) poison(num1 int) string {
	player := b.ActivePlayer
	if num1 == 0 {
		b.nextStep()
		return "操作成功，没有使用毒药。"
	}
	target := b.Seats[num1]
	if (target.HasLabel("恶魔")) {
		return "操作成功。"
	}
	b.meta["昨晚被毒"] = strconv.Itoa(num1)
	n, ok := player.Skills["女巫毒"]
	if (!ok || n < 1) {
		return "操作失败，您没有毒药。"
	} else {
		target.Label("被毒")
		b.Log(fmt.Sprintf("%d号女巫毒了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
	}
	return "操作成功，使用了一瓶毒药。"
}
func (b *Board) heal(num1 int) string {
	player := b.ActivePlayer
	nu, ok := b.meta["昨晚中刀"]
	if (!ok) {
		return "操作成功。没有使用解药。"
	}
	ns, _ := strconv.Atoi(nu)
	target := b.Seats[ns]
	if (target.HasLabel("女巫") && b.meta["女巫自救"] == "不能") {
		return "操作失败，女巫不可能自救。"
	}
	b.meta["昨晚被救"] = strconv.Itoa(ns)
	n, ok := player.Skills["女巫救"]
	if (!ok || n<1) {
		return "操作失败，您没有解药。"
	} else {
		target.Label("被救")
		b.Log(fmt.Sprintf("%d号女巫救了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
	}
	return "操作成功，使用了一瓶解药。"
}
func (b *Board) examine(num1 int) string {
	player := b.ActivePlayer
	if num1 == 0 || player.Seat == num1 {
		return "操作错误，不能验空号码或验自己。"
	}
	target := b.Seats[num1]
	b.meta["昨晚被验"] = strconv.Itoa(num1)
	_, ok := player.Skills["预言家验"]
	if (!ok) {
		return "操作失败，您没有验人功能。"
	} else {
		target.Label("被验")
		b.Log(fmt.Sprintf("%d号预言家验了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
		if target.IsGood() {
			return fmt.Sprintf("%d号玩家的身份是: 好人", target.Seat)
		} else {
			return fmt.Sprintf("%d号玩家的身份是: 狼人", target.Seat)
		}
	}
	return "操作成功。"
}

func (b *Board) charm(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	l, ok := b.meta["昨晚被连"]
	if !ok {
		l = "-1"
	}
	if (strconv.Itoa(num1) == l || num1 == player.Seat) {
		return "操作失败，不可以连续两晚连同一个人，也不可以连自己。"
	}
	b.meta["昨晚被连"] = strconv.Itoa(num1)
	n, ok := player.Skills["狼美人连"]
	if (!ok || n<1) {
		return "操作失败，您没有连人的能力。"
	} else {
		target.Label("被连")
		ll, _ := strconv.Atoi(l)
		if (ll > 0) {
			b.Seats[ll].DeLabel("被连")
		}
		b.Log(fmt.Sprintf("%d号狼美人连了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
	}
	return "操作成功。"
}

func (b *Board) qubit(num1, num2 int) string {
	player := b.ActivePlayer
	target1 := b.Seats[num1]
	target2 := b.Seats[num2]
	if (num1 == num2 || num1 == player.Seat || num2 == player.Seat) {
		return "操作失败，不可以连同一个人或丘比特自己。"
	}
	n, ok := player.Skills["丘比特连"]
	if (!ok || n<1) {
		return "操作失败，您没有丘比特的能力。"
	} else {
		target1.Label("情侣")
		target2.Label("情侣")
		b.Log(fmt.Sprintf("%d号丘比特连了%d号和%d号玩家作为情侣", player.Seat, target1.Seat, target2.Seat))
		b.nextStep()
	}	
	return "操作成功。"
}
func (b *Board) guard(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	l, ok := b.meta["昨晚被守"]
	if !ok {
		l = "-1"
	}
	if (strconv.Itoa(num1) == l && l != "0") {
		return "操作失败，不可以连续两晚守同一个人"
	}
	b.meta["昨晚被守"] = strconv.Itoa(num1)
	if num1 == 0 {
		b.nextStep()
		return "操作成功。"
	}
	n, ok := player.Skills["守卫守"]
	if (!ok || n<1) {
		return "操作失败，您没有守人的能力。"
	} else {
		target.Label("被守")
		ll, _ := strconv.Atoi(l)
		if (ll > 0) {
			b.Seats[ll].DeLabel("被守")
		}
		b.Log(fmt.Sprintf("%d号守卫守了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
	}
	return "操作成功。"
}
func (b *Board) endow(num1 int, skill string) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	n, ok := player.Skills["黑商给"]
	if (!ok || n < 1) {
		return "操作失败，您没有黑商给技能的能力。"
	} else {
		if target.IsGood() {
			target.Label("幸运儿")
			target.Skills[skill] = 1
			b.Log(fmt.Sprintf("%d号黑商给了%d号玩家‘%s’技能", player.Seat, num1, skill))
		} else {
			b.Log(fmt.Sprintf("%d号黑商给了%d号狼人‘%s’技能， 被反噬而死", player.Seat, num1, skill))
			//player.Skills["寿命"] = 0
			//b.Report(fmt.Sprintf("昨夜%d死亡。", player.Seat))
			player.Label("被反噬")
		}
		b.nextStep()
	}
	return "技能已经成功发出。如果是好人会收到技能，狼人则会令你受到反噬而死。"
}
func (b *Board) foxExamine(num1 int) string {
	player := b.ActivePlayer
//	target := b.Seats[num1]
	n, ok := player.Skills["狐狸验"]
	if (!ok || n<1) {
		return "操作失败，您没有狐狸验人的能力。"
	} else {
		b.Log(fmt.Sprintf("%d号狐狸验了%d号及其左右的身份", player.Seat, num1))

		/*
		左边第一个活着的人 a
		右边第一个活着的人 b
		if a == p return no_wolf
		if b == p return no_wolf
		if a == b return no_wolf
		if a <> b return bad(a) | bad(b) | bad(c)
		*/
		a1, a2 := 0, 0
		message := ""
		p := player.Seat
		for i := 0; i < b.SeatsCount; i ++ {
			p = p - 1
			if p == 0 {
				p = b.SeatsCount
			}
			if b.Seats[p].IsAlive() {
				a1 = p
			}
		}
		p = player.Seat
		for i := 0; i < b.SeatsCount; i ++ {
			p = p + 1
			if p > b.SeatsCount {
				p = 1
			}
			if b.Seats[p].IsAlive() {
				a2 = p
			}
		}
		if (a1 == a2 || a1 == p || a2 == p) {
			message = "您查验的人及其身边两个活着的人里面没有狼人。"
		} else {
			if (b.Seats[a1].IsGood() && b.Seats[a2].IsGood() && b.Seats[p].IsGood()) {
				message = "您查验的人及其身边两个活着的人里面没有狼人。"
			} else {
				message = "您查验的人及其身边两个活着的人里面有狼人。"

			}
		}
		b.nextStep()
		return message
	}
	return "操作成功。"
}
func (b *Board) freeze(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	f, ok := b.meta["昨晚被冻"]
	if (ok && f == strconv.Itoa(num1)) {
		return "操作失败，不能连续两晚冻同一个人。"
	}
	b.meta["昨晚被冻"] = strconv.Itoa(num1)
	n, ok := player.Skills["企鹅冻"]
	if (!ok || n<1) {
		return "操作失败，您没有冻人功能。"
	} else {
		//TODO: 在所有人操作之前验证有没有被冻，猎人状态也需要更新
		target.Label("被冻")
		b.Log(fmt.Sprintf("%d号企鹅冻了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
		return fmt.Sprintf("冷冻成功， 本轮%d号玩家将不能使用技能。", target.Seat)
	}	
}
func (b *Board) defamation(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	b.meta["昨晚被诽谤"] = strconv.Itoa(num1)
	n, ok := player.Skills["乌鸦诽谤"]
	if (!ok || n<1) {
		return "操作失败，您没有诽谤人功能。"
	} else {
		//TODO: 在所有人操作之前验证有没有被冻，猎人状态也需要更新
		target.Label("被诽谤")
		b.Log(fmt.Sprintf("%d号乌鸦诽谤了%d号玩家", player.Seat, target.Seat))
		b.Report(fmt.Sprintf("昨夜%d号玩家被诽谤，放逐投票时将多一票。", target.Seat))
		b.nextStep()
		return fmt.Sprintf("冷冻成功， 本轮%d号玩家将获得额外一张放逐票。", target.Seat)
	}
}

func (b *Board) generateBurgularCards() {
	//TODO： 从所有牌里选取两张作为盗贼的候选牌，分别为1号和2号
	//TODO： 生成板子之前，发现有盗贼则座位要少两个，先检查一下最后两张是否都是狼人牌，如果是，则无限shuffle，直到不是为止。
}

func (b *Board) burgular(selection string) string {
	//TODO： 选择一张牌作为身份牌，更新自己的身份
	return "TODO"
}
func (b *Board) swap(num1, num2 int) string {
	//TODO：魔术师的交换比较复杂，要记得天亮之前回复两人的号码
	return "TODO"
}
func (b *Board) sleep(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	l, ok := b.meta["昨晚被睡"]
	if !ok {
		l = "-1"
	}
	if num1 == player.Seat {
		return "操作失败，不可以睡自己。"
	}
	if (strconv.Itoa(num1) == l) {
		b.Log(fmt.Sprintf("%d号名媛连续两晚睡了%d号玩家致使其死亡。", player.Seat, num1))
		//target.Skills["寿命"] = 0
		//b.Report(fmt.Sprintf("昨夜%d死亡。", target.Seat))
		target.Label("被掏空")
		return fmt.Sprintf("%d号玩家因为连续两晚与你共度良宵精尽人亡。", num1)
	}
	b.meta["昨晚被睡"] = strconv.Itoa(num1)
	n, ok := player.Skills["名媛睡"]
	if (!ok || n<1) {
		return "操作失败，您没有睡人的能力。"
	} else {
		//TODO：计算死亡时要考虑被睡的情况
		target.Label("被睡")
		ll, _ := strconv.Atoi(l)
		if (ll > 0) {
			b.Seats[ll].DeLabel("被睡")
		}
		b.Log(fmt.Sprintf("%d号名媛睡了%d号玩家。", player.Seat, target.Seat))
		b.nextStep()
	}
	return "操作成功。"
}
//潜行者刺杀白天投的玩家
func (b *Board) assasin(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	n, ok := player.Skills["潜行者潜"]
	if (!ok || n<1) {
		return "操作失败，您没有刺杀功能。"
	} else {
		//TODO: 计算死亡时要考虑被潜行者处理的情况
		target.Label("被潜")
		b.Log(fmt.Sprintf("%d号潜行者刺杀了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
		return fmt.Sprintf("刺杀成功，%d号玩家已经死亡。", target.Seat)
	}
}
func (b *Board) forbidden(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	n, ok := player.Skills["禁言长老禁言"]
	if (!ok || n<1) {
		return "操作失败，您没有禁言功能。"
	} else {
		target.Label("被禁言")
		b.Log(fmt.Sprintf("%d号禁言长老禁言了%d号玩家", player.Seat, target.Seat))
		b.Report(fmt.Sprintf("昨夜%d号玩家被禁言，今天白天的发言将被跳过。", target.Seat))
		b.nextStep()
		return fmt.Sprintf("禁言成功，明天%d号玩家的发言将被跳过。", target.Seat)
	}
}
func (b *Board) psychicExamine(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	n, ok := player.Skills["通灵师验"]
	if (!ok || n<1) {
		return "操作失败，您没有通灵功能。"
	} else {
		b.Log(fmt.Sprintf("%d号通灵师验了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
		return fmt.Sprintf("%d号玩家的身份是%s。", target.Seat, target.Role)
	}
}
func (b *Board) deamonExamine(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	n, ok := player.Skills["恶魔验"]
	if (!ok || n<1) {
		return "操作失败，您没有恶魔验人功能。"
	} else {
		b.Log(fmt.Sprintf("%d号恶魔验了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
		if (target.IsGod()) {
			return fmt.Sprintf("%d号玩家的身份是神牌。", target.Seat)
		} else {
			return fmt.Sprintf("%d号玩家的身份不是神牌。", target.Seat)
		}
	}
}
func (b *Board) learn(num1 int) string {
	//TODO：比较复杂，身份需要复制，功能需要修改，如守卫功能要修改成守毒，女巫只选毒药
	return "TODO"
}
func (b *Board) guardPoison(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	l, ok := b.meta["昨晚被守毒"]
	if !ok {
		l = "-1"
	}
	if (strconv.Itoa(num1) == l) {
		return "操作失败，不可以连续两晚守毒同一个人"
	}
	b.meta["昨晚被守毒"] = strconv.Itoa(num1)
	n, ok := player.Skills["机械狼守毒"]
	if (!ok || n<1) {
		return "操作失败，您没有守毒的能力。"
	} else {
		//TODO: 计算死讯的时候考虑守毒
		target.Label("被守毒")
		ll, _ := strconv.Atoi(l)
		if (ll > 0) {
			b.Seats[ll].DeLabel("被守毒")
		}
		b.Log(fmt.Sprintf("%d号机械狼守毒了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
	}
	return "操作成功。"
}
func (b *Board) infect() string {
	player := b.ActivePlayer
	zd, ok := b.meta["昨晚中刀"]
	if (!ok) {
		return "操作失败，狼人昨晚没有刀人。"
	}
	ls, _ := strconv.Atoi(zd)
	target := b.Seats[ls]
	n, ok := player.Skills["种狼感染"]
	if (!ok || n<1) {
		return "操作失败，您没有种狼感染的能力。"
	} else {
		//TODO: 计算死讯的时候考虑感染，如果没有死且感染，需要将其转化为狼人，并且将能力转化为响应的能力
		target.Label("被感染")
		ll, _ := strconv.Atoi(zd)
		if (ll > 0) {
			b.Seats[ll].DeLabel("被感染")
		}
		b.Log(fmt.Sprintf("%d号种狼感染了%d号玩家", player.Seat, target.Seat))
		b.nextStep()
	}
	return "操作成功。"
}
func (b *Board) mix(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	if player.Seat == num1 || num1 == 0 {
		return "操作失败，不能混自己或混空号。"
	}
	n, ok := player.Skills["混血儿混"]
	if (!ok || n<1) {
		return "操作失败，您没有混血儿混人的能力。"
	} else {
		b.Log(fmt.Sprintf("%d号混血儿混了%d号玩家", player.Seat, target.Seat))
		b.meta["混血儿混"] = strconv.Itoa(num1)
		b.nextStep()
	}
	return "操作成功。"
}
func (b *Board) father(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	n, ok := player.Skills["野孩子混"]
	if (!ok || n<1) {
		return "操作失败，您没有野孩子认爹的能力。"
	} else {
		b.Log(fmt.Sprintf("%d号野孩子认了%d号玩家作为爸爸", player.Seat, target.Seat))
		b.meta["野孩子认"] = strconv.Itoa(num1)
		b.nextStep()
	}
	return "操作成功。"	
}
func (b *Board) fan(num1 int) string {
	player := b.ActivePlayer
	target := b.Seats[num1]
	idol, ok := player.Skills["迷妹迷"]
	if (ok) {
		return fmt.Sprintf("操作失败，您已经认定了%d号玩家作为偶像。", idol)
	} else {
		b.Log(fmt.Sprintf("%d号迷妹成为了%d号玩家的头号粉丝", player.Seat, target.Seat))
		b.meta["迷妹迷"] = strconv.Itoa(num1)
		b.nextStep()
	}
	return "操作成功。"
}
func (b *Board) idol() string {
	player := b.ActivePlayer
	idol, _ := strconv.Atoi(b.meta["迷妹迷"])
	target := b.Seats[idol]
	n, ok := player.Skills["迷妹粉"]
	if (!ok || n<1) {
		return "操作失败，您没有后援团功能。"
	} else {
		target.Label("被粉")
		b.Log(fmt.Sprintf("%d号迷妹粉了%d号玩家", player.Seat, target.Seat))
		b.Report(fmt.Sprintf("昨夜%d号玩家被迷妹后援团支持，放逐投票时将少一票。", target.Seat))
		b.nextStep()
		return fmt.Sprintf("支持偶像成功， 本轮%d号玩家将少获得一张放逐票。", target.Seat)
	}
}

func (b *Board) hate() string {
	player := b.ActivePlayer
	idol, _ := strconv.Atoi(b.meta["迷妹迷"])
	target := b.Seats[idol]
	n, ok := player.Skills["迷妹黑"]
	if (!ok || n<1) {
		return "操作失败，您没有黑人功能。"
	} else {
		target.Label("被黑")
		b.Log(fmt.Sprintf("%d号迷妹黑了%d号玩家", player.Seat, target.Seat))
		b.Report(fmt.Sprintf("昨夜%d号玩家被粉丝爆黑料，放逐投票时将多一票。", target.Seat))
		b.nextStep()
		return fmt.Sprintf("黑粉操作成功， 本轮%d号玩家将获得额外一张粉丝放逐票。", target.Seat)
	}
}