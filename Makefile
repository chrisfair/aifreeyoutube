# ==============================================================================
# AI-Free YouTube: Authentic Human Content Portal
# ==============================================================================
# Makefile for building, testing, developing, and running the application.
# ==============================================================================

.PHONY: all build build-backend build-frontend run dev test clean help

# Default target
all: build

help:
	@echo "AI-Free YouTube - Available Commands:"
	@echo "  make build          Build both Vue 3 frontend and Go backend"
	@echo "  make build-backend  Compile Go backend binary (backend/server)"
	@echo "  make build-frontend Build Vue 3 production bundle (frontend/dist)"
	@echo "  make run            Start the standalone production server on http://localhost:8080"
	@echo "  make test           Execute comprehensive Go test suite"
	@echo "  make clean          Clean compiled binaries and build artifacts"

# Build Vue 3 frontend and Go backend
build: build-frontend build-backend
	@echo "==> Build complete! Run 'make run' to launch the server on http://localhost:8080"

# Build Vue 3 static distribution files
build-frontend:
	@echo "==> Building Vue 3 frontend..."
	cd frontend && npm run build

# Compile Go backend server binary
build-backend:
	@echo "==> Compiling Go backend binary..."
	cd backend && go build -o server ./cmd/server

# Run the unified Go server (serves API and embedded frontend SPA)
run: build
	@echo "==> Launching AI-Free YouTube server on http://localhost:8080..."
	./backend/server

# Execute unit tests with verbose logging
test:
	@echo "==> Running backend test suite..."
	cd backend && go test -v ./...

# Remove build artifacts
clean:
	@echo "==> Cleaning build artifacts..."
	rm -f backend/server
	rm -rf frontend/dist
	@echo "==> Clean complete."
