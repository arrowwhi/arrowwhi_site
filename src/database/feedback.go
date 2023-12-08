package database

type FeedbackType struct {
	ID   uint   `gorm:"primaryKey;type:smallserial"`
	Type string `gorm:"type:varchar"`
}

type Feedback struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	TypeID      uint   `gorm:"not null;ForeignKey:FeedbackTypeID" json:"type_id"`
	Description string `gorm:"not null" json:"description"`
	Done        bool   `gorm:"default:false" json:"done"`
	Trashed     bool   `gorm:"default:false" json:"trashed"`
}

func (dbe *DbEngine) AddFeedback(input *Feedback) error {
	if err := dbe.db.Create(input).Error; err != nil {
		return err
	}
	return nil
}
