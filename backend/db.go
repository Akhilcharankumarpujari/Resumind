package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	sqldrvmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// validMySQLNetworks is the set of network types accepted by go-sql-driver/mysql.
// Any value outside this set in cfg.Net after ParseDSN indicates a malformed DSN.
var validMySQLNetworks = map[string]bool{
	"tcp":        true,
	"tcp4":       true,
	"tcp6":       true,
	"unix":       true,
	"unixgram":   true,
	"unixpacket": true,
}

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
// Required DSN format (go-sql-driver/mysql):
//
//	user:password@tcp(host:port)/dbname?parseTime=true[&tls=skip-verify]
//
// The "@tcp(...)" part is mandatory. Omitting it causes the driver to parse
// the username:password as the network protocol, which results in:
//
//	dial user:password: unknown network user:password
//
// For Aiven MySQL (custom CA, TLS required):
//
//	avnadmin:PASSWORD@tcp(host:22026)/defaultdb?parseTime=true&tls=skip-verify
//
// For local Docker MySQL (no TLS):
//
//	root:password@tcp(127.0.0.1:3306)/resumind?parseTime=true
func InitDB() {
	rawDSN := os.Getenv("MYSQL_DSN")
	if rawDSN == "" {
		log.Fatal(
			"MYSQL_DSN is not set.\n" +
				"Aiven format: user:password@tcp(host:port)/dbname?parseTime=true&tls=skip-verify\n" +
				"Local format: user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true",
		)
	}

	// ── Step 1: Parse DSN into a typed config struct ──────────────────────
	//
	// ParseDSN is the only correct way to read a DSN. It handles @, = and &
	// inside passwords correctly. Hand-written string splitting will break.
	cfg, parseErr := sqldrvmysql.ParseDSN(rawDSN)
	if parseErr != nil {
		log.Fatalf(
			"MYSQL_DSN is not valid: %v\n"+
				"Required format: user:password@tcp(host:port)/dbname?parseTime=true&tls=skip-verify",
			parseErr,
		)
	}

	// ── Step 2: Validate cfg.Net ──────────────────────────────────────────
	//
	// After ParseDSN the cfg.Net field MUST be one of the Go network types
	// (tcp, unix, …). If it is anything else — especially if it looks like
	// "user:password" — the DSN is missing the "@tcp(...)" part.
	//
	// Failure mode that occurred on Render:
	//   DSN:       avnadmin:PASS(host:22026)/defaultdb?tls=skip-verify
	//              ↑ missing "@tcp"
	//   Parsed as: cfg.Net="avnadmin:PASS"  cfg.User=""  cfg.Addr="host:22026"
	//   Result:    net.Dial("avnadmin:PASS", "host:22026") → unknown network
	//
	// The driver does NOT return an error from ParseDSN for this input; it
	// silently accepts "avnadmin:PASS" as a custom network type.
	if !validMySQLNetworks[cfg.Net] {
		// cfg.Net is not a Go network type — DSN is structurally broken.
		// Log the target addr and TLS config for diagnostics; never log Net
		// because it contains the credential in this malformed case.
		log.Fatalf(
			"MYSQL_DSN is structurally invalid: cfg.Net is not a recognized network type.\n"+
				"This means the DSN is missing '@tcp(...)'. The '@' and 'tcp' keyword are required.\n\n"+
				"  WRONG (missing @tcp):  user:password(host:port)/db?parseTime=true&tls=skip-verify\n"+
				"  CORRECT:               user:password@tcp(host:port)/db?parseTime=true&tls=skip-verify\n\n"+
				"Parsed addr=%q  tls=%q — update MYSQL_DSN in Render Environment Variables.",
			cfg.Addr, cfg.TLSConfig,
		)
	}

	// ── Step 3: Validate cfg.User is present ─────────────────────────────
	//
	// A well-formed DSN with "@tcp(...)" always populates cfg.User.
	// An empty User means the "@" separator was missing from the DSN.
	if cfg.User == "" {
		log.Fatal(
			"MYSQL_DSN is missing the username. The '@' separator between credentials\n" +
				"and the host is required: user:password@tcp(host:port)/db",
		)
	}

	// ── Step 4: Normalise TLS ─────────────────────────────────────────────
	//
	// tls=skip-verify → built-in go-sql-driver name, pass through unchanged.
	// tls=true        → uses system CA pool. Aiven's CA is NOT in the system
	//                   pool. Replace with skip-verify to allow encrypted
	//                   connections without CA chain verification.
	// tls="" / false  → no TLS (local Docker). Leave unchanged.
	if cfg.TLSConfig == "true" {
		cfg.TLSConfig = "skip-verify"
	}

	// ── Step 5: Rebuild DSN once from the validated config ────────────────
	//
	// FormatDSN serialises all struct fields into a canonical DSN string.
	// This is the only mutation of the DSN value — no string splitting,
	// no strings.ReplaceAll, no manual construction.
	dsn := cfg.FormatDSN()

	// ── Step 6: Log connection target — never credentials ─────────────────
	log.Printf("Connecting to MySQL: net=%s addr=%s db=%s tls=%s",
		cfg.Net, cfg.Addr, cfg.DBName, cfg.TLSConfig)

	// ── Step 7: Open GORM connection with retries ─────────────────────────
	//
	// gorm.Open(mysql.Open(dsn), ...) is the standard GORM+go-sql-driver
	// connection call. The complete DSN string is passed as one argument.
	// Credentials are NOT passed separately; they are embedded in dsn.
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

	// ── Step 8: Auto-migrate schema ───────────────────────────────────────
	if migrateErr := DB.AutoMigrate(&User{}, &Resume{}); migrateErr != nil {
		log.Fatalf("AutoMigrate failed: %v", migrateErr)
	}

	log.Printf("✅ Connected to MySQL at %s — schema migrated", cfg.Addr)
}
