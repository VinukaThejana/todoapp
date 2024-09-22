# TodoApp - Microservices with gRPC

TodoApp is a robust task management application built using a microservices architecture with gRPC for inter-service communication. This project demonstrates the implementation of a scalable and modular system for managing user authentication and todo items.

## Features

- User Authentication
  - Registration
  - Login
  - Access token refresh
  - Logout
- Todo Management
  - Create todos
  - Update todos
  - Delete todos

## Architecture

TodoApp consists of three microservices:

1. **Auth Service**: Handles user authentication and token management.
2. **Todo Service**: Manages todo items (create, update, delete).
3. **API Service**: Acts as the gateway, handling user requests and communicating with other microservices.

## Technologies Used

- gRPC: For efficient inter-service communication
- golang: For building the microservices
- sqlite: For storing user and todo data
- redis: For validating access tokens
- docker: For containerizing the microservices

## Getting Started

### Prerequisites

- `just` installed on your system (It it not required but it makes running the project easier)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/todoapp.git
   cd todoapp
   ```

2. Run each micorservice

   - Auth Service
     ```
     just auth
     ```

   - Todo Service
     ```
     just todo
     ```

   - API Service
     ```
     just run
     ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
