package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowMethods(t *testing.T) {
	t.Run("Allowed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		allowMethods([]string{http.MethodGet}, h).ServeHTTP(w, r)

		if got, want := w.Code, http.StatusOK; got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
	})

	t.Run("NotAllowed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPatch, "/", nil)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		allowMethods([]string{http.MethodGet}, h).ServeHTTP(w, r)

		if got, want := w.Code, http.StatusMethodNotAllowed; got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
	})
}

func TestValidatePath(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		for _, tc := range []struct {
			desc string
			path string
		}{
			{
				desc: "ValidType",
				path: "/characters",
			},
			{
				desc: "ValidTypeAndID",
				path: "/characters/123",
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

				validatePath(h).ServeHTTP(w, r)

				if got, want := w.Code, http.StatusOK; got != want {
					t.Errorf("got status %d, want %d", got, want)
				}
			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		for _, tc := range []struct {
			desc      string
			path      string
			errString string
		}{
			{
				desc:      "InvalidType",
				path:      "/foo",
				errString: "unknow type\n",
			},
			{
				desc:      "InvalidID",
				path:      "/characters/foo",
				errString: "invalid id\n",
			},
			{
				desc:      "EmptyID",
				path:      "/characters/",
				errString: "missing id\n",
			},
			{
				desc:      "TooMuch",
				path:      "/characters/123/foo",
				errString: "unsupported resource\n",
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, tc.path, nil)

				h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

				validatePath(h).ServeHTTP(w, r)

				if got, want := w.Code, http.StatusBadRequest; got != want {
					t.Fatalf("got status %d, want %d", got, want)
				}

				b, err := ioutil.ReadAll(w.Body)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if got, want := string(b), tc.errString; got != want {
					t.Errorf("got error %q, want %q", got, want)
				}
			})
		}
	})
}
