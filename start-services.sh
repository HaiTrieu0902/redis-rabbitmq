#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Starting Redis MQ Microservices${NC}"
echo -e "${GREEN}========================================${NC}"

# Check if Redis is running
echo -e "\n${YELLOW}Checking Redis...${NC}"
if docker ps | grep -q redis; then
    echo -e "${GREEN}✓ Redis is running${NC}"
else
    echo -e "${RED}✗ Redis is not running${NC}"
    echo -e "${YELLOW}Starting Redis...${NC}"
    cd "$(dirname "$0")"
    docker-compose up -d
    sleep 3
    echo -e "${GREEN}✓ Redis started${NC}"
fi

# Check PostgreSQL connection
echo -e "\n${YELLOW}Checking PostgreSQL...${NC}"
if psql -U postgres -d vuihoi -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PostgreSQL is accessible${NC}"
else
    echo -e "${RED}✗ Cannot connect to PostgreSQL${NC}"
    echo -e "${RED}Please make sure PostgreSQL is running${NC}"
    exit 1
fi

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}  All systems ready!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "${YELLOW}Starting services...${NC}"
echo -e ""
echo -e "${GREEN}Authentication Service:${NC} http://localhost:3001"
echo -e "${GREEN}User Service:${NC} http://localhost:3002"
echo -e ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo -e ""

# Start both services in background
cd "$(dirname "$0")/authenication"
npm run start:dev > ../logs/auth.log 2>&1 &
AUTH_PID=$!

cd "$(dirname "$0")/user"
npm run start:dev > ../logs/user.log 2>&1 &
USER_PID=$!

# Wait for services to start
sleep 5

echo -e "${GREEN}✓ Services started!${NC}"
echo -e ""
echo -e "Authentication Service PID: $AUTH_PID"
echo -e "User Service PID: $USER_PID"
echo -e ""
echo -e "View logs:"
echo -e "  Auth: tail -f logs/auth.log"
echo -e "  User: tail -f logs/user.log"
echo -e ""

# Wait for Ctrl+C
trap "echo -e '\n${YELLOW}Stopping services...${NC}'; kill $AUTH_PID $USER_PID; echo -e '${GREEN}Services stopped${NC}'; exit 0" INT

wait
