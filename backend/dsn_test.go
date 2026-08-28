package main

import (
	"fmt"
	"testing"

	sqldrvmysql "github.com/go-sql-driver/mysql"
)

func TestDSNParsing(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		wantNet     string
		wantTLS     string
		wantUser    string
		wantAddr    string
		wantDB      string
	}{
		{
			name:     "production Aiven DSN with tcp wrapper and tls=true",
			dsn:      "avnadmin:REDACTED@tcp(mysql-2fb71695-resumind19.h.aivencloud.com:22026)/defaultdb?parseTime=true&tls=true",
			wantNet:  "tcp",
			wantTLS:  "true",
			wantUser: "avnadmin",
			wantAddr: "mysql-2fb71695-resumind19.h.aivencloud.com:22026",
			wantDB:   "defaultdb",
		},
		{
			name:     "local Docker DSN without TLS",
			dsn:      "root:akhil123@tcp(127.0.0.1:3307)/resumind?parseTime=true",
			wantNet:  "tcp",
			wantTLS:  "",
			wantUser: "root",
			wantAddr: "127.0.0.1:3307",
			wantDB:   "resumind",
		},
		{
			name:     "DSN without tcp wrapper (common misconfiguration)",
			dsn:      "avnadmin:REDACTED@mysql-host.example.com:22026/defaultdb?parseTime=true&tls=true",
			wantNet:  "tcp", // our code sets this if empty
			wantTLS:  "true",
			wantUser: "avnadmin",
			wantDB:   "defaultdb",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := sqldrvmysql.ParseDSN(tc.dsn)
			if err != nil {
				t.Fatalf("ParseDSN failed: %v", err)
			}

			// Apply the same fix as InitDB
			if cfg.Net == "" {
				cfg.Net = "tcp"
			}

			if cfg.Net != tc.wantNet {
				t.Errorf("Net: got %q want %q", cfg.Net, tc.wantNet)
			}
			if cfg.TLSConfig != tc.wantTLS {
				t.Errorf("TLSConfig: got %q want %q", cfg.TLSConfig, tc.wantTLS)
			}
			if cfg.User != tc.wantUser {
				t.Errorf("User: got %q want %q", cfg.User, tc.wantUser)
			}
			if tc.wantDB != "" && cfg.DBName != tc.wantDB {
				t.Errorf("DBName: got %q want %q", cfg.DBName, tc.wantDB)
			}
			if tc.wantAddr != "" && cfg.Addr != tc.wantAddr {
				t.Errorf("Addr: got %q want %q", cfg.Addr, tc.wantAddr)
			}

			// Verify FormatDSN does not produce an unparseable DSN
			rebuilt := cfg.FormatDSN()
			cfg2, err := sqldrvmysql.ParseDSN(rebuilt)
			if err != nil {
				t.Errorf("FormatDSN produced unparseable DSN: %v", err)
			}
			if cfg2.Net != cfg.Net || cfg2.User != cfg.User || cfg2.DBName != cfg.DBName {
				t.Errorf("FormatDSN round-trip mismatch: net %q user %q db %q",
					cfg2.Net, cfg2.User, cfg2.DBName)
			}

			fmt.Printf("[%s] Net=%q User=%q Addr=%q DB=%q TLS=%q  rebuilt=%q\n",
				tc.name, cfg.Net, cfg.User, cfg.Addr, cfg.DBName, cfg.TLSConfig,
				rebuilt[:min(60, len(rebuilt))]+"...")
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
