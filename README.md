# upupaway - Up, up and away! - a cloud upload micro service
Up, up and Away (upupaway) is a cloud upload micro service. Files are streamed directly to the storage service without storing anything on the local server.

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
- Allow sending ACL on the request URL
- Allow setting S3 default ACL in config
- Support Azure Storage
- Support Digital Ocean Object Storage
