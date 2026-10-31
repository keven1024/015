package utils

import (
	"net"
	"net/mail"
	"net/url"
	"slices"
	"strings"

	"github.com/samber/lo"
)

func IsVaildEmail(email string) bool {
	parsedEmail, err := mail.ParseAddress(email)
	return err == nil && parsedEmail.Address == email
}

func IsVaildURL(rawURL string) (bool, *url.URL) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false, nil
	}
	return parsedURL.Scheme != "" &&
		parsedURL.Hostname() != "" &&
		slices.Contains(
			[]string{"http", "https"},
			strings.ToLower(parsedURL.Scheme),
		), parsedURL
}

func IsVaildIP(ip string) bool {
	ips, err := net.LookupIP(ip)
	if err != nil {
		return false
	}
	return lo.EveryBy(ips, func(ip net.IP) bool {
		return !ip.IsPrivate() &&
			!ip.IsLoopback() &&
			!ip.IsLinkLocalUnicast() &&
			!ip.IsUnspecified() &&
			!ip.IsMulticast()
	})
}
