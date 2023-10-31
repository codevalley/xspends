#!/bin/bash

MINIKUBE_URL=$1
SKIP_REGISTER=$2
SKIP_SOURCES=$3

if [ -z "$MINIKUBE_URL" ]; then
  echo "Usage: $0 <MINIKUBE_URL> [skip_register] [skip_sources]"
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

if [ -z "$SKIP_SOURCES" ]; then
     # Create a source
     echo -e "\n\nCreating a source..."
     response=$(curl -s -X POST "$MINIKUBE_URL/sources" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $token" \
          -d '{
               "name": "Bank Savings",
               "type": "SAVINGS",
               "balance": 20000.0
          }')
     sourceID=$(echo $response | jq -r .id)
     echo $response

     # Create a category
     echo -e "\n\nCreating a category..."
     response=$(curl -s -X POST "$MINIKUBE_URL/categories" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $token" \
          -d '{
               "name": "Groceries",
               "description": "Grocery related transactions",
               "icon": "shopping-cart"
          }')
     categoryID=$(echo $response | jq -r .id)
fi
# Create a transaction
echo -e "\n\nCreating a transaction... source: $sourceID,category: $categoryID"
curl -X POST "$MINIKUBE_URL/transactions" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $token" \
     -d '{
          "amount": 100.50,
          "type": "expense",
          "description": "My transaction",
          "source_id": '"$sourceID"',
          "category_id": '"$categoryID"'
         }'
     txnID=$(echo $response | jq -r .id)
# Fetch transactions
echo -e "\n\nFetching transactions..."
curl -X GET "$MINIKUBE_URL/transactions" \
     -H "Authorization: Bearer $token"

# Update a transaction (assuming transaction ID is 1 for simplicity)
echo -e "\n\nUpdating a transaction..."
curl -X PUT "$MINIKUBE_URL/transactions/484552817248305200" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $token" \
     -d '{
          "amount": 110.75,
          "type": "expense",
          "description": "Newer transaction",
          "source_id": '"$sourceID"',
          "category_id": '"$categoryID"',
          "tags": ["movy", "knight"]
         }'

# Get transaction
echo -e "\n\nFetching transactions..."
curl -X GET "$MINIKUBE_URL/transactions/484552817248305200" \
     -H "Authorization: Bearer $token"

echo -e "\n\nDone testing."
