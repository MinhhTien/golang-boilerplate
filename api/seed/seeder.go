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
		models.User{
			Username: "stev",
			Email:    "stev@gmail.com",
			Password: "stev123",
		   },
		   models.User{
			Username: "martin",
			Email:    "martin@gmail.com",
			Password: "martin123",
		   },
	}
}
