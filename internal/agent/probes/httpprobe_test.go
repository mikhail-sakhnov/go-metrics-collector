package probes

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

func TestHTTPProbe(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("panic") != "" {
			panic("this panic would be intercepted by the http server")
		}
		w.WriteHeader(418)
		w.Write([]byte("It works!"))
		time.Sleep(time.Millisecond * 100)
	}))
	defer ts.Close()
	t.Run("http_probe_response_status", func(t *testing.T) {
		msg, err := HTTPProbe(context.Background(), "unittest", ts.URL, nil)
		require.NoError(t, err)
		require.Equal(t, 418, msg.HTTPCode)
	})
	t.Run("http_probe_duration_test", func(t *testing.T) {
		msg, err := HTTPProbe(context.Background(), "unittest", ts.URL, nil)
		require.NoError(t, err)
		require.True(t, msg.ResponseTime > 100)
	})

	t.Run("http_probe_if_content_found", func(t *testing.T) {
		msg, err := HTTPProbe(context.Background(), "unittest", ts.URL, regexp.MustCompile(".*works.*"))
		require.NoError(t, err)
		require.True(t, msg.ContentFound)
	})

	t.Run("http_probe_if_content_not_found", func(t *testing.T) {
		msg, err := HTTPProbe(context.Background(), "unittest", ts.URL, regexp.MustCompile(".*hello.*"))
		require.NoError(t, err)
		require.False(t, msg.ContentFound)
	})

	t.Run("http_probe_error_case", func(t *testing.T) {
		msg, err := HTTPProbe(context.Background(), "unittest", ts.URL+"?panic=1", regexp.MustCompile(".*hello.*"))
		require.Error(t, err)
		require.True(t, errors.Is(err, io.EOF))
		require.False(t, msg.ContentFound)
	})
}
