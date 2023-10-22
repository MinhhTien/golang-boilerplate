package seed

import (
	"log"
	"todolist/api/models"

	"github.com/jinzhu/gorm"
)

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.ToDo{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("Cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.ToDo{}).Error
	if err != nil {
		log.Fatalf("Cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.ToDo{}).AddForeignKey("author_id", "user(id)", "cascade", "cascade").Error
	if err != nil {
		// Log the error if it occurs
		log.Fatalf("Attaching foreign key error: %v", err)
	}

	users := []models.User{
		{
			Username: "stev",
			Email:    "stev@gmail.com",
			Password: "stev123",
		},
		{
			Username: "martin",
			Email:    "martin@gmail.com",
			Password: "martin123",
		},
	}

	todos := []models.ToDo{
		{
			Title:   "Dance class",
			Content: "I have to find dance class near me",
		},
		{
			Title:   "Coding",
			Content: "I have to learn coding",
		},
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			// Log the error if it occurs
			log.Fatalf("Cannot seed users table: %v", err)
		}
		todos[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.ToDo{}).Create(&todos[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed todos table: %v", err)
		}
	}
}
