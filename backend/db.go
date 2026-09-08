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

// User is the authenticated user account model.
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GoogleID  string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"google_id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Picture   string    `gorm:"type:varchar(1024)" json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Resume is a stored resume analysis record model.
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

// Feedback is the AI-generated resume analysis result structure.
type Feedback struct {
	OverallScore int          `json:"overallScore"`
	ATS          SectionScore `json:"ATS"`
	ToneAndStyle SectionScore `json:"toneAndStyle"`
	Content      SectionScore `json:"content"`
	Structure    SectionScore `json:"structure"`
	Skills       SectionScore `json:"skills"`
}

// SectionScore holds a category score and actionable tips.
type SectionScore struct {
	Score int   `json:"score"`
	Tips  []Tip `json:"tips"`
}

// Tip is a single improvement suggestion.
type Tip struct {
	Type        string `json:"type"`
	Tip         string `json:"tip"`
	Explanation string `json:"explanation,omitempty"`
}

// AfterFind deserialises FeedbackRaw into Feedback after a GORM query.
func (r *Resume) AfterFind(tx *gorm.DB) error {
	if r.FeedbackRaw != "" {
		var f Feedback
		if err := json.Unmarshal([]byte(r.FeedbackRaw), &f); err == nil {
			r.Feedback = &f
		}
	}
	return nil
}

// InitDB connects to PostgreSQL using the DATABASE_URL environment variable and runs GORM AutoMigrate.
//
// Expected DSN formats:
//
//	Render / Cloud PostgreSQL (SSL enabled):
//	  postgres://user:password@host:5432/dbname?sslmode=require
//
//	Local Development / Docker PostgreSQL (SSL disabled):
//	  postgres://postgres:password@127.0.0.1:5432/resumind?sslmode=disable
func InitDB() {
	rawDSN := os.Getenv("DATABASE_URL")
	if rawDSN == "" {
		rawDSN = os.Getenv("POSTGRES_DSN")
	}
	if rawDSN == "" {
		rawDSN = os.Getenv("MYSQL_DSN")
	}
	if rawDSN == "" {
		log.Fatal(
			"DATABASE_URL environment variable is not set.\n" +
				"Production format: postgres://user:password@host:5432/dbname?sslmode=require\n" +
				"Local format:      postgres://postgres:password@127.0.0.1:5432/resumind?sslmode=disable",
		)
	}

	// Sanitize / ensure valid scheme for GORM PostgreSQL driver
	dsn := rawDSN
	if strings.HasPrefix(dsn, "mysql://") {
		log.Fatal("MYSQL_DSN is obsolete. Please update environment variable to DATABASE_URL with postgres:// scheme.")
	}

	targetLog := safeRedactPostgresDSN(dsn)
	log.Printf("Connecting to PostgreSQL at %s", targetLog)

	var dbErr error
	for attempt := 1; attempt <= 10; attempt++ {
		DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if dbErr == nil {
			break
		}
		log.Printf("PostgreSQL connect attempt %d/10 failed: %v — retrying in 3 s...", attempt, dbErr)
		time.Sleep(3 * time.Second)
	}
	if dbErr != nil {
		log.Fatalf("Failed to connect to PostgreSQL after 10 attempts: %v", dbErr)
	}

	// Auto-migrate database tables
	if migrateErr := DB.AutoMigrate(&User{}, &Resume{}); migrateErr != nil {
		log.Fatalf("AutoMigrate failed: %v", migrateErr)
	}

	log.Printf("✅ Connected to PostgreSQL (%s) — schema migrated successfully", targetLog)
}

// safeRedactPostgresDSN returns host:port/dbname with credentials masked for safe logging.
func safeRedactPostgresDSN(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "invalid-dsn"
	}
	host := u.Host
	path := strings.TrimPrefix(u.Path, "/")
	scheme := u.Scheme
	return scheme + "://***@" + host + "/" + path
}
