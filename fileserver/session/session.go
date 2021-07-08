package session

import (
	"fmt"
	"sync"
)

//var GoSessionManager

var provides = make(map[string]SessionProvider)

func init() {

}

//Register register SessionProvider
func Register(name string, provider SessionProvider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, ok := provides[name]; ok {
		panic("session: Register called twice for provider")
	}
	provides[name] = provider
}

func NewSessionManager(provideName, cookieName string, maxLifeTime int64) (*SessionManager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &SessionManager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

type SessionManager struct {
	cookieName  string
	lock        sync.Mutex
	provider    SessionProvider
	maxLifeTime int64
}

//interface

type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // back current sessionID
}

type SessionProvider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}
