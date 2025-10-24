# Project Learning & CV Insights

## Project Overview

This repository contains the source code for SnippetBox, a full-stack web application designed for creating, managing, and sharing text-based snippets. The application provides a secure and intuitive platform for users to manage their code snippets or notes, with features including user authentication, session management, and a clean, user-friendly interface. The project is built entirely in Go, following best practices for web development and application security.

## Tech Stack and Key Technologies

*   **Language:** Go (Golang)
*   **Database:** MySQL
*   **Web Framework:** While not a full framework, the project leverages a combination of Go's standard `net/http` library with high-performance routing via `julienschmidt/httprouter`.
*   **Frontend:** Server-side rendered HTML templates using Go's `html/template` package.
*   **Styling:** Plain CSS for styling the user interface.
*   **Session Management:** `alexedwards/scs` for secure, database-backed session management.
*   **Security:** TLS for encrypted communication (HTTPS), password hashing with `golang.org/x/crypto`.

## Notable Libraries

*   **`julienschmidt/httprouter`:** A high-performance, lightweight HTTP router that supports URL parameters and scales efficiently. Chosen for its speed and minimal overhead.
*   **`alexedwards/scs/v2`:** A session management library that provides secure, flexible, and easy-to-use session handling. It is configured with a `mysqlstore` to persist sessions in the database.
*   **`go-sql-driver/mysql`:** The official MySQL driver for Go, enabling seamless and efficient communication with the MySQL database.
*   **`go-playground/form/v4`:** A versatile library for decoding form data from HTTP requests into Go structs, simplifying data validation and processing.
*   **`justinas/alice`:** A lightweight library for chaining middleware, making it easy to compose and apply middleware to HTTP handlers in a clean and readable way.
*   **`golang.org/x/crypto`:** Provides robust cryptographic primitives, used in this project for securely hashing and verifying user passwords.

## Major Achievements and Skills Demonstrated

*   **Full-Stack Application Development:** Designed, developed, and deployed a complete web application from scratch using Go for both the backend and frontend templating.
*   **Secure User Authentication System:** Implemented a robust user authentication and authorization system, including password hashing, session management, and protected routes.
*   **Database Design and Management:** Designed a relational database schema for snippets and users, and managed database interactions using prepared statements to prevent SQL injection attacks.
*   **RESTful API Design Principles:** Structured the application's routes and handlers following RESTful principles for managing snippet and user resources.
*   **Custom Middleware Implementation:** Developed a chain of custom middleware for logging, security headers, and session handling to process requests efficiently and securely.
*   **Configuration Management:** Implemented a flexible configuration system using command-line flags to manage application settings like database connections and server addresses.
*   **HTTPS and TLS Implementation:** Configured and deployed the server with TLS to ensure secure, encrypted communication between the client and server.
*   **Test-Driven Development (TDD):** Wrote unit and integration tests for handlers and middleware to ensure code quality and maintainability.

## Skills Gained/Reinforced

*   **Go Programming:** Advanced proficiency in Go, including concurrency, error handling, and standard library usage.
*   **Backend Web Development:** Deep understanding of building web applications, handling HTTP requests, and managing server-side logic.
*   **Database Management (MySQL):** Skills in schema design, querying, and interacting with a relational database in a production environment.
*   **Application Security:** Practical experience in implementing security best practices, including password hashing, session security, and preventing common vulnerabilities.
*   **Software Architecture:** Experience in designing and structuring a web application with a clear separation of concerns (e.g., handlers, models, views).
*   **API Design:** Proficiency in designing and implementing clean and intuitive APIs.
*   **Testing:** Competence in writing effective unit and integration tests to ensure application reliability.
