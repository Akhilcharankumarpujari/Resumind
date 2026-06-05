package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB


type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GoogleID  string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"google_id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Picture   string    `gorm:"type:varchar(1024)" json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type Resume struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	CompanyName string    `gorm:"type:varchar(255);not null" json:"companyName"`
	JobTitle    string    `gorm:"type:varchar(255);not null" json:"jobTitle"`
	JobDesc     string    `gorm:"type:text" json:"jobDescription"`
	FilePath    string    `gorm:"type:varchar(512);not null" json:"resumePath"`
	FeedbackRaw string    `gorm:"type:longtext" json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`


	Feedback *Feedback `gorm:"-" json:"feedback,omitempty"`
}


type Feedback struct {
	OverallScore int          `json:"overallScore"`
	ATS          SectionScore `json:"ATS"`
	ToneAndStyle SectionScore `json:"toneAndStyle"`
	Content      SectionScore `json:"content"`
	Structure    SectionScore `json:"structure"`
	Skills       SectionScore `json:"skills"`
}

type SectionScore struct {
	Score int   `json:"score"`
	Tips  []Tip `json:"tips"`
}

type Tip struct {
	Type        string `json:"type"`
	Tip         string `json:"tip"`
	Explanation string `json:"explanation,omitempty"`
}


func (r *Resume) AfterFind(tx *gorm.DB) error {
	if r.FeedbackRaw != "" {
		var f Feedback
		if err := json.Unmarshal([]byte(r.FeedbackRaw), &f); err == nil {
			r.Feedback = &f
		}
	}
	return nil
}


func InitDB() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN environment variable is not set. Example: user:password@tcp(127.0.0.1:3306)/resumind?parseTime=true")
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
	}


	err = DB.AutoMigrate(&User{}, &Resume{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("✅ Connected to MySQL and migrated schema successfully")
}
