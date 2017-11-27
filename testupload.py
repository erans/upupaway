import os
import sys
import requests
import mimetypes
import argparse

server_url = "http://localhost:8000"

upload_file = None

def prepare_upload(server):
        r = requests.get("{0}/prepare".format(server_url))
        if r.status_code != 200:
            raise Exception("Failed to fetch an Upload ID")

        upload_id = r.json().get("result", {}).get("uploadId", None)
        if upload_id is None:
            raise Exception("Failed to get a valid Upload ID")

        return upload_id

def preform_upload(server, server_path, upload_id, upload_file):
    mime_type = mimetypes.guess_type(upload_file)
    if not mime_type:
        mime_type = "application/octet-stream"
    else:
        mime_type = mime_type[0]

    r = requests.post("{0}{1}?uid={2}".format(server, server_path, upload_id), files={ "file" : (os.path.basename(upload_file), open(upload_file, "rb"), mime_type) })
    if r.status_code != 200:
        raise Exception("Upload failed. Status={0}  Reason={1}".format(r.status_code, r.text))

    return r.json()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Test Uploads to Up, Up and Away (upupaway) upload micro service")
    parser.add_argument("--server", help="Server URL (default: http://localhost:8000)", default="http://localhost:8000")
    parser.add_argument("file", help="File to upload")
    parser.add_argument("path", help="Server Upload Path")

    args = parser.parse_args()

    upload_id = prepare_upload(args.server)
    response = preform_upload(args.server, args.path, upload_id, args.file)
    if response and response.get("status", "") == "ok":
        print "URL: {0}".format(response.get("data",{}).get("url"))

    print "Done"
