package models

// 文件管理
type File struct {
	BaseModel
	Name      string `gorm:"type:varchar(255)" json:"name"`  // 文件名
	Path      string `gorm:"type:varchar(255)" json:"path"`  // 路径
	Ext       string `gorm:"type:varchar(50)" json:"ext"`    // 扩展名
	Type      int    `gorm:"type:tinyint(1)" json:"type"`    // 文件类型 1.图片类 2.压缩包
	UserId    uint   `gorm:"type:int(11)" json:"user_id"`    // 用户id
	ArticleId uint   `gorm:"type:int(11)" json:"article_id"` // 文章id
	// 关系模型
	User     User
	Articles []Article `gorm:"many2many:article_anthologys"`
}
