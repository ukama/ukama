#!/usr/bin/env python3
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.
#
import argparse
import http.server
import os
import socketserver
from urllib.parse import urlparse

class Handler(http.server.BaseHTTPRequestHandler):
    repo = None

    def do_GET(self):
        parts = [p for p in urlparse(self.path).path.split('/') if p]
        if len(parts) == 5 and parts[0] == 'v1' and parts[1] == 'apps' and parts[4] == 'pkg':
            app = parts[2]
            tag = parts[3]
            pkg = os.path.join(self.repo, f"{app}-{tag}.tar.gz")
            if os.path.exists(pkg):
                self.send_response(200)
                self.send_header('Content-Type', 'application/gzip')
                self.send_header('Content-Length', str(os.path.getsize(pkg)))
                self.end_headers()
                with open(pkg, 'rb') as f:
                    self.wfile.write(f.read())
                return
        self.send_response(404)
        self.end_headers()

    def log_message(self, fmt, *args):
        print("mock_wimc:", fmt % args, flush=True)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--host', default='127.0.0.1')
    ap.add_argument('--port', type=int, required=True)
    ap.add_argument('--repo', required=True)
    args = ap.parse_args()
    Handler.repo = args.repo
    class Server(socketserver.ThreadingMixIn, http.server.HTTPServer):
        daemon_threads = True
    with Server((args.host, args.port), Handler) as httpd:
        print(f"mock_wimc listening on {args.host}:{args.port} repo={args.repo}", flush=True)
        httpd.serve_forever()

if __name__ == '__main__':
    main()
