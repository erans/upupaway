# upupaway - Up, up and away! - a cloud upload micro service

Written by Eran Sandler &copy; 2017

Up, up and Away (upupaway) is a cloud upload micro service. Files are streamed directly to the storage service bucket without storing anything on the local server.

## Benefits
- Small footprint
- Streams directly to the cloud storage service without storing anything locally
- Supports multiple storage services
- Maps an upload path to a specific bucket or bucket location which enables using the same service for multiple sites/services (multi-tenancy)
- Each upload resides in its own upload namespace (upload ID)
- Combine with [Thumbla](https://github.com/erans/thumbla) for a complete file upload + image processing solution

## Supported Storage Services
- [Google Storage](https://cloud.google.com/storage/)
- [AWS S3](https://aws.amazon.com/s3/)

## TODO
- Define a way to protect calls to /preapre
- Allow sending ACL on the request URL
- Allow setting S3 default ACL in config
- Support Azure Storage
- Support Digital Ocean Object Storage

## Usage - Configuration
See `config-example.yml` for an example of the configuration.

In general, the rule is that you first configure a bucket under the `buckets` section
and then configure a path under the `paths` section that uses a specific bucket.

## Usage - Client
Before uploads can be made to any of the paths configured under the `paths` section of the configuration a call to `/prepare` must be made, usually from the server side.

The result of will be an `UploadID` which is used to perform the upload.

To perform the upload, simply do a `POST` to any of the paths configured under the `paths` section of the configuration file with a query string parameter named `uid` for example, if the `uid` is "AABBCC" the upload should be POSTed to `/u/p1?uid=AABBCC`

Suggested usage would be to make a call to `/prepare`, get the UploadID, generate the upload URL and put it as the action of a FORM element that has an input type=file in it.
