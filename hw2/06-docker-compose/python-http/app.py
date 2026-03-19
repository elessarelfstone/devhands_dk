from flask import Flask, jsonify, request
import os
import requests

app = Flask(__name__)

# Get Go service URL from environment
GO_SERVICE = os.getenv("GO_SERVICE", "http://go-http:8080")

@app.route('/document/<int:doc_id>')
def get_document(doc_id):
    try:
        # Call Go service
        response = requests.get(f"{GO_SERVICE}/document/{doc_id}", timeout=5)
        response.raise_for_status()
        return jsonify(response.json())
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500

@app.route('/health')
def health():
    return "OK", 200

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
