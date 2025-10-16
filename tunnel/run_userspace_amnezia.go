//go:build amnezia

package tunnel

import (
	"golang.org/x/sys/windows/svc"
	"golang.zx2c4.com/wireguard/windows/conf"
)

// maybeRunUserspace runs the userspace tunnel service if built with 'amnezia' tag.
func maybeRunUserspace(confPath, serviceName string) (bool, error) {
	return true, svc.Run(serviceName, &userspaceTunnelService{confPath})
}
