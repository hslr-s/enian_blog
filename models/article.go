package models

import (
	"enian_blog/lib/cmn"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 文章表
type Article struct {
	BaseModel
	PageLimitStruct
	Title         string `gorm:"type:varchar(100)"` // 标题
	Content       string `gorm:"type:longtext"`     // 文章原始内容数据
	ContentRender string `gorm:"type:longtext"`     // 发布的渲染内容
	Visit         int    `gorm:"type:int(11)"`      // 访问次数
	Status        int    `gorm:"type:tinyint(1)"`   // 状态:0.私有 1.公开
	Description   string `gorm:"type:varchar(300)"` // 描述
	Editor        int    `gorm:"type:tinyint(1)"`   // 编辑器类型
	// AutoRelease      int       `gorm:"type:tinyint(1)"`   // 自动发布 0.否 1.是
	UserId            uint      `gorm:""`              // 用户id
	SaveTime          time.Time `gorm:"type:datetime"` // 最后保存时间
	ReleaseUpdateTime time.Time `gorm:"type:datetime"` // 发布更新时间
	ReleaseTime       time.Time `gorm:"type:datetime"` // 首次发布时间

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

// 根据专栏获取文章
// article_condition: user_id keyword only_release status
func (m *Article) GetListByAnthologyId(page, limit int, AnthologyId uint, article_condition cmn.Mss) (articleList []Article, count int64, err error) {
	offset, limit := calcPage(page, limit)
	mAnthology := Anthology{}
	mAnthology.ID = AnthologyId
	db := Db.Model(&mAnthology)
	db_count := Db.Model(&mAnthology)
	if v, ok := article_condition["user_id"]; ok {
		db = db.Where("article.user_id=?", v)
		db_count = db_count.Where("article.user_id=?", v)
	}

	if v, ok := article_condition["keyword"]; ok {
		db = db.Where("article.title LIKE ? OR content LIKE ?", "%"+v+"%", "%"+v+"%")
		db_count = db_count.Where("article.title LIKE ? OR content LIKE ?", "%"+v+"%", "%"+v+"%")
	}

	// 仅发布
	if _, ok := article_condition["only_release"]; ok {
		db = db.Where("article.release_time is not null")
		db_count = db_count.Where("article.release_time is not null")
	}

	// 公开状态
	if v, ok := article_condition["status"]; ok {
		db = db.Where("article.status= ?", v)
		db_count = db_count.Where("article.status= ?", v)
	}

	err = db.Preload("User").Preload("Tags").
		Where("article.release_time is not null").
		Offset(offset).Limit(limit).
		Association("Articles").
		Find(&articleList)
	count = db_count.Preload("User").Preload("Tags").
		Association("Articles").Count()
	return
}

// 根标签获取文章
func (m *Article) GetListByTagId(page, limit int, tagId uint, article_condition cmn.Mss) (articleList []Article, count int64, err error) {
	offset, limit := calcPage(page, limit)
	mTag := Tag{}
	mTag.ID = tagId
	db := Db.Model(&mTag)
	db_count := Db.Model(&mTag)
	if v, ok := article_condition["user_id"]; ok {
		db = db.Where("article.user_id=?", v)
		db_count = db_count.Where("article.user_id=?", v)
	}

	if v, ok := article_condition["keyword"]; ok {
		db = db.Where("article.title LIKE ? OR content LIKE ?", "%"+v+"%", "%"+v+"%")
		db_count = db_count.Where("article.title LIKE ? OR content LIKE ?", "%"+v+"%", "%"+v+"%")
	}

	// 仅发布
	if _, ok := article_condition["only_release"]; ok {
		db = db.Where("article.release_time is not null")
		db_count = db_count.Where("article.release_time is not null")
	}

	// 公开状态
	if v, ok := article_condition["status"]; ok {
		db = db.Where("article.status= ?", v)
		db_count = db_count.Where("article.status= ?", v)
	}
	err = db.Preload("User").Preload("Tags").
		Offset(offset).Limit(limit).
		Association("Articles").
		Find(&articleList)

	count = db_count.Preload("User").Preload("Tags").
		Association("Articles").Count()
	return
}

// 按页获取文章列表(筛选条件)
// user_id keyword only_release status
func (m *Article) GetListByCondition(page, limit int, condition cmn.Mss, order string) (articleList []Article, count int64, err error) {
	db := Db.Order(order).Preload("Tags").Preload("User")
	// if v, ok := condition["anthologys"]; ok {
	// 	db = db.Debug().Preload("Anthologys", func(db *gorm.DB) *gorm.DB {
	// 		return db.Where("id in (?)", v)
	// 	})
	// }

	if v, ok := condition["user_id"]; ok {
		db = db.Where("user_id=?", v)
	}

	if v, ok := condition["keyword"]; ok {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+v+"%", "%"+v+"%")
	}

	// 仅发布
	if _, ok := condition["only_release"]; ok {
		db = db.Where("release_time is not null")
	}

	// 公开状态
	if v, ok := condition["status"]; ok {
		db = db.Where("status= ?", v)
	}

	offset, limit := calcPage(page, limit)
	err = db.Model(&Article{}).Offset(offset).Limit(limit).Find(&articleList).Offset(-1).Limit(-1).Count(&count).Error
	return
}

// 按分类和用户id获取文章列表
func (m *Article) GetListByAnthologyIdAndUserId(page, limit int, keyword string, AnthologyId uint, user_id uint) (Articles []Article, count int64) {
	db := Db.Order("save_time DESC,created_at DESC")
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
		db.Preload("User").Offset(offset).Limit(limit)
		db.Preload("User").Find(&Articles).Offset(-1).Limit(-1).Count(&count)
	} else {
		err = db.Model(&mAnthology).Preload("User").Offset(offset).Limit(limit).Association("Articles").Find(&Articles)
		count = db.Model(&mAnthology).Preload("User").Association("Articles").Count()
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
	err = Db.Preload("Tags").Preload("User").Where("status=1").First(&article).Error
	return
}

// 获取详情和标签
func (m *Article) GetInfo(id uint) (article Article, err error) {
	err = Db.Where("id=? ", id).First(&article).Error
	return
}

// 获取文章配置
func (m *Article) GetConfig(article_id uint) (article_config Article, err error) {
	mArticle := Article{}
	mArticle.ID = article_id
	tags := []Tag{}
	anthologys := []Anthology{}
	err = Db.Model(&mArticle).Association("Tags").Find(&tags)
	err = Db.Model(&mArticle).Association("Anthologys").Find(&anthologys)

	article_config.Tags = tags
	article_config.Anthologys = anthologys
	return
}

// 添加一个
func (m *Article) AddOne(article Article) (info Article, err error) {
	if article.ReleaseTime.IsZero() {
		err = Db.Omit("ReleaseUpdateTime,ReleaseTime").Create(&article).Error
	} else {
		err = Db.Create(&article).Error
	}
	return article, err
}

// 更新一个文章
func (m *Article) UpdateByArticleId(article_id uint, updateData map[string]interface{}) (err error) {
	// 添加关联
	article := Article{}
	article.ID = article_id
	m.ClearTagAll(article_id)
	m.ClearAnthologyAll(article_id)

	tags, _ := updateData["tags"].([]Tag)
	anthologys, _ := updateData["anthologys"].([]Anthology)
	// 删除原始结构
	delete(updateData, "tags")
	delete(updateData, "anthologys")

	Db.Model(&article).Association("Tags").Append(tags)
	Db.Model(&article).Association("Anthologys").Append(anthologys)

	err = Db.Model(&Article{}).Where("id=?", article_id).Updates(updateData).Error
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

func (m *Article) AddAnthology(article_id uint, anthology_id uint) (err error) {
	article := Article{}
	anthology := Anthology{}
	article.ID = article_id
	anthology.ID = anthology_id
	err = Db.Debug().Model(&article).Association("Anthologys").Append(&anthology)
	return
}

// 访问增加X
func (m *Article) VisitSetInc(article_id uint, num int) error {
	if err := Db.Model(&Article{}).Where("id=?", article_id).Update("visit", gorm.Expr("visit+?", num)).Error; err != nil {
		return err
	}
	return nil
}
