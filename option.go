package wpasupplicant

import (
	"net"
)

// Option is the most commonly practice to configure the struct.
type Option func(conn *unixgram) error

func CustomUnixgram(local, remote string) Option {
	return func(conn *unixgram) (err error) {
		conn.local = local
		conn.remote = remote
		conn.conn, err = net.DialUnix("unixgram",
			&net.UnixAddr{Name: local, Net: "unixgram"},
			&net.UnixAddr{Name: remote, Net: "unixgram"})
		return err
	}
}
