package geoip

import (
	"testing"
)

func initOrSkip(t *testing.T) {
	t.Helper()
	if err := Init(); err != nil {
		t.Skipf("skipping: geoip database unavailable: %v", err)
	}
}

func TestGetIpGeoInfo_Cloudflare(t *testing.T) {
	initOrSkip(t)

	ip := "223.167.150.96"
	city := GetIpGeoInfo(ip)
	if city == nil {
		t.Fatalf("GetIpGeoInfo(%q) returned nil, expected a result", ip)
	}

	t.Logf("Country:   %s (%s)", city.Country.Country.Names.English, city.Country.Country.ISOCode)
	t.Logf("Emoji: %s", city.Emoji)
}

func TestGetIpGeoInfo_EmptyIP(t *testing.T) {
	initOrSkip(t)

	if city := GetIpGeoInfo(""); city != nil {
		t.Errorf("GetIpGeoInfo(\"\") = %+v, want nil", city)
	}
}

func TestGetIpGeoInfo_InvalidIP(t *testing.T) {
	initOrSkip(t)

	if city := GetIpGeoInfo("not-an-ip"); city != nil {
		t.Errorf("GetIpGeoInfo(\"not-an-ip\") = %+v, want nil", city)
	}
}

func TestGetIpGeoInfo_UninitializedReader(t *testing.T) {
	geoipReader = nil
	if city := GetIpGeoInfo("103.21.244.12"); city != nil {
		t.Errorf("expected nil when reader is uninitialized, got %+v", city)
	}
}
