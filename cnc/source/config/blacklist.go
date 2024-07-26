package config

import (
	"encoding/binary"
	"net"
)

func (blacklist *blacklist) IsBlacklisted(target uint32, netmask uint8) bool {
	for _, t := range blacklist.Blacklist {
		// Parse prefix
		ip := net.ParseIP(t.Prefix).To4() // Ensure IPv4 address
		if ip == nil {
			continue // Skip invalid prefixes
		}
		iWhitelistPrefix := binary.BigEndian.Uint32(ip)

		if netmask > t.Netmask { // Whitelist is less specific than attack target
			if Netshift(iWhitelistPrefix, t.Netmask) == Netshift(target, t.Netmask) {
				return true
			}
		} else if netmask < t.Netmask { // Attack target is less specific than whitelist
			if (target >> netmask) == (iWhitelistPrefix >> netmask) {
				return true
			}
		} else { // Both target and whitelist have the same prefix
			if iWhitelistPrefix == target {
				return true
			}
		}
	}

	return false
}

func Netshift(prefix uint32, netmask uint8) uint32 {
	return uint32(prefix >> (32 - netmask))
}
