package commands

import (
	"reflect"
	"testing"
)

func TestNewCredentialKDCs(t *testing.T) {
	kdcs := NewCredentialKDCs("kdc1", "kdc2")
	if !reflect.DeepEqual(kdcs.KDC, []string{"kdc1", "kdc2"}) {
		t.Errorf("expected KDCs to be [kdc1 kdc2], got %v", kdcs.KDC)
	}
}

func TestNewCredentialKey(t *testing.T) {
	key := NewCredentialKey("phrase", "priv", "pub")
	if key.Phrase != "phrase" || key.Private != "priv" || key.Public != "pub" {
		t.Errorf("unexpected key values: %+v", key)
	}
}

func TestNewCredentialPrivacy(t *testing.T) {
	privacy := NewCredentialPrivacy("algo", "pw")
	if privacy.Algorithm != "algo" || privacy.Password != "pw" {
		t.Errorf("unexpected privacy values: %+v", privacy)
	}
}
