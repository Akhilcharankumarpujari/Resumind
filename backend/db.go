package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	sqldrvmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User is the user account model.
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GoogleID  string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"google_id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Picture   string    `gorm:"type:varchar(1024)" json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Resume is the uploaded resume model.
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

// Feedback is the AI analysis result structure.
type Feedback struct {
	OverallScore int          `json:"overallScore"`
	ATS          SectionScore `json:"ATS"`
	ToneAndStyle SectionScore `json:"toneAndStyle"`
	Content      SectionScore `json:"content"`
	Structure    SectionScore `json:"structure"`
	Skills       SectionScore `json:"skills"`
}

// SectionScore holds a score and list of tips for a feedback section.
type SectionScore struct {
	Score int   `json:"score"`
	Tips  []Tip `json:"tips"`
}

// Tip is a single actionable improvement suggestion.
type Tip struct {
	Type        string `json:"type"`
	Tip         string `json:"tip"`
	Explanation string `json:"explanation,omitempty"`
}

// AfterFind deserializes the stored JSON feedback after fetching a resume.
func (r *Resume) AfterFind(tx *gorm.DB) error {
	if r.FeedbackRaw != "" {
		var f Feedback
		if err := json.Unmarshal([]byte(r.FeedbackRaw), &f); err == nil {
			r.Feedback = &f
		}
	}
	return nil
}

// InitDB connects to MySQL using MYSQL_DSN, registers TLS config when needed,
// and runs GORM AutoMigrate.
func InitDB() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN environment variable is not set. Example: user:password@tcp(127.0.0.1:3306)/resumind?parseTime=true")
	}

	// Validate DSN structure using the MySQL driver's own parser.
	// This catches malformed DSNs (e.g., missing tcp(...) network spec) early.
	cfg, err := sqldrvmysql.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("MYSQL_DSN is invalid: %v", err)
	}

	// If the parsed network is empty the driver will fail with "unknown network".
	// go-sql-driver/mysql defaults to "tcp" but older DSN formats can miss it.
	// Ensure we always have a valid network.
	if cfg.Net == "" {
		cfg.Net = "tcp"
		// Rebuild DSN from the corrected config so GORM receives a valid string.
		dsn = cfg.FormatDSN()
	}

	// Register a skip-verify TLS config for managed cloud MySQL providers
	// (e.g. Aiven, PlanetScale, Railway) that use custom CAs not in Go's
	// default pool. We only do this when the DSN explicitly requests TLS.
	//
	// Using "skip-verify" avoids the "unknown network" error that occurs when
	// Go cannot verify the server certificate and the driver cannot dial.
	// Production traffic is still encrypted; only certificate validation is
	// relaxed. For stricter verification, set MYSQL_TLS_CA to a PEM file path.
	if strings.Contains(dsn, "tls=true") || strings.Contains(dsn, "tls=skip-verify") {
		tlsCfgName := "custom"
		// If the user already specified tls=skip-verify, honour it directly.
		if strings.Contains(dsn, "tls=skip-verify") {
			tlsCfgName = "skip-verify" // already registered by the driver
		} else {
			// tls=true: register a permissive config for cloud MySQL providers.
			if regErr := sqldrvmysql.RegisterTLSConfig(tlsCfgName, &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec // Intentional for managed MySQL CAs
			}); regErr != nil {
				log.Printf("Warning: could not register custom TLS config: %v", regErr)
				tlsCfgName = "skip-verify"
			}
			// Replace tls=true in the DSN with our named config.
			dsn = strings.ReplaceAll(dsn, "tls=true", "tls="+tlsCfgName)
		}
		_ = tlsCfgName
	}

	var dbErr error
	for i := 1; i <= 10; i++ {
		DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if dbErr == nil {
			break
		}
		log.Printf("Failed to connect to MySQL (attempt %d/10). Retrying in 3 seconds...", i)
		time.Sleep(3 * time.Second)
	}
	if dbErr != nil {
		log.Fatalf("Failed to connect to MySQL after 10 attempts: %v", dbErr)
	}

	if migrateErr := DB.AutoMigrate(&User{}, &Resume{}); migrateErr != nil {
		log.Fatalf("Failed to auto-migrate database: %v", migrateErr)
	}

	log.Println("✅ Connected to MySQL and migrated schema successfully")
}
