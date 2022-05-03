package models

import (
	"enian_blog/lib/cmn"
)

// 专栏
type Anthology struct {
	BaseModel
	Title          string `gorm:"type:varchar(20)" json:"title"`         // 标题
	Golbal_open    int    `gorm:"type:tinyint(1)" json:"golbal_open"`    // 公开到全局：1.公开（管理员可以添加专栏展示在团队首页） 2.不公开（专栏仅显示在个人首页）
	Accept_article int    `gorm:"type:tinyint(1)" json:"accept_article"` // 接收文章：1.接受 2. 不接受（新文章不可再向此专栏推送，已推送的不受影响） 3.需要审核
	UserId         uint   `gorm:"type:int(11)" json:"user_id"`           // 用户id
	Description    string `gorm:"type:varchar(300)" json:"description"`  // 描述
	// AnthologyRelArticleId uint

	// 关系模型
	User     User
	Articles []Article `gorm:"many2many:article_anthologys"`
}

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
func (m *Anthology) GetListByIds(ids string) (list []Anthology, err error) {
	if ids != "" {
		err = Db.Preload("User").Where("id in (" + ids + ")").Order("field(id ," + ids + ")").Find(&list).Error
	}
	return
}

// 获取专栏列表
func (m *Anthology) GetListByIdsUint(ids []uint) (list []Anthology, err error) {
	err = Db.Preload("User").Where("id in ?", ids).Find(&list).Error
	return
}

// 修改一个专栏
func (m *Anthology) Edit(newAnthology Anthology) (returnData Anthology, err error) {
	if newAnthology.ID == 0 {
		err = Db.Create(&newAnthology).Error
	} else {
		// fmt.Println("更新记录", newAnthology)
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

// 查找关联
func (m *Anthology) FindRelation(anthologyId, articleId uint) bool {
	mAntology := Anthology{}
	mAntology.ID = anthologyId
	count := Db.Model(&mAntology).Where("article_id = ?", articleId).Association("Articles").Count()
	if count > 0 {
		return true
	} else {
		return false
	}
}
