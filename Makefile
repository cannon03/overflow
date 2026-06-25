.PHONY: all native wasm web clean serve netlify

BINARY   = overflow
WASM_OUT = web/main.wasm
WASM_JS  = web/wasm_exec.js

# ──────────────────────────────────────────────
# Default: build native binary
# ──────────────────────────────────────────────
all: native

native:
	go build -o $(BINARY) ./cmd/game/

# ──────────────────────────────────────────────
# WASM binary for browser deployment
# ──────────────────────────────────────────────
wasm: $(WASM_OUT) $(WASM_JS)

$(WASM_OUT):
	GOOS=js GOARCH=wasm go build -o $(WASM_OUT) ./cmd/game/

$(WASM_JS):
	@GOROOT=$$(go env GOROOT); \
	if [ -f "$$GOROOT/misc/wasm/wasm_exec.js" ]; then \
		cp "$$GOROOT/misc/wasm/wasm_exec.js" $(WASM_JS); \
	elif [ -f "$$GOROOT/lib/wasm/wasm_exec.js" ]; then \
		cp "$$GOROOT/lib/wasm/wasm_exec.js" $(WASM_JS); \
	else \
		echo "Error: wasm_exec.js not found in GOROOT ($$GOROOT)"; \
		exit 1; \
	fi

# ──────────────────────────────────────────────
# Build everything for web deployment
# ──────────────────────────────────────────────
web: wasm

# ──────────────────────────────────────────────
# Serve the web build locally (requires python3)
# ──────────────────────────────────────────────
serve: web
	@echo "  → Starting server at http://localhost:8080"
	@cd web && python3 -m http.server 8080

# ──────────────────────────────────────────────
# Netlify: build wasm + copy web assets
# ──────────────────────────────────────────────
netlify: web

# ──────────────────────────────────────────────
# Clean build artifacts
# ──────────────────────────────────────────────
clean:
	rm -f $(BINARY) $(WASM_OUT)
