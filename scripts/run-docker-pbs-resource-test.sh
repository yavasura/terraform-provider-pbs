#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
COMPOSE_FILE="${SCRIPT_DIR}/docker-compose.pbs-smoke.yml"
TEST_PATTERN="${1:-^TestUsersDockerSmoke$}"
COMPOSE_CMD=()

PBS_ENDPOINT="${PBS_ENDPOINT:-https://127.0.0.1:8007}"
PBS_USERNAME="${PBS_USERNAME:-admin@pbs}"
PBS_PASSWORD="${PBS_PASSWORD:-pbspbs}"
PBS_INSECURE="${PBS_INSECURE:-true}"
PBS_TEST_USER_REALM="${PBS_TEST_USER_REALM:-pbs}"
GOCACHE="${GOCACHE:-/tmp/gocache-docker-pbs-resource-test}"

detect_compose_cmd() {
    if command -v docker-compose >/dev/null 2>&1; then
        COMPOSE_CMD=(docker-compose)
        return
    fi

    if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
        COMPOSE_CMD=(docker compose)
        return
    fi

    echo "Neither docker-compose nor 'docker compose' is available." >&2
    exit 1
}

cleanup() {
    if [ ${#COMPOSE_CMD[@]} -eq 0 ]; then
        return
    fi

    "${COMPOSE_CMD[@]}" -f "${COMPOSE_FILE}" logs --no-color || true
    "${COMPOSE_CMD[@]}" -f "${COMPOSE_FILE}" down -v
}

wait_for_pbs() {
    local max_attempts=60
    local attempt=1

    until curl -skf \
        --data-urlencode "username=${PBS_USERNAME}" \
        --data-urlencode "password=${PBS_PASSWORD}" \
        "${PBS_ENDPOINT}/api2/json/access/ticket" >/dev/null; do
        if [ "${attempt}" -ge "${max_attempts}" ]; then
            echo "PBS did not become ready after ${max_attempts} attempts." >&2
            return 1
        fi

        echo "Waiting for PBS login endpoint... (${attempt}/${max_attempts})"
        sleep 5
        attempt=$((attempt + 1))
    done
}

cd "${PROJECT_ROOT}"

if [ ! -x "${PROJECT_ROOT}/terraform-provider-pbs" ]; then
    echo "Provider binary not found at ${PROJECT_ROOT}/terraform-provider-pbs. Run 'go build -o terraform-provider-pbs .' first." >&2
    exit 1
fi

detect_compose_cmd
trap cleanup EXIT INT TERM

"${COMPOSE_CMD[@]}" -f "${COMPOSE_FILE}" up -d
wait_for_pbs

export PBS_ENDPOINT
export PBS_USERNAME
export PBS_PASSWORD
export PBS_INSECURE
export PBS_TEST_USER_REALM
export GOCACHE

go test ./test/integration -run "${TEST_PATTERN}" -count=1 -v -timeout 30m
