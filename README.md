# gofetch: High-Performance Local and Web Search Engine

![gofetch architecture](https://user-images.githubusercontent.com/12345678/123456789-abcdef.png)  <!-- Placeholder for a future architecture diagram -->

`gofetch` is a powerful, high-performance search engine written in Go, designed to index and search content from local file systems and (in the future) the web. It provides a clean RESTful API for querying and a minimalist web interface for a seamless user experience.

This project is perfect for developers who want to build their own search solutions, learn about information retrieval concepts, or integrate full-text search capabilities into their Go applications.

## Project Vision

The goal of `gofetch` is to be a lightweight yet powerful search engine that is:

- **Fast and Efficient:** Built with Go's concurrency model to handle indexing and searching with high throughput.
- **Easy to Use:** Simple to configure, deploy, and integrate via a clean API.
- **Extensible:** Designed with a modular architecture to allow for future expansion, such as adding new data sources or ranking algorithms.
- **Educational:** Provides a clear and well-documented codebase for those interested in learning about search engine internals.

## Features at a Glance

- **Dual Indexing:** Supports indexing of local directories (text and markdown files). Web crawling capabilities are planned for a future release.
- **Relevance-Based Ranking:** Implements the TF-IDF (Term Frequency-Inverse Document Frequency) algorithm to deliver relevance-ranked search results.
- **Advanced Text Analysis:** Features a sophisticated text analysis pipeline including:
    - **Tokenization:** Breaks down text into individual words or terms.
    - **Normalization:** Converts text to a consistent case (lowercase).
    - **Stop Word Filtering:** Removes common words to improve index quality.
    - **Stemming:** Reduces words to their root form using Snowball stemmers.
- **Multi-Language Support:** Out-of-the-box support for **English** and **Spanish**.
- **Robust Persistence:** Utilizes MongoDB for scalable and reliable storage of the search index and document metadata.
- **Concurrent by Design:** Leverages Go's goroutines to perform indexing and searching operations concurrently, maximizing performance.
- **Containerized:** Comes with `Dockerfile` and `docker-compose.yaml` for easy, reproducible deployments.

## Getting Started

### Prerequisites

- Go (1.18 or higher)
- MongoDB
- Docker & Docker Compose (Recommended)

### 1. Clone the Repository

```sh
git clone https://github.com/TonyGLL/gofetch
cd gofetch
```

### 2. Configuration

`gofetch` uses a `config.yaml` file for configuration, supplemented by environment variables.

**A. Create a `config.yaml` file:**

You can copy the example file:

```sh
cp config.yaml.example config.yaml
```

**B. Configure your environment:**

The following variables can be set in your `config.yaml` or as environment variables:

| Variable            | Description                                | Default                      |
| ------------------- | ------------------------------------------ | ---------------------------- |
| `MONGODB_URI`       | MongoDB connection string.                 | `mongodb://localhost:27017`  |
| `DB_NAME`           | The name of the database.                  | `gofetch`                    |
| `ANALYZER_LANGUAGE` | Language for text analysis (`english` or `spanish`). | `english`                    |
| `INDEXER_PATH`      | The directory path to index.               | `./data`                     |
| `SERVER_PORT`       | The port for the API server.               | `8080`                       |

### 3. Build and Run with Docker (Recommended)

The simplest way to get `gofetch` running is with Docker Compose.

```sh
docker-compose up --build
```

This command will:
1.  Build the `gofetch` Docker image.
2.  Start the search server and a MongoDB container.
3.  Mount the source code for live-reloading on changes.

The API will be available at `http://localhost:8080`.

### 4. Build and Run Manually

**A. Install Dependencies:**

```sh
go mod download
```

**B. Run the Indexer:**

To populate the search index, run the indexer and point it to the directory you want to scan.

```sh
go run cmd/indexer/main.go --path=./data
```

**C. Run the API Server:**

Once indexing is complete, start the server.

```sh
go run cmd/server/main.go
```

The server will be available at `http://localhost:8080`.

## Usage

### Search UI

Navigate to `http://localhost:8080` in your browser to use the simple web interface for searching.

### API Endpoints

#### Search for Documents

-   **Endpoint:** `/api/v1/search`
-   **Method:** `GET`
-   **Query Parameters:**
    -   `q` (string, required): The search query.
-   **Example Request:**

    ```sh
    curl "http://localhost:8080/api/v1/search?q=hello+world"
    ```

-   **Example Success Response (`200 OK`):**

    ```json
    [
        {
            "Score": 0.6931471805599453,
            "Document": {
                "ID": "635f2a4b2f8e4b1e3e3e3e3e",
                "Title": "Hello World",
                "URL": "data/hello.txt",
                "IndexedAt": "2023-10-27T10:00:00Z",
                "Content": "This is the content of the hello world file."
            }
        }
    ]
    ```

## Project Structure

The project follows a standard Go layout to maintain a clean and scalable architecture.

```
.
├── cmd/                # Application entry points
│   ├── indexer/        # Main package for the indexer binary
│   └── server/         # Main package for the API server binary
├── internal/           # Private application logic
│   ├── analysis/       # Text analysis (tokenization, stemming, etc.)
│   ├── builder/        # Dependency injection builders
│   ├── config/         # Configuration management (Viper)
│   ├── indexer/        # Core indexing logic and pipeline
│   ├── ranking/        # Search result ranking algorithms (TF-IDF)
│   ├── search/         # Core search logic
│   ├── server/         # Web server, handlers, and routing
│   └── storage/        # MongoDB interaction and data models
├── pkg/                # Public libraries (currently none)
├── ui/                 # Frontend static files (HTML, JS)
├── .github/            # GitHub Actions workflows
├── Makefile            # Development and automation commands
└── Dockerfile          # Docker build configuration
```

## Development

### Makefile Commands

The `Makefile` contains helpful commands for development:

-   `make lint`: Run the Go linter to check for code style and errors.
-   `make test`: Run all unit and integration tests.
--   `make test-coverage`: Generate a test coverage report.
-   `make build-indexer`: Compile the indexer binary.
-   `make run-indexer`: Run the indexer on the default `data/` directory.
-   `make watch`: Start the API server in development mode with live reloading (`air`).

### Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/my-new-feature`).
3.  Commit your changes (`git commit -am 'Add some feature'`).
4.  Push to the branch (`git push origin feature/my-new-feature`).
5.  Open a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
