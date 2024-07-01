import argparse
import os
import tempfile

import handler_util
from flask import Flask, jsonify, request

app = Flask(__name__)

exec_path = {
    "signing": "/usr/local/bin/signing",
    "auth_passcode_2PC": "/app/mpcauth/build/bin/auth_passcode_2PC",
    "auth_passcode_3PC": "/app/mpcauth/build/bin/auth_passcode_3PC",
    "aes_ctr": "/app/mpcauth/build/bin/aes_ctr",
    "pir": "/app/pir/bazel-bin/server_handle_pir_requests_bin",
}


@app.route("/", methods=["POST"])
def handler():
    event = request.json

    storage_name = os.environ.get("STORAGE")
    bucket_name = "flock-storage"

    response, status_code = handler_util.handler_body(
        event, bucket_name, storage_name, exec_path
    )
    return jsonify(response), status_code


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-p", "--port", type=int, default=443, help="port number to run the server on"
    )
    parser.add_argument(
        "-s",
        "--storage_name",
        default="azure",
        help="which storage type you are using (e.g. aws, gcp, azure, local)",
    )
    args = parser.parse_args()
    port = args.port

    if os.environ.get("STORAGE") is None:
        os.environ["STORAGE"] = args.storage_name

    if port == 443:
        cert = os.environ.get("PARTY_CERT")
        key = os.environ.get("PARTY_KEY")

        cert_file_path = '/app/certs/client.pem'
        key_file_path = '/app/certs/client.key'
        with open(cert_file_path, 'w') as cert_file:
            cert_file.write(cert)
        with open(key_file_path, 'w') as key_file:
            key_file.write(key)

        app.run(
            debug=True,
            host="0.0.0.0",
            port=443,
            ssl_context=(cert_file_path, key_file_path),
        )
    else:
        app.run(debug=True, host="0.0.0.0", port=port)
