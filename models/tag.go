package models

import (
	"enian_blog/lib/cmn"
)

// 标签
type Tag struct {
	BaseModel
	Title  string `gorm:"type:varchar(10)"` // 标题
	Status int    `gorm:"type:tinyint(1)"`  // 隐藏:0.隐藏 1.显示
	UserId uint
	// TagRelArticleId uint

	// 关联模型
	User     User
	Articles []*Article `gorm:"many2many:article_tags"`
	// TagRelArticle []TagRelArticle `gorm:"-"`
}

func (m *Tag) GetAll() (list []Tag, err error) {
	err = Db.Where("status=1").Find(&list).Error
	return
}

// func (m *Tag) GetMoreById(ids string) (list []Tag, err error) {
// 	err = Db.Where("status=1").Find(&list).Error
// 	return
// }

func (m *Tag) GetOneById(id uint) (info Tag, err error) {
	err = Db.Where("status=1").Where("id=?", id).Find(&info).Error
	return
}

// 条件 title（模糊）
func (m *Tag) GetList(condition cmn.Msi) (list []Tag, err error) {
	db := Db.Model(&Tag{})
	if v, ok := condition["title"]; ok {
		db = db.Where("title like ?", "%"+cmn.InterfaceToString(v)+"%")
	}
	err = db.Where("status=1").Find(&list).Error
	return
}

// 创建标签
func (m *Tag) CreateOne(title string, user_id uint) (info Tag, err error) {
	newTag := Tag{
		Title:  title,
		UserId: user_id,
		Status: 1,
	}
	err = Db.Create(&newTag).Error
	info = newTag
	return
}

// 获取一个按标题（精准）
func (m *Tag) GetOneByTitle(title string) (info Tag, err error) {
	err = Db.Where("title=?", title).First(&info).Error
	return info, err
}

// // 获取文章下的标签
// func (m *Tag) GetListByArticleId(articleId uint) (list []TagRelArticle, err error) {
// 	err = Db.Model(&TagRelArticle{}).Preload("Tag").Where("article_id=?", articleId).Find(&list).Error
// 	return
// }
