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
  response=$(curl -s -X POST "$MINIKUBE_URL/auth/register" \
       -H "Content-Type: application/json" \
       -d '{
            "username": "testyuser",
            "password": "testypass",
            "email": "testy@example.com"
           }')
  accessToken=$(echo $response | jq -r .access_token)
  refreshToken=$(echo $response | jq -r .refresh_token)

  if [ -z "$accessToken" ] || [ -z "$refreshToken" ]; then
    echo "User registration failed with response $response"
    exit 1
  fi
  echo $response
fi

# Login and get token
echo -e "\n\nLogging in..."
response=$(curl -s -X POST "$MINIKUBE_URL/auth/login" \
             -H "Content-Type: application/json" \
             -d '{
                  "username": "testyuser",
                  "password": "testypass"
                 }')
accessToken=$(echo $response | jq -r .access_token)
refreshToken=$(echo $response | jq -r .refresh_token)

if [ -z "$accessToken" ] || [ -z "$refreshToken" ]; then
  echo "Failed to get tokens from login response"
  exit 1
fi
echo $response

echo "Refresh token:"
echo $refreshToken
# Refresh token
echo -e "\n\nRefreshing token..."
response=$(curl -s -X POST "$MINIKUBE_URL/auth/refresh" \
             -H "Content-Type: application/json" \
             -d '{
                  "refresh_token": "'"$refreshToken"'"
                 }')
newAccessToken=$(echo $response | jq -r .access_token)
newRefreshToken=$(echo $response | jq -r .refresh_token)

if [ -z "$newAccessToken" ] || [ -z "$newRefreshToken" ]; then
  echo "Failed to get tokens from refresh response"
  exit 1
fi
accessToken=$newAccessToken
refreshToken=$newRefreshToken
echo $response

# Create a source
if [ -z "$SKIP_SOURCES" ]; then
     echo -e "\n\nCreating a source..."
     response=$(curl -s -X POST "$MINIKUBE_URL/sources" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $accessToken" \
          -d '{
               "name": "Bank Current",
               "type": "SAVINGS",
               "balance": 10000.0
          }')
     sourceID=$(echo $response | jq -r .id)
     echo $response

     # Create a category
     echo -e "\n\nCreating a category..."
     response=$(curl -s -X POST "$MINIKUBE_URL/categories" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $accessToken" \
          -d '{
               "name": "Food",
               "description": "Food related transactions",
               "icon": "food"
          }')
     categoryID=$(echo $response | jq -r .id)
     echo $response
fi

# Create a transaction
echo -e "\n\nCreating a transaction... source: $sourceID,category: $categoryID"
response=$(curl -s -X POST "$MINIKUBE_URL/transactions" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $accessToken" \
     -d '{
          "amount": 90.50,
          "type": "expense",
          "description": "Food transaction",
          "source_id": '"$sourceID"',
          "category_id": '"$categoryID"'
         }')
echo $response

# Fetch transactions
echo -e "\n\nFetching transactions..."
response=$(curl -s -X GET "$MINIKUBE_URL/transactions" \
     -H "Authorization: Bearer $accessToken")
txnID=$(echo $response | jq -r '.[0].id')
echo $txnID

# Update a transaction
echo -e "\n\nUpdating a transaction..."
response=$(curl -s -X PUT "$MINIKUBE_URL/transactions/$txnID" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $accessToken" \
     -d '{
          "amount": 70.75,
          "type": "expense",
          "description": "Tasty transaction",
          "source_id": '"$sourceID"',
          "category_id": '"$categoryID"',
          "tags": ["eat", "binge"]
         }')
echo $response

# Get transaction
echo -e "\n\nFetching transactions..."
response=$(curl -s -X GET "$MINIKUBE_URL/transactions/$txnID" \
     -H "Authorization: Bearer $accessToken")
echo $response


# Fetch transactions
echo -e "\n\nTest KV..."
response=$(curl -s -X GET "$MINIKUBE_URL/health")

echo $response 
echo -e "\n\nDone testing."