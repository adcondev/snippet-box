# Learning: Snippet Box Project

## Project Overview
Snippet Box is a secure, database-driven web application designed for creating, viewing, and managing text-based snippets (primarily code). It serves as a centralized repository where users can sign up, log in, and save snippets with expiration dates. The project demonstrates a robust implementation of a web server in Go, featuring user authentication, session management, and secure data handling.

## Tech Stack and Key Technologies
*   **Language:** Go (Golang) v1.22+
*   **Database:** MySQL (Relational Database)
*   **Web Framework/Router:** `net/http` standard library with `julienschmidt/httprouter` for efficient routing.
*   **Frontend:** Server-Side Rendered HTML5 Templates, CSS3.
*   **Automation:** GNU Make (`Makefile`) for build scripts and task automation.
*   **CI/CD:** GitHub Actions for automated testing and release management.
*   **Version Control:** Git.

## Notable Libraries
*   **`github.com/julienschmidt/httprouter`**: A high-performance, lightweight request router. Solved the need for efficient pattern matching and parameter extraction in URLs.
*   **`github.com/justinas/alice`**: A middleware chaining library. Simplified the management of middleware stacks (e.g., logging, security, panic recovery) by allowing them to be composed cleanly.
*   **`github.com/alexedwards/scs/v2`**: A session management library. Handled the complexity of secure session storage and persistence using a MySQL backend.
*   **`github.com/go-playground/form/v4`**: A form decoding library. Automates the parsing of HTML form data into Go structs, reducing boilerplate code.
*   **`golang.org/x/crypto`**: Used for `bcrypt` password hashing, ensuring secure user authentication.
*   **`github.com/go-sql-driver/mysql`**: The official MySQL driver for Go's `database/sql` package.

## Major Achievements and Skills Demonstrated
*   **Backend Architecture**: Designed and implemented a modular web application structure separating concerns between handlers, models (database logic), and middleware.
*   **Database Design**: Created a normalized MySQL database schema for Users and Snippets, including proper indexing for performance optimization.
*   **Secure Authentication**: Implemented a complete user authentication system with secure password hashing (bcrypt), session management, and protection against common vulnerabilities (CSRF, session hijacking).
*   **Middleware Implementation**: Developed custom middleware for request logging, panic recovery, and enforcing security headers (CSP, X-Frame-Options, XSS Protection).
*   **CI/CD Pipeline**: Configured a GitHub Actions workflow to automatically run tests on pull requests and generate releases with conventional changelogs.
*   **Testing**: Wrote unit and integration tests for handlers, models, and middleware to ensure code reliability.
*   **TLS/SSL**: Configured the server to run securely over HTTPS using TLS 1.2/1.3 with modern cipher suites.

## Skills Gained/Reinforced
*   **Go Web Development**: Deepened understanding of `net/http`, handlers, and the Go context package.
*   **RESTful Principles**: Experience in designing resource-oriented routes and handling HTTP methods correctly.
*   **SQL & Database Management**: Writing raw SQL queries and managing database connections in Go.
*   **Web Security**: Practical application of security best practices (HTTPS, Secure Cookies, OWASP headers).
*   **DevOps & Automation**: Writing Makefiles and configuring CI/CD workflows.
*   **Concurrent Programming**: Understanding Go's concurrency model in the context of handling web requests.
