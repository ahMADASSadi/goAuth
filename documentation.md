# goAuth Project Guide

## 1. Running Locally

### **Prerequisites**
- Go 1.24+ installed
- [Make](https://www.gnu.org/software/make/) installed
- [Git](https://git-scm.com/) installed

### **Setup Steps**

1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd goAuth/src
   ```

2. **Configure environment variables:**
   - Copy .env.example to `.env` and adjust as needed:
     ```sh
     cp .env.example .env
     ```
   - Example `.env`:
     ```
     PORT=8000
     DB_URL="./test.sqlite3"
     ACCESS_EXPIRY="24h"
     SECRET_KEY="Some secret key"
     ```

3. **Install dependencies:**
   ```sh
   go mod download
   ```

4. **Build and run the application:**
   ```sh
   make build
   make run
   ```
   Or, for live reload (requires [air](https://github.com/air-verse/air)):
   ```sh
   make watch
   ```

5. **Access the API:**
   - The server runs at `http://localhost:8000/api/v1`

---

## 2. Running with Docker

### **Prerequisites**
- [Docker](https://www.docker.com/) installed
- [Docker Compose](https://docs.docker.com/compose/) installed

### **Steps**

1. **Build and start the container:**
   ```sh
   docker-compose up --build
   ```
   - The app will be available at `http://localhost:8000/api/v1`

2. **Environment Variables:**
   - Edit `src/.env` or pass environment variables via Docker Compose if needed.

---

## 3. API Documentation

### **Swagger/OpenAPI**

- **Swagger UI:**  
  After running the app, access API docs at:  
  `http://localhost:8000/swagger/index.html`

- **Endpoints Overview:**

  | Method | Endpoint                | Description                       |
  |--------|-------------------------|-----------------------------------|
  | POST   | `/api/v1/auth/request`  | Request OTP for phone number      |
  | POST   | `/api/v1/auth/verify`   | Verify OTP and get access token   |
  | GET    | `/api/v1/users/:id`     | Get user by ID                    |
  | GET    | `/api/v1/users`         | List users (pagination supported) |

- **Example Requests:**  
  See [src/requests/client.http](src/requests/client.http) for ready-to-use HTTP requests.

---

## 4. Database Choice Justification

### **Current Choice: SQLite**

- **File:** [src/internal/database/database.go](src/internal/database/database.go)
- **Configured via:** `DB_URL` in `.env` (default: `./test.sqlite3`)

#### **Why SQLite?**
- **Simplicity:** No server setup required; runs as a local file.
- **Zero Configuration:** Ideal for development and testing.
- **Portability:** The database is a single file, easy to move or reset.
- **Performance:** Fast for small to medium workloads and prototyping.

#### **Production Considerations**
- For production, switching to a more robust DBMS (e.g., PostgreSQL or MySQL) for:
  - Better concurrency
  - Scalability
  - Advanced features (backups, replication, etc.)

- The codebase is structured to allow easy migration to another database by changing the GORM driver and updating the `DB_URL`.

---

**References:**
- [src/internal/database/database.go](src/internal/database/database.go)
- [src/Makefile](src/Makefile)
- [docker-compose.yml](docker-compose.yml)
- [dockerFile](dockerFile)
- [src/requests/client.http](src/requests/client.http)

---

**Tip:**  
For further details on API request/response formats, see the Swagger UI or the handler comments in [src/internal/server/api](src/internal/server/api).   - The app will be available at `http://localhost:8000/api/v1`

1. **Environment Variables:**
   - Edit `src/.env` or pass environment variables via Docker Compose if needed.

2. **Request Docs:**
   - Simply just navigate to the `src/requests` and use the `client.http` if you have any http client installed.