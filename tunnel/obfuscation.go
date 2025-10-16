/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2025.
 */

package tunnel

import "golang.zx2c4.com/wireguard/windows/conf"

// needsUserspaceObfuscation returns true if the configuration enables any
// AmneziaWG obfuscation parameters that require userspace backend.
func needsUserspaceObfuscation(c *conf.Config) bool {
	ci := c.Interface
	if ci.JunkPacketCount > 0 || ci.JunkPacketMinSize > 0 || ci.JunkPacketMaxSize > 0 {
		return true
	}
	if ci.InitPacketJunkSize > 0 || ci.ResponsePacketJunkSize > 0 || ci.CookieReplyPacketJunkSize > 0 || ci.TransportPacketJunkSize > 0 {
		return true
	}
	if ci.InitPacketMagicHeader > 0 || ci.ResponsePacketMagicHeader > 0 || ci.UnderloadPacketMagicHeader > 0 || ci.TransportPacketMagicHeader > 0 {
		return true
	}
	if len(ci.IPackets) > 0 || len(ci.JPackets) > 0 || ci.ITime > 0 {
		return true
	}
	return false
}
