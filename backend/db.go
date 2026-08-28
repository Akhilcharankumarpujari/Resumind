package main

import (
	"encoding/json"
	"fmt"
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

// InitDB opens a GORM/MySQL connection using the MYSQL_DSN environment variable.
//
// ── DSN format (go-sql-driver/mysql) ──────────────────────────────────────────
//
//	user:password@tcp(host:port)/dbname?parseTime=true[&tls=skip-verify]
//
// ── TLS for managed providers (Aiven, PlanetScale, Railway, …) ───────────────
//
//	Use tls=skip-verify in MYSQL_DSN.
//	"skip-verify" is a built-in go-sql-driver/mysql TLS config that enables
//	full encryption while skipping CA certificate chain validation.
//	This is necessary because managed MySQL providers use custom CAs that
//	are not present in the Go runtime's system certificate pool.
//
//	DO NOT use tls=true with Aiven — it will fail with "unknown network" or
//	"x509: certificate signed by unknown authority".
//
// ── Local Docker ──────────────────────────────────────────────────────────────
//
//	Omit the tls parameter entirely: ...?parseTime=true
func InitDB() {
	rawDSN := os.Getenv("MYSQL_DSN")
	if rawDSN == "" {
		log.Fatal(
			"MYSQL_DSN environment variable is not set.\n" +
				"For Aiven: user:password@tcp(host:port)/dbname?parseTime=true&tls=skip-verify\n" +
				"For local: user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true",
		)
	}

	// ── Parse DSN → structured config ─────────────────────────────────────
	// Using ParseDSN (the driver's own parser) is the only safe way to
	// inspect and modify a DSN. Never split the raw string manually — the
	// password may contain characters like @, =, & that would break any
	// hand-written parser.
	cfg, parseErr := sqldrvmysql.ParseDSN(rawDSN)
	if parseErr != nil {
		log.Fatalf(
			"MYSQL_DSN is not a valid go-sql-driver/mysql DSN: %v\n"+
				"For Aiven: user:password@tcp(host:port)/dbname?parseTime=true&tls=skip-verify",
			parseErr,
		)
	}

	// ── Ensure network is set ─────────────────────────────────────────────
	// ParseDSN sets cfg.Net to "tcp" by default when a host:port address is
	// present. Guard against the empty-string edge case so we never pass a
	// blank network to net.Dial (which would surface as "unknown network").
	if cfg.Net == "" {
		cfg.Net = "tcp"
	}

	// ── TLS normalisation ─────────────────────────────────────────────────
	// "skip-verify"  → already a valid built-in name; pass through as-is.
	// "true"         → would use the system CA pool, which does not contain
	//                   Aiven's CA. Replace with "skip-verify" so the driver
	//                   enables encryption without failing CA verification.
	//                   The TLS config name "skip-verify" is pre-registered
	//                   by go-sql-driver/mysql and requires no additional
	//                   RegisterTLSConfig call.
	// ""  / "false"  → no TLS (local dev). Leave unchanged.
	if cfg.TLSConfig == "true" {
		cfg.TLSConfig = "skip-verify"
	}

	// ── Rebuild the DSN from the (possibly modified) config ───────────────
	// FormatDSN serialises all struct fields back into a canonically valid
	// DSN string. This is the only place the DSN is mutated.
	dsn := cfg.FormatDSN()

	// ── Log connection target (host/port only — never credentials) ────────
	log.Printf("Connecting to MySQL at %s (net=%s, tls=%s)", cfg.Addr, cfg.Net, cfg.TLSConfig)

	// ── Open connection with retries ──────────────────────────────────────
	var dbErr error
	for attempt := 1; attempt <= 10; attempt++ {
		DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if dbErr == nil {
			break
		}
		log.Printf("MySQL connect attempt %d/10 failed: %v — retrying in 3 s...", attempt, dbErr)
		time.Sleep(3 * time.Second)
	}
	if dbErr != nil {
		log.Fatalf("Failed to connect to MySQL after 10 attempts: %v", dbErr)
	}

	// ── Auto-migrate schema ───────────────────────────────────────────────
	if migrateErr := DB.AutoMigrate(&User{}, &Resume{}); migrateErr != nil {
		log.Fatalf("AutoMigrate failed: %v", migrateErr)
	}

	log.Printf("✅ Connected to MySQL at %s and migrated schema successfully", cfg.Addr)
}

// safeRedactDSN returns a DSN-like string with the password replaced by ***
// for use in diagnostic messages. Never log the real DSN.
func safeRedactDSN(cfg *sqldrvmysql.Config) string {
	return fmt.Sprintf("%s:***@%s(%s)/%s", cfg.User, cfg.Net, cfg.Addr, cfg.DBName)
}
