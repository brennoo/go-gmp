#!/bin/sh
# Stops the stack brought up by hack/gvmd-up.sh.
#
# By default this preserves volumes (Postgres data and the synced feeds),
# which is what saves the next 20-45 minute first boot. Pass --purge to
# also delete them.
#
# Usage:
#   hack/gvmd-down.sh
#   hack/gvmd-down.sh --purge
set -eu

. "$(dirname "$0")/compose-lib.sh"

purge=0
for arg in "$@"; do
	case "$arg" in
	--purge)
		purge=1
		;;
	*)
		log "usage: $0 [--purge]"
		exit 2
		;;
	esac
done

if [ "$purge" -eq 1 ]; then
	log "Stopping the stack and deleting all volumes."
	log "This throws away the feed sync (20-45 minutes to redo) and the Postgres data."
	dc down -v
else
	dc down
	log "Volumes preserved (Postgres data, synced feeds). Pass --purge to delete them too."
fi
