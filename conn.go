package wpasupplicant

// Conn is a connection to wpa_supplicant over one of its communication
// channels.
type Conn interface {
	// Close closes the unixgram connection
	Close() error

	// Ping tests the connection.  It returns nil if wpa_supplicant is
	// responding.
	Ping() error

	// AddNetwork creates an empty network configuration. Returns the network
	// ID.
	AddNetwork() (int, error)

	// SetNetwork configures a network property. Returns error if the property
	// configuration failed.
	// Value's type must one of int, string and []byte. The int type always shown
	// without double quotes. The string type always shown with double quotes except
	// the variable name is key_mgmt. The []byte type may only uses for ssid, maybe
	// useful when it contains non-ascii encoded chars.
	SetNetwork(networkID int, field string, value interface{}) error

	// EnableNetwork enables a network. Returns error if the command fails.
	EnableNetwork(int) error

	// EnableAllNetworks enables all configured networks. Returns error if the command fails.
	EnableAllNetworks() error

	// SelectNetwork selects a network (and disables the others).
	SelectNetwork(int) error

	// DisableNetwork disables a network.
	DisableNetwork(int) error

	// RemoveNetwork removes a network from the configuration.
	RemoveNetwork(int) error

	// RemoveAllNetworks removes all networks (basically running `REMOVE_NETWORK all`).
	// Returns error if command fails.
	RemoveAllNetworks() error

	// SaveConfig stores the current network configuration to disk.
	SaveConfig() error

	// Reconfigure sends a RECONFIGURE command to the wpa_supplicant. Returns error when
	// command fails.
	Reconfigure() error

	// Reassociate sends a REASSOCIATE command to the wpa_supplicant. Returns error when
	// command fails.
	Reassociate() error

	// Reconnect sends a RECONNECT command to the wpa_supplicant. Returns error when
	// command fails.
	Reconnect() error

	// ListNetworks returns the currently configured networks.
	ListNetworks() ([]ConfiguredNetwork, error)

	// Status returns current wpa_supplicant status
	Status() (StatusResult, error)

	// Scan triggers a new scan. Returns error if the wpa_supplicant does not
	// return OK.
	Scan() error

	// ScanResult returns the latest scanning results.  It returns a slice
	// of scanned BSSs, and/or a slice of errors representing problems
	// communicating with wpa_supplicant or parsing its output.
	ScanResults() ([]ScanResult, []error)

	EventQueue() chan WPAEvent
}
