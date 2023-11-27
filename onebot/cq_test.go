package onebot

import (
	"testing"

	"github.com/shynome/fuck-qq/onebot/msg"
)

var testMsgRaw = "[CQ:at,qq=all] xxx 开播啦\n直播间标题\nhttps://live.bilibili.com/24393\n[CQ:image,file=https://yameng.remoon.cn/_app/immutable/assets/left-girl.384ee8d5.png]"

func TestParseBBCode(t *testing.T) {
	r := msg.ParseString(testMsgRaw)
	t.Log(r)
}
