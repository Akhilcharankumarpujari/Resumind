package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
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
	FeedbackRaw string    `gorm:"type:text" json:"-"`
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

// InitDB connects to PostgreSQL using DATABASE_URL and auto-migrates the schema.
//
// Local:      DATABASE_URL=postgres://postgres:password@127.0.0.1:5432/resumind?sslmode=disable
// Production: DATABASE_URL=postgres://user:password@host:5432/dbname?sslmode=require
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal(
			"DATABASE_URL environment variable is not set.\n" +
				"Local format:      postgres://postgres:password@127.0.0.1:5432/resumind?sslmode=disable\n" +
				"Production format: postgres://user:password@host:5432/dbname?sslmode=require",
		)
	}

	log.Printf("Connecting to PostgreSQL at %s", safeRedactPostgresDSN(dsn))

	var err error
	for i := 1; i <= 10; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to PostgreSQL database (attempt %d/10): %v. Retrying in 3 seconds...", i, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL database after 10 attempts: %v", err)
	}

	err = DB.AutoMigrate(&User{}, &Resume{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL and migrated schema successfully")
}

// safeRedactPostgresDSN masks credentials in a postgres:// URL for safe log output.
func safeRedactPostgresDSN(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "invalid-dsn"
	}
	return u.Scheme + "://***@" + u.Host + "/" + strings.TrimPrefix(u.Path, "/")
}
