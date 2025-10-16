package tunnel

// maybeRunUserspace is a stub when built without the 'amnezia' tag.
// It never handles the tunnel and always falls back to the kernel driver path.
func maybeRunUserspace(confPath, serviceName string) (bool, error) {
	return false, nil
}
