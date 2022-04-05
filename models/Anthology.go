package models

import (
	"enian_blog/lib/cmn"
	"fmt"
)

// 专栏
type Anthology struct {
	BaseModel
	Title string `gorm:"type:varchar(20)" json:"title"` // 标题
	// Title_en          string `gorm:"type:varchar(50)"` // 标题英文标识
	Order          int    `gorm:"type:int(11)" json:"order"`             // 排序
	Golbal_open    int    `gorm:"type:tinyint(1)" json:"golbal_open"`    // 公开到全局：1.公开（管理员可以添加专栏展示在团队首页） 2.不公开（专栏仅显示在个人首页）
	Accept_article int    `gorm:"type:tinyint(1)" json:"accept_article"` // 接收文章：1.接受 2. 不接受（新文章不可再向此专栏推送，已推送的不受影响） 3.需要审核
	UserId         uint   `gorm:"type:int(11)" json:"user_id"`           // 用户id
	Description    string `gorm:"type:varchar(300)" json:"description"`  // 描述
	// AnthologyRelArticleId uint

	// 关系模型
	User     User
	Articles []Article `gorm:"many2many:article_anthologys"`
}

// 首页专栏(用户首页，站首页)
// type HomeAnthology struct {
// 	BaseModel
// 	Order       int `gorm:"type:int(11)" json:"order"` // 排序
// 	UserId      uint
// 	AnthologyID uint
// }

// // 专栏关联文章
// type AnthologyRelArticle struct {
// 	Model
// 	TagID       uint
// 	AnthologyID uint
// 	Tag         []Tag       `gorm:"-"`
// 	Anthology   []Anthology `gorm:"-"`
// }

func (m *Anthology) GetMoreById(ids string) (list []Anthology, err error) {
	err = Db.Where("id in (" + ids + ")").Find(&list).Error
	return
}

// 获取详情根据一个
func (m *Anthology) GetInfoById(id uint) (info Anthology, err error) {
	err = Db.Preload("User").Where("id = ?", id).First(&info).Error
	return
}

// 获取专栏列表
// 支持条件 uer_id golbal_open Accept_article exclude_user_id
func (m *Anthology) GetList(condition cmn.Msi) (list []Anthology, err error) {
	db := Db.Model(&Anthology{}).Order("updated_at DESC")
	if v, ok := condition["user_id"]; ok {
		db = db.Where("user_id=?", v)
	}
	if v, ok := condition["exclude_user_id"]; ok {
		db = db.Where("user_id not in ?", v)
	}

	if v, ok := condition["golbal_open"]; ok {
		db = db.Where("golbal_open=?", v)
	}
	if v, ok := condition["accept_article"]; ok {
		db = db.Where("accept_article in ?", v)
	}

	err = db.Preload("User").Find(&list).Error
	return
}

// 获取专栏列表
// 支持条件 uer_id golbal_open Accept_article
// func (m *Anthology) GetList(condition Anthology) (list []Anthology, err error) {
// 	db := Db.Debug().Model(&Anthology{}).Order("updated_at DESC")
// 	if condition.UserId != 0 {
// 		db = db.Where("user_id=?", condition.UserId)
// 	}
// 	if condition.Golbal_open != 0 {
// 		db = db.Where("golbal_open=?", condition.Golbal_open)
// 	}
// 	if condition.Accept_article != 0 {
// 		db = db.Where("accept_article=?", condition.Accept_article)
// 	}

// 	err = db.Preload("User").Find(&list).Error
// 	return
// }

// 获取专栏列表
// 支持条件 uer_id golbal_open Accept_article
func (m *Anthology) GetListByIds(ids string) (list []Anthology, err error) {

	err = Db.Debug().Preload("User").Where("id in (" + ids + ")").Order("field(id ," + ids + ")").Find(&list).Error
	return
}

// 修改一个专栏
func (m *Anthology) Edit(newAnthology Anthology) (returnData Anthology, err error) {
	if newAnthology.ID == 0 {
		err = Db.Create(&newAnthology).Error
	} else {
		fmt.Println("更新记录", newAnthology)
		err = Db.Model(Anthology{}).Where("id=?", newAnthology.ID).Updates(map[string]interface{}{
			"user_id":        newAnthology.UserId,
			"title":          newAnthology.Title,
			"description":    newAnthology.Description,
			"accept_article": newAnthology.Accept_article,
			"golbal_open":    newAnthology.Golbal_open,
		}).Error
	}
	returnData = newAnthology
	return returnData, err
}

// 删除专栏用id
func (m *Anthology) DeleteById(ids []int) (err error) {
	err = Db.Delete(&Anthology{}, ids).Error
	return
}