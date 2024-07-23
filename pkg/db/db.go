package db

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"zerotier-webhook/pkg/models"
)

type Database interface {
	CreateEvent(event *models.Event) error
	GetEvents(networkID, memberID, userID string) ([]models.Event, error)
}

type GormDatabase struct {
	db *gorm.DB
}

func NewDatabase(dbPath string) (Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Event{}); err != nil {
		return nil, err
	}

	return &GormDatabase{db: db}, nil
}

func (g *GormDatabase) CreateEvent(event *models.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	return g.db.Create(event).Error
}

func (g *GormDatabase) GetEvents(networkID, memberID, userID string) ([]models.Event, error) {
	var events []models.Event
	tx := g.db.Model(&models.Event{})

	if networkID != "" {
		tx = tx.Where("network_id = ?", networkID)
	}
	if memberID != "" {
		tx = tx.Where("member_id = ?", memberID)
	}
	if userID != "" {
		tx = tx.Where("user_id = ?", userID)
	}

	if err := tx.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
