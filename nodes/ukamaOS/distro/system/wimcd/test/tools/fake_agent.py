#!/usr/bin/env python3

from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import os
import shutil
import sys
import urllib.request

PORT = int(sys.argv[1]) if len(sys.argv) > 1 else 19081
FIXTURE_DIR = os.environ.get('WIMC_FIXTURE_TARBALLS',
                             'test/fixtures/tarballs')
PKG_DIR = os.environ.get('WIMC_TEST_PKG_DIR', '/ukama/apps/pkgs')

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/v1/ping':
            self.send_response(200)
            self.end_headers()
            return

        if self.path == '/v1/version':
            data = b'test-agent\n'
            self.send_response(200)
            self.send_header('Content-Type', 'text/plain')
            self.send_header('Content-Length', str(len(data)))
            self.end_headers()
            self.wfile.write(data)
            return

        self.send_response(404)
        self.end_headers()

    def do_POST(self):
        length = int(self.headers.get('Content-Length', '0'))
        raw = self.rfile.read(length)
        req = json.loads(raw.decode('utf-8'))

        # wimc.d currently serializes fetch requests as a flat JSON object.
        name = req['name']
        tag = req['tag']
        cb_url = req['callback_url']
        uuid = req['uuid']

        src = os.path.join(FIXTURE_DIR, f'{name}_{tag}.tar.gz')
        dst = os.path.join(PKG_DIR, f'{name}_{tag}.tar.gz')

        os.makedirs(PKG_DIR, exist_ok=True)
        shutil.copyfile(src, dst)

        update = {
            'type_update': {
                'uuid': uuid,
                'total_kilobytes': 1,
                'transfer_kilobytes': 1,
                'transfer_state': 'done',
                'void': dst
            }
        }

        data = json.dumps(update).encode('utf-8')
        request = urllib.request.Request(
            cb_url,
            data=data,
            method='PUT',
            headers={'Content-Type': 'application/json'}
        )
        urllib.request.urlopen(request, timeout=5).read()

        self.send_response(202)
        self.end_headers()

    def log_message(self, fmt, *args):
        return

HTTPServer(('127.0.0.1', PORT), Handler).serve_forever()
