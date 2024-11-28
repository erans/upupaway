# upupaway - Your Cloud Upload Solution Made Simple

Created by Eran Sandler &copy; 2018

![Thumbla](examples/img/upupaway-logo.png)

Upupaway is a powerful yet lightweight cloud upload microservice that seamlessly streams files directly to your favorite cloud storage providers. No local storage needed - just pure efficiency.

## Why Choose Upupaway?
- Lightning fast with minimal resource footprint
- Direct streaming to cloud storage - zero local storage overhead
- Flexible multi-cloud support for all major providers
- True multi-tenant architecture - perfect for managing multiple sites/services
- Secure upload namespacing with unique upload IDs
- Seamless integration with [Thumbla](https://github.com/erans/thumbla) for a complete media handling solution

## Enterprise-Ready Cloud Storage Support
- [Google Cloud Storage](https://cloud.google.com/storage/)
- [Amazon S3](https://aws.amazon.com/s3/)
- [Azure Storage](https://azure.microsoft.com/en-us/services/storage/)
- [DigitalOcean Spaces](https://www.digitalocean.com/products/object-storage/)

## Roadmap
- Enhanced security for /prepare endpoint - <b>NOTE: <i>prepare endpoint currently requires additional security implementation</i>.</b>
- ACL configuration via request URLs
- Configurable default ACLs for S3 and Azure
- Advanced metadata management per bucket
- Dynamic metadata via request URLs

## Getting Started - Configuration
Check out our [`config-example.yml`](https://github.com/erans/upupaway/blob/master/config-example.yml) for a quick start guide.

The configuration is elegantly organized into 3 key sections:
- `buckets`: Define your cloud storage configurations with custom paths and settings
- `paths`: Map your upload endpoints to specific storage buckets
- `storage`: Configure your token storage backend:
  - `inmemory`: Perfect for development and testing
  - `redis`: Enterprise-ready Redis integration
  - `memcache`: High-performance memcached support

  Supported storages are:
  - `inmemory` - used for development purposes to temporarily store valid upload IDs during development
  - `redis` - used to store upload IDs in Redis
  - `memcache` - used to store upload IDs in memcached

## Usage - Client
Before uploads can be made to any of the paths configured under the `paths` section of the configuration a call to `/prepare` must be made, usually from the server side.

The result of will be an `UploadID` which is used to perform the upload.

## Security Best Practices
- Implement authentication for the /prepare endpoint through your preferred security layer
- Leverage cloud-native security roles instead of storing access keys in config files
- Use direct key storage only for testing or special use cases

## Kubernetes Deployment
Deploy like a pro using Kubernetes configmaps:

```
apiVersion: v1
kind: ReplicationController
metadata:
  name: upupaway
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: upupaway
    spec:
      containers:
      - name: upupaway
        image: erans/upupaway:latest
        volumeMounts:
        -
          name: config-volume
          mountPath: /etc/config
        env:
          -
            name: UUA_CFG
            value: "/etc/config/upupaway.yml"
          -
            name: UUA_PORT
            value: "8000"
        ports:
        - containerPort: 8000

      volumes:
        - name: config-volume
          configMap:
            name: upupaway-config
```
