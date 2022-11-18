package dns_cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
	"regexp"

	"github.com/getstackhead/stackhead/system"
)

func constructRecord(domain string) cloudflare.DNSRecord {
	dnsTarget := system.Context.TargetHost.String()
	ipv4RegEx := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	dnsType := "AAAA"
	if len(ipv4RegEx.FindStringIndex(dnsTarget)) > 0 {
		dnsType = "A"
	}
	return cloudflare.DNSRecord{
		Type:    dnsType,
		Name:    GetSubdomain(domain),
		Content: dnsTarget,
	}
}
