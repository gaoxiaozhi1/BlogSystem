package message_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

type MessageRequest struct {
	SendUserID uint   `json:"send_user_id" binding:"required"` // 发送人id
	RecvUserID uint   `json:"recv_user_id" binding:"required"` // 接收人id
	Content    string `json:"content" binding:"required"`      // 消息内容
}

// MessageCreateView 当前用户发送消息
func (MessageApi) MessageCreateView(c *gin.Context) {
	var cr MessageRequest
	// SendUserID就是当前登录人的id
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	var sendUser, recvUser models.UserModel
	err = global.DB.Take(&sendUser, cr.SendUserID).Error
	if err != nil {
		res.FailWithMessage("发送人信息不存在", c)
		return
	}
	err = global.DB.Take(&recvUser, cr.RecvUserID).Error
	if err != nil {
		res.FailWithMessage("接收人信息不存在", c)
		return
	}

	err = global.DB.Create(&models.MessageModel{
		SendUserID:       cr.SendUserID,
		SendUserNickName: sendUser.NickName,
		SendUserAvatar:   sendUser.Avatar,
		RevUserID:        cr.RecvUserID,
		RevUserNickName:  recvUser.NickName,
		RevUserAvatar:    recvUser.Avatar,
		IsRead:           false,
		Content:          cr.Content,
	}).Error

	if err != nil {
		res.FailWithMessage("消息发送失败", c)
		return
	}
	res.OKWithMessage("消息发送成功", c)
	return
}
