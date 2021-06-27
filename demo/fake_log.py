#!/usr/bin/env python3

from datetime import datetime
import json
import copy
import os


from flask import Flask, request
from werkzeug.serving import run_simple

app = Flask(__name__)


LINE = json.loads(
    """
{
  "method": "GET",
  "path": "/factory/gitlab-py.git/info/refs",
  "format": "*/*",
  "controller": "Repositories::GitHttpController",
  "action": "info_refs",
  "status": 200,
  "time": "2021-06-24T17:26:44.000Z",
  "params": [
    {
      "key": "service",
      "value": "git-upload-pack"
    },
    {
      "key": "repository_path",
      "value": "factory/gitlab-py.git"
    }
  ],
  "remote_ip": "78.40.125.12",
  "user_id": 3,
  "username": "bdenard",
  "ua": "gitlab-runner 13.12.0 linux/amd64",
  "correlation_id": "01F829RKZS28Y8Q7JKGE4XSTXH",
  "meta.user": "bdenard",
  "meta.project": "factory/gitlab-py",
  "meta.root_namespace": "factory",
  "meta.caller_id": "Repositories::GitHttpController#info_refs",
  "meta.remote_ip": "78.40.125.12",
  "meta.feature_category": "source_code_management",
  "meta.client_id": "user/33",
  "redis_calls": 1,
  "redis_duration_s": 0.000425,
  "redis_read_bytes": 109,
  "redis_write_bytes": 44,
  "redis_cache_calls": 1,
  "redis_cache_duration_s": 0.000425,
  "redis_cache_read_bytes": 109,
  "redis_cache_write_bytes": 44,
  "db_count": 7,
  "db_write_count": 0,
  "db_cached_count": 1,
  "cpu_s": 0.042879,
  "db_duration_s": 0.00661,
  "view_duration_s": 0.00044,
  "duration_s": 0.03741
}
"""
)

log = open(os.getenv("LOG_PATH", "/tmp/glr.log"), 'w+')


@app.route("/")
def hello_world():
    v = copy.deepcopy(LINE)
    v["time"] = datetime.now().strftime("%Y-%m-%dT%H:%M:%S.000Z")
    v["remote_ip"] = request.access_route[0]
    v["meta.remote_ip"] = v["remote_ip"]
    v["ua"] = request.headers["user-agent"]
    json.dump(v, log, sort_keys=True, separators=(',', ':'))
    log.write("\n")
    log.flush()
    return "<p>Hello, World!</p>"


if __name__ == "__main__":
    run_simple("0.0.0.0", 5000, app)
