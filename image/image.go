package image

import "os"

const (
	DefaultBusyBoxImage = "docker.io/busybox:latest"
	DefaultSleeperImage = "docker.io/ubuntu:latest"
	EnvBusyBoxImage     = "DFA_BUSYBOX_IMAGE"
	EnvSleeperImage     = "DFA_SLEEPER_IMAGE"
)

func BusyBox() string {
	if img, isSet := os.LookupEnv(EnvBusyBoxImage); isSet {
		return img
	}

	return DefaultBusyBoxImage
}

func Sleeper() string {
	if img, isSet := os.LookupEnv(EnvSleeperImage); isSet {
		return img
	}

	return DefaultSleeperImage
}
