package models

import (
	"fmt"

	"gorm.io/gorm"
)

// 文章表
type Article struct {
	BaseModel
	PageLimitStruct
	Title         string `gorm:"type:varchar(30)"`  // 标题
	Content       string `gorm:"type:text"`         // 文章原始内容数据
	ContentRender string `gorm:"type:text"`         // 文章渲染数据
	Visit         int    `gorm:"type:int(11)"`      // 访问次数
	Status        int    `gorm:"type:tinyint(1)"`   // 状态:0.私有 1.公开
	Description   string `gorm:"type:varchar(300)"` // 标题
	Editor        int    `gorm:"type:tinyint(1)"`   // 编辑器类型
	AutoRelease   int    `gorm:"type:tinyint(1)"`   // 自动发布 0.否 1.是
	UserId        uint   `gorm:""`                  // 用户id

	// 关联
	User User
	// 关系
	Tags       []Tag       `gorm:"many2many:article_tags"`
	Anthologys []Anthology `gorm:"many2many:article_anthologys"`
	// TagRelArticle []TagRelArticle `gorm:"-"`
}

//
func (m *Article) Group() *Article {
	return m
}

// 按页获取文章列表
func (m *Article) GetList(page, limit int) ([]Article, int64) {

	db := Db.Order("updated_at Desc")
	var count int64

	offset, limit := calcPage(page, limit)

	articleList := []Article{}

	// err := db.Debug().Table("article").Joins("LEFT JOIN `blog_user` ON blog_article.u_id=blog_user.id").Offset(offset).Limit(limit).Find(&articleList).Offset(-1).Limit(-1).Count(&count).Error
	err := db.Debug().Model(&Article{}).Preload("User").Preload("Tags").Offset(offset).Limit(limit).Find(&articleList).Offset(-1).Limit(-1).Count(&count).Error

	// fmt.Println(articleList)
	if err != nil {
		return nil, 0
	} else {
		return articleList, count
	}
}

// 根据专栏获取文章
func (m *Article) GetListByAnthologyId(page, limit int, AnthologyId uint) (articleList []Article, count int64) {
	offset, limit := calcPage(page, limit)
	mAnthology := Anthology{}
	mAnthology.ID = AnthologyId
	err := Db.Model(&mAnthology).
		Preload("User").Preload("Tags").
		Offset(offset).Limit(limit).
		Association("Articles").
		Find(&articleList)
		// Offset(-1).Limit(-1).
		// Count(&count)
	count = Db.Model(&mAnthology).
		Preload("User").
		Association("Articles").Count()
	_ = err
	return
}

// 根标签获取文章
func (m *Article) GetListByTagId(page, limit int, tagId uint) (articleList []Article, count int64) {
	offset, limit := calcPage(page, limit)
	mTag := Tag{}
	mTag.ID = tagId
	err := Db.Model(&mTag).
		Preload("User").Preload("Tags").
		Offset(offset).Limit(limit).
		Association("Articles").
		Find(&articleList)
		// Offset(-1).Limit(-1).
		// Count(&count)
	count = Db.Model(&mTag).
		Preload("User").
		Association("Articles").Count()
	_ = err
	return
}

// 按页获取文章列表
func (m *Article) GetListByCondition(page, limit int, keyword string, AnthologyIds []uint, article Article) (articleList []Article, count int64) {
	db := Db.Debug().Order("updated_at Desc")
	if article.UserId != 0 {
		db = db.Where("user_id=?", article.UserId)
	}

	if len(AnthologyIds) != 0 {
		db = db.Preload("Article.Anthologys", func(db *gorm.DB) *gorm.DB {
			return db.Where("id in ?", AnthologyIds)
		})
		// db = db.Joins("Anthologys", "id in ?", AnthologyIds)
	}

	if keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	offset, limit := calcPage(page, limit)
	err := db.Model(&Article{}).Offset(offset).Limit(limit).Find(&articleList).Offset(-1).Limit(-1).Count(&count).Error
	fmt.Println("查询错误", err)
	return
	// fmt.Println("查询错误", err)
	// if err != nil {
	// 	fmt.Println("22")
	// 	return nil, 0
	// } else {
	// 	fmt.Println("11---", count)
	// 	return
	// }
}

// 按分类和用户id获取文章列表
func (m *Article) GetListByAnthologyIdAndUserId(page, limit int, keyword string, AnthologyId uint, user_id uint) (Articles []Article, count int64) {
	db := Db.Order("updated_at Desc")
	mAnthology := Anthology{}
	mAnthology.ID = AnthologyId
	offset, limit := calcPage(page, limit)
	if keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if user_id != 0 {
		db = db.Where("user_id=?", user_id)
	}

	var (
		err error
	)
	if AnthologyId == 0 {
		db.Offset(offset).Limit(limit)
		db.Find(&Articles).Offset(-1).Limit(-1).Count(&count)
	} else {
		err = db.Model(&mAnthology).Offset(offset).Limit(limit).Association("Articles").Find(&Articles)
		count = db.Model(&mAnthology).Association("Articles").Count()
	}

	if err != nil {
		return nil, 0
	} else {
		return
	}
}

func (m *Article) GetContentById(id int) {
	articleList := Article{}
	Db.Preload("TagRelArticle").Where("id=?", id).Find(&articleList)
	fmt.Println(articleList)
}

// 获取用户的文章列表
func (m *Article) GetListByUserIdAndPage(page, limit, user_id int) (articleList []Article, count int64, err error) {
	offset, limit := calcPage(page, limit)
	err = Db.Preload("Tags").Where("user_id=?", user_id).Offset(offset).Limit(limit).Find(&articleList).Offset(-1).Limit(-1).Count(&count).Error
	return
}

// 获取专栏下的所有文章
func (m *Article) GetListByGroupId(user_id int) {

}

// 获取详情和标签
func (m *Article) GetInfoAndTag(id uint) (article Article, err error) {
	// tagList := []Tag{}
	// article := Article{}
	article.ID = id
	// article.Status = 2
	// err = Db.Model(&article).Association("Tags").Find(&tagList)
	// err = Db.Model(&article).Association("Tags").Find(&tagList)
	err = Db.Debug().Preload("Tags").Preload("User").Where("status=1").First(&article).Error
	return
}

// 获取详情和标签
func (m *Article) GetInfo(id uint) (article Article, err error) {
	err = Db.Debug().Where("id=? ", id).First(&article).Error
	return
}

// 获取文章配置
func (m *Article) GetConfig(article_id uint) (article_config Article, err error) {
	mArticle := Article{}
	mArticle.ID = article_id
	tags := []Tag{}
	Anthologys := []Anthology{}
	err = Db.Model(&mArticle).Association("Tags").Find(&tags)
	err = Db.Model(&mArticle).Association("Anthologys").Find(&Anthologys)
	article_config.Tags = tags
	article_config.Anthologys = Anthologys
	return
}

// 添加一个
func (m *Article) AddOne(article Article) (info Article, err error) {
	err = Db.Create(&article).Error
	// // 添加关联
	// Db.Debug().Model(&article).Association("Tags").Append(tags)
	// Db.Debug().Model(&article).Association("Anthologys").Append(anthologys)
	return article, err
}

// 更新一个文章
// 参数-updateData 支持：title,content,content_render,visit,status,description,editor,auto_release
func (m *Article) UpdateByArticleId(article_id uint, updateData map[string]interface{}) (err error) {

	sqlUpdateData := map[string]interface{}{}

	sqlUpdateData["title"] = updateData["title"]
	sqlUpdateData["content"] = updateData["content"]
	sqlUpdateData["content_render"] = updateData["content_render"]
	sqlUpdateData["visit"] = updateData["visit"]
	sqlUpdateData["status"] = updateData["status"]
	sqlUpdateData["description"] = updateData["description"]
	sqlUpdateData["editor"] = updateData["editor"]
	sqlUpdateData["auto_release"] = updateData["auto_release"]

	db := Db.Debug().Model(&Article{}).Where("id=?", article_id)
	// 添加关联
	article := Article{}
	article.ID = article_id
	m.ClearTagAll(article_id)
	m.ClearAnthologyAll(article_id)

	tags, _ := updateData["tags"].([]Tag)
	anthologys, _ := updateData["anthologys"].([]Anthology)

	Db.Debug().Model(&article).Association("Tags").Append(tags)
	Db.Debug().Model(&article).Association("Anthologys").Append(anthologys)
	err = db.Updates(sqlUpdateData).Error
	return err
}

// 清空全部标签的关联
func (m *Article) ClearTagAll(article_id uint) (info Article, err error) {
	article := Article{}
	article.ID = article_id
	err = Db.Model(&article).Association("Tags").Clear()
	return article, err
}

// 清空全部专栏关联
func (m *Article) ClearAnthologyAll(article_id uint) (info Article, err error) {
	article := Article{}
	article.ID = article_id
	err = Db.Model(&article).Association("Anthologys").Clear()
	return article, err
}

// 根据主键删除文章
func (m *Article) DeleteByArticleId(article_id uint) (err error) {
	err = Db.Delete(&Article{}, article_id).Error
	return err
}
