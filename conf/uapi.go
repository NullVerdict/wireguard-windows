/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.
 */

package conf

import (
	"fmt"
	"strings"
)

// ToUAPI builds a UAPI configuration string compatible with userspace backends.
// This includes AmneziaWG obfuscation parameters (J*, S*, H*, i*, j*, itime).
func (conf *Config) ToUAPI() (uapi string, dnsErr error) {
	var output strings.Builder
	output.WriteString(fmt.Sprintf("private_key=%s\n", conf.Interface.PrivateKey.HexString()))

	if conf.Interface.ListenPort > 0 {
		output.WriteString(fmt.Sprintf("listen_port=%d\n", conf.Interface.ListenPort))
	}

	if conf.Interface.JunkPacketCount > 0 {
		output.WriteString(fmt.Sprintf("jc=%d\n", conf.Interface.JunkPacketCount))
	}
	if conf.Interface.JunkPacketMinSize > 0 {
		output.WriteString(fmt.Sprintf("jmin=%d\n", conf.Interface.JunkPacketMinSize))
	}
	if conf.Interface.JunkPacketMaxSize > 0 {
		output.WriteString(fmt.Sprintf("jmax=%d\n", conf.Interface.JunkPacketMaxSize))
	}
	if conf.Interface.InitPacketJunkSize > 0 {
		output.WriteString(fmt.Sprintf("s1=%d\n", conf.Interface.InitPacketJunkSize))
	}
	if conf.Interface.ResponsePacketJunkSize > 0 {
		output.WriteString(fmt.Sprintf("s2=%d\n", conf.Interface.ResponsePacketJunkSize))
	}
	if conf.Interface.CookieReplyPacketJunkSize > 0 {
		output.WriteString(fmt.Sprintf("s3=%d\n", conf.Interface.CookieReplyPacketJunkSize))
	}
	if conf.Interface.TransportPacketJunkSize > 0 {
		output.WriteString(fmt.Sprintf("s4=%d\n", conf.Interface.TransportPacketJunkSize))
	}
	if conf.Interface.InitPacketMagicHeader > 0 {
		output.WriteString(fmt.Sprintf("h1=%d\n", conf.Interface.InitPacketMagicHeader))
	}
	if conf.Interface.ResponsePacketMagicHeader > 0 {
		output.WriteString(fmt.Sprintf("h2=%d\n", conf.Interface.ResponsePacketMagicHeader))
	}
	if conf.Interface.UnderloadPacketMagicHeader > 0 {
		output.WriteString(fmt.Sprintf("h3=%d\n", conf.Interface.UnderloadPacketMagicHeader))
	}
	if conf.Interface.TransportPacketMagicHeader > 0 {
		output.WriteString(fmt.Sprintf("h4=%d\n", conf.Interface.TransportPacketMagicHeader))
	}
	for key, value := range conf.Interface.IPackets {
		output.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	for key, value := range conf.Interface.JPackets {
		output.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	if conf.Interface.ITime > 0 {
		output.WriteString(fmt.Sprintf("itime=%d\n", conf.Interface.ITime))
	}

	if len(conf.Peers) > 0 {
		output.WriteString("replace_peers=true\n")
	}
	for _, peer := range conf.Peers {
		output.WriteString(fmt.Sprintf("public_key=%s\n", peer.PublicKey.HexString()))
		if !peer.PresharedKey.IsZero() {
			output.WriteString(fmt.Sprintf("preshared_key=%s\n", peer.PresharedKey.HexString()))
		}
		if !peer.Endpoint.IsEmpty() {
			var resolvedIP string
			resolvedIP, dnsErr = resolveHostname(peer.Endpoint.Host)
			if dnsErr != nil {
				return
			}
			resolvedEndpoint := Endpoint{resolvedIP, peer.Endpoint.Port}
			output.WriteString(fmt.Sprintf("endpoint=%s\n", resolvedEndpoint.String()))
		}
		output.WriteString(fmt.Sprintf("persistent_keepalive_interval=%d\n", peer.PersistentKeepalive))
		if len(peer.AllowedIPs) > 0 {
			output.WriteString("replace_allowed_ips=true\n")
			for _, address := range peer.AllowedIPs {
				output.WriteString(fmt.Sprintf("allowed_ip=%s\n", address.String()))
			}
		}
	}
	return output.String(), nil
}
