package image

import "os"

const (
	defaultBusyBoxImage = "docker.io/busybox:latest"
	defaultSleeperImage = "docker.io/ubuntu:latest"
	envBusyBoxImage     = "DFA_BUSYBOX_IMAGE"
	envSleeperImage     = "DFA_SLEEPER_IMAGE"
)

func BusyBox() string {
	if img, isSet := os.LookupEnv(envBusyBoxImage); isSet {
		return img
	}

	return defaultBusyBoxImage
}

func Sleeper() string {
	if img, isSet := os.LookupEnv(envSleeperImage); isSet {
		return img
	}

	return defaultSleeperImage
}
