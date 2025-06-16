package db

import "time"

type SearchSetting struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	SearchOn bool      `json:"searchOn"`
	AddNew   bool      `json:"addNew"`
	Amount   uint      `json:"amount"`
	UpdateAt time.Time `json:"updatedAt"`
}

func (s *SearchSetting) Get() error {
	err := DBConn.Where("id = 1").First(s).Error
	return err
}
func (s *SearchSetting) Update() error {
	tx := DBConn.Select("search_on", "add_new", "amount", "updatedAt").Where("id = 1").Updates(s)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
