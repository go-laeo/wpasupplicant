package wpasupplicant

import "net"

// ScanResult is a scanned BSS.
type ScanResult interface {
	// BSSID is the MAC address of the BSS.
	BSSID() net.HardwareAddr

	// SSID is the SSID of the BSS.
	SSID() string

	// Frequency is the frequency, in Mhz, of the BSS.
	Frequency() int

	// RSSI is the received signal strength, in dB, of the BSS.
	RSSI() int

	// Flags is an array of flags, in string format, returned by the
	// wpa_supplicant SCAN_RESULTS command.  Future versions of this code
	// will parse these into something more meaningful.
	Flags() []string
}

// scanResult is a package-private implementation of ScanResult.
type scanResult struct {
	bssid     net.HardwareAddr
	ssid      string
	frequency int
	rssi      int
	flags     []string
}

func (r *scanResult) BSSID() net.HardwareAddr { return r.bssid }
func (r *scanResult) SSID() string            { return r.ssid }
func (r *scanResult) Frequency() int          { return r.frequency }
func (r *scanResult) RSSI() int               { return r.rssi }
func (r *scanResult) Flags() []string         { return r.flags }
