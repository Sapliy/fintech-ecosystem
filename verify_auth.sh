#!/bin/bash
set -e

BASE_URL="http://localhost:8081"
EMAIL="test_refresh_$(date +%s)@example.com"
PASSWORD="password123"

echo "1. Registering user $EMAIL..."
curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" > register.json

echo " - Registration response:"
cat register.json
echo ""

echo "2. Logging in..."
curl -v -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" \
  -c cookies.txt > login.json 2> login_headers.txt

echo " - Login response:"
cat login.json
echo ""
echo " - Cookies received:"
cat cookies.txt
echo ""

# Extract Refresh Token from cookies
REFRESH_COOKIE=$(grep "refreshtoken" cookies.txt | awk '{print $7}')

if [ -z "$REFRESH_COOKIE" ]; then
  echo "ERROR: Refreshtoken cookie not found!"
  exit 1
fi

echo "3. Testing Refresh..."
curl -v -s -X POST "$BASE_URL/refresh" \
  -b cookies.txt \
  > refresh.json 2> refresh_headers.txt

echo " - Refresh response:"
cat refresh.json
echo ""

TOKEN=$(grep -o '"token":"[^"]*' refresh.json | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "ERROR: New token not found in refresh response"
    exit 1
fi

echo "SUCCESS: Got new token: ${TOKEN:0:20}..."
