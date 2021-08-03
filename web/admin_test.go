package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/factorysh/gitlab-log-reader/rg"
	"github.com/factorysh/gitlab-log-reader/state"
	"github.com/stretchr/testify/assert"
)

func TestAdminAPI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := &API{
		rg:      rg.NewRG(nil, state.NewState(ctx, 3*time.Hour)),
		handler: Admin,
	}
	err := a.rg.ProcessLine(fmt.Sprintf(`{"method":"GET","path":"/factory/gitlab-py.git/info/refs","format":"*/*","controller":"Repositories::GitHttpController","action":"info_refs","status":200,"time":"%s","params":[{"key":"service","value":"git-upload-pack"},{"key":"repository_path","value":"factory/gitlab-py.git"}],"remote_ip":"78.40.125.12","user_id":3,"username":"bdenard","ua":"gitlab-runner 13.12.0 linux/amd64","correlation_id":"01F829RKZS28Y8Q7JKGE4XSTXH","meta.user":"bdenard","meta.project":"factory/gitlab-py","meta.root_namespace":"factory","meta.caller_id":"Repositories::GitHttpController#info_refs","meta.remote_ip":"78.40.125.12","meta.feature_category":"source_code_management","meta.client_id":"user/33","redis_calls":1,"redis_duration_s":0.000425,"redis_read_bytes":109,"redis_write_bytes":44,"redis_cache_calls":1,"redis_cache_duration_s":0.000425,"redis_cache_read_bytes":109,"redis_cache_write_bytes":44,"db_count":7,"db_write_count":0,"db_cached_count":1,"cpu_s":0.042879,"db_duration_s":0.00661,"view_duration_s":0.00044,"duration_s":0.03741}`, time.Now().Format(rg.TimeFormat)))
	assert.NoError(t, err)
	ts := httptest.NewServer(a)
	defer ts.Close()
	c := ts.Client()
	// invalid url, 404
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	r, err := c.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, r.StatusCode)
	// valid url, 200
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", ts.URL, "allowlist"), nil)
	r, err = c.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	data, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	var state []StateResult
	err = json.Unmarshal(data, &state)
	assert.NoError(t, err)
	assert.Len(t, state, 1)
}
