import os
import sys
import requests
import mimetypes

server_url = "http://localhost:8000"

upload_file = None

if len(sys.argv) > 2:
    server_url = sys.argv[1]
    upload_file = sys.argv[2]
else:
    upload_file = sys.argv[1]

upload_id = requests.get("{0}/prepare".format(server_url)).json()["result"]["uploadId"]
print "Got Upload ID: {0}".format(upload_id)

mime_type = mimetypes.guess_type(upload_file)
print "Uploading {0}   Content-Type: {1}   Size: {2}".format(upload_file, mime_type[0], os.path.getsize(upload_file))

r = requests.post("{0}/u/p2/?uid={1}".format(server_url, upload_id), files={ "file" : (os.path.basename(upload_file), open(upload_file, "rb"), mime_type[0]) })
print r.status_code, r.text
