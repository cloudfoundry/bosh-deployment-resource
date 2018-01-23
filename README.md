# BOSH Deployment Resource

A resource that will deploy releases and stemcells using the [BOSH CLI v2](https://bosh.io/docs/cli-v2.html). 

## Differences from original BOSH Deployment Resource

The original [BOSH Deployment Resource](https://github.com/concourse/bosh-deployment-resource)
uses the Ruby CLI and does not support newer BOSH features.

### Breaking Changes

* This resource requires that the target director's SSL certificate is trusted. If the director's certificate is not
 already trusted by normal root authorities, a custom CA certificate must be provided.

## Adding to your pipeline

To use the BOSH Deployment Resource, you must declare it in your pipeline as a resource type:

```
resource_types:
- name: bosh-deployment
  type: docker-image
  source:
    repository: cloudfoundry/bosh-deployment-resource
```

## Source Configuration

* `deployment`: *Required.* The name of the deployment.
* `target`: *Optional.* The address of the BOSH director which will be used for the deployment. If omitted, target_file
  must be specified via out parameters, as documented below.
* `client`: *Required.* The username or UAA client ID for the BOSH director.
* `client_secret`: *Required.* The password or UAA client secret for the BOSH director.
* `ca_cert`: *Optional.* CA certificate used to validate SSL connections to Director and UAA. If omitted, the director's
  certificate must be already trusted.
* `jumpbox_url`: *Optional.* The URL, including port, of the jumpbox. If set, `jumpbox_ssh_key` must also be set. If omitted,
  the BOSH director will be dialed directly.
* `jumpbox_ssh_key`: *Optional.* The private key of the jumpbox. If set, `jumpbox_url` must also be set.
* `vars_store`: *Optional.* Configuration for a persisted variables store. Currently only the Google Cloud Storage (GCS)
  provider is supported. `json_key` must be the the JSON key for your service account. Example:

  ```
  provider: gcs
  config:
    bucket: my-bucket
    file_name: path/to/vars-store.yml
    json_key: "{\"type\": \"service_account\"}"
  ```

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
normal parameters for `put`, the following parameters can be provided to redefine the source:

* `source_file`: *Optional.* Path to a file containing a YAML or JSON source config. This allows the target to be determined
  at runtime, e.g. by acquiring a BOSH lite instance using the
  [Pool resource](https://github.com/concourse/pool-resource). The content of the `source_file` should have the same
  structure as the source configuration for the resource itself. The `source_file` will be merged into the exist source
  configuration.

_Notes_:
 - `target` must **ONLY** be configured via the `source_file` otherwise the implicit `get` will fail after the `put`.
 - This is only supported for a `put`.


#### Example

```
- put: staging
  params:
    source_file: path/to/sourcefile
```

Sample source file:

```
{
    "target": "dynamic-director.example.com",
    "client_secret": "generated-secret",
    "vars_store": {
        "config": {
            "bucket": "my-bucket"
        }
    }
}
```

## Behaviour

### `in`: Download information about a BOSH deployment

This will download the deployment manifest. It will place two files in the target directory:

- `manifest.yml`: The deployment manifest
- `target`: The hostname for the director

_Note_: Only the most recent version is fetchable

#### Parameters

* `compiled_releases`: *Optional.* List of compiled releases to download. Deployment can only have one stemcell.

``` yaml
- get: staging
  params:
    compiled_releases:
    - name: release-one
    - name: release-two
```

### `out`: Deploy or Delete a BOSH deployment (defaults to deploy)

This will upload any given stemcells and releases, lock them down in the
deployment manifest and then deploy.

#### Parameters

* `manifest`: *Required.* Path to a BOSH deployment manifest file.

* `stemcells`: *Optional.* An array of globs that should point to where the
  stemcells used in the deployment can be found. Stemcell entries in the
  manifest with version 'latest' will be updated to the actual provided
  stemcell versions.

* `releases`: *Optional.* An array of globs that should point to where the
  releases used in the deployment can be found.

* `vars`: *Optional.* A collection of variables to be set in the deployment manifest.

* `vars_files`: *Optional.* A collection of vars files to be interpolated into the deployment manifest.

* `ops_files`: *Optional.* A collection of ops files to be applied over the deployment manifest.

* `cleanup`: *Optional* An boolean that specifies if a bosh cleanup should be
  run after deployment. Defaults to false.

* `no_redact`: *Optional* Removes redacted from Bosh output. Defaults to false.

* `dry_run`: *Optional* Shows the deployment diff without running a deploy. Defaults to false.

* `recreate`: *Optional* Recreate all VMs in deployment. Defaults to false.

* `target_file`: *Optional.* Path to a file containing a BOSH director address.
  This allows the target to be determined at runtime, e.g. by acquiring a BOSH
  lite instance using the [Pool
  resource](https://github.com/concourse/pool-resource).

  If both `target_file` and `target` are specified, `target_file` takes
  precedence.

* `delete.enabled`: *Optional* Deletes the configured deployment instead of doing a deploy.

* `delete.force`: *Optional* Defaults to `false`. Asks bosh to ignore errors when deleting the configured deployment.


``` yaml
# Deploy
- put: staging
  params:
    manifest: path/to/manifest.yml
    stemcells:
    - path/to/stemcells-*
    releases:
    - path/to/releases-*
    vars:
      enable_ssl: true
      domains: ["example.com", "example.net"]
      smtp:
        server: example.com
        port: 25

# Delete
- put: staging
  params:
    delete:
      enabled: true
      force: true
```
