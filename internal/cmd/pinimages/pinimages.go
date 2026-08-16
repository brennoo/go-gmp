// Command pinimages parses and rewrites the "image:" lines of a Docker
// Compose file, pinning each to a resolved content digest (repo@sha256:...)
// and, on a later run, re-resolving already-pinned lines so a pin can move
// with its tag instead of being permanently frozen. See hack/pin-images.sh,
// the thin wrapper around this command.
package main

import (
	"fmt"
	"strings"

	"github.com/brennoo/go-gmp/internal/imageref"
)

// commentUntagged is the comment written for a repository with no tag at
// all — matching hack/pin-images.sh's historical wording, so existing pins
// stay byte-identical across the shell-to-Go rewrite.
const commentUntagged = "untagged, resolved from the implicit :latest"

const imagePrefix = "    image: "

// parseImageLine extracts an "    image: ..." line's repository and, for an
// unpinned line, its tag (empty if none), or for an already-pinned line, its
// trailing recovery comment (empty if the line has none).
//
// A tag or repository is taken as everything after the reference's first
// ':', matching hack/docker-compose.yml's own documented assumption that
// none of its registries are host:port references — this package does not
// generalize that assumption, only carries it forward.
func parseImageLine(line string) (repo, tagOrComment string, pinned, ok bool) {
	rest, isImage := strings.CutPrefix(line, imagePrefix)
	if !isImage {
		return "", "", false, false
	}

	if repo, afterDigest, ok := strings.Cut(rest, "@sha256:"); ok {
		_, comment, _ := strings.Cut(afterDigest, " # ")
		return repo, comment, true, true
	}

	if repo := imageref.Repository(rest); repo != rest {
		return repo, strings.TrimPrefix(rest[len(repo):], ":"), false, true
	}
	return rest, "", false, true
}

// originalRef reconstructs the reference an already-pinned line was pinned
// from, given its repository and recovery comment, so it can be
// re-resolved. It is the inverse of the comment refreshLine's unpinned
// branch writes, and fails closed if comment is missing or does not match
// either shape that branch ever produces.
func originalRef(repo, comment string) (string, error) {
	if comment == commentUntagged {
		return repo, nil
	}
	if tag, ok := strings.CutPrefix(comment, "was "); ok && tag != "" {
		return repo + ":" + tag, nil
	}
	return "", fmt.Errorf("comment %q does not record what this image was pinned from; fix it by hand or re-run against the untagged/tagged source", comment)
}

// refreshLine resolves and (re)pins line, which must be one "\n"-free line
// of a Compose file. Lines that are not an "image:" entry pass through
// unchanged. An unpinned line is resolved and pinned, exactly as
// hack/pin-images.sh always did. An already-pinned line is re-resolved
// against the reference it was originally pinned from — recovered from its
// own repository and recovery comment — and rewritten only if the digest
// changed; when it hasn't, the reconstructed line is byte-identical to line,
// so a no-op run produces no diff. A pinned line whose comment cannot be
// recovered is an error rather than a silently frozen pin.
func refreshLine(line string, resolve func(ref string) (string, error)) (string, error) {
	repo, tagOrComment, pinned, ok := parseImageLine(line)
	if !ok {
		return line, nil
	}

	if !pinned {
		ref, comment := repo, commentUntagged
		if tagOrComment != "" {
			ref, comment = repo+":"+tagOrComment, "was "+tagOrComment
		}
		digest, err := resolve(ref)
		if err != nil {
			return "", fmt.Errorf("resolving %s: %w", ref, err)
		}
		return fmt.Sprintf("%s%s@sha256:%s # %s", imagePrefix, repo, digest, comment), nil
	}

	ref, err := originalRef(repo, tagOrComment)
	if err != nil {
		return "", fmt.Errorf("%s: %w", strings.TrimSpace(line), err)
	}
	digest, err := resolve(ref)
	if err != nil {
		return "", fmt.Errorf("re-resolving %s: %w", ref, err)
	}
	return fmt.Sprintf("%s%s@sha256:%s # %s", imagePrefix, repo, digest, tagOrComment), nil
}
