package wpasupplicant

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// ParseError is returned when we can't parse the wpa_supplicant response.
// Some functions may return multiple ParseErrors.
type ParseError struct {
	// Line is the line of output from wpa_supplicant which we couldn't
	// parse.
	Line string

	// Err is any nested error.
	Err error
}

func (err *ParseError) Error() string {
	b := &bytes.Buffer{}
	b.WriteString("failed to parse wpa_supplicant response")

	if err.Line != "" {
		fmt.Fprintf(b, ": %q", err.Line)
	}

	if err.Err != nil {
		fmt.Fprintf(b, ": %s", err.Err.Error())
	}

	return b.String()
}

func parseListNetworksResult(resp io.Reader) (res []ConfiguredNetwork, err error) {
	s := bufio.NewScanner(resp)
	if !s.Scan() {
		return nil, &ParseError{}
	}

	networkIDCol, ssidCol, bssidCol, flagsCol, maxCol := -1, -1, -1, -1, -1
	for n, col := range strings.Split(s.Text(), " / ") {
		switch col {
		case "network id":
			networkIDCol = n
		case "ssid":
			ssidCol = n
		case "bssid":
			bssidCol = n
		case "flags":
			flagsCol = n
		}

		maxCol = n
	}

	for s.Scan() {
		ln := s.Text()
		fields := strings.Split(ln, "\t")
		if len(fields) < maxCol {
			return nil, &ParseError{Line: ln}
		}

		var networkID string
		if networkIDCol != -1 {
			networkID = fields[networkIDCol]
		}

		var ssid string
		if ssidCol != -1 {
			ssid = fields[ssidCol]
		}

		var bssid string
		if bssidCol != -1 {
			bssid = fields[bssidCol]
		}

		var flags []string
		if flagsCol != -1 {
			if len(fields[flagsCol]) >= 2 && fields[flagsCol][0] == '[' && fields[flagsCol][len(fields[flagsCol])-1] == ']' {
				flags = strings.Split(fields[flagsCol][1:len(fields[flagsCol])-1], "][")
			}
		}

		res = append(res, &configuredNetwork{
			networkID: networkID,
			ssid:      ssid,
			bssid:     bssid,
			flags:     flags,
		})
	}

	return res, nil
}

func parseStatusResults(resp io.Reader) (StatusResult, error) {
	s := bufio.NewScanner(resp)

	res := &statusResult{}

	for s.Scan() {
		ln := s.Text()
		fields := strings.Split(ln, "=")
		if len(fields) != 2 {
			continue
		}

		switch fields[0] {
		case "wpa_state":
			res.wpaState = fields[1]
		case "key_mgmt":
			res.keyMgmt = fields[1]
		case "ip_address":
			res.ipAddr = fields[1]
		case "ssid":
			res.ssid = fields[1]
		case "address":
			res.address = fields[1]
		case "bssid":
			res.bssid = fields[1]
		case "freq":
			res.freq = fields[1]
		}
	}

	return res, nil
}

// parseScanResults parses the SCAN_RESULTS output from wpa_supplicant.  This
// is split out from ScanResults() to make testing easier.
func parseScanResults(resp io.Reader) (res []ScanResult, errs []error) {
	// In an attempt to make our parser more resilient, we start by
	// parsing the header line and using that to determine the column
	// order.
	s := bufio.NewScanner(resp)
	if !s.Scan() {
		errs = append(errs, &ParseError{})
		return
	}
	bssidCol, freqCol, rssiCol, flagsCol, ssidCol, maxCol := -1, -1, -1, -1, -1, -1
	for n, col := range strings.Split(s.Text(), " / ") {
		switch col {
		case "bssid":
			bssidCol = n
		case "frequency":
			freqCol = n
		case "signal level":
			rssiCol = n
		case "flags":
			flagsCol = n
		case "ssid":
			ssidCol = n
		}
		maxCol = n
	}

	var err error
	for s.Scan() {
		ln := s.Text()
		fields := strings.Split(ln, "\t")
		if len(fields) < maxCol {
			errs = append(errs, &ParseError{Line: ln})
			continue
		}

		var bssid net.HardwareAddr
		if bssidCol != -1 {
			if bssid, err = net.ParseMAC(fields[bssidCol]); err != nil {
				errs = append(errs, &ParseError{Line: ln, Err: err})
				continue
			}
		}

		var freq int
		if freqCol != -1 {
			if freq, err = strconv.Atoi(fields[freqCol]); err != nil {
				errs = append(errs, &ParseError{Line: ln, Err: err})
				continue
			}
		}

		var rssi int
		if rssiCol != -1 {
			if rssi, err = strconv.Atoi(fields[rssiCol]); err != nil {
				errs = append(errs, &ParseError{Line: ln, Err: err})
				continue
			}
		}

		var flags []string
		if flagsCol != -1 {
			if len(fields[flagsCol]) >= 2 && fields[flagsCol][0] == '[' && fields[flagsCol][len(fields[flagsCol])-1] == ']' {
				flags = strings.Split(fields[flagsCol][1:len(fields[flagsCol])-1], "][")
			}
		}

		var ssid string
		if ssidCol != -1 {
			ssid = fields[ssidCol]
		}

		res = append(res, &scanResult{
			bssid:     bssid,
			frequency: freq,
			rssi:      rssi,
			flags:     flags,
			ssid:      ssid,
		})
	}

	return
}
