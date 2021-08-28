package wpasupplicant

import "os"

func createLocalPath(iface string) (string, error) {
	local, err := os.CreateTemp(os.TempDir(), "wpa_"+iface+"_*")
	if err != nil {
		return "", err
	}

	local.Close()
	os.Remove(local.Name())

	return local.Name(), nil
}
