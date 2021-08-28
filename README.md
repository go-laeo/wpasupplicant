# wpasupplicant

![build.yaml](https://github.com/go-laeo/wpasupplicant/actions/workflows/build.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/go-laeo/wpasupplicant.svg)](https://pkg.go.dev/github.com/go-laeo/wpasupplicant) ![golangci.yaml](https://github.com/go-laeo/wpasupplicant/actions/workflows/golangci-lint.yaml/badge.svg)

Golang interface for talking to wpa_supplicant.

## Install

```
go get github.com/go-laeo/wpasupplicant
```

## Example

```
// Prints the BSSID (MAC address) and SSID of each access point in range:
w, err := wpasupplicant.Connect(context.Background(), "wlan0")
if err != nil {
	panic(err)
}
for _, bss := range w.ScanResults() {
	fmt.Fprintf("%s\t%s\n", bss.BSSID(), bss.SSID())
}
```

## License

Three-clause BSD.  See LICENSE.txt.

## Thanks

* [@Dave Pifke](https://github.com/dpifke/golang-wpasupplicant)
