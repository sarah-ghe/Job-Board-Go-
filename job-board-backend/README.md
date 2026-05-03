# Job Board API (Go)

## 🚀 Overview

This is a RESTful Job Board API built with Go.
It allows users to register, authenticate, and manage job postings.

The project follows a clean and scalable architecture:

```
Handler → Service → Repository → Database
```

---

## ✨ Features

* User registration and login
* JWT-based authentication
* Protected routes using middleware
* Full CRUD operations for jobs
* Multi-user system (jobs belong to users)
* Authorization: users can only modify their own jobs
* Secure password hashing
* Clean architecture with separation of concerns
* DTOs to prevent sensitive data exposure

---

## 🧠 Architecture

The application is structured into layers:

* **Handlers** → Handle HTTP requests and responses
* **Services** → Contain business logic
* **Repositories** → Interact with the database
* **Middleware** → Handle authentication and request validation

---

## 🔐 Authentication

Authentication is implemented using JWT.

Flow:

1. User logs in with email and password
2. Server verifies credentials
3. A JWT token is generated
4. The client sends the token in the `Authorization` header

```
Authorization: Bearer <token>
```

---

## 📦 Project Structure

```
backend/
│
├── handlers/
├── services/
├── repositories/
├── middleware/
├── models/
├── utils/
├── config/
└── main.go
```

---

## 🛠️ Technologies Used

* Go
* Gorilla Mux
* PostgreSQL
* JWT (authentication)
* Bcrypt (password hashing)

---

## 📌 API Endpoints

### Public Routes

* `POST /register` → Register a new user
* `POST /login` → Login and receive JWT

---

### Protected Routes

* `GET /me` → Get current user
* `POST /jobs` → Create job
* `GET /my-jobs` → Get user’s jobs
* `PUT /jobs/{id}` → Update job (owner only)
* `DELETE /jobs/{id}` → Delete job (owner only)

---

### Public Jobs

* `GET /jobs` → Get all jobs (no authentication required)

---

## 🔒 Security

* Passwords are hashed using bcrypt
* JWT tokens are used for authentication
* Sensitive fields (like passwords) are never exposed in API responses
* Authorization ensures users can only access their own data

---

