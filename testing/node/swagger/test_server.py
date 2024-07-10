from flask import Flask, jsonify, request, abort
import json
import yaml
import argparse
from openapi_spec_validator import validate_spec

app = Flask(__name__)

def load_spec(spec_path):
    with open(spec_path, 'r') as file:
        if spec_path.endswith('.json'):
            spec_dict = json.load(file)
        elif spec_path.endswith('.yaml') or spec_path.endswith('.yml'):
            spec_dict = yaml.safe_load(file)
        else:
            raise ValueError("Unsupported file format. Use JSON or YAML.")
    validate_spec(spec_dict)
    return spec_dict

def create_mock_endpoints(spec):
    base_path = spec.get('basePath', '')
    for path, methods in spec['paths'].items():
        for method, details in methods.items():
            endpoint = base_path + path
            if '{' in endpoint:
                endpoint = endpoint.replace('{', '<').replace('}', '>')
            create_mock_endpoint(endpoint, method, details)

def create_mock_endpoint(endpoint, method, details):
    responses = details.get('responses', {})

    def handler():
        response_code = 200
        response_description = "Mock response"
        if '200' in responses:
            response_code = 200
            response_description = responses['200'].get('description', 'OK')
        elif '404' in responses:
            response_code = 404
            response_description = responses['404'].get('description', 'Not Found')

        return jsonify({'message': response_description}), response_code

    app.add_url_rule(endpoint, endpoint + '_' + method,
                     handler, methods=[method.upper()])

@app.errorhandler(405)
def method_not_allowed(error):
    return jsonify({'message': 'Method Not Allowed'}), 405

@app.errorhandler(404)
def not_found(error):
    return jsonify({'message': 'Not Found'}), 404

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description='Run a mock server based on a Swagger spec file.')
    parser.add_argument('spec_path', type=str,
                        help='Path to the Swagger spec file')
    args = parser.parse_args()

    spec = load_spec(args.spec_path)

    # Get host from the specification, default to localhost:5000 if not defined
    host = spec.get('host', 'localhost:5000')
    port = int(host.split(':')[-1])  # Extract port number from the host

    create_mock_endpoints(spec)

    app.run(host='0.0.0.0', port=port)
