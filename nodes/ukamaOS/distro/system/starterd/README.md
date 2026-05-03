starter and wimc contract flow:

starter.d -> POST /v1/apps/<name>/<tag>
starter.d -> GET  /v1/apps/<name>/<tag>/status
wimc.d    -> downloads only the tar.gz
shared    -> /ukama/apps/pkgs/<name>_<tag>.tar.gz
starter.d -> install/unpack/switch/start/commit
