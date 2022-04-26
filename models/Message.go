package models

import (
	"encoding/json"
	"enian_blog/lib/cmn"
)

// 消息
type Message struct {
	BaseModel
	Title       string `gorm:"type:varchar(50)"`   // 标题
	Read        int    `gorm:"type:tinyint(1)"`    // 是否已读:0.否 1.是
	Content     string `gorm:"type:varchar(2000)"` // 消息内容
	MessageType int    `gorm:"type:tinyint(1)"`    // 消息类型:1.通知 2.审核 3.私信
	FromUserId  uint   `gorm:"type:int(11)"`       // 消息来自用户。发信人（0 代表系统消息）
	ExtendParam string `gorm:"type:varchar(2000)"` // 扩展参数，json格式
	ToUserId    uint   `gorm:"type:int(11)"`       // 收信人
	// TagRelArticleId uint

	// 关联模型
	// User User
}

// 创建一条消息
func (m *Message) CreateOneMessage(title, content string, ExtendParam cmn.Msi, MessageType int, toUserId, fromUserId uint) (Message, error) {
	paramMsi, _ := json.Marshal(ExtendParam)
	// if jsonErr != nil {
	// 	fmt.Println("json 解析 错误")
	// }
	message := Message{
		Title:       title,
		Content:     content,
		MessageType: MessageType,
		ToUserId:    toUserId,
		FromUserId:  fromUserId,
		ExtendParam: string(paramMsi),
		Read:        0,
	}

	err := Db.Create(&message).Error
	return message, err
}

// 获取列表
// from_user_id message_type to_user_id
func (m *Message) GetListByCondition(page, limit int, condition cmn.Msi) (messageList []Message, count int64, err error) {
	// mAnthology := Message{}
	offset, limit := calcPage(page, limit)
	db := Db.Model(&Message{}).Order("created_at DESC")

	if v, ok := condition["from_user_id"]; ok {
		db = db.Where("from_user_id =?", v)
	}

	if v, ok := condition["to_user_id"]; ok {
		db = db.Where("to_user_id =?", v)
	}

	if v, ok := condition["message_type"]; ok {
		db = db.Where("message_type =?", v)
	}

	db.Offset(offset).Limit(limit)
	err = db.Find(&messageList).Offset(-1).Limit(-1).Count(&count).Error

	return
}

// 获取消息详情
func (m *Message) GetInfoById(id uint) (info Message, extend_param cmn.Msi, err error) {
	err = Db.Where("id=?", id).First(&info).Error
	err = json.Unmarshal([]byte(info.ExtendParam), &extend_param)
	return
}

// 修改消息
// key 可选：read extend_param
func (m *Message) Update(id uint, msi cmn.Msi) (info Message, err error) {
	update := map[string]interface{}{}
	if v, ok := msi["read"]; ok {
		update["read"] = v
	}

	if v, ok := msi["extend_param"]; ok {
		paramMsi, _ := json.Marshal(v)
		update["extend_param"] = string(paramMsi)
	}

	err = Db.Model(&Message{}).Where("id=?", id).Updates(&update).Error
	return
}
