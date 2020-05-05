package cas

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/net/publicsuffix"
)

func TestAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		if r.RequestURI != "/login" {
			t.Errorf("expected request URL to be /login, got %s", r.RequestURI)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("expected parsable form, got err %s", err)
		}

		username := r.PostForm.Get("username")
		if username != "someUser" {
			t.Errorf("expected username to be 'someUser', got '%s'", username)
		}

		password := r.PostForm.Get("password")
		if password != "somePass" {
			t.Errorf("expected password to be 'somePass', got '%s'", password)
		}

		service := r.PostForm.Get("service")
		if service != "https://spinup.example.com/login" {
			t.Errorf("expected service to be 'https://spinup.example.com/login', got '%s'", service)
		}
	}))
	defer ts.Close()

	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	casUrl = ts.URL
	if err := Auth("someUser", "somePass", "https://spinup.example.com/login", httpClient); err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
}
