# wpasupplicant

![build.yaml](https://github.com/go-laeo/wpasupplicant/actions/workflows/build.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/go-laeo/wpasupplicant.svg)](https://pkg.go.dev/github.com/go-laeo/wpasupplicant) ![golangci.yaml](https://github.com/go-laeo/wpasupplicant/actions/workflows/golangci-lint.yaml/badge.svg)

Golang interface for talking to wpa_supplicant.

At the moment, this simply provides an interface for fetching wifi scan
results.  More functionality (probably) coming soon.

## Example

```
import (
	"fmt"

	"github.com/go-laeo/wpasupplicant"
)

// Prints the BSSID (MAC address) and SSID of each access point in range:
w, err := wpasupplicant.Unixgram("wlan0")
if err != nil {
	panic(err)
}
for _, bss := range w.ScanResults() {
	fmt.Fprintf("%s\t%s\n", bss.BSSID(), bss.SSID())
}
```

## Downloading

If you use this library in your own code, please use the canonical URL in your
Go code, instead of Github:

```
go get github.com/go-laeo/wpasupplicant
```

Or (until I finish setting up the self-hosted repository):

```
# From the root of your project:
git submodule add https://github.com/dpifke/golang-wpasupplicant vendor/github.com/go-laeo/wpasupplicant
```

Then:

```
import (
        "github.com/go-laeo/wpasupplicant"
)
```

As opposed to the pifke.org URL, I make no guarantee this Github repository
will exist or be up-to-date in the future.

## Documentation

Available on [godoc.org](https://godoc.org/github.com/go-laeo/wpasupplicant).

## License

Three-clause BSD.  See LICENSE.txt.

Contact me if you want to use this code under different terms.

## Thanks

* [@Dave Pifke](https://github.com/dpifke/golang-wpasupplicant)
