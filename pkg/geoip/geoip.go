package geoip

import (
	"embed"
	"net/netip"

	"github.com/enescakir/emoji"
	"github.com/oschwald/geoip2-golang/v2"
)

//go:embed resource/*.mmdb
var dbFS embed.FS

var geoipReader *geoip2.Reader

func Init() error {
	data, err := dbFS.ReadFile("resource/GeoLite2-Country.mmdb")
	if err != nil {
		return err
	}

	geoipReader, err = geoip2.OpenBytes(data)
	return err
}

type IpGeoInfo struct {
	Country *geoip2.Country
	Emoji   string
}

func GetIpGeoInfo(ip string) *IpGeoInfo {
	if geoipReader == nil || ip == "" {
		return nil
	}

	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil
	}

	country, err := geoipReader.Country(ipAddr)
	if err != nil || country.Country.ISOCode == "" {
		return nil
	}

	emoji, err := emoji.CountryFlag(country.Country.ISOCode)
	if err != nil {
		return nil
	}

	return &IpGeoInfo{
		Country: country,
		Emoji:   emoji.String(),
	}
}
