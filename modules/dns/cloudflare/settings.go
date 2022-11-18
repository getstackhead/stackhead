package dns_cloudflare

type ModuleSettings struct {
	ApiToken        string `json:"api_token"`
	DisableSafeMode bool   `json:"disable_safemode"`
}
