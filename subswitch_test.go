package subswitch

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubDomain(t *testing.T) {
	tests := []struct {
		host, expected string
	}{
		{"domain.ru", ""},
		{"mydomain.ru", ""},
		{"DOMAIN.com", ""},
		{"domain.com:8000", ""},
		{"my_domain-2123.ru:12323", ""},
		{"domain.ru.com", "domain"},
		{"go.go.go.go", "go.go"},
		{"go.go", ""},
		{"go.go:2323", ""},
		{"go.go-go.go:2323", "go"},
		{"wefkw.wefkwf3", ""},
		{"fast:213", ""},
		{"main.site.ru", "main"},
		{"fff.fwef3.we32r.wef23r.wef23.d2w1", "fff.fwef3.we32r.wef23r"},
	}

	for _, test := range tests {
		actual := subDomain(test.host)

		if actual != test.expected {
			t.Errorf("subDomain(%s) failed: expected %s, actual %s", test.host, test.expected, actual)
		}
	}
}

// getTestHandler returns a http.HandlerFunc for testing http middleware
func getTestHandler(output string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, output)
	}

	return http.HandlerFunc(fn)
}

func TestNew(t *testing.T) {
	domain := "domain.ru"

	subdomains := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"subdomain",
		"good.wefewf",
		"f23123",
	}

	m := make(map[string]http.Handler)

	for _, subdomain := range subdomains {
		m[subdomain] = getTestHandler(subdomain)
	}

	h := getTestHandler(domain)
	instance := New(h, m)

	for _, subdomain := range subdomains {
		_, has := instance.subdomains[subdomain]

		if !has {
			t.Errorf("Subdomain \"%s\" wasn't found in SubSwitch.subdomains.", subdomain)
		}
	}
}

func TestSubSwitch_ServeHTTP(t *testing.T) {
	domain := "domain.ru"

	subdomains := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"subdomain",
		"good.wefewf",
		"f23123",
	}

	m := make(map[string]http.Handler)

	for _, subdomain := range subdomains {
		m[subdomain] = getTestHandler(subdomain)
	}

	h := getTestHandler(domain)
	ss := New(h, m)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	ss.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ss.ServeHTTP failed: status expected %d, actual %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != domain {
		t.Errorf("ss.ServeHTTP failed: body expected %s, actual %s", domain, rr.Body.String())
	}

	for _, subdomain := range subdomains {
		rr = httptest.NewRecorder()
		req.Host = subdomain + "." + domain

		ss.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("ss.ServeHTTP failed: status expected %d, actual %d", http.StatusOK, rr.Code)
		}

		if rr.Body.String() != subdomain {
			t.Errorf("ss.ServeHTTP failed: body expected %s, actual %s", subdomain, rr.Body.String())
		}
	}
}
