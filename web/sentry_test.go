package web

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/factorysh/gitlab-log-reader/metrics"
	"github.com/factorysh/gitlab-log-reader/rg"
	"github.com/factorysh/gitlab-log-reader/state"
	"github.com/factorysh/postier/pkg/postester"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/stretchr/testify/assert"
)

func TestSentryCapture(t *testing.T) {
	// https://github.com/factorysh/postier
	// start postier testing server
	pt, err := postester.StartTesting()
	assert.NoError(t, err, "Error when starting postier test server")
	defer pt.Cleanup()

	// init sentry client
	err = sentry.Init(sentry.ClientOptions{
		Dsn: fmt.Sprintf("http://3a29ce59d51c4d55a139ccb0c366aaaa@%s/0", pt.URL),
	})
	assert.NoError(t, err, "Error on Sentry init")

	// create admin API server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := &API{
		rg:      rg.NewRG(nil, state.NewState(ctx, 3*time.Hour, metrics.Collector)),
		handler: Admin,
		sentryHandler: sentryhttp.New(sentryhttp.Options{
			Repanic: false,
		}),
	}
	assert.NoError(t, err)
	// start admin server using httptest
	ts := httptest.NewServer(a)
	defer ts.Close()
	c := ts.Client()
	// resquest a specific hidden endpoint to trigger a panic
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/sentry/hidden/test", ts.URL), nil)
	_, err = c.Do(req)
	assert.NoError(t, err)

	// wait for postier to save things
	time.Sleep(50 * time.Millisecond)
	assert.Len(t, pt.History().FilterURL("/api/0/store/"), 1)
}
