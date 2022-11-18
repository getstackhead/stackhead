package dns_cloudflare

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	logger "github.com/sirupsen/logrus"
	"robpike.io/filter"

	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

func (m Module) Destroy(_modulesSettings interface{}) error {
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

	for _, domain := range domains {
		zoneName := GetDomain(domain.Domain)
		zoneId, err := api.ZoneIDByName(zoneName)
		if err != nil {
			return fmt.Errorf("unable to find Zone ID for name \"" + zoneName + "\": " + err.Error())
		}
		dnsRecord := constructRecord(domain.Domain)
		dnsRecord.Name = domain.Domain // when looking up the records name is the full domain name
		if moduleSettings.DisableSafeMode {
			dnsRecord.Content = ""
			dnsRecord.Type = ""
		}
		// API scope: dns_records:read
		foundRecords, err := api.DNSRecords(ctx, zoneId, dnsRecord)
		if err != nil {
			return fmt.Errorf("unable to find DNS record: " + err.Error())
		}
		if len(foundRecords) == 0 {
			if moduleSettings.DisableSafeMode {
				logger.Warning("no DNS records were found for name \"" + dnsRecord.Name + "\"")
			} else {
				logger.Warning("no DNS " + dnsRecord.Type + " record was found for name \"" + dnsRecord.Name + "\" and content \"" + dnsRecord.Content + "\"")
			}
			return nil
		}
		if !moduleSettings.DisableSafeMode && len(foundRecords) > 1 {
			return fmt.Errorf("more than 1 DNS record was found")
		}

		// API scope: dns_records:edit
		for _, record := range foundRecords {
			logger.Info("Removing DNS " + record.Type + " record for name \"" + record.Name + "\" with content \"" + record.Content + "\"")
			if err := api.DeleteDNSRecord(ctx, zoneId, record.ID); err != nil {
				return err
			}
		}
	}
	return nil
}
