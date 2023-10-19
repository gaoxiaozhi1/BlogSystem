package message_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/utils/jwts"
	"time"
)

type Message struct {
	SendUserID       uint      `json:"send_user_id"`       // 发送人id
	SendUserNickName string    `json:"send_user_nickname"` // 发送人昵称
	SendUserAvatar   string    `json:"send_user_avatar"`   // 发送人头像
	RevUserID        uint      `json:"rev_user_id"`        // 接收人id
	RevUserNickName  string    `json:"rev_user_nickname"`  // 接收人昵称
	RevUserAvatar    string    `json:"rev_user_avatar"`    // 接收人头像
	Content          string    `json:"content"`
	CreatedAt        time.Time `json:"created_at"`    // 最新的消息时间
	MessageCount     int       `json:"message_count"` // 消息条数
}

type MessageGroup map[uint]*Message

// MessageListView 类似于QQ聊天页面，展示消息数量和最新的一条消息，没有分页
func (MessageApi) MessageListView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var messageGroup = MessageGroup{}
	var messageList []models.MessageModel
	var messages []Message

	// 按升序找到我和其他人的聊天记录，其中我和某个人的聊天记录，对应的id和肯定相同，以此区分我与不同人的聊天
	global.DB.Order("created_at asc").
		Find(&messageList, "send_user_id = ? or rev_user_id = ?", claims.UserID, claims.UserID)

	for _, model := range messageList {
		message := Message{
			SendUserID:       model.SendUserID,
			SendUserNickName: model.SendUserNickName,
			SendUserAvatar:   model.SendUserAvatar,
			RevUserID:        model.RevUserID,
			RevUserNickName:  model.RevUserNickName,
			RevUserAvatar:    model.RevUserAvatar,
			Content:          model.Content,
			CreatedAt:        model.CreatedAt,
			MessageCount:     1,
		}
		idNum := model.SendUserID + model.RevUserID
		val, ok := messageGroup[idNum]
		if !ok {
			// 不存在
			messageGroup[idNum] = &message
			continue
		}
		// 存在
		message.MessageCount = val.MessageCount + 1
		messageGroup[idNum] = &message
	}

	for _, model := range messageGroup {
		messages = append(messages, *model)
	}
	res.OKWithData(messages, c)
	return
}
