# GoAuth - An authentication app which handles register & login via RESTful API

I made this project while learning golang.This backend app is my first project, that I made using **GoFiber**(A backend framework for Go).

## Requirement

- Go
- PostgreSQL

## Starting the backend server

**Setting up files:**

```
    git clone https://github.com/imrishabk/goauth.git
    cd goauth
    mv example.env .env  # Edit your env file to setup database. (Make sure to have postgres and access to database.)
```

**Running Server**

```
    go run cmd/main.go
```

OR

```
    go build cmd/main.go
    ./goauth
```

## End Points

### 1. Register

```
    POST /api/auth/register
```

### 2. Login

```
    POST /api/auth/login
```

### 3. Default

```
    GET /
    POST /
```

Just returns simple message with goauth working..

## How to send request ?

### 1. Register

```json
    POST /api/auth/register
    
    Body:

    {
        "username": "",
        "display_name": "",
        "password": "",
        "email": "",
        "phone_number": ""
    }
```

Fill the fields properly & Send Request.

### Status Response

**1. 400 - Bad Request**

- Invalid JSON request
- Invalid email
- Invalid username
- Email already taken
- Username already taken

**2. 500 - Internal Server Error**

- Failed to check if email already exists
- Failed to check if username already exists
- If failed to hash the password.
- Failed to write user into the database

**3. 201 - Status Created**

- If user is successfully registered.

### 2. Login

```json
    POST /api/auth/register

    Body:

    {
        "identity": "", 
        "password": ""
    }
```

*Identity field should either contain email or username*
Fill the field properly.

### Status Response

**1. 400 - Bad Request**

- Invalid JSON request

**2. 500 - Internal Server Error**

- Failed to retreive info from database
- Failed to generate JWT Token (JWT are not implemented yet. It just generates token for now)

**3. 401 - Status Unauthorized**

- Entered identity or password doesn't match

**4. 200 - OK**

- Returns a JWT Token for session based authentication. (JWT are not implemented yet. It just generates token for now)

## Packages Used

- GoFiber (The core backend framework)
- GORM (For Object-Relation-Mapping & also postgres driver)
- godotenv (To load environment variables from .env file)

## Future Plans

- Implement Middlewares for protecting routes.
- Implement JWT for session based authentication.

## License

This project is licensed under the [Apache License 2.0](LICENSE).

© 2025 Rishab Karki
