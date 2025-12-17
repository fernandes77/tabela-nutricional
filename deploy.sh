#!/usr/bin/env bash
set -euo pipefail

echo "Starting deployment..."

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

SERVICE_NAME="tabela-nutricional"
SYSTEMD_PATH="/etc/systemd/system/${SERVICE_NAME}.service"
DEPLOY_DIR="/opt/tabela-nutricional"

require_tool() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: $1 is required but not found in PATH" >&2
    exit 1
  fi
}

require_tool pnpm
require_tool go
require_tool curl

if [[ "$ROOT" != "$DEPLOY_DIR" ]]; then
  echo "Warning: repository root ($ROOT) does not match expected deploy dir ($DEPLOY_DIR)."
  echo "The systemd service points to $DEPLOY_DIR. Make sure files are deployed there."
fi

echo "Building Go backend..."
cd "$ROOT"
export GOCACHE=/tmp/go-build-cache
export GOMODCACHE=/tmp/go-mod-cache
export GOTOOLCHAIN=auto
CGO_ENABLED=0 go build -o "$ROOT/tabela-nutricional" .
chmod +x "$ROOT/tabela-nutricional"

echo "Checking for database file..."
if [[ ! -f "$ROOT/tabelas-nutricionais.db" ]]; then
  echo "Error: Database file not found at $ROOT/tabelas-nutricionais.db" >&2
  exit 1
fi

echo "Starting API server for build..."
cd "$ROOT"
export PORT=8085
"$ROOT/tabela-nutricional" &
API_PID=$!

# Wait for server to be ready
echo "Waiting for API server to start..."
for i in {1..30}; do
  if curl -s http://localhost:8085/api/foods >/dev/null 2>&1; then
    echo "API server is ready"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "Error: API server failed to start" >&2
    kill $API_PID 2>/dev/null || true
    exit 1
  fi
  sleep 1
done

echo "Installing frontend dependencies..."
cd "$ROOT/web"
export CI=true
if [[ -f pnpm-lock.yaml ]]; then
  pnpm install --frozen-lockfile
else
  pnpm install
fi

echo "Building frontend..."
export ASTRO_TELEMETRY_DISABLED=1
export API_URL=http://localhost:8085
pnpm build || {
  echo "Frontend build failed, stopping API server..."
  kill $API_PID 2>/dev/null || true
  exit 1
}
cd "$ROOT"

echo "Stopping API server..."
kill $API_PID 2>/dev/null || true
wait $API_PID 2>/dev/null || true

echo "Copying systemd service file..."
sudo cp "$ROOT/app.service" "$SYSTEMD_PATH"

echo "Reloading systemd..."
sudo systemctl daemon-reload

echo "Restarting service ${SERVICE_NAME}..."
sudo systemctl restart "${SERVICE_NAME}"

echo "Checking service status..."
sudo systemctl status "${SERVICE_NAME}" --no-pager

echo "Deployment completed successfully!"

