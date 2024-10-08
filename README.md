# chat.gg-api

# TODO (pending)
- Deploy on AWS
- Setup CI/CD pipelines

This repository contains the backend for a chat application built using Golang, Gin, and PostgreSQL with GORM as the ORM. It handles user authentication, message routing, contact management, and real-time communication via WebSockets.

# Tech Stack
- Golang: Backend language
- Gin Framework: Web framework for routing and middleware
- PostgreSQL: Relational database for storing persistent data
- GORM: Golang ORM library to interact with PostgreSQL
- Gorilla WebSocket: For implementing websocket protocol
- JWT (JSON Web Tokens): For secure authentication and authorization

# API Endpoints

- POST /auth/register: Register a new user
- POST /auth/login: Login an existing user
- GET /contacts: Retrieve contacts registered in the app
- WebSocket for real-time messaging: ws://localhost:8080/ws
