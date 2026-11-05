package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/batijo/poll-scraper/config"
)

const (
	testUser     = "Maistobankas"
	testPassword = "MB123"
	testToken    = "secret-token"
)

type authServer struct {
	username   string
	password   string
	loginPosts int
	valid      bool
}

func newAuthServer(t *testing.T, username, password string) (*httptest.Server, *authServer) {
	t.Helper()

	state := &authServer{username: username, password: password}
	mux := http.NewServeMux()
	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || !state.valid || cookie.Value != "valid" {
			http.Redirect(w, r, "/Login?ReturnUrl=%2Fdata", http.StatusFound)
			return
		}
		writeDataPage(w, r)
	})
	mux.HandleFunc("/public", writeDataPage)
	mux.HandleFunc("/Login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderLoginPage(w)
			return
		}

		state.loginPosts++
		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse login form: %v", err)
		}
		if r.PostForm.Get("token") != testToken ||
			r.PostForm.Get("Login1$UserName") != state.username ||
			r.PostForm.Get("Login1$Password") != state.password ||
			r.PostForm.Get("Login1$Button1") != "Prisijungti" {
			renderLoginPage(w)
			return
		}

		state.valid = true
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "valid", Path: "/", HttpOnly: true})
		http.Redirect(w, r, r.URL.Query().Get("ReturnUrl"), http.StatusFound)
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, state
}

func renderLoginPage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<html><body>
		<form method="post" action="/Login?ReturnUrl=%%2Fdata">
			<input type="hidden" name="token" value="%s">
			<input type="text" name="Login1$UserName" value="placeholder">
			<input type="password" name="Login1$Password" value="placeholder">
			<input type="submit" name="Login1$Button1" value="Prisijungti">
		</form>
	</body></html>`, testToken)
}

func writeDataPage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, `<html><body><table><tbody>
		<tr><td class="pdg">Item1</td><td class="pdg">100</td></tr>
	</tbody></table></body></html>`)
}

func TestSession_LoginOnRedirect(t *testing.T) {
	ts, state := newAuthServer(t, testUser, testPassword)
	session := NewSession(config.Link{Name: testUser, Password: testPassword})

	data := session.Scrape(ts.URL+"/data", false)

	if len(data) != 1 {
		t.Fatalf("got %d items, want 1", len(data))
	}
	if data[0].Name != "Item1" || data[0].Value != "100" {
		t.Errorf("data[0] = {%q, %q}, want {%q, %q}", data[0].Name, data[0].Value, "Item1", "100")
	}
	if state.loginPosts != 1 {
		t.Errorf("login posts = %d, want 1", state.loginPosts)
	}
}

func TestSession_NoLoginWhenNotRequired(t *testing.T) {
	ts, state := newAuthServer(t, testUser, testPassword)
	session := NewSession(config.Link{Name: testUser, Password: testPassword})

	data := session.Scrape(ts.URL+"/public", false)

	if len(data) != 1 {
		t.Fatalf("got %d items, want 1", len(data))
	}
	if state.loginPosts != 0 {
		t.Errorf("login posts = %d, want 0", state.loginPosts)
	}
}

func TestSession_FallbackWithoutCredentials(t *testing.T) {
	ts, state := newAuthServer(t, testUser, testPassword)
	session := NewSession(config.Link{})

	if session != nil {
		t.Fatal("NewSession() should return nil when credentials are missing")
	}

	data := session.Scrape(ts.URL+"/data", false)

	if len(data) != 0 {
		t.Errorf("got %d items, want 0", len(data))
	}
	if state.loginPosts != 0 {
		t.Errorf("login posts = %d, want 0", state.loginPosts)
	}
}

func TestNewSession_RequiresBothCredentials(t *testing.T) {
	cases := []config.Link{
		{},
		{Name: "user"},
		{Password: "pass"},
	}
	for _, creds := range cases {
		if NewSession(creds) != nil {
			t.Errorf("NewSession(%+v) should return nil", creds)
		}
	}
	if NewSession(config.Link{Name: "user", Password: "pass"}) == nil {
		t.Error("NewSession() returned nil for complete credentials")
	}
}

func TestSession_ReloginAfterExpiry(t *testing.T) {
	ts, state := newAuthServer(t, testUser, testPassword)
	session := NewSession(config.Link{Name: testUser, Password: testPassword})

	if data := session.Scrape(ts.URL+"/data", false); len(data) != 1 {
		t.Fatalf("first scrape got %d items, want 1", len(data))
	}

	state.valid = false
	if data := session.Scrape(ts.URL+"/data", false); len(data) != 1 {
		t.Fatalf("second scrape got %d items, want 1", len(data))
	}
	if state.loginPosts != 2 {
		t.Errorf("login posts = %d, want 2", state.loginPosts)
	}
}

func TestSession_RejectedLogin(t *testing.T) {
	ts, state := newAuthServer(t, testUser, testPassword)
	session := NewSession(config.Link{Name: testUser, Password: "wrong"})

	data := session.Scrape(ts.URL+"/data", false)

	if len(data) != 0 {
		t.Errorf("got %d items, want 0", len(data))
	}
	if state.loginPosts != 1 {
		t.Errorf("login posts = %d, want 1", state.loginPosts)
	}
}

func TestNewSessions_MapsCredentialsByIndex(t *testing.T) {
	links := []config.Link{
		{URL: "https://a.example", Name: "first", Password: "one"},
		{URL: "https://b.example", Name: "second", Password: "two"},
		{URL: "https://c.example"},
	}

	sessions := NewSessions(links)

	if len(sessions) != len(links) {
		t.Fatalf("got %d sessions, want %d", len(sessions), len(links))
	}
	if sessions[0] == nil || sessions[0].creds.Name != "first" {
		t.Errorf("session 0 user = %v, want first", sessions[0])
	}
	if sessions[1] == nil || sessions[1].creds.Name != "second" {
		t.Errorf("session 1 user = %v, want second", sessions[1])
	}
	if sessions[2] != nil {
		t.Errorf("session 2 = %v, want nil", sessions[2])
	}
}
