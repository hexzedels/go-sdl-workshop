//go:build !exploit && !docker

// Package verify provides static fix checks over the source code.
// Run with: go test -v ./verify/
package verify

import (
	"os"
	"strings"
	"testing"
)

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	return string(b)
}

// ---------- Round 1 ----------

func TestRound1_SQLInjection(t *testing.T) {
	src := mustRead(t, "../internal/store/store.go")

	if strings.Contains(src, `LIKE '%"`) || strings.Contains(src, `LIKE '%%%s`) {
		t.Error("Round 1.1: title LIKE is still being built with string concatenation — use a ? parameter")
	}

	if !strings.Contains(src, "LIKE ?") {
		t.Error("Round 1.1: title search must use `LIKE ?` with a bound parameter")
	}

	if !strings.Contains(src, "owner_id = ?") {
		t.Error("Round 1.1: search must still filter by owner_id using a parameter placeholder")
	}
}

func TestRound1_PathTraversal(t *testing.T) {
	src := mustRead(t, "../internal/handler/files.go")

	if strings.Contains(src, "filepath.HasPrefix") {
		t.Error("Round 1.3: filepath.HasPrefix is deprecated since Go 1.19 — use filepath.Abs + filepath.Rel (or strings.HasPrefix on absolute paths)")
	}

	if !strings.Contains(src, "filepath.Abs") {
		t.Error("Round 1.3: fix must canonicalise paths via filepath.Abs before comparing")
	}

	hasRel := strings.Contains(src, "filepath.Rel")
	hasStringsPrefix := strings.Contains(src, "strings.HasPrefix")
	if !hasRel && !hasStringsPrefix {
		t.Error("Round 1.3: fix must reject paths outside the upload dir (use filepath.Rel, or strings.HasPrefix on absolute paths)")
	}
}

func TestRound1_JWTValidation(t *testing.T) {
	src := mustRead(t, "../internal/auth/jwt.go")

	if strings.Contains(src, "ParseWithClaims") && strings.Contains(src, ", _ :=") {
		t.Error("Round 1.4: the error from jwt.ParseWithClaims must not be discarded")
	}
	if strings.Contains(src, "_ = token") {
		t.Error("Round 1.4: the parsed token must be inspected, not discarded")
	}

	if !strings.Contains(src, "err != nil") {
		t.Error("Round 1.4: fix must branch on err != nil and return an error")
	}
	if !strings.Contains(src, "token.Valid") {
		t.Error("Round 1.4: fix must check token.Valid")
	}
}

// ---------- Round 2 ----------

func TestRound2_Typosquat(t *testing.T) {
	logutils, err := os.ReadFile("../pkg/logutils/logutils.go")
	if err == nil {
		src := string(logutils)
		if strings.Contains(src, "/tmp/.env_dump") || strings.Contains(src, "os.Environ()") {
			t.Error("Round 2.1: the typosquat package still exfiltrates env vars — remove it")
		}
	}

	mainSrc := mustRead(t, "../cmd/main.go")
	if strings.Contains(mainSrc, `"github.com/hexzedels/gosdlworkshop/pkg/logutils"`) {
		t.Error("Round 2.1: cmd/main.go still imports the typosquat package")
	}
}

func TestRound2_WeakRandom(t *testing.T) {
	src := mustRead(t, "../internal/token/token.go")

	if strings.Contains(src, `"math/rand"`) {
		t.Error("Round 2.2: token.go still imports math/rand — switch to crypto/rand")
	}
	if !strings.Contains(src, "crypto/rand") {
		t.Error("Round 2.2: token.go must use crypto/rand")
	}
}

func TestRound2_VulnerableDep(t *testing.T) {
	src := mustRead(t, "../go.mod")
	if strings.Contains(src, "golang.org/x/text v0.3.7") {
		t.Error("Round 2.3: vulnerable golang.org/x/text v0.3.7 is still in go.mod")
	}
}

// ---------- Round 3 ----------

func TestRound3_GoroutineLeak(t *testing.T) {
	src := mustRead(t, "../internal/handler/webhook.go")

	if !strings.Contains(src, "context.WithTimeout") && !strings.Contains(src, "http.Client{") {
		t.Error("Round 3.1: webhook dispatch needs a per-request timeout (context.WithTimeout or http.Client{Timeout})")
	}

	hasSem := strings.Contains(src, "chan struct{}") || strings.Contains(src, "semaphore") || strings.Contains(src, "errgroup")
	if !hasSem {
		t.Error("Round 3.1: webhook dispatch needs a concurrency limit (buffered channel, semaphore, or errgroup)")
	}
}

func TestRound3_RaceCondition(t *testing.T) {
	src := mustRead(t, "../internal/auth/session.go")

	if !strings.Contains(src, "sync.RWMutex") && !strings.Contains(src, "sync.Mutex") && !strings.Contains(src, "sync.Map") {
		t.Error("Round 3.2: session store must use sync.Mutex/RWMutex or sync.Map")
	}
}

func TestRound3_TimingAttack(t *testing.T) {
	src := mustRead(t, "../internal/auth/apikey.go")

	if strings.Contains(src, "provided == configured") || strings.Contains(src, "provided != configured") {
		t.Error("Round 3.3: direct string equality leaks timing — use crypto/subtle")
	}
	if strings.Contains(src, "time.Sleep") {
		t.Error("Round 3.3: artificial per-byte sleep must be removed")
	}
	if !strings.Contains(src, "subtle.ConstantTimeCompare") {
		t.Error("Round 3.3: fix must use subtle.ConstantTimeCompare")
	}
}

// ---------- Round 4 ----------

func TestRound4_HardcodedSecrets(t *testing.T) {
	dockerfile := mustRead(t, "../Dockerfile")

	if strings.Contains(dockerfile, "COPY . .") {
		ignore, err := os.ReadFile("../.dockerignore")
		if err != nil || !strings.Contains(string(ignore), "config.yaml") {
			t.Error("Round 4.1: .dockerignore must exclude config.yaml when Dockerfile uses COPY . .")
		}
	}

	if strings.Contains(dockerfile, "ENV SIGNING_TOKEN") ||
		strings.Contains(dockerfile, "ENV JWT_SECRET") ||
		strings.Contains(dockerfile, "ENV API_KEY") {
		t.Error("Round 4.1: secrets must not be frozen into image layers via ENV directives")
	}
}

