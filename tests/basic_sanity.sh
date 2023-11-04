#!/bin/bash

MINIKUBE_URL=$1
SKIP_REGISTER=$2
SKIP_SOURCES=$3

if [ -z "$MINIKUBE_URL" ]; then
  echo "Usage: $0 <MINIKUBE_URL> [skip_register] [skip_sources]"
  exit 1
fi

# Register a user
if [ -z "$SKIP_REGISTER" ]; then
  echo "Registering a new user..."
  response=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$MINIKUBE_URL/auth/register" \
       -H "Content-Type: application/json" \
       -d '{
            "username": "testuser",
            "password": "testpass",
            "email": "test@example.com"
           }')
  if [ "$response" -ne 200 ]; then
    echo "User registration failed with HTTP status $response"
    exit 1
  fi
fi

# Login and get token
echo -e "\n\nLogging in..."
response=$(curl -s -X POST "$MINIKUBE_URL/auth/login" \
             -H "Content-Type: application/json" \
             -d '{
                  "username": "testuser",
                  "password": "testpass"
                 }')
token=$(echo $response | jq -r .token)

if [ -z "$token" ]; then
  echo "Failed to get a token from login response"
  exit 1
fi

# Create a source
if [ -z "$SKIP_SOURCES" ]; then
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
     echo $response
fi

# Create a transaction
echo -e "\n\nCreating a transaction... source: $sourceID,category: $categoryID"
response=$(curl -s -X POST "$MINIKUBE_URL/transactions" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $token" \
     -d '{
          "amount": 100.50,
          "type": "expense",
          "description": "My transaction",
          "source_id": '"$sourceID"',
          "category_id": '"$categoryID"'
         }')
echo $response

# Fetch transactions
echo -e "\n\nFetching transactions..."
response=$(curl -s -X GET "$MINIKUBE_URL/transactions" \
     -H "Authorization: Bearer $token")
txnID=$(echo $response | jq -r '.[0].id')
echo $txnID

# Update a transaction
echo -e "\n\nUpdating a transaction..."
response=$(curl -s -X PUT "$MINIKUBE_URL/transactions/$txnID" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $token" \
     -d '{
          "amount": 110.75,
          "type": "expense",
          "description": "Newer transaction",
          "source_id": '"$sourceID"',
          "category_id": '"$categoryID"',
          "tags": ["movie", "night"]
         }')
echo $response

# Get transaction
echo -e "\n\nFetching transactions..."
response=$(curl -s -X GET "$MINIKUBE_URL/transactions/$txnID" \
     -H "Authorization: Bearer $token")
echo $response


# Fetch transactions
echo -e "\n\nTest KV..."
response=$(curl -s -X GET "$MINIKUBE_URL/health")

echo $response 
echo -e "\n\nDone testing."