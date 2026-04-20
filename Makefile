.PHONY: build run run-race clean test verify-all verify-round1 verify-round2 verify-round3 verify-round4 verify-setup docker

# Build the workshop binary
build:
	go build -o workshop ./cmd

# Run the server with local-dev placeholder secrets. These are NOT the values
# the hosted CTFd target uses — exploit against your local server to practise
# the attack, then submit the flag you capture from the hosted target in CTFd.
# Override by exporting any of CTF_FLAG / SIGNING_TOKEN / JWT_SECRET / API_KEY
# in your shell before running `make run`.
run:
	CTF_FLAG="$${CTF_FLAG:-LOCAL_DEV_ONLY}" \
	SIGNING_TOKEN="$${SIGNING_TOKEN:-LOCAL_DEV_ONLY}" \
	JWT_SECRET="$${JWT_SECRET:-local-dev-jwt}" \
	API_KEY="$${API_KEY:-local-dev-api-key}" \
	go run ./cmd --config ./config.yaml

# Run with race detector
run-race:
	go run -race .

# Clean build artifacts
clean:
	rm -f workshop
	rm -f workshop.db

# Run all standard tests
test:
	go test ./...

# Verify all rounds (used after applying fixes)
verify-all: verify-round1 verify-round2 verify-round3 verify-round4
	@echo ""
	@echo "All rounds verified!"

# Verify individual rounds
verify-round1:
	@echo "=== Round 1: SAST ==="
	go test -v -run TestRound1 ./verify/

verify-round2:
	@echo "=== Round 2: SCA ==="
	go test -v -run TestRound2 ./verify/

verify-round3:
	@echo "=== Round 3: Runtime ==="
	go test -v -run TestRound3 ./verify/

verify-round4:
	@echo "=== Round 4: Container ==="
	go test -v -run TestRound4 ./verify/

# Pre-workshop environment check
verify-setup:
	bash scripts/verify-setup.sh

# Docker build and run
docker:
	docker compose up --build app
