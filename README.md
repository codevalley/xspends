# XSpends: Personal Expense Management API Service

## Overview
XSpends is a RESTful API server designed for personal expense management. It offers functionalities like user authentication, transaction recording, fund source tracking, and spend categorization.

## Features
- **User Authentication**: Secure user registration and login.
- **Transaction Management**: Record and manage financial transactions.
- **Fund Source Management**: Keep track of different financial sources.
- **Expense Categorization**: Organize expenses with categories and tags.

## Getting Started

### Prerequisites
- Docker
- Helm (for Kubernetes deployments)
- Minikube (for local Kubernetes setup)

### Installation
For detailed installation instructions, refer to [install.md](install.md). This includes steps for Minikube setup, TiDB setup using Helm, and deploying XSpends.

### Accessing the Service
Use `minikube service xspends-service --url` to get the service URL for accessing the API.

## Usage
Make RESTful API requests to the provided service URL. Example endpoints:
- `/auth/register` (POST): Register a new user.
- `/auth/login` (POST): Login for existing users.
- `/transactions` (POST/GET): Manage transactions.

Refer to the API documentation (if available) for detailed usage.

## Development
- **Local Development**: Use `docker-compose up`.
- **Testing**: Execute `tests/basic_sanity.sh` for basic checks.
- **Rebuild and Redeploy**: Instructions available in [install.md](install.md).

## Contributing
Contributions to XSpends are welcome! Please read our contributing guidelines for details on how to contribute.

## License
This project is licensed under the MIT License - see the `LICENSE` file for details.
