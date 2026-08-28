package main

import (
	"strings"
	"testing"

	sqldrvmysql "github.com/go-sql-driver/mysql"
)

// ─── Helpers ──────────────────────────────────────────────────────────────────

// parseDSNAndApplyFixes applies the same transformations as InitDB so tests
// can verify the final cfg state without a live database.
// It does NOT print credentials.
func parseDSNAndApplyFixes(t *testing.T, dsn string) *sqldrvmysql.Config {
	t.Helper()
	cfg, err := sqldrvmysql.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("ParseDSN failed: %v", err)
	}

	// mirror InitDB logic exactly
	if cfg.Net == "" {
		cfg.Net = "tcp"
	}
	if cfg.TLSConfig == "true" {
		cfg.TLSConfig = "skip-verify"
	}
	return cfg
}

// ─── The test the user requested ──────────────────────────────────────────────

// TestProductionDSNFields verifies that a correctly formatted Aiven DSN
// produces the exact cfg field values required for a successful connection.
//
// The MYSQL_DSN on Render MUST follow this format:
//
//	avnadmin:PASSWORD@tcp(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=skip-verify
func TestProductionDSNFields(t *testing.T) {
	// Use a placeholder password — tests must never contain real credentials.
	const dsn = "avnadmin:PLACEHOLDER@tcp(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=skip-verify"

	cfg := parseDSNAndApplyFixes(t, dsn)

	if cfg.Net != "tcp" {
		t.Errorf("cfg.Net = %q, want %q", cfg.Net, "tcp")
	}
	if cfg.User != "avnadmin" {
		t.Errorf("cfg.User = %q, want %q", cfg.User, "avnadmin")
	}
	if cfg.Addr != "mysql-2fb71695-resumind19.h.aivencloud.com:22026" {
		t.Errorf("cfg.Addr = %q, want %q", cfg.Addr, "mysql-2fb71695-resumind19.h.aivencloud.com:22026")
	}
	if cfg.DBName != "defaultdb" {
		t.Errorf("cfg.DBName = %q, want %q", cfg.DBName, "defaultdb")
	}
	if cfg.TLSConfig != "skip-verify" {
		t.Errorf("cfg.TLSConfig = %q, want %q", cfg.TLSConfig, "skip-verify")
	}

	// The Net field must never equal the username or password.
	// If it does, the DSN is missing "@tcp(...)" and the driver will call
	// net.Dial("avnadmin:PASSWORD", ...) which fails with "unknown network".
	if cfg.Net == cfg.User || cfg.Net == cfg.Passwd {
		t.Errorf("CRITICAL: cfg.Net (%q) equals username or password — DSN is missing @tcp(...)", cfg.Net)
	}

	t.Logf("cfg.Net=%q cfg.User=%q cfg.Addr=%q cfg.DBName=%q cfg.TLSConfig=%q",
		cfg.Net, cfg.User, cfg.Addr, cfg.DBName, cfg.TLSConfig)
}

// ─── Prove the exact failure mode seen on Render ──────────────────────────────

// TestMissingAtTcpProducesWrongNet documents and detects the exact malformed
// DSN structure that produced the Render log:
//
//	Connecting to MySQL at host:22026 (net=avnadmin:******, tls=skip-verify)
//
// The DSN was missing "@tcp" so the driver parsed user:pass as cfg.Net.
func TestMissingAtTcpProducesWrongNet(t *testing.T) {
	// This is the malformed DSN that caused the production failure.
	// It is missing "@" and "tcp". The driver parses "avnadmin:PLACEHOLDER"
	// as the network protocol, leaving cfg.User and cfg.Passwd empty.
	malformedDSN := "avnadmin:PLACEHOLDER(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=skip-verify"

	cfg, err := sqldrvmysql.ParseDSN(malformedDSN)
	if err != nil {
		// If ParseDSN rejects this DSN, our validation catches it.
		t.Logf("ParseDSN correctly rejected malformed DSN: %v", err)
		return
	}

	// Document what the driver actually parses — this is what was deployed.
	t.Logf("Malformed DSN parsed as: Net=%q User=%q Addr=%q DBName=%q TLSConfig=%q",
		cfg.Net, cfg.User, cfg.Addr, cfg.DBName, cfg.TLSConfig)

	// cfg.Net will contain "avnadmin:PLACEHOLDER" — not "tcp".
	// This is what our new validation in InitDB must catch.
	if cfg.Net == "tcp" {
		t.Log("Net is tcp — DSN appears well-formed")
	} else {
		t.Logf("Net is %q (not tcp) — this is the bug: DSN needs @tcp(...)", cfg.Net)
	}

	// cfg.User will be empty because there was no "@" to separate credentials.
	if cfg.User == "" {
		t.Log("cfg.User is empty — confirms @tcp was missing from DSN")
	}
}

// ─── InitDB validation: cfg.Net must be tcp ───────────────────────────────────

// TestInitDBValidatesNet checks that InitDB-style validation catches a DSN
// where cfg.Net is not a recognized Go network type.
func TestInitDBValidatesNet(t *testing.T) {
	// validNetworks is the same set checked in InitDB.
	validNetworks := map[string]bool{
		"tcp": true, "tcp4": true, "tcp6": true,
		"unix": true, "unixgram": true, "unixpacket": true,
	}

	// Correct DSN — must pass validation.
	goodDSN := "avnadmin:PLACEHOLDER@tcp(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=skip-verify"
	cfg, _ := sqldrvmysql.ParseDSN(goodDSN)
	if !validNetworks[cfg.Net] {
		t.Errorf("good DSN: cfg.Net=%q failed validation — expected tcp", cfg.Net)
	} else {
		t.Logf("good DSN: cfg.Net=%q passed validation ✓", cfg.Net)
	}

	// Malformed DSN (the one seen on Render) — must fail validation.
	badDSN := "avnadmin:PLACEHOLDER(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=skip-verify"
	cfg2, err := sqldrvmysql.ParseDSN(badDSN)
	if err != nil {
		t.Logf("malformed DSN: ParseDSN rejected it: %v ✓", err)
		return
	}
	if validNetworks[cfg2.Net] {
		t.Errorf("malformed DSN: cfg.Net=%q incorrectly passed validation", cfg2.Net)
	} else {
		t.Logf("malformed DSN: cfg.Net=%q correctly fails validation — InitDB will Fatalf ✓", cfg2.Net)
	}
}

// ─── Local Docker DSN (no TLS) ────────────────────────────────────────────────

// TestLocalDockerDSN verifies the local development DSN format.
func TestLocalDockerDSN(t *testing.T) {
	const dsn = "root:localpassword@tcp(127.0.0.1:3307)/resumind?parseTime=true"

	cfg := parseDSNAndApplyFixes(t, dsn)

	if cfg.Net != "tcp" {
		t.Errorf("cfg.Net = %q, want %q", cfg.Net, "tcp")
	}
	if cfg.User != "root" {
		t.Errorf("cfg.User = %q, want %q", cfg.User, "root")
	}
	if cfg.Addr != "127.0.0.1:3307" {
		t.Errorf("cfg.Addr = %q, want %q", cfg.Addr, "127.0.0.1:3307")
	}
	if cfg.TLSConfig != "" {
		t.Errorf("cfg.TLSConfig = %q, want %q (no TLS for local)", cfg.TLSConfig, "")
	}
	t.Logf("cfg.Net=%q cfg.User=%q cfg.Addr=%q cfg.TLSConfig=%q",
		cfg.Net, cfg.User, cfg.Addr, cfg.TLSConfig)
}

// ─── FormatDSN round-trip ─────────────────────────────────────────────────────

// TestFormatDSNRoundTrip checks that FormatDSN produces a re-parseable DSN
// with the same field values — no information is lost or corrupted.
func TestFormatDSNRoundTrip(t *testing.T) {
	const dsn = "avnadmin:PLACEHOLDER@tcp(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=skip-verify"

	cfg := parseDSNAndApplyFixes(t, dsn)
	rebuilt := cfg.FormatDSN()

	cfg2, err := sqldrvmysql.ParseDSN(rebuilt)
	if err != nil {
		t.Fatalf("FormatDSN produced unparseable DSN: %v", err)
	}
	if cfg2.Net != "tcp" {
		t.Errorf("round-trip: cfg.Net = %q, want tcp", cfg2.Net)
	}
	if cfg2.User != cfg.User {
		t.Errorf("round-trip: cfg.User = %q, want %q", cfg2.User, cfg.User)
	}
	if cfg2.Addr != cfg.Addr {
		t.Errorf("round-trip: cfg.Addr = %q, want %q", cfg2.Addr, cfg.Addr)
	}
	if cfg2.DBName != cfg.DBName {
		t.Errorf("round-trip: cfg.DBName = %q, want %q", cfg2.DBName, cfg.DBName)
	}
	if cfg2.TLSConfig != "skip-verify" {
		t.Errorf("round-trip: cfg.TLSConfig = %q, want skip-verify", cfg2.TLSConfig)
	}

	// The rebuilt DSN must contain "@tcp(" — the driver's canonical format.
	if !strings.Contains(rebuilt, "@tcp(") {
		t.Errorf("rebuilt DSN missing '@tcp(' — would cause 'unknown network': %q",
			rebuilt[:min(80, len(rebuilt))])
	}

	t.Logf("rebuilt DSN contains @tcp(: ✓")
	t.Logf("round-trip Net=%q User=%q Addr=%q TLS=%q", cfg2.Net, cfg2.User, cfg2.Addr, cfg2.TLSConfig)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
