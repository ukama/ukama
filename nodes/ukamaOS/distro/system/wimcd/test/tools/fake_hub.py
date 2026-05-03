#!/usr/bin/env python3

from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import sys

PORT = int(sys.argv[1]) if len(sys.argv) > 1 else 19080

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path.startswith('/v1/hub/apps/example') or \
           self.path.startswith('/v1/apps/example'):
            body = {
                'name': 'example',
                'artifacts': [
                    {
                        'version': 'v1-abc',
                        'formats': [
                            {
                                'type': 'chunk',
                                'url': '/fake/example/v1-abc.caidx',
                                'created_at': '2026-01-01T00:00:00Z',
                                'extra_info': {'chunks': '/fake/chunks/'}
                            }
                        ]
                    },
                    {
                        'version': 'v1-xyz',
                        'formats': [
                            {
                                'type': 'chunk',
                                'url': '/fake/example/v1-xyz.caidx',
                                'created_at': '2026-01-01T00:00:00Z',
                                'extra_info': {'chunks': '/fake/chunks/'}
                            }
                        ]
                    }
                ]
            }
            data = json.dumps(body).encode('utf-8')
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.send_header('Content-Length', str(len(data)))
            self.end_headers()
            self.wfile.write(data)
            return

        self.send_response(404)
        self.end_headers()

    def log_message(self, fmt, *args):
        return

HTTPServer(('127.0.0.1', PORT), Handler).serve_forever()
