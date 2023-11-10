
# XSpends API Specification

## 1. Register User

- **Endpoint**: `/auth/register`
- **Method**: POST
- **Description**: Register a new user account.
- **Request Format**:
  ```json
  {
    "username": "newuser",
    "password": "password123",
    "email": "newuser@example.com"
  }
  ```
- **Response Format**:
  ```json
  {
    "user_id": "12345",
    "username": "newuser",
    "email": "newuser@example.com"
  }
  ```
- **Error Response**: (if user already exists or input is invalid)
  ```json
  {
    "error": "user already exists or invalid input"
  }
  ```

## 2. Login User

- **Endpoint**: `/auth/login`
- **Method**: POST
- **Description**: Log in an existing user.
- **Request Format**:
  ```json
  {
    "username": "existinguser",
    "password": "password123"
  }
  ```
- **Response Format**:
  ```json
  {
    "token": "jwt-token-here",
    "user_id": "12345"
  }
  ```
- **Error Response**: (if credentials are incorrect)
  ```json
  {
    "error": "invalid username or password"
  }
  ```

## 3. Refresh Token

- **Endpoint**: `/auth/refresh`
- **Method**: POST
- **Description**: Refresh the authentication token.
- **Request Format**:
  ```json
  {
    "refresh_token": "existing-refresh-token"
  }
  ```
- **Response Format**:
  ```json
  {
    "new_token": "new-jwt-token-here",
    "new_refresh_token": "new-refresh-token"
  }
  ```
- **Error Response**: (if the refresh token is invalid or expired)
  ```json
  {
    "error": "invalid or expired refresh token"
  }
  ```

## 4. Logout User

- **Endpoint**: `/auth/logout`
- **Method**: POST
- **Description**: Log out the current user.
- **Request Format**: No body required (Authorization header with token is needed)
- **Response Format**:
  ```json
  {
    "message": "successfully logged out"
  }
  ```
- **Error Response**: (if not logged in or invalid token)
  ```json
  {
    "error": "user not logged in or invalid token"
  }
  ```
Continuing with the API specification for the `/sources` endpoints based on the analysis of the `routes.go` and corresponding handler files in the `xspends` project:

---

## 2. List Sources

- **Endpoint**: `/sources`
- **Method**: GET
- **Description**: Retrieve a list of all financial sources for the authenticated user.
- **Request Format**: No body required (Authorization header with token is needed)
- **Response Format**:
  ```json
  [
    {
      "source_id": 1,
      "source_name": "Bank Account",
      "balance": 1000.00,
      "type": "SAVINGS"
      // other source details
    },
    // ... other sources
  ]
  ```
- **Error Response**: (e.g., if the user is not authenticated)
  ```json
  {
    "error": "unauthorized access"
  }
  ```

## 2. Create Source

- **Endpoint**: `/sources`
- **Method**: POST
- **Description**: Create a new financial source for the authenticated user.
- **Request Format**:
  ```json
  {
    "source_name": "New Source",
    "balance": 500.00,
    "type": "CREDIT"
    // other source details
  }
  ```
- **Response Format**:
  ```json
  {
    "source_id": 2,
    "source_name": "New Source",
    "balance": 500.00,
    "type": "CREDIT"
    // other source details
  }
  ```
- **Error Response**: (e.g., if input data is invalid)
  ```json
  {
    "error": "invalid input data"
  }
  ```

## 3. Get Source

- **Endpoint**: `/sources/:id`
- **Method**: GET
- **Description**: Retrieve details of a specific financial source by ID.
- **Request Format**: Source ID in URL path
- **Response Format**:
  ```json
  {
    "source_id": 1,
    "source_name": "Bank Account",
    "balance": 1000.00,
    "type": "SAVINGS"
    // other source details
  }
  ```
- **Error Response**: (e.g., if the source is not found)
  ```json
  {
    "error": "source not found"
  }
  ```

## 4. Update Source

- **Endpoint**: `/sources/:id`
- **Method**: PUT
- **Description**: Update an existing financial source.
- **Request Format**:
  ```json
  {
    "source_name": "Updated Source",
    "balance": 1200.00,
    "type": "CREDIT"
    // other updated details
  }
  ```
- **Response Format**:
  ```json
  {
    "source_id": 1,
    "source_name": "Updated Source",
    "balance": 1200.00,
    "type": "CREDIT"
    // other updated details
  }
  ```
- **Error Response**: (e.g., if the update fails)
  ```json
  {
    "error": "update failed"
  }
  ```

## 5. Delete Source

- **Endpoint**: `/sources/:id`
- **Method**: DELETE
- **Description**: Delete a specific financial source by ID.
- **Request Format**: Source ID in URL path
- **Response Format**:
  ```json
  {
    "message": "source deleted successfully"
  }
  ```
- **Error Response**: (e.g., if the source is not found)
  ```json
  {
    "error": "source not found"
  }
  ```
---


## 1. List Categories

- **Endpoint**: `/categories`
- **Method**: GET
- **Description**: Retrieve a list of all categories for the authenticated user.
- **Request Format**: No body required (Authorization header with token is needed)
- **Response Format**:
  ```json
  [
    {
      "category_id": 1,
      "name": "Groceries",
      "description": "Grocery shopping",
      "icon": "shopping-cart"
      // other category details
    },
    // ... other categories
  ]
  ```
- **Error Response**: (e.g., unauthorized access)
  ```json
  {
    "error": "unauthorized access"
  }
  ```

## 2. Create Category

- **Endpoint**: `/categories`
- **Method**: POST
- **Description**: Create a new category for the authenticated user.
- **Request Format**:
  ```json
  {
    "name": "Utilities",
    "description": "Monthly bills and utilities",
    "icon": "utilities-icon"
    // other category details
  }
  ```
- **Response Format**:
  ```json
  {
    "category_id": 2,
    "name": "Utilities",
    "description": "Monthly bills and utilities",
    "icon": "utilities-icon"
    // other category details
  }
  ```
- **Error Response**: (e.g., invalid input data)
  ```json
  {
    "error": "invalid input data"
  }
  ```

## 3. Get Category

- **Endpoint**: `/categories/:id`
- **Method**: GET
- **Description**: Retrieve details of a specific category by ID.
- **Request Format**: Category ID in URL path
- **Response Format**:
  ```json
  {
    "category_id": 1,
    "name": "Groceries",
    "description": "Grocery shopping",
    "icon": "shopping-cart"
    // other category details
  }
  ```
- **Error Response**: (e.g., category not found)
  ```json
  {
    "error": "category not found"
  }
  ```

## 4. Update Category

- **Endpoint**: `/categories/:id`
- **Method**: PUT
- **Description**: Update an existing category.
- **Request Format**:
  ```json
  {
    "name": "Updated Category",
    "description": "Updated description",
    "icon": "updated-icon"
    // other updated details
  }
  ```
- **Response Format**:
  ```json
  {
    "category_id": 1,
    "name": "Updated Category",
    "description": "Updated description",
    "icon": "updated-icon"
    // other updated details
  }
  ```
- **Error Response**: (e.g., update failed)
  ```json
  {
    "error": "update failed"
  }
  ```

## 5. Delete Category

- **Endpoint**: `/categories/:id`
- **Method**: DELETE
- **Description**: Delete a specific category by ID.
- **Request Format**: Category ID in URL path
- **Response Format**:
  ```json
  {
    "message": "category deleted successfully"
  }
  ```
- **Error Response**: (e.g., category not found)
  ```json
  {
    "error": "category not found"
  }
  ```
---

## 1. List Tags

- **Endpoint**: `/tags`
- **Method**: GET
- **Description**: Retrieve a list of all tags for the authenticated user.
- **Request Format**: No body required (Authorization header with token is needed)
- **Response Format**:
  ```json
  [
    {
      "tag_id": 1,
      "name": "Food",
      // other tag details
    },
    // ... other tags
  ]
  ```
- **Error Response**: (e.g., unauthorized access)
  ```json
  {
    "error": "unauthorized access"
  }
  ```

## 2. Create Tag

- **Endpoint**: `/tags`
- **Method**: POST
- **Description**: Create a new tag for the authenticated user.
- **Request Format**:
  ```json
  {
    "name": "Entertainment",
    // other tag details
  }
  ```
- **Response Format**:
  ```json
  {
    "tag_id": 2,
    "name": "Entertainment",
    // other tag details
  }
  ```
- **Error Response**: (e.g., invalid input data)
  ```json
  {
    "error": "invalid input data"
  }
  ```

## 3. Get Tag

- **Endpoint**: `/tags/:id`
- **Method**: GET
- **Description**: Retrieve details of a specific tag by ID.
- **Request Format**: Tag ID in URL path
- **Response Format**:
  ```json
  {
    "tag_id": 1,
    "name": "Food",
    // other tag details
  }
  ```
- **Error Response**: (e.g., tag not found)
  ```json
  {
    "error": "tag not found"
  }
  ```

## 4. Update Tag

- **Endpoint**: `/tags/:id`
- **Method**: PUT
- **Description**: Update an existing tag.
- **Request Format**:
  ```json
  {
    "name": "Updated Tag",
    // other updated details
  }
  ```
- **Response Format**:
  ```json
  {
    "tag_id": 1,
    "name": "Updated Tag",
    // other updated details
  }
  ```
- **Error Response**: (e.g., update failed)
  ```json
  {
    "error": "update failed"
  }
  ```

## 5. Delete Tag

- **Endpoint**: `/tags/:id`
- **Method**: DELETE
- **Description**: Delete a specific tag by ID.
- **Request Format**: Tag ID in URL path
- **Response Format**:
  ```json
  {
    "message": "tag deleted successfully"
  }
  ```
- **Error Response**: (e.g., tag not found)
  ```json
  {
    "error": "tag not found"
  }
  ```
---

## 1. List Transactions

- **Endpoint**: `/transactions`
- **Method**: GET
- **Description**: Retrieve a paginated list of transactions for the authenticated user.
- **Request Parameters**:
  - `page`: Page number for pagination (optional, default is 1).
  - `limit`: Number of transactions per page (optional, default is 10).
- **Request Format**: Query parameters for pagination.
- **Response Format**:
  ```json
  {
    "transactions": [
      {
        "transaction_id": 1,
        "amount": 50.00,
        "type": "Expense",
        "source_id": 1,
        "category_id": 1,
        "tags": [1, 2],
        "description": "Grocery shopping",
        "date": "2023-01-01T00:00:00Z"
        // other transaction details
      },
      // ... other transactions
    ],
    "page": 1,
    "limit": 10,
    "total_pages": 5,
    "total_transactions": 50
  }
  ```
- **Error Response**: (e.g., unauthorized access)
  ```json
  {
    "error": "unauthorized access"
  }
  ```

## 2. Create Transaction

- **Endpoint**: `/transactions`
- **Method**: POST
- **Description**: Create a new transaction for the authenticated user.
- **Request Format**:
  ```json
  {
    "amount": 100.00,
    "type": "Income",
    "source_id": 1,
    "category_id": 1,
    "tags": [1, 3],
    "description": "Salary",
    "date": "2023-01-15T00:00:00Z"
    // other transaction details
  }
  ```
- **Response Format**:
  ```json
  {
    "transaction_id": 2,
    "amount": 100.00,
    "type": "Income",
    // other transaction details
  }
  ```
- **Error Response**: (e.g., invalid input data)
  ```json
  {
    "error": "invalid input data"
  }
  ```

## 3. Get Transaction

- **Endpoint**: `/transactions/:id`
- **Method**: GET
- **Description**: Retrieve details of a specific transaction by ID.
- **Request Format**: Transaction ID in URL path
- **Response Format**:
  ```json
  {
    "transaction_id": 1,
    "amount": 50.00,
    "type": "Expense",
    // other transaction details
  }
  ```
- **Error Response**: (e.g., transaction not found)
  ```json
  {
    "error": "transaction not found"
  }
  ```

## 4. Update Transaction

- **Endpoint**: `/transactions/:id`
- **Method**: PUT
- **Description**: Update an existing transaction.
- **Request Format**:
  ```json
  {
    "amount": 75.00,
    "type": "Expense",
    // other updated details
  }
  ```
- **Response Format**:
  ```json
  {
    "transaction_id": 1,
    "amount": 75.00,
    "type": "Expense",
    // other updated details
  }
  ```
- **Error Response**: (e.g., update failed)
  ```json
  {
    "error": "update failed"
  }
  ```

## 5. Delete Transaction

- **Endpoint**: `/transactions/:id`
- **Method**: DELETE
- **Description**: Delete a specific transaction by ID.
- **Request Format**: Transaction ID in URL path
- **Response Format**:
  ```json
  {
    "message": "transaction deleted successfully"
  }
  ```
- **Error Response**: (e.g., transaction not found)
  ```json
  {
    "error": "transaction not found"
  }
  ```

