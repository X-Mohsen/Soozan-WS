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

var (
	userChannels = struct {
		sync.RWMutex
		channels map[float64]*Channel
	}{
		channels: make(map[float64]*Channel),
	}
)

/*
Retrieves or creates a private channel for a given userID.
Used after upgrading connection
*/
func GetOrCreateUserChannel(userID float64) *Channel {
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

// Returns list of connections based on userID
func GetUserConnections(userID float64) ([]net.Conn, bool) {
	userChannels.RLock()
	ch, exists := userChannels.channels[userID]
	userChannels.RUnlock()

	if !exists {
		return nil, false
	}

	ch.mu.Lock()
	connections := make([]net.Conn, 0, len(ch.connections))
	for conn := range ch.connections {
		connections = append(connections, conn)
	}

	return connections, true
}

// Previous function but for multiple users (cleaner for loop in main)
func GetMultipleUserConnections(userIDs [2]float64) []net.Conn {
	connCollection := make([]net.Conn, 0)

	for i := range userIDs {
		userConns, exists := GetUserConnections(userIDs[i])
		if exists {
			connCollection = append(connCollection, userConns...)
		}
	}

	return connCollection
}
