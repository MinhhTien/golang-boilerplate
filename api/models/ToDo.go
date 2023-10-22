package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type ToDo struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:100;not null;unique" json:"title"`
	Content   string    `gorm:"text;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null;" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (t *ToDo) Prepare() {
	t.Title = html.EscapeString(strings.TrimSpace(t.Title))
	t.Content = html.EscapeString(strings.TrimSpace(t.Content))
	t.Author = User{}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *ToDo) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)
	if t.Title == "" {
		err = errors.New("Required Title")
		errorMessages["Required_title"] = err.Error()
	}
	if t.Content == "" {
		err = errors.New("Required Content")
		errorMessages["Required_content"] = err.Error()
	}
	if t.AuthorID < 1 {
		err = errors.New("Required Author")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}

func (t *ToDo) SaveToDo(db *gorm.DB) (*ToDo, error) {
	var err error
	err = db.Debug().Model(&ToDo{}).Create(&t).Error
	if err != nil {
		return &ToDo{}, err
	}
	return t, nil
}

func (t *ToDo) UpdateAToDo(db *gorm.DB) (*ToDo, error) {
	var err error
	err = db.Debug().Model(&ToDo{}).Where("id = ?", t.ID).Updates(ToDo{Title: t.Title, Content: t.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &ToDo{}, err
	}
	
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &ToDo{}, err
		}
	}
	return t, nil
}

func (t *ToDo) DeleteATodo(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&ToDo{}).Where("id = ?", t.ID).Take(&ToDo{}).Delete(&ToDo{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (t *ToDo) FindUserToDos(db *gorm.DB, uid uint32) (*[]ToDo, error) {
	var err error
	todos := []ToDo{}
	err = db.Debug().Model(&ToDo{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&todos).Error
	if err != nil {
		return &[]ToDo{}, err
	}
	if len(todos) > 0 {
		for i, _ := range todos {
			err := db.Debug().Model(&User{}).Where("id = ?", todos[i].AuthorID).Take(&todos[i].Author).Error
			if err != nil {
				return &[]ToDo{}, err
			}
		}
	}
	return &todos, nil
}  