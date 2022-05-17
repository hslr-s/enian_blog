package models

// 内容排序,首页专栏，用户首页专栏，置顶文章等
type ContentSort struct {
	ID        uint
	ContentID uint `gorm:"type:bigint(20)"` // 内容id，专栏id，文章id...
	SortID    uint `gorm:"type:bigint(20)"` // 排序id
	TypeID    uint `gorm:"type:int(11)"`    // 类型：1.首页专栏 2.用户专栏排序
}

// 获取首页转栏的列表
func (m *ContentSort) GetHomeAnthologyList() (anthologyList []Anthology) {
	Db.Joins("INNER JOIN content_sort ON content_sort.content_id=anthology.id").
		Where("content_sort.type_id=?", 1).
		Order("content_sort.sort_id").
		Find(&anthologyList)
	return
}

// 获取首页转栏的列表
// func (m *ContentSort) SetHomeAnthologyList() (anthologyList []Anthology) {

// }
