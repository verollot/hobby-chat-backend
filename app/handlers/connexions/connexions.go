package connexions

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Connexions struct {
	connxs map[uuid.UUID]*websocket.Conn
	mutex  sync.RWMutex
}

var ErrUserNotAvailable error = errors.New("user not available")
var ErrUserAlreadyConnected error = errors.New("user already connected")
var ErrUserNotConnected error = errors.New("user not connected")

func NewConnexionsMap() *Connexions {
	return &Connexions{connxs: make(map[uuid.UUID]*websocket.Conn)}
}

func (c *Connexions) AddConn(id uuid.UUID, conn *websocket.Conn) error {
	fmt.Println(c.connxs)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.connxs[id]
	if ok {
		return fmt.Errorf("%w : [%s]", ErrUserAlreadyConnected, id.String())
	}
	c.connxs[id] = conn
	return nil
}

func (c *Connexions) RemoveConn(id uuid.UUID) error {
	fmt.Println(c.connxs)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.connxs[id]
	if !ok {
		return fmt.Errorf("%w : [%s]", ErrUserNotConnected, id.String())
	}
	delete(c.connxs, id)
	return nil
}

func (c *Connexions) GetConn(id uuid.UUID) (*websocket.Conn, error) {
	fmt.Println(c.connxs)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	conn, ok := c.connxs[id]
	if !ok {
		return nil, fmt.Errorf("%w : [%s]", ErrUserNotAvailable, id.String())
	}
	return conn, nil
}
