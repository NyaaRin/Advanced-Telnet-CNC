package flood

import (
	"sync"
	"time"
)

type FloodLimiter struct {
	mu      sync.Mutex
	Entries map[Host]*RateLimitEntry
}

type Host struct {
	Target  uint32
	Netmask uint8
}

type RateLimitEntry struct {
	Count     int
	LastReset time.Time
}

const (
	MaxFloods     = 5
	ResetInterval = 30 * time.Minute
)

func NewFloodLimiter() *FloodLimiter {
	return &FloodLimiter{
		Entries: make(map[Host]*RateLimitEntry),
	}
}

func (f *FloodLimiter) Add(host uint32, netmask uint8) {
	f.mu.Lock()
	defer f.mu.Unlock()

	h := Host{Target: host, Netmask: netmask}
	entry, ok := f.Entries[h]
	if !ok {
		entry = &RateLimitEntry{}
		f.Entries[h] = entry
	}

	entry.Count++
	entry.LastReset = time.Now()
}

func (f *FloodLimiter) IsLimited(target uint32, netmask uint8) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	h := Host{Target: target, Netmask: netmask}
	entry, ok := f.Entries[h]
	if !ok {
		return false
	}

	if time.Since(entry.LastReset) >= ResetInterval {
		entry.Count = 0
		entry.LastReset = time.Now()
		return false
	}

	return entry.Count >= MaxFloods
}
