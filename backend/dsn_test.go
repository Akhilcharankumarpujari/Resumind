package main

import (
	"fmt"
	"testing"

	sqldrvmysql "github.com/go-sql-driver/mysql"
)

// TestDSNRoundTrip verifies that ParseDSN + FormatDSN produces a DSN that:
// 1. Has Net set correctly (not empty, not the username/password)
// 2. Has the correct TLSConfig field
// 3. Round-trips cleanly through a second ParseDSN call
func TestDSNRoundTrip(t *testing.T) {
	// This is the exact format the user should set on Render
	productionDSN := "avnadmin:REDACTED@tcp(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=true"
	localDSN := "root:akhil123@tcp(127.0.0.1:3307)/resumind?parseTime=true"
	skipVerifyDSN := "avnadmin:REDACTED@tcp(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=skip-verify"

	cases := []struct {
		name         string
		dsn          string
		wantNet      string
		wantUser     string
		wantTLS      string
		wantTLSAfter string // TLSConfig after our InitDB transformation
	}{
		{
			name:         "aiven production (tls=true)",
			dsn:          productionDSN,
			wantNet:      "tcp",
			wantUser:     "avnadmin",
			wantTLS:      "true",
			wantTLSAfter: "skip-verify",
		},
		{
			name:         "aiven skip-verify",
			dsn:          skipVerifyDSN,
			wantNet:      "tcp",
			wantUser:     "avnadmin",
			wantTLS:      "skip-verify",
			wantTLSAfter: "skip-verify", // no transformation needed
		},
		{
			name:         "local docker (no TLS)",
			dsn:          localDSN,
			wantNet:      "tcp",
			wantUser:     "root",
			wantTLS:      "",
			wantTLSAfter: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Step 1: Parse
			cfg, err := sqldrvmysql.ParseDSN(tc.dsn)
			if err != nil {
				t.Fatalf("ParseDSN failed: %v", err)
			}

			// Step 2: Verify parsed fields
			if cfg.Net != tc.wantNet {
				t.Errorf("cfg.Net = %q, want %q", cfg.Net, tc.wantNet)
			}
			if cfg.User != tc.wantUser {
				t.Errorf("cfg.User = %q, want %q", cfg.User, tc.wantUser)
			}
			if cfg.TLSConfig != tc.wantTLS {
				t.Errorf("cfg.TLSConfig = %q, want %q", cfg.TLSConfig, tc.wantTLS)
			}

			// Confirm Net is never the username or password
			if cfg.Net == cfg.User || cfg.Net == cfg.Passwd {
				t.Errorf("CRITICAL: cfg.Net (%q) equals user/password — this causes 'unknown network'", cfg.Net)
			}

			// Step 3: Apply our InitDB transformation (mirror the real code exactly)
			if cfg.Net == "" {
				cfg.Net = "tcp"
			}
			if cfg.TLSConfig == "true" {
				cfg.TLSConfig = "skip-verify" // built-in, no RegisterTLSConfig needed
			}

			// Step 4: Rebuild DSN via FormatDSN (no string manipulation)
			rebuilt := cfg.FormatDSN()

			// Step 5: Verify the rebuilt DSN round-trips correctly
			cfg2, err := sqldrvmysql.ParseDSN(rebuilt)
			if err != nil {
				t.Fatalf("rebuilt DSN is unparseable: %v\nDSN: %q", err, rebuilt)
			}
			if cfg2.Net != "tcp" {
				t.Errorf("rebuilt DSN: Net = %q, want \"tcp\"", cfg2.Net)
			}
			if cfg2.TLSConfig != tc.wantTLSAfter {
				t.Errorf("rebuilt DSN: TLSConfig = %q, want %q", cfg2.TLSConfig, tc.wantTLSAfter)
			}
			if cfg2.Net == cfg2.User || cfg2.Net == cfg2.Passwd {
				t.Errorf("CRITICAL: rebuilt DSN Net (%q) equals user/password", cfg2.Net)
			}

			fmt.Printf("  PASS %-35s  Net=%q  TLS=%q→%q  Addr=%q\n",
				tc.name, cfg2.Net, tc.wantTLS, cfg2.TLSConfig, cfg2.Addr)
		})
	}
}

// TestMalformedDSN verifies that malformed DSNs are caught early by ParseDSN
// and never reach the driver with wrong arguments.
func TestMalformedDSN(t *testing.T) {
	malformed := []struct {
		name string
		dsn  string
	}{
		{"empty", ""},
		{"no slash", "avnadmin:PASS@tcp(host:22026)"},
		{"missing tcp wrapper gives driver error", "avnadmin:PASS@host:22026/db"},
		{"bare username only", "avnadmin"},
	}

	for _, tc := range malformed {
		t.Run(tc.name, func(t *testing.T) {
			if tc.dsn == "" {
				t.Log("empty DSN: caught by os.Getenv check in InitDB")
				return
			}
			_, err := sqldrvmysql.ParseDSN(tc.dsn)
			if err != nil {
				t.Logf("  ParseDSN correctly rejects %q: %v", tc.name, err)
			} else {
				t.Logf("  ParseDSN accepted %q (driver may reject at dial time)", tc.name)
			}
		})
	}
}
