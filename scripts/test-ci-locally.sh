#!/bin/bash
# test-ci-locally.sh
# This script simulates the GitHub Actions CI environment locally
# to help validate that tests will pass in CI.

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Testing Dkron CI Setup Locally${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if a port is in use
port_in_use() {
    lsof -i ":$1" >/dev/null 2>&1 || netstat -an | grep -q ":$1.*LISTEN" 2>/dev/null
}

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

if ! command_exists docker; then
    echo -e "${RED}Error: Docker is not installed${NC}"
    exit 1
fi

if ! command_exists go; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker is installed${NC}"
echo -e "${GREEN}✓ Go is installed${NC}"

# Check Go version
GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓ Go version: ${GO_VERSION}${NC}"
echo ""

# Check if Mailpit is already running
echo -e "${YELLOW}Checking Mailpit status...${NC}"
MAILPIT_CONTAINER="dkron-ci-mailpit"

if docker ps -a --format '{{.Names}}' | grep -q "^${MAILPIT_CONTAINER}$"; then
    echo -e "${YELLOW}Removing existing Mailpit container...${NC}"
    docker rm -f ${MAILPIT_CONTAINER} >/dev/null 2>&1
fi

# Check if ports are available
if port_in_use 1025; then
    echo -e "${RED}Error: Port 1025 is already in use${NC}"
    echo -e "${YELLOW}Stop the service using port 1025 or stop existing Mailpit${NC}"
    exit 1
fi

if port_in_use 8025; then
    echo -e "${RED}Error: Port 8025 is already in use${NC}"
    echo -e "${YELLOW}Stop the service using port 8025 or stop existing Mailpit${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Ports 1025 and 8025 are available${NC}"
echo ""

# Start Mailpit (simulating GitHub Actions service container)
echo -e "${YELLOW}Starting Mailpit service container...${NC}"
docker run -d \
    --rm \
    --name ${MAILPIT_CONTAINER} \
    -p 1025:1025 \
    -p 8025:8025 \
    axllent/mailpit >/dev/null

# Wait for Mailpit to be ready
echo -e "${YELLOW}Waiting for Mailpit to be ready...${NC}"
sleep 2

# Check if Mailpit is responding
if ! nc -z localhost 1025 2>/dev/null && ! timeout 1 bash -c "</dev/tcp/localhost/1025" 2>/dev/null; then
    echo -e "${RED}Error: Mailpit failed to start${NC}"
    docker logs ${MAILPIT_CONTAINER}
    docker stop ${MAILPIT_CONTAINER}
    exit 1
fi

echo -e "${GREEN}✓ Mailpit is running${NC}"
echo -e "${GREEN}  - SMTP: localhost:1025${NC}"
echo -e "${GREEN}  - Web UI: http://localhost:8025${NC}"
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}Cleaning up...${NC}"
    docker stop ${MAILPIT_CONTAINER} >/dev/null 2>&1 || true
    echo -e "${GREEN}✓ Mailpit stopped${NC}"
}

# Set trap to cleanup on exit
trap cleanup EXIT INT TERM

# Run tests (simulating GitHub Actions test step)
echo -e "${YELLOW}Running tests (matching CI configuration)...${NC}"
echo -e "${YELLOW}Command: go test -v -timeout 200s -coverprofile=coverage.txt ./...${NC}"
echo ""

if go test -v -timeout 200s -coverprofile=coverage.txt ./...; then
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}✓ All tests passed!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${GREEN}Your tests should pass in GitHub Actions.${NC}"

    if [ -f coverage.txt ]; then
        echo ""
        echo -e "${YELLOW}Coverage report generated: coverage.txt${NC}"
        # Show coverage summary if go tool cover is available
        if command_exists "go"; then
            COVERAGE=$(go tool cover -func=coverage.txt | grep total | awk '{print $3}')
            echo -e "${GREEN}Total coverage: ${COVERAGE}${NC}"
        fi
    fi

    echo ""
    echo -e "${YELLOW}View captured emails at: http://localhost:8025${NC}"
    echo -e "${YELLOW}Press Ctrl+C to stop Mailpit and exit${NC}"

    # Keep script running so user can view Mailpit UI
    read -p "Press Enter to stop Mailpit and exit..."

    exit 0
else
    echo ""
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}✗ Tests failed!${NC}"
    echo -e "${RED}========================================${NC}"
    echo ""
    echo -e "${RED}Fix the failing tests before pushing to GitHub.${NC}"
    echo ""
    echo -e "${YELLOW}View captured emails at: http://localhost:8025${NC}"
    echo -e "${YELLOW}Mailpit logs:${NC}"
    docker logs ${MAILPIT_CONTAINER}

    exit 1
fi
