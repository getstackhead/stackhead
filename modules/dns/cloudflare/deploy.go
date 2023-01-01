package dns_cloudflare

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"robpike.io/filter"

	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

func GetDomain(domain string) string {
	var re = regexp.MustCompile(`(?P<Sub>.*)\.(?P<Main>.*\..*)`)
	match := re.FindStringSubmatch(domain)
	return match[re.SubexpIndex("Main")]
}

func GetSubdomain(domain string) string {
	var re = regexp.MustCompile(`(?P<Sub>.*)\.(?P<Main>.*\..*)`)
	match := re.FindStringSubmatch(domain)
	subIndex := re.SubexpIndex("Sub")
	if len(match) <= subIndex {
		// no subdomain found
		return "@"
	}
	return match[subIndex]
}

func (m Module) Deploy(_modulesSettings interface{}) error {
	moduleSettings, err := system.UnpackModuleSettings[ModuleSettings](_modulesSettings)
	if err != nil {
		return fmt.Errorf("unable to load module settings: " + err.Error())
	}
	if len(moduleSettings.ApiToken) == 0 {
		return fmt.Errorf("missing Cloudflare API token. Supply one in module settings.")
	}
	// Construct a new API object using a global API key
	api, err := cloudflare.NewWithAPIToken(moduleSettings.ApiToken)
	if err != nil {
		return err
	}

	ctx := context.Background()
	domains := filter.Choose(system.Context.Project.Domains, func(d project.Domains) bool {
		return d.DNS.Provider == "cloudflare"
	}).([]project.Domains)

	proxied := false
	for _, domain := range domains {
		zoneName := GetDomain(domain.Domain)
		zoneId, err := api.ZoneIDByName(zoneName)
		if err != nil {
			return fmt.Errorf("unable to find Zone ID for name \"" + zoneName + "\": " + err.Error())
		}
		dnsRecord := constructRecord(domain.Domain)
		dnsRecord.Proxied = &proxied
		// API scope: dns_records:edit
		_, err = api.CreateDNSRecord(ctx, zoneId, dnsRecord)
		// todo: add to created resources
		if err != nil {
			if strings.Contains(err.Error(), "Record already exists.") {
				return nil
			}
			return fmt.Errorf("Cloudflare error: " + err.Error())
		}
	}
	return nil
}
