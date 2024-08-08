from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route('/node/v1/notify', methods=['POST'])
def forwarded_alert():
    data = request.json
    print("Received forwarded alert:")
    print(data)
    return jsonify({"status": "received"}), 200

if __name__ == '__main__':
    app.run(port=18300)
