#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

check_status() {
    local expected=$1
    local actual=$2
    local test_name=$3
    
    if [ "$actual" -eq "$expected" ]; then
        echo -e "${GREEN}✓${NC} $test_name (status: $actual)"
        return 0
    else
        echo -e "${RED}✗${NC} $test_name (expected: $expected, got: $actual)"
        return 1
    fi
}

check_server() {
    local port=$1
    if ! curl -s http://localhost:$port/_facade/health > /dev/null 2>&1; then
        echo -e "${RED}✗${NC} Server is not running on port $port"
        echo "Please start the server first:"
        echo "  make run                    # for example config"
        echo "  ./build/facade -c configs/testing.yaml -v   # for testing config"
        exit 1
    fi
}

test_example_config() {
    local port=8080
    local base_url="http://localhost:$port"
    
    echo -e "${YELLOW}Testing Example Config (port $port)...${NC}"
    check_server $port
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/_facade/health)
    check_status 200 $status "Health check"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/users)
    check_status 200 $status "Users list (schema generation)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/users/123)
    check_status 200 $status "User by ID (template)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" -X POST $base_url/users \
        -H "Content-Type: application/json" \
        -d '{"name": "Test User", "email": "test@example.com"}')
    check_status 201 $status "Create user (POST)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/_facade/state)
    check_status 200 $status "Check state"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" "$base_url/products?category=electronics")
    check_status 200 $status "Products with condition (electronics)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/products)
    check_status 200 $status "Products default"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" "$base_url/search?q=test&limit=10")
    check_status 200 $status "Search with query params"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/files/test.txt)
    check_status 200 $status "Files endpoint"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/health)
    check_status 200 $status "Custom health endpoint"
}

test_testing_config() {
    local port=8081
    local base_url="http://localhost:$port"
    
    echo -e "${YELLOW}Testing Testing Config (port $port)...${NC}"
    check_server $port
    
    status=$(curl -s -o /dev/null -w "%{http_code}" -X POST $base_url/api/payments \
        -H "Content-Type: application/json" \
        -d '{"amount": "500", "card_number": "4111111111111111"}')
    check_status 200 $status "Successful payment"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" -X POST $base_url/api/payments \
        -H "Content-Type: application/json" \
        -d '{"amount": "1500", "card_number": "4111111111111111"}')
    check_status 200 $status "High amount payment"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" -X POST $base_url/api/payments \
        -H "Content-Type: application/json" \
        -d '{"amount": "100", "card_number": "4000000000000002"}')
    check_status 402 $status "Declined payment"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" -X POST $base_url/api/payments \
        -H "Content-Type: application/json" \
        -d '{"amount": "100", "card_number": "1234567890123456"}')
    check_status 400 $status "Invalid card"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/orders/1/status)
    check_status 200 $status "Order status (pending)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/orders/2/status)
    check_status 200 $status "Order status (shipped)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/orders/3/status)
    check_status 200 $status "Order status (delivered)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/orders/999/status)
    check_status 404 $status "Order status (not found)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/error/500)
    check_status 500 $status "Error simulation (500)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/error/404)
    check_status 404 $status "Error simulation (404)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/error/401)
    check_status 401 $status "Error simulation (401)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/error/200)
    check_status 200 $status "Error simulation (200)"
    
    status=$(curl -s -o /dev/null -w "%{http_code}" $base_url/api/slow)
    check_status 200 $status "Slow endpoint"
}

main() {
    echo -e "${YELLOW}Facade API Test Suite${NC}"
    echo "================================"
    
    case "${1:-both}" in
        "example")
            test_example_config
            ;;
        "testing")
            test_testing_config
            ;;
        "both"|"")
            test_example_config
            echo ""
            test_testing_config
            ;;
        *)
            echo "Usage: $0 [example|testing|both]"
            exit 1
            ;;
    esac
    
    echo ""
    echo -e "${GREEN}Tests completed!${NC}"
}

if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is required but not installed${NC}"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}Warning: jq not found. Install it for better JSON output${NC}"
fi

main "$@"
