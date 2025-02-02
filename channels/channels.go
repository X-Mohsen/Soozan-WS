package channels

import (
	"net"
	"sync"
)

type Channel struct {
	connections map[net.Conn]bool
	mu          sync.Mutex
}

func (c *Channel) Add(conn net.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connections[conn] = true
}

func (c *Channel) Remove(conn net.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.connections, conn)
}

// userChannels holds a mapping from userID to their private Channel.
var (
	userChannels = struct {
		sync.RWMutex
		channels map[float64]*Channel
	}{
		channels: make(map[float64]*Channel),
	}
)

// getUserChannel retrieves or creates a private channel for a given userID.
func GetUserChannel(userID float64) *Channel {
	userChannels.RLock()
	ch, exists := userChannels.channels[userID]
	userChannels.RUnlock()

	if !exists {
		userChannels.Lock()
		// Double-check in case it was created in the meantime.
		ch, exists = userChannels.channels[userID]
		if !exists {
			ch = &Channel{
				connections: make(map[net.Conn]bool),
			}
			userChannels.channels[userID] = ch
		}
		userChannels.Unlock()
	}

	return ch
}
