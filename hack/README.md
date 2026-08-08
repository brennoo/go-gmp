# Local Greenbone stack

`docker-compose.yml` brings up a Greenbone Community Edition stack (`gvmd`,
PostgreSQL, the vulnerability feed containers, etc.) for manual and
integration testing against a real GMP server.

Prefer `hack/gvmd-up.sh` and `hack/gvmd-down.sh` at the repo root over
running the commands below by hand — they wrap this exact sequence, get the
working directory right, and additionally create the admin user and set the
feed import owner (see "Admin user and feed import owner" below). What
follows documents what those scripts do and why, for when something goes
wrong and you need to run a step manually.

The commands below are written to run from this `hack/` directory. They also
work from anywhere with an explicit `-f path/to/hack/docker-compose.yml`: the
CA bind mounts are relative to the compose file (`./CA`), which Compose
resolves against the file's own location rather than your shell's current
directory, and the project name is pinned by `name:` at the top of the file.
The scripts pass `-f` for exactly this reason.

Compose ships as either the `docker compose` plugin or the standalone
`docker-compose` binary. Substitute whichever you have for `docker compose`
below; the scripts detect it for you.

## Generate certificates (first run only)

`CA/*.pem` and `CA/private/*.pem` are gitignored; a fresh clone has none.

Do **not** use `docker-compose exec -it gvmd gvm-manage-certs -af` — the
previous version of this doc recommended it, but it races `gvmd`'s
`restart: on-failure` and requires a full 20-45 minute first boot to
complete before there is a running container to exec into at all.

Instead, generate certificates with a disposable, dependency-free container
before starting anything else:

```bash
cd hack
docker compose run --rm --no-deps --entrypoint gvm-manage-certs gvmd -af
```

`--no-deps` means only the `gvmd` image is used — Postgres and the five
one-shot feed-data loaders never start — and the CA bind mounts still apply,
so the generated PEMs land on the host under `hack/CA/`.

**Unconfirmed, verify by actually running the above:** `gvmd` (the daemon)
runs as uid `gvmd` inside the container and must be able to read its own
generated private key. If the container above runs as a different user
(root, most likely) and leaves `CA/private/*.pem` owned by that user with
mode `0600`, `gvmd` will fail to start with a permission error visible in
`docker compose logs gvmd`. If that happens, regenerate as the current host
user and loosen permissions:

```bash
rm -rf CA
docker compose run --rm --no-deps --user "$(id -u):$(id -g)" \
  --entrypoint gvm-manage-certs gvmd -af
chmod -R a+r CA
```

`hack/gvmd-up.sh` does exactly this when run with
`GVMD_FIX_CERT_OWNERSHIP=1`.

## Start gvmd

```bash
cd hack
docker compose up -d gvmd
```

`gvmd`'s own `depends_on` already names the minimum dependency set —
`pg-gvm` plus the five one-shot data loaders (`scap-data`, `cert-bund-data`,
`dfn-cert-data`, `data-objects`, `report-formats`) — so this starts nothing
scanner-side. To also start `vulnerability-tests` (needed for `get_nvts` to
return anything), add it explicitly: `docker compose up -d gvmd
vulnerability-tests`.

There are no healthchecks anywhere in `docker-compose.yml`. Readiness has to
be polled; see "Waiting for it to be ready" below.

## Admin user and feed import owner

Nothing in the compose stack creates a GMP user or configures a feed import
owner. Both are required before the server is useful, and both are handled
by `hack/gvmd-up.sh` — do this by hand only if you're not using the
script.

`gvmd`'s CLI supports these as one-shot commands run inside the container
(they talk to Postgres directly and do not need the TLS listener up):

```bash
# Probe first: --create-user fails if the user already exists (a warm
# volume), and --user/--new-password fails if it does not exist yet (a cold
# volume).
docker compose exec -T -u gvmd gvmd gvmd --get-users

# If "admin" was not in that list:
docker compose exec -T -u gvmd gvmd gvmd --create-user=admin

# Always safe, whichever branch above ran:
docker compose exec -T -u gvmd gvmd gvmd --user=admin --new-password=<password>

# Feed import owner — without this, get_configs, get_port_lists,
# get_report_formats and get_scanners all silently return empty:
docker compose exec -T -u gvmd gvmd gvmd --get-users --verbose   # note admin's UUID
docker compose exec -T -u gvmd gvmd gvmd --modify-setting 78eceaec-3385-11ea-b237-28d24461215b --value <admin UUID>
```

## Waiting for it to be ready

The feed sync (SCAP, CERT, NVT data) that starts on first boot takes 20-45
minutes and there is no healthcheck for it. `<get_feeds/>` is gvmd's own
direct answer to whether it has finished: every feed should have a
non-empty `<version>` and none should have `<currently_syncing>` set.
`internal/gmptest`'s readiness ladder (behind `-tags integration`, see
`internal/gmptest/live.go`) polls exactly this, plus `<get_version/>` and
`<authenticate/>` as earlier gates, with an hour-long budget:

```bash
go test -tags integration ./internal/gmptest/... \
  -gmp.pass=<password> -gmp.servername=<CN, see below>
```

## Connecting over TLS

Port `9390` requires TLS 1.2 or 1.3 (`GVMD_ARGS` sets
`--gnutls-priorities=SECURE256:-VERS-TLS-ALL:+VERS-TLS1.2:+VERS-TLS1.3`),
which this library's default `tls.Config` already satisfies.

The server certificate's Subject Alternative Name very likely does not
include `127.0.0.1`, since `gvm-manage-certs` sets the CN from the
generating container's hostname rather than anything related to how you'll
dial it. Verification against the CA generated above therefore needs an
explicit server name matching that CN:

```bash
openssl x509 -in CA/servercert.pem -noout -subject
```

and pass the resulting CN to `connections.WithServerName(...)` (or
`-gmp.servername` for the test suite above).

### Validate the certs

```bash
docker compose exec -it gvmd gvm-manage-certs -V
```

This one is fine to run against the long-lived container — it only reads
files that already exist, so it does not race the restart-on-failure loop
the way generating certs against a not-yet-started container would.

## Stopping

```bash
cd hack
docker compose down       # keeps volumes: Postgres data and the synced feeds
docker compose down -v    # also deletes them, throwing away the feed sync
```

or `hack/gvmd-down.sh` / `hack/gvmd-down.sh --purge` from the repo
root, which do the same `cd` and the same distinction.
