#!/bin/bash

MINIKUBE_URL=$1
SKIP_REGISTER=$2

if [ -z "$MINIKUBE_URL" ]; then
  echo "Usage: $0 <MINIKUBE_URL> [skip_register]"
  exit 1
fi

if [ -z "$SKIP_REGISTER" ]; then
  # Register a user
  echo "Registering a new user..."
  curl -X POST "$MINIKUBE_URL/register" \
       -H "Content-Type: application/json" \
       -d '{
            "username": "testuser",
            "password": "testpass",
            "email": "test@example.com"
           }'
fi

# Login and get token
echo -e "\n\nLogging in..."
response=$(curl -s -X POST "$MINIKUBE_URL/login" \
             -H "Content-Type: application/json" \
             -d '{
                  "username": "testuser",
                  "password": "testpass"
                 }')
token=$(echo $response | jq -r .token)

# Create a transaction
echo -e "\n\nCreating a transaction..."
curl -X POST "$MINIKUBE_URL/transactions" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $token" \
     -d '{
          "title": "Test Transaction",
          "amount": 100.50,
          "type": "expense",
          "category": "Groceries",
          "tags": ["food", "essentials"]
         }'

# Fetch transactions
echo -e "\n\nFetching transactions..."
curl -X GET "$MINIKUBE_URL/transactions" \
     -H "Authorization: Bearer $token"

# Update a transaction (assuming transaction ID is 1 for simplicity)
echo -e "\n\nUpdating a transaction..."
curl -X PUT "$MINIKUBE_URL/transactions/1" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $token" \
     -d '{
          "title": "Updated Transaction",
          "amount": 110.75,
          "type": "expense",
          "category": "Entertainment",
          "tags": ["movie", "night out"]
         }'

echo -e "\n\nDone testing."
