package msgdao

import (
	"strconv"
	"testing"
	"time"
)

func TestChatMsgDao_GetRecentChatMessagesBySessionID(t *testing.T) {
	for i := 0; i < 100; i++ {
		ms, err := instance.GetChatMessagesBySession(2, 1, 0, 10)
		if err != nil {
			t.Error(err)
		}
		if len(ms) == 0 {
			break
		}
		for _, m := range ms {
			t.Log(m)
		}
	}
}

func TestChatMsgDao_GetRecentChatMessages(t *testing.T) {
	messages, err := instance.GetRecentChatMessages(1, 1637650186)
	if err != nil {
		t.Error(err)
	}
	for _, message := range messages {
		t.Log(message)
	}
}

func TestChatMsgDao_GetOfflineMessage(t *testing.T) {
	m, err := GetOfflineMessage(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(m)
}

func TestChatMsgDao_AddOfflineMessage(t *testing.T) {
	err := AddOfflineMessage(1, 4)
	if err != nil {
		t.Error(err)
	}
}

func TestChatMsgDao_DelOfflineMessage(t *testing.T) {
	err := DelOfflineMessage(1, []int64{1, 2, 3, 4})
	if err != nil {
		t.Error(err)
	}
}

func TestAddChatMessage(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, err := AddChatMessage(&ChatMessage{
			MID:       int64(2000 + i),
			SessionID: "2_1",
			CliSeq:    int64(i),
			From:      1,
			To:        2,
			Type:      1,
			SendAt:    time.Now().Unix() - int64(60*i),
			CreateAt:  time.Now().Unix(),
			Content:   "Hello-" + strconv.FormatInt(int64(i*100), 10),
		})
		if err != nil {
			t.Error(err)
		}
	}
}

func TestChatMsgDao_AddOrUpdateChatMessage(t *testing.T) {

	message, err := AddChatMessage(&ChatMessage{
		MID:       15,
		SessionID: "2_1",
		CliSeq:    2,
		From:      2,
		To:        1,
		Type:      1,
		SendAt:    time.Now().Unix(),
		CreateAt:  time.Now().Unix(),
		Content:   "hello",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func TestChatMsgDao_GetChatMessageSeqAfter(t *testing.T) {
	after, err := GetChatMessageMidAfter(1, 2, 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(after)
}
