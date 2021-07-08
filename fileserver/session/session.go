package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var provides = make(map[string]SessionProvider)

func init() {
}

type SessionManager struct {
	cookieName  string
	lock        sync.Mutex
	provider    SessionProvider
	maxLifeTime int64
}

func NewSessionManager(provideName, cookieName string, maxLifeTime int64) (*SessionManager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	sessionManager := &SessionManager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}
	go func() {
		sessionManager.GC()
	}()
	return sessionManager, nil
}

func (manager *SessionManager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) Session {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ := manager.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(manager.maxLifeTime),
		}
		http.SetCookie(w, &cookie)
		return session
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ := manager.provider.SessionRead(sid)
		return session
	}
}

func (manager *SessionManager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	sid, _ := url.QueryUnescape(cookie.Value)
	manager.provider.SessionDestroy(sid)
	expireCookie := http.Cookie{
		Name:     manager.cookieName,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1,
	}
	http.SetCookie(w, &expireCookie)
}

func (manager *SessionManager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() {
		manager.GC()
	})
}

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
