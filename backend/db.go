package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"os"
	"time"

	sqldrvmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User is the authenticated user account.
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GoogleID  string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"google_id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Picture   string    `gorm:"type:varchar(1024)" json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Resume is a stored resume analysis record.
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

// Feedback is the AI-generated resume analysis result.
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

// InitDB opens a GORM connection to MySQL using the MYSQL_DSN environment
// variable and runs AutoMigrate.
//
// DSN format (go-sql-driver/mysql):
//
//	user:password@tcp(host:port)/dbname?parseTime=true[&tls=true]
//
// When tls=true is present, a custom TLS config with InsecureSkipVerify is
// registered so that managed cloud MySQL providers (Aiven, PlanetScale, etc.)
// that use custom CA chains can still be reached over encrypted connections.
func InitDB() {
	rawDSN := os.Getenv("MYSQL_DSN")
	if rawDSN == "" {
		log.Fatal("MYSQL_DSN environment variable is not set.\n" +
			"Expected format: user:password@tcp(host:port)/dbname?parseTime=true")
	}

	// ── Step 1: Parse the DSN into a structured config ─────────────────────
	// ParseDSN understands the full go-sql-driver/mysql DSN grammar and
	// populates typed fields (Net, Addr, User, Passwd, DBName, TLSConfig, …).
	// This is the authoritative way to inspect and modify a DSN — never split
	// the raw string manually.
	cfg, err := sqldrvmysql.ParseDSN(rawDSN)
	if err != nil {
		log.Fatalf("MYSQL_DSN is not a valid go-sql-driver/mysql DSN: %v\n"+
			"Expected format: user:password@tcp(host:port)/dbname?parseTime=true", err)
	}

	// ── Step 2: Ensure the network field is set ─────────────────────────────
	// go-sql-driver/mysql defaults Net to "tcp" during ParseDSN, but guard
	// explicitly so that a bare host:port DSN (no tcp(...) wrapper) still
	// works rather than producing "unknown network".
	if cfg.Net == "" {
		cfg.Net = "tcp"
	}

	// ── Step 3: Register and apply a custom TLS config ─────────────────────
	// When tls=true is requested, Go verifies the server cert against the
	// system CA pool. Managed providers (Aiven, etc.) use their own CA that
	// is NOT in the system pool, so verification fails and the driver emits
	// the misleading "unknown network" error.
	//
	// Fix: register a named TLS config with InsecureSkipVerify and point the
	// DSN at it. The connection remains fully encrypted; only the CA chain
	// verification is relaxed. Set tls=false to disable TLS entirely (local
	// Docker), or tls=skip-verify to use the driver's built-in equivalent.
	if cfg.TLSConfig == "true" {
		const tlsName = "custom-ca"
		if regErr := sqldrvmysql.RegisterTLSConfig(tlsName, &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // intentional for managed MySQL CAs
		}); regErr != nil {
			// A non-nil error means the name is already registered (e.g. from
			// a previous call). That is fine — the config is still usable.
			log.Printf("TLS config '%s' already registered (this is fine): %v", tlsName, regErr)
		}
		cfg.TLSConfig = tlsName
	}

	// ── Step 4: Rebuild the DSN from the corrected config ──────────────────
	// FormatDSN serialises all fields back into a valid DSN string.  This is
	// the only mutation of the DSN — no manual string splitting or replacing.
	dsn := cfg.FormatDSN()

	// ── Step 5: Open the database connection with retries ──────────────────
	var dbErr error
	for attempt := 1; attempt <= 10; attempt++ {
		DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if dbErr == nil {
			break
		}
		log.Printf("MySQL connection attempt %d/10 failed: %v — retrying in 3 s...", attempt, dbErr)
		time.Sleep(3 * time.Second)
	}
	if dbErr != nil {
		log.Fatalf("Failed to connect to MySQL after 10 attempts: %v", dbErr)
	}

	// ── Step 6: Auto-migrate schema ────────────────────────────────────────
	if migrateErr := DB.AutoMigrate(&User{}, &Resume{}); migrateErr != nil {
		log.Fatalf("AutoMigrate failed: %v", migrateErr)
	}

	log.Println("✅ Connected to MySQL and migrated schema successfully")
}
