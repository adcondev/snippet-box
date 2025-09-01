# Snippet Box 📦

A modern, fast code snippet sharing platform built with Go. Share, organize, and collaborate on code snippets.

## Features

- 🔐 **Authentication**: User authentication and session management.
- 🔗 **Share Links**: Public links. 
- ⚡ **Fast**: Go backend with optimized database queries. Dynamic HTML templates for user interface.
- 🎯 **API**: RESTful API for integrations. CRUD operations for code snippets. Efficient request routing.

## Tech Stack

- **Backend**: Go 1.20+
- **Database**: PostgreSQL
- **Frontend**: HTML/CSS/Vanilla JS

## Getting Started

### Prerequisites

- Go
- MySQL

### Installation

1. Clone the repository
```sh
git clone https://github.com/consdotpy/snippetbox.git
```

2. Install dependencies
```sh
go mod download
```

3. Set up your MySQL database using the provided SQL scripts in the `sql/` directory. Be sure to change credentials in the `sql/create_user.sql` script before running it.

```sh
sudo mysql
source sql/*.sql
```


## Usage

Check `.env.example` for environment variables that need to be set. Then, run the application using `Makefile`:

```sh
make run
```

## Testing

To run quality control checks, execute:

```sh
make audit
```

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

## License


Distributed under the MIT License. See `LICENSE` for more information.

