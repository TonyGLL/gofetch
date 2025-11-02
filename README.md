# gofetch: A Local and Web Search Engine

`gofetch` is a high-performance search engine written in Go, designed to index content from both local files and websites. The system provides relevance-based search through a RESTful API and a minimalist web user interface.

## Key Features

- **Dual Indexing:** Support for indexing local directories (`.txt`, `.md`) and websites (future feature).
- **Relevance-Based Search:** Implementation of ranking algorithms like TF-IDF.
- **Robust Persistence:** Uses MongoDB to store the inverted index and metadata.
- **Concurrent Architecture:** Efficient design that leverages Go's goroutines for indexing.
- **Professional Operability:** Easy to configure, deploy, and monitor using Docker and GitHub Actions.

## Getting Started

### Prerequisites

- Go (1.18 or higher)
- MongoDB
- Docker and Docker Compose (for containerized development)

### Installation

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/TonyGLL/gofetch
    cd gofetch
    ```

2.  **Install dependencies:**
    ```sh
    go mod download
    ```

### Configuration

The application uses environment variables for configuration. You can set them in a `.env` file or export them in your shell.

-   `MONGODB_URI`: The MongoDB connection string. Default: `mongodb://localhost:27017`
-   `DB_NAME`: The name of the database to use. Default: `gofetch`

## Usage

### 1. Indexing Content

The indexer crawls a specified directory, processes the files, and builds the search index.

To run the indexer:

```sh
go run cmd/indexer/main.go --path=./data
```

This will index all `.txt` and `.md` files in the `data` directory.

### 2. Running the Search Server

The search server provides an API to query the indexed content.

To start the server:

```sh
go run cmd/server/main.go
```

The server will be available at `http://localhost:8080`.

### 3. Using the Search UI

Once the server is running, you can use the web interface to search:

1.  Open your browser and navigate to `http://localhost:8080/`.
2.  The project's `index.html` file will be served, allowing you to enter search queries.

## API Endpoints

### Search

-   **Endpoint:** `/api/v1/search`
-   **Method:** `GET`
-   **Query Parameters:**
    -   `q`: The search query.
-   **Example:**
    ```sh
    curl "http://localhost:8080/api/v1/search?q=hello+world"
    ```
-   **Response:**
    ```json
    [
        {
            "docID": "635f2a4b2f8e4b1e3e3e3e3e",
            "title": "Hello World",
            "url": "data/hello.txt"
        }
    ]
    ```

## Project Structure

```
.
├── cmd/
│   ├── indexer/  # Indexer application
│   └── server/   # Search server application
├── internal/
│   ├── analysis/ # Text analysis (tokenization, stemming)
│   ├── crawler/  # Web crawler (future feature)
│   ├── indexer/  # Indexing logic
│   ├── ranking/  # Search result ranking (TF-IDF)
│   ├── search/   # Search logic
│   ├── server/   # Web server implementation (handlers, router)
│   └── storage/  # MongoDB storage
├── pkg/
│   └── ...       # Shared packages
├── ui/
│   ├── index.html # Frontend search page
│   └── script.js  # Frontend JavaScript
├── Makefile      # Development commands
└── Dockerfile    # Docker configuration
```

## Development with Docker

The project includes a `docker-compose.yaml` file for a consistent development environment.

1.  **Start the environment:**
    ```sh
    docker-compose up --build
    ```
    This will build the Go application image, start the application and MongoDB containers, and mount the source code for live reloading.

2.  **Access services:**
    -   **API Server:** `http://localhost:8080`
    -   **MongoDB:** `mongodb://admin:password@localhost:27017`

3.  **Stop the environment:**
    ```sh
    docker-compose down
    ```

## Makefile Commands

The `Makefile` provides useful commands for common development tasks:

-   `make lint`: Run the linter to analyze the code.
-   `make test`: Run all project tests.
-   `make test-coverage`: Generate a test coverage report.
-   `make build-indexer`: Compile the indexer binary.
-   `make run-indexer`: Run the indexer on the `data/` directory.
-   `make watch`: Start the application in development mode with live reloading (`air`).
