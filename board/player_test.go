package board

import (
	"fmt"
	"testing"
)

func TestNewPlayer(t *testing.T) {

	var allRoles = []string{
		"预言家","女巫","猎人","通灵师","白痴","守卫",
		"熊","乌鸦","企鹅","黑商","盗贼","魔术师",
		"名媛","骑士","潜行者","丘比特","村民","高级村民",
		"长老","老流氓","两姐妹","三兄弟","白狼王","狼美人",
		"狼枪","恶魔","种狼","隐狼","机械狼","狼人",
		"狼兄","狼弟","混血儿","野孩子","迷妹","小女孩",
	}
	for _, role := range allRoles {
		p := new(Player)
		p.New(role)
		good, god := p.IsGood(), p.IsGod()
		inwwteam := p.InWerewolfTeam()
		var isgood, isgod, openeyes string
		if good {
			isgood = "好人"
		} else {
			isgood = "狼人"
		}
		if god {
			isgod = "是神"
		} else {
			isgod = "不是神"
		}
		if inwwteam {
			openeyes = "跟狼队睁眼"
		} else {
			openeyes = "不跟狼队睁眼"
		}
		fmt.Printf("身份：%s, %s %s %s    %t   %v\n", p.Role, isgood, isgod, openeyes, p.HasLabel("狼人"), p)
	}
}