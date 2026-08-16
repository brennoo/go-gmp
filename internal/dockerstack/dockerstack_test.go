package dockerstack

import "testing"

func TestParsePSEntries(t *testing.T) {
	t.Run("single running container", func(t *testing.T) {
		out := []byte(`{"Service":"gvmd","ID":"72ddd84224c7","State":"running"}` + "\n")
		entries, err := parsePSEntries(out)
		if err != nil {
			t.Fatalf("parsePSEntries: %v", err)
		}
		if len(entries) != 1 {
			t.Fatalf("got %d entries, want 1", len(entries))
		}
		if entries[0].ID != "72ddd84224c7" || entries[0].Service != "gvmd" || entries[0].State != "running" {
			t.Fatalf("unexpected entry: %+v", entries[0])
		}
	})

	t.Run("empty output means no containers", func(t *testing.T) {
		entries, err := parsePSEntries([]byte("\n"))
		if err != nil {
			t.Fatalf("parsePSEntries: %v", err)
		}
		if len(entries) != 0 {
			t.Fatalf("got %d entries, want 0", len(entries))
		}
	})

	t.Run("multiple lines are JSON lines, not a JSON array", func(t *testing.T) {
		out := []byte(
			`{"Service":"gvmd","ID":"aaa","State":"running"}` + "\n" +
				`{"Service":"gvmd","ID":"bbb","State":"running"}` + "\n",
		)
		entries, err := parsePSEntries(out)
		if err != nil {
			t.Fatalf("parsePSEntries: %v", err)
		}
		if len(entries) != 2 {
			t.Fatalf("got %d entries, want 2", len(entries))
		}
	})

	t.Run("malformed line is an error", func(t *testing.T) {
		if _, err := parsePSEntries([]byte("not json")); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestParseRepoDigests(t *testing.T) {
	t.Run("one digest", func(t *testing.T) {
		digests, err := parseRepoDigests([]byte(
			`["registry.community.greenbone.net/community/gvmd@sha256:cb82c501ca93809405a8d0ed4c21846703466802ce6462536e985625f102b497"]`,
		))
		if err != nil {
			t.Fatalf("parseRepoDigests: %v", err)
		}
		if len(digests) != 1 || digests[0] == "" {
			t.Fatalf("unexpected digests: %+v", digests)
		}
	})

	t.Run("empty array means no resolvable digest", func(t *testing.T) {
		digests, err := parseRepoDigests([]byte("[]"))
		if err != nil {
			t.Fatalf("parseRepoDigests: %v", err)
		}
		if len(digests) != 0 {
			t.Fatalf("got %d digests, want 0", len(digests))
		}
	})

	t.Run("malformed JSON is an error", func(t *testing.T) {
		if _, err := parseRepoDigests([]byte("not json")); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestSelectDigest(t *testing.T) {
	t.Run("one match", func(t *testing.T) {
		digests := []string{"gvmd@sha256:aaa", "other@sha256:bbb"}
		got, err := selectDigest(digests, "gvmd")
		if err != nil {
			t.Fatalf("selectDigest: %v", err)
		}
		if got != "gvmd@sha256:aaa" {
			t.Fatalf("got %q, want gvmd@sha256:aaa", got)
		}
	})

	t.Run("no match is an error", func(t *testing.T) {
		if _, err := selectDigest([]string{"other@sha256:bbb"}, "gvmd"); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})

	t.Run("two matches for the same repository is ambiguous", func(t *testing.T) {
		digests := []string{"gvmd@sha256:aaa", "gvmd@sha256:bbb"}
		if _, err := selectDigest(digests, "gvmd"); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})

	t.Run("result is independent of RepoDigests order", func(t *testing.T) {
		forward := []string{"other@sha256:bbb", "gvmd@sha256:aaa"}
		reversed := []string{"gvmd@sha256:aaa", "other@sha256:bbb"}
		got1, err := selectDigest(forward, "gvmd")
		if err != nil {
			t.Fatalf("selectDigest: %v", err)
		}
		got2, err := selectDigest(reversed, "gvmd")
		if err != nil {
			t.Fatalf("selectDigest: %v", err)
		}
		if got1 != got2 {
			t.Fatalf("got %q and %q, want the same result regardless of order", got1, got2)
		}
	})
}
