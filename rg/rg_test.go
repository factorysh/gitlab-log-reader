package rg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/factorysh/gitlab-log-reader/state"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func TestTimeFormat(t *testing.T) {
	ts, err := time.Parse(timeFormat, "2021-06-13T20:26:56.186Z")
	assert.NoError(t, err)
	assert.Equal(t, 2021, ts.Year())
}

func TestRG(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rg := &RG{
		state:  state.NewState(ctx, 3*time.Hour),
		parser: &fastjson.Parser{},
	}
	err := rg.ProcessLine(`{"method":"GET","path":"/ndupond/eBPF_demo/tree/0db310cfd5486f1b7c63c5e5b543dc921ec34c30/demo-wordpress/wordpress/wp-includes/Text/Diff","format":"html","controller":"Projects::TreeController","action":"show","status":200,"time":"2021-06-13T20:26:56.186Z","params":[{"key":"namespace_id","value":"ndubouilh"},{"key":"project_id","value":"eBPF_demo"},{"key":"id","value":"0db310cfd5486f1b7c63c5e5b543dc921ec34c30/demo-wordpress/wordpress/wp-includes/Text/Diff"}],"remote_ip":"192.99.5.48","user_id":null,"username":null,"ua":"Mozilla/5.0 (compatible; MJ12bot/v1.4.8; http://mj12bot.com/)","correlation_id":"01F83GWKM1AE0W2BMPHJBX3S4B","meta.project":"ndupond/eBPF_demo","meta.root_namespace":"ndupond","meta.caller_id":"Projects::TreeController#show","meta.remote_ip":"192.99.5.48","meta.feature_category":"source_code_management","meta.client_id":"ip/192.99.5.48","gitaly_calls":2,"gitaly_duration_s":0.009537,"redis_calls":13,"redis_duration_s":0.003593,"redis_read_bytes":1372,"redis_write_bytes":2328,"redis_cache_calls":13,"redis_cache_duration_s":0.003593,"redis_cache_read_bytes":1372,"redis_cache_write_bytes":2328,"db_count":11,"db_write_count":0,"db_cached_count":1,"cpu_s":0.168202,"db_duration_s":0.00702,"view_duration_s":0.11873,"duration_s":0.16273}`)
	assert.NoError(t, err)
	_, ok := rg.state.Get(state.Key{"factory/gitlab-py", "192.99.5.48", ""})
	assert.False(t, ok)
	err = rg.ProcessLine(fmt.Sprintf(`{"method":"GET","path":"/factory/gitlab-py.git/info/refs","format":"*/*","controller":"Repositories::GitHttpController","action":"info_refs","status":200,"time":"%s","params":[{"key":"service","value":"git-upload-pack"},{"key":"repository_path","value":"factory/gitlab-py.git"}],"remote_ip":"78.40.125.12","user_id":3,"username":"bdenard","ua":"gitlab-runner 13.12.0 linux/amd64","correlation_id":"01F829RKZS28Y8Q7JKGE4XSTXH","meta.user":"bdenard","meta.project":"factory/gitlab-py","meta.root_namespace":"factory","meta.caller_id":"Repositories::GitHttpController#info_refs","meta.remote_ip":"78.40.125.12","meta.feature_category":"source_code_management","meta.client_id":"user/33","redis_calls":1,"redis_duration_s":0.000425,"redis_read_bytes":109,"redis_write_bytes":44,"redis_cache_calls":1,"redis_cache_duration_s":0.000425,"redis_cache_read_bytes":109,"redis_cache_write_bytes":44,"db_count":7,"db_write_count":0,"db_cached_count":1,"cpu_s":0.042879,"db_duration_s":0.00661,"view_duration_s":0.00044,"duration_s":0.03741}`, time.Now().Format(timeFormat)))
	assert.NoError(t, err)
	v, ok := rg.state.Get(
		state.Key{
			"factory/gitlab-py",
			"78.40.125.12",
			"",
		})
	fmt.Println(v)
	assert.True(t, ok)
}
