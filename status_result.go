package wpasupplicant

type StatusResult interface {
	WPAState() string
	KeyMgmt() string
	IPAddr() string
	SSID() string
	Address() string
	BSSID() string
	Freq() string
}

type statusResult struct {
	wpaState string
	keyMgmt  string
	ipAddr   string
	ssid     string
	address  string
	bssid    string
	freq     string
}

func (s *statusResult) WPAState() string { return s.wpaState }
func (s *statusResult) KeyMgmt() string  { return s.keyMgmt }
func (s *statusResult) IPAddr() string   { return s.ipAddr }
func (s *statusResult) SSID() string     { return s.ssid }
func (s *statusResult) Address() string  { return s.address }
func (s *statusResult) BSSID() string    { return s.bssid }
func (s *statusResult) Freq() string     { return s.freq }
