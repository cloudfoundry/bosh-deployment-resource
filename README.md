- client: Username or UAA client
- client-secret: Password or UAA client secret

# BOSH Deployment Resource

A resource that will deploy releases and stemcells using the [BOSH CLI v2](https://bosh.io/docs/cli-v2.html). 

## Differences from original BOSH Deployment Resource

The original [BOSH Deployment Resource](https://github.com/concourse/bosh-deployment-resource)
uses the Ruby CLI and does not support newer BOSH features.

### Breaking Changes

* This resource requires that the target director's SSL certificate is trusted. If the director's certificate is not
 already trusted by normal root authorities, a custom CA certificate must be provided.

## Source Configuration

* `deployment`: *Required.* The name of the deployment.
* `target`: *Optional.* The address of the BOSH director which will be used for the deployment. If omitted, target_file
  must be specified via out parameters, as documented below.
* `client`: *Required.* The UAA client ID for the BOSH director.
* `client_secret`: *Required.* The UAA client secret for the BOSH director.
* `ca_cert`: *Required.* CA certificate used to validate SSL connections to Director and UAA.

### Example

``` yaml
- name: staging
  type: bosh-deployment
  source:
    deployment: staging-deployment-name
    target: https://bosh.example.com:25555
    client: admin
    client_secret: admin
    ca_cert: "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"
```

### Dynamic Source Configuration

Sometimes source configuration cannot be known ahead of time, such as when a BOSH director is created as part of your
pipeline. In these scenarios, it is helpful to be able to have a dynamic source configuration. In addition to the
normal parameters for `get` and `put`, the following parameters can be provided to redefine the source:

* `target_file`: *Optional.* Path to a file containing a BOSH director address. This allows the target to be determined
  at runtime, e.g. by acquiring a BOSH lite instance using the
  [Pool resource](https://github.com/concourse/pool-resource).

#### Example

```
- put: staging
  params:
    target_file: path/to/target_file
```

## Behaviour

### `in`: Deploy a BOSH deployment

This will download the deployment manifest. It will place two files in the target directory:

- `manifest.yml`: The deployment manifest
- `target`: The hostname for the director

_Note_: Only the most recent version is fetchable

### `out`: Deploy a BOSH deployment

This will deploy the deployment provided.

#### Parameters

* `manifest`: *Required.* Path to a BOSH deployment manifest file.

``` yaml
- put: staging
  params:
    manifest: path/to/manifest.yml
```