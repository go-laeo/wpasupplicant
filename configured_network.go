package wpasupplicant

// ConfiguredNetwork is a configured network (from LIST_NETWORKS)
type ConfiguredNetwork interface {
	NetworkID() string
	SSID() string
	BSSID() string
	Flags() []string
}

type configuredNetwork struct {
	networkID string
	ssid      string
	bssid     string // Since bssid can be any
	flags     []string
}

func (r *configuredNetwork) NetworkID() string { return r.networkID }
func (r *configuredNetwork) BSSID() string     { return r.bssid }
func (r *configuredNetwork) SSID() string      { return r.ssid }
func (r *configuredNetwork) Flags() []string   { return r.flags }
