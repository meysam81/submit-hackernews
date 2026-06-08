package hackernews

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meysam81/submit-hackernews/internal/logger"
)

func testClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	log := logger.New(false)
	client, err := New(&log, WithBaseURL(baseURL))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return client
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantErr  error
		wantUser string
	}{
		{name: "success", body: "<html>welcome</html>", wantUser: "alice"},
		{name: "bad login", body: "<html>Bad login.</html>", wantErr: ErrInvalidCredentials},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotForm string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/login" {
					t.Errorf("unexpected path %q", r.URL.Path)
				}
				b, _ := io.ReadAll(r.Body)
				gotForm = string(b)
				http.SetCookie(w, &http.Cookie{Name: "user", Value: "session"})
				_, _ = io.WriteString(w, tt.body)
			}))
			defer srv.Close()

			client := testClient(t, srv.URL)
			err := client.Login(context.Background(), "alice", "hunter2")

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if want := "acct=alice"; !strings.Contains(gotForm, want) {
				t.Errorf("login form %q missing %q", gotForm, want)
			}
		})
	}
}

func TestSubmit(t *testing.T) {
	var (
		sawCookie bool
		postForm  string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login":
			http.SetCookie(w, &http.Cookie{Name: "user", Value: "session"})
			_, _ = io.WriteString(w, "ok")
		case "/submitlink":
			if _, err := r.Cookie("user"); err == nil {
				sawCookie = true
			}
			_, _ = io.WriteString(w, `<form>
				<input type="hidden" name="fnid" value="tok-42">
				<input type="hidden" name="fnop" value="submit-page">
			</form>`)
		case "/r":
			b, _ := io.ReadAll(r.Body)
			postForm = string(b)
			w.WriteHeader(http.StatusFound)
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
		}
	}))
	defer srv.Close()

	client := testClient(t, srv.URL)
	if err := client.Login(context.Background(), "alice", "pw"); err != nil {
		t.Fatalf("login: %v", err)
	}
	if err := client.Submit(context.Background(), "My Title", "https://example.com"); err != nil {
		t.Fatalf("submit: %v", err)
	}

	if !sawCookie {
		t.Error("session cookie was not carried to the submit form request")
	}
	for _, want := range []string{"fnid=tok-42", "fnop=submit-page", "title=My+Title", "url=https%3A%2F%2Fexample.com"} {
		if !strings.Contains(postForm, want) {
			t.Errorf("submit POST %q missing %q", postForm, want)
		}
	}
}

func TestSubmitMissingFormFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "<form></form>")
	}))
	defer srv.Close()

	client := testClient(t, srv.URL)
	err := client.Submit(context.Background(), "t", "u")
	if !errors.Is(err, ErrFormFieldMissing) {
		t.Fatalf("err = %v, want ErrFormFieldMissing", err)
	}
}
