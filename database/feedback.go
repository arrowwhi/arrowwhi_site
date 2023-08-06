package database

type FeedbackType struct {
	ID   uint   `gorm:"primaryKey;type:smallserial"`
	Type string `gorm:"type:varchar"`
}

type Feedback struct {
	ID          uint   `gorm:"primaryKey"`
	TypeID      uint   `gorm:"not null;ForeignKey:FeedbackTypeID"`
	Description string `gorm:"not null"`
	Done        bool   `gorm:"default:false"`
	Trashed     bool   `gorm:"default:false"`
}

func (dbe *DbEngine) AddFeedback(input *Feedback) error {
	if err := dbe.db.Create(input).Error; err != nil {
		return err
	}
	return nil
}
