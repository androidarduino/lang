package board

import ()

type Player struct {
	Seat   int
	Labels []string
	Skills map[string]int
	Role   string
	Nick   string
}

// Add label to a player
func (p *Player) Label(l string) {
	p.Labels = append(p.Labels, l)
}

// Remove label from a player
func (p *Player) DeLabel(l string) {
	newLabels := []string{}
	for _, i := range p.Labels {
		if i != l {
			newLabels = append(newLabels, i)
		}
	}
	p.Labels = newLabels
}

// Create a new player with constructor values
func (p *Player) Create(Labels []string, Skills map[string]int) Player {
	p.Labels = Labels
	p.Skills = Skills
	return *p
}

//是否被贴了某个标签
func (p Player) HasLabel(l string) bool {
	for _, e := range p.Labels {
		if l == e {
			return true
		}
	}
	return false
}

//是否是好人
func (p Player) IsGood() bool {
	if !p.HasLabel("狼人") {
		return true
	}
	return false
}

//是否是神
func (p Player) IsGod() bool {
	return p.HasLabel("神")
}

//是否活着
func (p Player) IsAlive() bool {
	return p.Skills["寿命"] > 0
}

//是否能进行某些操作
func (p Player) CanDo(a string) bool {
	_, ok := p.Skills[a]
	return ok
}

func (p Player) InWerewolfTeam() bool {
	return (p.HasLabel("狼人") && !p.HasLabel("狼队不睁眼")) || (p.HasLabel("狼队睁眼"))
}

// Create a new player of given type
func (p *Player) New(playerType string) Player {
	if playerType == "预言家" {
		p.Create([]string{"预言家", "好人", "神"}, map[string]int{"预言家验": 99})
	}
	if playerType == "女巫" {
		p.Create([]string{"女巫", "好人", "神"}, map[string]int{"女巫毒": 1, "女巫救": 1})
	}
	if playerType == "猎人" {
		p.Create([]string{"猎人", "好人", "神"}, map[string]int{"查看状态": 99})
	}
	if playerType == "白痴" {
		p.Create([]string{"白痴", "好人", "神"}, map[string]int{})
	}
	if playerType == "守卫" {
		p.Create([]string{"守卫", "好人", "神"}, map[string]int{"守卫守": 99})
	}
	if playerType == "熊" {
		p.Create([]string{"熊", "好人", "神"}, map[string]int{})
	}
	if playerType == "乌鸦" {
		p.Create([]string{"乌鸦", "好人", "神"}, map[string]int{"乌鸦诽谤": 99})
	}
	if playerType == "企鹅" {
		p.Create([]string{"企鹅", "好人", "神"}, map[string]int{"企鹅冻": 99})
	}
	if playerType == "黑商" {
		p.Create([]string{"黑商", "好人", "神"}, map[string]int{"黑商给": 1})
	}
	if playerType == "盗贼" {
		p.Create([]string{"盗贼"}, map[string]int{"盗贼选": 1, "女巫救": 1})
	}
	if playerType == "魔术师" {
		p.Create([]string{"魔术师", "好人", "神"}, map[string]int{"魔术师交换": 1})
	}
	if playerType == "名媛" {
		p.Create([]string{"名媛", "好人", "神"}, map[string]int{"名媛睡": 1})
	}
	if playerType == "骑士" {
		p.Create([]string{"骑士", "好人", "神"}, map[string]int{})
	}
	if playerType == "潜行者" {
		p.Create([]string{"潜行者", "好人", "神"}, map[string]int{"潜行者暗杀": 1})
	}
	if playerType == "丘比特" {
		p.Create([]string{"丘比特", "好人", "神"}, map[string]int{"丘比特连": 1})
	}
	if playerType == "村民" {
		p.Create([]string{"村民", "好人", "民"}, map[string]int{})
	}
	if playerType == "高级村民" {
		p.Create([]string{"高级村民", "好人", "民"}, map[string]int{})
	}
	if playerType == "长老" {
		p.Create([]string{"长老", "好人", "民"}, map[string]int{"寿命": 2})
	}
	if playerType == "老流氓" {
		p.Create([]string{"老流氓", "好人", "民"}, map[string]int{"抗伤": 1})
	}
	if playerType == "两姐妹" {
		p.Create([]string{"两姐妹", "好人", "民"}, map[string]int{})
	}
	if playerType == "三兄弟" {
		p.Create([]string{"三兄弟", "好人", "民"}, map[string]int{})
	}
	if playerType == "白狼王" {
		p.Create([]string{"白狼王", "狼人"}, map[string]int{})
	}
	if playerType == "狼美人" {
		p.Create([]string{"狼美人", "狼人"}, map[string]int{"狼美人连": 99, "狼人刀": 99})
	}
	if playerType == "狼枪" {
		p.Create([]string{"狼枪", "狼人"}, map[string]int{"查看状态": 99})
	}
	if playerType == "恶魔" {
		p.Create([]string{"恶魔", "狼人"}, map[string]int{"恶魔验": 99})
	}
	if playerType == "种狼" {
		p.Create([]string{"种狼", "狼人"}, map[string]int{"种狼种": 1})
	}
	if playerType == "隐狼" {
		p.Create([]string{"隐狼"}, map[string]int{})
	}
	if playerType == "机械狼" {
		p.Create([]string{"机械狼", "好人"}, map[string]int{"机械狼学": 1})
	}
	if playerType == "狼人" {
		p.Create([]string{"狼人"}, map[string]int{"狼人刀": 99})
	}
	if playerType == "狼兄" {
		p.Create([]string{"狼兄", "狼人"}, map[string]int{"狼人刀": 99})
	}
	if playerType == "狼弟" {
		p.Create([]string{"狼弟", "狼人", "狼队不睁眼"}, map[string]int{})
	}
	if playerType == "混血儿" {
		p.Create([]string{"混血儿", "好人"}, map[string]int{"混血儿混": 1})
	}
	if playerType == "野孩子" {
		p.Create([]string{"野孩子", "好人"}, map[string]int{"野孩子混": 1})
	}
	if playerType == "迷妹" {
		p.Create([]string{"迷妹", "好人"}, map[string]int{"迷妹迷": 1, "迷妹粉": 1, "迷妹黑": 1})
	}
	if playerType == "小女孩" {
		p.Create([]string{"小女孩", "好人", "狼队睁眼", "民"}, map[string]int{})
	}
	if playerType == "通灵师" {
		p.Create([]string{"通灵师", "好人", "神"}, map[string]int{"通灵师验": 1})
	}
	if playerType == "守墓人" {
		p.Create([]string{"守墓人", "好人", "民"}, map[string]int{})
	}
	p.Role = playerType
	p.Label("所有人")
	p.Skills["寿命"] = 1
	p.Skills["抗伤"] = 0
	return *p
}
