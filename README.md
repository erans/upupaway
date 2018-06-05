# upupaway - Up, up and away! - a cloud upload micro service

Written by Eran Sandler &copy; 2018

![Thumbla](examples/img/upupaway-logo.png)

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
- [Azure Storage](https://azure.microsoft.com/en-us/services/storage/)
- [DigitalOcean Spaces](https://www.digitalocean.com/products/object-storage/)

## TODO
- Define a way to protect calls to /prepare - <b>NOTE: <i>prepare is currently completely unprotected and can be called by anyone</i>.</b>
- Allow sending ACL on the request URL
- Allow setting S3 default ACL in config
- Allow setting Azure default ACL in config
- Allow setting default meta-data per bucket in config
- Allow sending meta-data on the request URL

## Usage - Configuration
See [`config-example.yml`](https://github.com/erans/upupaway/blob/master/config-example.yml) for an example of the configuration.

The configuration file has 3 major sections:
- `buckets` where the different bucket configurations are defined per cloud storage service. If you use the same bucket but with different sub-paths you will need to configure different buckets - each with its own name and path.
- `paths` where the different paths that files can be uploaded (POSTed) to. Each such path contains the bucket name (defined in the `buckets` section) that files should be uploaded to.
- `storage` where different storage configurations are used to store the UploadID tokens and the amount of time a specific token is available after generating it with a GET call to `/prepare`.

  Supported storages are:
  - `inmemory` - used for development purposes to temporarily store valid upload IDs during development
  - `redis` - used to store upload IDs in Redis
  - `memcache` - used to store upload IDs in memcached

## Usage - Client
Before uploads can be made to any of the paths configured under the `paths` section of the configuration a call to `/prepare` must be made, usually from the server side.

The result of will be an `UploadID` which is used to perform the upload.

To perform the upload, simply do a `POST` to any of the paths configured under the `paths` section of the configuration file with a query string parameter named `uid` for example, if the `uid` is "AABBCC" the upload should be POSTed to `/u/p1?uid=AABBCC`

Suggested usage would be to make a call to `/prepare`, get the UploadID, generate the upload URL and put it as the action of a FORM element that has an input type=file in it.

## Security Considerations
- Access to /prepare method is not protected. The micro service assumes something else will perform the check before that. This can be a server-to-server call that will send the upload ID only to authenticated users or some kind of an API proxy that checks for proper keys.
- Don't store access keys in the configuration file (unless you really have to). Most cloud services such as AWS, GCP and Azure have ways of defining security roles for instances. It is recommended to use that instead of setting the access keys in the configuration. The configuration supports this feature for testing purposes or for cases where the micro service itself runs outside of a certain cloud environment (or all cloud environments).

## Running Under Kubernetes
- The best way to run the mico service under Kubernetes with custom configuration is to update the configuration file as a configmap:
```
kubectl create configmap upupaway-config --from-file=upupaway.yml
```

You can then mount `upupaway-config` as a volume inside your container and point to it using an environment varaible `UUA_CFG`. For example:
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

The above configuration will mount `upupaway-config` map onto `/etc/config` inside the container. The environment variable `UUA_CFG` points to the config file `upupaway.yml` under `/etc/config` (the mounted volume).
