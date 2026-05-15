package portal

import (
	"net/http"
	"net/url"
	"testing"
)

// MockRedis is not needed if we test the interface logic
// But for now, let's just test the domain filtering and cookie merging logic
// if we can refactor RedisCookieJar to be more testable.

func TestIsPortalDomain(t *testing.T) {
	jar := &RedisCookieJar{}
	
	tests := []struct {
		host     string
		expected bool
	}{
		{"tiemchung.vncdc.gov.vn", true},
		{"google.com", false},
		{"localhost", false},
	}

	for _, tt := range tests {
		if jar.isPortalDomain(tt.host) != tt.expected {
			t.Errorf("isPortalDomain(%s) = %v, expected %v", tt.host, !tt.expected, tt.expected)
		}
	}
}

func TestCookieJarLogic(t *testing.T) {
	// Since we can't easily mock go-redis without extra deps, 
	// we'll focus on testing the client integration
}
