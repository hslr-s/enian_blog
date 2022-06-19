package models

// 友情链接
type FriendLink struct {
	BaseModel
	Link  string `gorm:"type:varchar(2000)"` // 链接地址
	Title string `gorm:"type:varchar(100)"`  // 标题
	Sort  int    `gorm:"type:int(11)"`       // 排序数字
}

// 添加
func (m *FriendLink) AddOne(title, link string, sort int) (FriendLink, error) {
	f := FriendLink{
		Title: title,
		Link:  link,
		Sort:  sort,
	}
	err := Db.Create(&f).Error
	return f, err
}

// 修改
func (m *FriendLink) UpdateOne(id uint, title, link string, sort int) (FriendLink, error) {
	f := FriendLink{
		Title: title,
		Link:  link,
		Sort:  sort,
	}
	err := Db.Where("id", id).Updates(&f).Error
	return f, err
}

// 删除
func (m *FriendLink) DeleteOne(id uint) error {
	err := Db.Delete(m, id).Error
	return err
}

// 查列表
func (m *FriendLink) GetList(sortAsc bool) ([]FriendLink, error) {
	list := []FriendLink{}
	db := Db.Model(m)
	if sortAsc {
		db = db.Order("sort asc,id asc")
	} else {
		db = db.Order("sort Desc,id asc")
	}
	err := db.Find(&list).Error
	return list, err
}
