#!/bin/sh
# Shared setup for hack/gvmd-up.sh and hack/gvmd-down.sh. Not executable on its
# own; both scripts dot-source it.
#
# Defines:
#   log()   - write a line to stderr
#   dc()    - run docker compose against hack/docker-compose.yml
#   HACK_DIR - absolute path to this directory

HACK_DIR="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="$HACK_DIR/docker-compose.yml"

log() {
	printf '%s\n' "$*" >&2
}

# Docker Compose ships in two forms and neither is universal: `docker compose`
# is the v2 CLI plugin, `docker-compose` is the standalone binary. A machine may
# have either, and hardcoding one means the scripts simply do not run for
# roughly half of users.
if docker compose version >/dev/null 2>&1; then
	COMPOSE_CMD="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
	COMPOSE_CMD="docker-compose"
else
	log "Neither 'docker compose' nor 'docker-compose' is available."
	log "Install Docker Compose: https://docs.docker.com/compose/install/"
	exit 1
fi

# dc runs compose with an explicit -f, so these scripts work from any working
# directory. The compose file's bind mounts use paths relative to itself
# (./CA), which Compose resolves against the file's own directory rather than
# the caller's cwd — so there is nothing to cd into and no way to end up
# mounting the wrong CA.
#
# COMPOSE_OVERRIDE_FILE, if set, is layered on top via a second -f — used by
# hack/schema-drift.sh to boot gvmd from a freshly-resolved candidate image
# without touching this checked-in compose file.
dc() {
	# shellcheck disable=SC2086 # COMPOSE_CMD is one or two words, deliberately unquoted.
	if [ -n "${COMPOSE_OVERRIDE_FILE:-}" ]; then
		$COMPOSE_CMD -f "$COMPOSE_FILE" -f "$COMPOSE_OVERRIDE_FILE" "$@"
	else
		$COMPOSE_CMD -f "$COMPOSE_FILE" "$@"
	fi
}
