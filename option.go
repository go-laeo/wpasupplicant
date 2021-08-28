package wpasupplicant

import (
	"net"
	"os"
)

// Option is the most commonly practice to configure the struct.
type Option func(conn *unixgram) error

func CustomUnixgram(local, remote string) Option {
	return func(conn *unixgram) (err error) {
		// We have to ensure the exist connection is properly
		// closed and the local socket file also removed.
		if conn.conn != nil {
			conn.conn.Close()
			os.Remove(conn.local)
		}

		conn.local = local
		conn.remote = remote
		conn.conn, err = net.DialUnix("unixgram",
			&net.UnixAddr{Name: local, Net: "unixgram"},
			&net.UnixAddr{Name: remote, Net: "unixgram"})
		return err
	}
}
