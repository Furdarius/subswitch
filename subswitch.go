package subswitch

import (
	"net/http"
)

// SubSwitcher allow separate mux instance for subdomains
type SubSwitcher struct {
	next       http.Handler
	subdomains map[string]http.Handler
}

func (s *SubSwitcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := s.subdomains[subDomain(r.Host)]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		s.next.ServeHTTP(w, r)
	}
}

// New return instance of SubSwitcher
func New(h http.Handler, subdomains map[string]http.Handler) *SubSwitcher {
	return &SubSwitcher{h, subdomains}
}

func subDomain(host string) string {
	for i, cnt := len(host)-1, 2; i >= 0 && cnt > 0; i-- {
		if host[i] == '.' {
			cnt--
		}

		if cnt == 0 {
			return host[:i]
		}
	}

	return ""
}
