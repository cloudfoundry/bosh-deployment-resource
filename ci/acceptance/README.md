# Instant-BOSH Acceptance Prototype

This fork-only prototype replaces the external BOSH Director with a disposable
[instant-bosh](https://github.com/rkoster/instant-bosh) Docker Director. It
does not require Concourse, cloud credentials, or a registry push. It does not
publish images or replace the existing release gates.

## Run

Use a dedicated Linux amd64 host with a local Docker daemon and at least 12 GiB
free Docker disk space. instant-bosh starts a privileged Director container,
mounts the Docker socket, and publishes ports 25555, 2222, 8443, and 8081 on the
host. It also reserves subnet `10.245.0.0/16`. Do not run it from an untrusted
checkout or on a shared host. The preflight refuses to continue when instant-bosh
resources already exist, so it does not destroy another environment.

```bash
bash ci/acceptance/run.sh
```

The test builds the instant-bosh CLI from pinned commit `e1f2a43` (the v0.3
release predates the `docker` subcommands) and pins bosh CLI v7.11.0, the Director
image digest, and the Dummy BOSH Release v2 checksum. The Director image does not
expose a revision label, so the digest was pinned from the current `latest` and
`main` tags and is not proven to match that commit. A Docker Desktop host on
macOS is not a supported target.

The `Docker CPI acceptance` workflow runs on pushes to
`ci/docker-cpi-acceptance` in `benjaminguttmann-avtq/bosh-deployment-resource`.
Commit and push the workflow and scripts to start the first run. Actions must be
enabled in the fork. `workflow_dispatch` is also supported after the workflow is
registered on the default branch.

```bash
gh run list --repo benjaminguttmann-avtq/bosh-deployment-resource \
  --workflow acceptance.yml --branch ci/docker-cpi-acceptance
gh run watch RUN_ID --repo benjaminguttmann-avtq/bosh-deployment-resource
```

## Assertions

- Build the candidate check, in, and out binaries from this checkout.
- Start instant-bosh and use its generated CA, client secret, and jumpbox key.
- Invoke the candidate image's out entrypoint to upload a release and deploy one
  instance using the light stemcell uploaded by instant-bosh, variables, and an
  ops file.
- Assert one running instance, a stable check result, in's version, downloaded
  manifest and target, and the interpolated marker value.
- Invoke out with deletion enabled and verify the deployment is absent.
- Destroy only the instant-bosh environment started by this run. Downloaded
  Docker images and build caches are retained.

The legacy Concourse pipeline is intentionally retained. This prototype does not
cover Concourse scheduling/caching, GCP-specific behavior, or the production
jumpbox. It uses a local instant-bosh jumpbox.

## Validation Status

The candidate image builds successfully. The earlier create-env prototype did not
pass on Docker Desktop and has been replaced. The instant-bosh version has not
yet been executed on a supported Linux host: the Director startup, resource
tunnel through the jumpbox, light stemcell upload, Dummy release deployment, and
cleanup are unverified. Run it on the fork before relying on it as a release gate.
Runtime credentials are temporary and are not uploaded as CI artifacts.
