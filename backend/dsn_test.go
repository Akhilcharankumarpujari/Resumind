package main

import (
	"net/url"
	"strings"
	"testing"
)

// TestPostgresDSNParsing tests various PostgreSQL connection URL formats (cloud, local, SSL modes).
func TestPostgresDSNParsing(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		wantScheme  string
		wantHost    string
		wantPath    string
		wantSSLMode string
	}{
		{
			name:        "Cloud Production PostgreSQL (SSL required)",
			dsn:         "postgres://resumind_user:SecretPass123@postgres.production.aws.internal:5432/resumind_db?sslmode=require",
			wantScheme:  "postgres",
			wantHost:    "postgres.production.aws.internal:5432",
			wantPath:    "/resumind_db",
			wantSSLMode: "require",
		},
		{
			name:        "Local Development Docker PostgreSQL (SSL disabled)",
			dsn:         "postgres://postgres:rootpassword@127.0.0.1:5432/resumind?sslmode=disable",
			wantScheme:  "postgres",
			wantHost:    "127.0.0.1:5432",
			wantPath:    "/resumind",
			wantSSLMode: "disable",
		},
		{
			name:        "PostgreSQL URL with postgresql:// scheme",
			dsn:         "postgresql://admin:pass@localhost:5432/resumind_dev",
			wantScheme:  "postgresql",
			wantHost:    "localhost:5432",
			wantPath:    "/resumind_dev",
			wantSSLMode: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.Parse(tc.dsn)
			if err != nil {
				t.Fatalf("Failed to parse DSN %q: %v", tc.dsn, err)
			}

			if u.Scheme != tc.wantScheme {
				t.Errorf("Scheme = %q, want %q", u.Scheme, tc.wantScheme)
			}
			if u.Host != tc.wantHost {
				t.Errorf("Host = %q, want %q", u.Host, tc.wantHost)
			}
			if u.Path != tc.wantPath {
				t.Errorf("Path = %q, want %q", u.Path, tc.wantPath)
			}
			if tc.wantSSLMode != "" {
				sslMode := u.Query().Get("sslmode")
				if sslMode != tc.wantSSLMode {
					t.Errorf("sslmode = %q, want %q", sslMode, tc.wantSSLMode)
				}
			}

			// Ensure safe redaction helper does not leak passwords
			redacted := safeRedactPostgresDSN(tc.dsn)
			if strings.Contains(redacted, "SecretPass123") || strings.Contains(redacted, "rootpassword") || strings.Contains(redacted, "pass") {
				t.Errorf("safeRedactPostgresDSN leaked password: %q", redacted)
			}
		})
	}
}

// TestSafeRedactPostgresDSN ensures password redaction works safely on invalid or empty input.
func TestSafeRedactPostgresDSN(t *testing.T) {
	cases := []struct {
		input    string
		contains string
	}{
		{"postgres://user:supersecret@host:5432/db", "postgres://***@"},
		{"invalid dsn syntax", "invalid-dsn"},
	}

	for _, c := range cases {
		res := safeRedactPostgresDSN(c.input)
		if !strings.Contains(res, c.contains) {
			t.Errorf("safeRedactPostgresDSN(%q) = %q, expected to contain %q", c.input, res, c.contains)
		}
	}
}
