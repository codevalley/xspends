# XSpends API Documentation

## Overview
This document outlines the available API endpoints for the XSpends personal expense management application.

### Base URL
`http://localhost:8080` (Adjust according to deployment)

---

## Authentication Endpoints

### 1. Register User
- **Endpoint**: `/auth/register`
- **Method**: POST
- **Description**: Registers a new user with email and password.
- **Request Format**:
  ```json
  {
    "email": "user@example.com",
    "password": "yourpassword"
  }
  ```
- **Response Format**:
  ```json
  {
    "user_id": 1,
    "email": "user@example.com"
  }
  ```

### 2. Login User
- **Endpoint**: `/auth/login`
- **Method**: POST
- **Description**: Logs in an existing user with email and password.
- **Request Format**:
  ```json
  {
    "email": "user@example.com",
    "password": "yourpassword"
  }
  ```
- **Response Format**:
  ```json
  {
    "token": "jwt-token-here"
  }
  ```

---

## User Profile Management

---

## Fund Source Management

### 1. List Sources
- **Endpoint**: `/sources` 
- **Method**: GET
- **Description**: Retrieves all sources for the authenticated user.
- **Request Format**: N/A (Authorization token required)
- **Response Format**:
  ```json
  [
    {
      "source_id": 1,
      "source_name": "Bank Account",
      // other source details
    },
    // ... other sources
  ]
  ```

---

## Health Check Endpoint

### 1. Check Health
- **Endpoint**: `/health`
- **Method**: GET
- **Description**: Checks the health status of the application.
- **Response Format**:
  ```json
  {
    "status": "UP"
  }
  ```

---

## Spend Tag Management

### 1. List Tags
- **Endpoint**: `/tags/list`
- **Method**: GET
- **Description**: Retrieves a list of all tags for the authenticated user.
- **Request Format**: N/A (Authorization token required)
- **Response Format**:
  ```json
  [
    {
      "tag_id": 1,
      "tag_name": "Groceries",
      // other tag details
    },
    // ... other tags
  ]
  ```

### 2. Create Tag
- **Endpoint**: `/tags/create`
- **Method**: POST
- **Description**: Creates a new tag for the authenticated user.
- **Request Format**:
  ```json
  {
    "tag_name": "New Tag",
    // other tag details
  }
  ```
- **Response Format**:
  ```json
  {
    "tag_id": 2,
    "tag_name": "New Tag",
    // other tag details
  }
  ```

### 3. Update Tag
- **Endpoint**: `/tags/update`
- **Method**: PUT
- **Description**: Updates an existing tag for the authenticated user.
- **Request Format**:
  ```json
  {
    "tag_id": 2,
    "tag_name": "Updated Tag",
    // other updated details
  }
  ```
- **Response Format**:
  ```json
  {
    "tag_id": 2,
    "tag_name": "Updated Tag",
    // other updated details
  }
  ```

### 4. Delete Tag
- **Endpoint**: `/tags/delete/{id}`
- **Method**: DELETE
- **Description**: Deletes a specific tag by ID for the authenticated user.
- **Request Format**: N/A (Tag ID in URL)
- **Response Format**:
  ```json
  {
    "message": "tag deleted successfully"
  }
  ```

---

## Tags 

### 1. List Transaction Tags
- **Endpoint**: `/transaction/tags/list/{transaction_id}`
- **Method**: GET
- **Description**: Retrieves a list of tags associated with a specific transaction.
- **Request Format**: N/A (Transaction ID in URL)
- **Response Format**:
  ```json
  [
    {
      "tag_id": 1,
      "tag_name": "Groceries",
      // other tag details
    },
    // ... other tags
  ]
  ```

### 2. Add Tag to Transaction
- **Endpoint**: `/transaction/tags/add`
- **Method**: POST
- **Description**: Adds a tag to a specific transaction.
- **Request Format**:
  ```json
  {
    "transaction_id": 123,
    "tag_id": 1
  }
  ```
- **Response Format**:
  ```json
  {
    "message": "tag added successfully to the transaction"
  }
  ```

### 3. Remove Tag from Transaction
- **Endpoint**: `/transaction/tags/remove`
- **Method**: DELETE
- **Description**: Removes a tag from a specific transaction.
- **Request Format**:
  ```json
  {
    "transaction_id": 123,
    "tag_id": 1
  }
  ```
- **Response Format**:
  ```json
  {
    "message": "tag removed successfully from the transaction"
  }
  ```
