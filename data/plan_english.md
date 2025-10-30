# Project Plan: Local Search Engine with Web Crawler in Go

**Version:** 1.0
**Date:** 2023-10-27
**Author:** Project Team

## 1. Definition and Project Vision

### 1.1. Vision
To build a high-performance search engine, written in Go, capable of indexing content from local files and websites. The system will offer relevance-based searches through a RESTful API and a minimalist web user interface. The project will be self-contained, easy to deploy with Docker, and designed to be extensible.

### 1.2. Key Objectives
- **Dual Indexing:** Support indexing of local directories (`.txt`, `.md`, `.html` files) and websites through an integrated crawler.
- **Relevance-Based Search:** Implement standard ranking algorithms like TF-IDF (initially) and BM25 (as an enhancement).
- **Robust Persistence:** Use MongoDB to store the inverted index, document metadata, and index statistics.
- **Concurrent Architecture:** Design a system that leverages Go's goroutines for efficient and parallel indexing and crawling.
- **Modern Interface:** Provide a REST API for programmatic interaction and a simple web UI for manual searching.
- **Professional Operability:** Ensure the project is easy to configure, deploy, monitor, and maintain in different environments (development, production).

### 1.3. MVP (Minimum Viable Product) Scope
- **Indexer:** Indexing of local files (`.txt`, `.md`).
- **Crawler:** Crawling of a single domain (without following external links) respecting `robots.txt`.
- **Text Analysis:** Tokenization, lowercasing, and basic stopword removal.
- **Inverted Index:** Storage of the index in MongoDB.
- **Ranking:** Implementation of TF-IDF.
- **API:** Endpoints for indexing a directory, starting a crawl, and searching (`/index`, `/crawl`, `/search`).
- **UI:** A single HTML page with a search field and a list of results (title and path/URL).
- **Deployment:** A `docker-compose.yml` file to run the application and the MongoDB database in a development environment.

### 1.4. Out of MVP Scope (Potential Future Enhancements)
- BM25 ranking.
- Stemming and lemmatization.
- Exact phrase searches and boolean operators (AND, OR, NOT).
- Term highlighting in results.
- Advanced pagination and facets.
- Indexing of complex formats (PDF, DOCX).
- Advanced UI with filters and previews.
- Admin panel for managing the index.
- Scaling to multiple nodes.

## 2. Architecture and Technology Stack

### 2.1. Component Flow Diagram
```
      +------------------+     +-----------------+      +-----------------+
      | Local File       | --> |                 |      |                 |
      | Ingestor         |     |                 |      |                 |
      +------------------+     |                 |      |   Indexer       | --> | MongoDB |
                               |   Analyzer      |      |   (Writer)      |     | (Inverted Index,
      +------------------+     |   (Tokenizer,   | -->  |                 |     |  Documents, Stats)
      | Web Crawler      | --> |    Stopwords)   |      |                 |     +-----------------+
      +------------------+     +-----------------+      +-----------------+              ^
                                                                                         |
+-------------+      +-----------------+      +---------------------+      +-----------------+
|   User      | -->  |   Web UI        | -->  |     API Server      | -->  |    Searcher     |
|             |      |   (HTML/JS)     |      |  (Go, net/http)     |      |    (Ranker)     |
+-------------+      +-----------------+      +---------------------+      +-----------------+
```

### 2.2. Technology Stack
- **Language:** Go (v1.20+)
- **Database:** MongoDB (v6.0+)
- **API Framework:** `net/http` (Go standard library, to minimize dependencies).
- **MongoDB Driver:** `go.mongodb.org/mongo-driver`
- **Crawler/HTML Parser:** `golang.org/x/net/html` and `github.com/temoto/robotstxt`
- **Containerization:** Docker & Docker Compose
- **CI/CD:** GitHub Actions
- **Observability:** `prometheus/client_golang` for metrics, `uber-go/zap` for structured logging.

### 2.3. Data Modeling in MongoDB

#### Collection: `documents`
Stores metadata for each indexed document.
```json
{
  "_id": "ObjectId('...')", // Unique document ID
  "source_type": "file" | "web",
  "path": "/path/to/document.txt", // Optional, if file
  "url": "https://example.com/page",  // Optional, if web
  "title": "Document Title",
  "length": 152, // Number of tokens in the body
  "indexed_at": "ISODate('...')"
}
```

#### Collection: `inverted_index`
The heart of the engine. Each document in this collection represents a term and contains its postings list.
**Consideration:** Documents in MongoDB have a 16MB limit. For very common terms, the postings list could exceed this limit. For this project, we assume this will not be an issue. An advanced solution would be to segment the postings list.
```json
{
  "_id": "golang", // The term (token) itself
  "doc_frequency": 42, // How many documents contain this term (df)
  "postings": [
    {
      "doc_id": "ObjectId('...')", // Reference to the `documents` collection
      "term_frequency": 5, // How many times the term appears in this doc (tf)
      "positions": [12, 45, 78, 101, 134] // Positions for snippets/phrases
    },
    // ... more postings
  ]
}
```

#### Collection: `index_stats`
A single document to store global index statistics.
```json
{
  "_id": "global_stats",
  "total_docs": 15000,
  "avg_doc_length": 250.5
}
```

## 3. Detailed Implementation Plan (Phases and Tasks)

### Phase 0: Foundation and Environment Setup (Duration: ~2 days)
- **Task 0.1:** Initialize the Git repository with `.gitignore` and `README.md`.
- **Task 0.2:** Define the project directory structure (`/cmd`, `/internal`, `/pkg`, `/api`, `/ui`).
- **Task 0.3:** Initialize the Go module (`go mod init`).
- **Task 0.4:** Create `docker-compose.yml` to run the Go service and a MongoDB instance.
- **Task 0.5:** Configure a linter (`golangci-lint`) and formatter.
- **Task 0.6:** Create a `Makefile` skeleton for common tasks (`build`, `run`, `test`, `lint`).
- **Task 0.7:** Set up a basic CI pipeline in GitHub Actions that runs `lint` and `test` on every push.

### Phase 1: Indexing and Analysis Core (Duration: ~1 week)
- **Task 1.1: Text Analysis Module (`/internal/analysis`)**
    - Sub-task 1.1.1: Implement a `Tokenizer` that converts text into a slice of tokens (words).
    - Sub-task 1.1.2: Add filters: lowercasing and punctuation removal.
    - Sub-task 1.1.3: Implement a `stopwords` filter (with a configurable list in English/Spanish).
    - Sub-task 1.1.4: Write unit tests for the analyzer.
- **Task 1.2: MongoDB Connection and Abstraction (`/internal/storage`)**
    - Sub-task 1.2.1: Implement the connection to MongoDB and its lifecycle management (connect/disconnect).
    - Sub-task 1.2.2: Define the Go structs that map to the MongoDB collections.
    - Sub-task 1.2.3: Create a wrapper with methods for `AddDocument`, `GetDocument`, `UpdateTermPostings`, etc.
- **Task 1.3: Local File Indexer (`/internal/indexer`)**
    - Sub-task 1.3.1: Implement logic to recursively traverse a directory and read files (`.txt`, `.md`).
    - Sub-task 1.3.2: Create the indexing pipeline: `Read file -> Analyze text -> Generate Postings -> Write to MongoDB`.
    - Sub-task 1.3.3: Use goroutines and a worker pool to parallelize file analysis.
    - Sub-task 1.3.4: Implement the write logic for MongoDB (create/update documents and terms). Use bulk operations for efficiency.
- **Task 1.4: CLI (`/cmd/indexer`)**
    - Sub-task 1.4.1: Create a simple CLI command (using `flag` or `cobra`) that takes a directory path and runs the indexing process.

### Phase 2: Search and API (Duration: ~1 week)
- **Task 2.1: Search Logic (`/internal/searcher`)**
    - Sub-task 2.1.1: Implement the logic to analyze a search query using the same `Analyzer` from Phase 1.
    - Sub-task 2.1.2: Develop the function to retrieve `postings lists` from MongoDB for each query term.
    - Sub-task 2.1.3: Implement the **TF-IDF** ranking algorithm.
        - `TF` is already in the postings.
        - `IDF` is calculated using `total_docs` (from `index_stats`) and `doc_frequency` (from `inverted_index`).
    - Sub-task 2.1.4: Implement score accumulation for documents and result sorting. Use a heap to efficiently get the top-K results.
- **Task 2.2: API Server (`/cmd/server` and `/internal/api`)**
    - Sub-task 2.2.1: Set up the HTTP server using the standard `net/http` package. This involves creating an `http.ServeMux` to register routes and starting the server with `http.ListenAndServe`.
    - Sub-task 2.2.2: Register the `GET /api/v1/search` route using `http.HandleFunc`. The handler will parse query parameters (e.g., `?q=...`).
    - Sub-task 2.2.3: Connect the route handler with the `searcher` logic implemented in Task 2.1.
    - Sub-task 2.2.4: Define the JSON response structs and use the `encoding/json` package to serialize the response.
    - Sub-task 2.2.5: Implement request logging and centralized error handling by creating a middleware wrapper for `http.Handler` to keep the code clean.
- **Task 2.3: Minimal UI (`/ui`)**
    - Sub-task 2.3.1: Create an `index.html` file with a search form.
    - Sub-task 2.3.2: Write a vanilla JavaScript `script.js` that calls the `/api/v1/search` endpoint and renders the results in a list.
    - Sub-task 2.3.3: Serve the static files using `http.FileServer` from the Go server.

### Phase 3: Web Crawler Integration (Duration: ~1 week)
- **Task 3.1: Crawler Component (`/internal/crawler`)**
    - Sub-task 3.1.1: Implement an HTTP fetcher to download the content of a URL.
    - Sub-task 3.1.2: Integrate a library to parse `robots.txt` and respect its rules.
    - Sub-task 3.1.3: Implement a "politeness" mechanism (delay between requests to the same domain).
    - Sub-task 3.1.4: Use `golang.org/x/net/html` to extract the title and body text from HTML pages.
    - Sub-task 3.1.5: Implement logic to extract and enqueue new links from the same domain for recursive crawling.
- **Task 3.2: Pipeline Integration**
    - Sub-task 3.2.1: Connect the crawler's output (extracted text) to the existing `Analyzer` and `Indexer`.
    - Sub-task 3.2.2: Ensure the `documents` model in MongoDB correctly handles URLs.
    - Sub-task 3.2.3: Use a map or a MongoDB structure to avoid visiting duplicate URLs (`visited_urls`).
- **Task 3.3: API for the Crawler**
    - Sub-task 3.3.1: Register the `POST /api/v1/crawl` route in the `net/http` server. The handler will read the request body to get the starting URL (`seed_url`).
    - Sub-task 3.3.2: Run the crawling process in a background goroutine to avoid blocking the API response.

### Phase 4: Enhancements and Refinement (Duration: ~1.5 weeks)
- **Task 4.1: BM25 Ranking**
    - Sub-task 4.1.1: Add `avg_doc_length` to the `index_stats` collection.
    - Sub-task 4.1.2: Implement the BM25 formula in the `searcher`.
    - Sub-task 4.1.3: Modify the search handler to allow selecting the ranking engine via a query param: `GET /search?q=...&ranker=bm25`.
- **Task 4.2: Snippets and Highlighting**
    - Sub-task 4.2.1: Implement a function that uses the `positions` from the postings to find the best text fragment containing the search terms.
    - Sub-task 4.2.2: Add a `snippet` field to the API response.
    - Sub-task 4.2.3: Update the UI to display the snippet.
- **Task 4.3: Index Management**
    - Sub-task 4.3.1: Implement the logic to delete a document from the index.
    - Sub-task 4.3.2: Create the `DELETE /api/v1/documents/{id}` endpoint. Extracting the ID from the path (e.g., `/api/v1/documents/12345`) will require manual parsing of `r.URL.Path` in the handler.
    - Sub-task 4.3.3: Implement a `GET /api/v1/stats` endpoint to view index statistics.

### Phase 5: Production and Operability (Duration: ~1 week)
- **Task 5.1: Configuration**
    - Sub-task 5.1.1: Externalize configuration (server port, MongoDB connection string, etc.) using environment variables or a configuration file (`viper`).
- **Task 5.2: Observability**
    - Sub-task 5.2.1: Integrate structured logging (e.g., `zap`) throughout the application.
    - Sub-task 5.2.2: Expose metrics in Prometheus format on a `/metrics` endpoint (e.g., search latency, number of indexed documents).
- **Task 5.3: Docker for Production**
    - Sub-task 5.3.1: Create a multi-stage `Dockerfile` to build a static, lightweight Go binary in a minimal final image (e.g., `scratch` or `alpine`).
- **Task 5.4: Documentation**
    - Sub-task 5.4.1: Improve the `README.md` with instructions for installation, configuration, and API usage.
    - Sub-task 5.4.2: Add code comments (`godoc`).
- **Task 5.5: CI/CD Refinement**
    - Sub-task 5.5.1: Extend the GitHub Actions pipeline to build and push the Docker image to a registry (e.g., Docker Hub, GHCR) on Git tags.

## 4. Testing Strategy

- **Unit Tests:** Each package (`analysis`, `storage`, `ranker`) must have solid unit test coverage. Pure functions and edge cases will be tested.
- **Integration Tests:** Tests will be created that spin up a Dockerized MongoDB instance to test the full indexing and search flow on a small corpus of documents.
- **End-to-End (E2E) Tests:** (Optional) A script or framework can be used to launch the application with `docker-compose` and make real HTTP requests to the API, validating the responses.

## 5. Risks and Mitigations

- **Risk 1:** MongoDB's 16MB document limit for postings lists of very frequent terms.
  - **Mitigation:** For the MVP, document this limitation. In the future, implement postings list segmentation (e.g., `term_golang_1`, `term_golang_2`).
- **Risk 2:** The web crawler may be blocked or considered abusive.
  - **Mitigation:** Implement a descriptive `User-Agent`, strictly respect `robots.txt`, and ensure the "politeness" `delay` is conservative by default.
- **Risk 3:** Database performance degrades with a large index.
  - **Mitigation:** Ensure proper indexes are created on MongoDB collections. Use `bulk writes` for indexing. Optimize search queries.
- **Risk 4:** Complexity in concurrency management may introduce race conditions.
  - **Mitigation:** Use safe Go concurrency patterns (channels, mutexes where necessary) and conduct concurrency testing. A single-writer model for the index can simplify this initially.

## 6. Deployment Strategy

The deployment strategy is divided into three environments to ensure a robust and secure software lifecycle.

### 6.1. Development Environment (Local)
- **Objective:** Allow developers to work quickly and in isolation.
- **Setup:** The `docker-compose.yml` file defined in Phase 0 will be used.
- **Components:**
  - A service for the Go application, with the source code mounted as a volume to enable hot-reloading (using tools like `air` or `CompileDaemon`).
  - A service for the MongoDB database, with its data persisted in a Docker volume to avoid losing the index between restarts.
- **Workflow:** The developer clones the repository, runs `docker-compose up`, and can start coding and testing changes in real-time.

### 6.2. Staging Environment (Testing)
- **Objective:** An environment identical to production for integration testing, E2E, and validation of new features before release.
- **Setup:**
  - The Go application is deployed as a Docker image built by the CI/CD pipeline. The source code is not mounted.
  - The database will be a separate MongoDB instance (it could be another container or a free instance on a service like MongoDB Atlas).
  - Configuration (connection strings, log levels, etc.) is injected via environment variables.
- **Workflow:** Whenever a new feature is merged into the `develop` or `main` branch, the CI/CD pipeline builds the image and automatically deploys it to this environment. User Acceptance Testing (UAT) would be performed here.

### 6.3. Production Environment
- **Objective:** Serve the application to end-users reliably, scalably, and securely.
- **Setup:**
  - **Go Application:** Deployed as a Docker container (using the optimized production image) in an orchestrator like **Kubernetes (K8s)**, **AWS ECS**, or a PaaS platform like **Google Cloud Run** or **Heroku**. A minimum of 2 replicas is recommended for high availability.
  - **MongoDB Database:** A managed service like **MongoDB Atlas** will be used. This delegates the responsibility for backups, scaling, security, and maintenance to a specialized provider.
  - **Networking:** The service will be behind a load balancer. API access can be protected with an API Gateway if needed.
  - **Configuration:** All configurations and secrets (API keys, connection strings) will be managed securely through the orchestrator's mechanisms (e.g., K8s Secrets) or a secrets management service.

## 7. Maintenance and Operations

A plan to ensure the system continues to operate optimally after deployment.

### 7.1. Monitoring and Alerting
- **Key Metrics to Monitor (via Prometheus):**
  - **API Latency:** p95 and p99 latency of the `/search` endpoint.
  - **Error Rate:** Percentage of 5xx responses from the API.
  - **Indexing Throughput:** Documents indexed per second.
  - **Crawler Status:** Number of URLs in the queue, HTTP error rate.
  - **Resource Usage:** CPU and memory of the application container.
  - **Database Health:** Query response times, active connections (provided by MongoDB Atlas).
- **Alerts (Configured in Alertmanager/Grafana):**
  - Alert if p99 search latency exceeds 500ms.
  - Alert if the error rate exceeds 1% for 5 minutes.
  - Alert if the indexing or crawling process fails repeatedly.
  - Alert if CPU or memory usage reaches 85% of the limit.

### 7.2. Backups and Disaster Recovery
- **Database:** By using MongoDB Atlas, automatic backups and point-in-time recovery (PITR) will be configured. This greatly simplifies disaster recovery.
- **Application:** The application is "stateless," so recovery simply involves re-deploying the functional Docker image from the container registry.

### 7.3. Update Process
- **Application:** Updates will be performed using "rolling updates" in the container orchestrator to ensure zero downtime.
- **Index:** In case of changes requiring a full re-index (e.g., a new text analyzer), a "blue/green index" strategy will be followed:
  1. A new set of collections is created in MongoDB (e.g., `documents_v2`, `inverted_index_v2`).
  2. A re-indexing job is launched to populate these new collections.
  3. Once complete, the application is reconfigured (via an environment variable) to point to the new collections.
  4. The old collections are archived and then deleted.

## 8. Summary Schedule and Milestones

| Phase                                   | Estimated Duration | Key Deliverable                                            |
| --------------------------------------- | ------------------ | ---------------------------------------------------------- |
| **Phase 0: Foundation and Setup**       | 2 days             | Repository, CI, and functional local development environment. |
| **Phase 1: Indexing and Analysis Core** | 1 week             | CLI to index local files into MongoDB.                     |
| **Phase 2: Search and API**             | 1 week             | Functional `/search` API with TF-IDF ranking and basic web UI. |
| **Phase 3: Web Crawler Integration**    | 1 week             | `/crawl` API to index a website.                           |
| **Phase 4: Enhancements & Refinement**  | 1.5 weeks          | BM25 ranking, snippets, and index management endpoints.    |
| **Phase 5: Production & Operability**   | 1 week             | Production Dockerfile, metrics, logging, and documentation. |
| **Total Estimated**                     | **~ 5-6 weeks**    |                                                            |