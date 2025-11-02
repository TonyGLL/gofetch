# gofetch: Motor de Búsqueda Local y Web

`gofetch` es un motor de búsqueda de alto rendimiento escrito en Go, diseñado para indexar contenido tanto de archivos locales como de sitios web. El sistema ofrece búsquedas por relevancia a través de una API RESTful y una futura interfaz de usuario web minimalista.

## Objetivos Clave

- **Indexación Dual:** Soporte para indexar directorios locales (`.txt`, `.md`) y sitios web.
- **Búsqueda por Relevancia:** Implementación de algoritmos de ranking como TF-IDF y BM25.
- **Persistencia Robusta:** Uso de MongoDB para almacenar el índice invertido y los metadatos.
- **Arquitectura Concurrente:** Diseño eficiente que aprovecha las goroutines de Go.
- **Operatividad Profesional:** Fácil de configurar, desplegar y monitorear usando Docker y GitHub Actions.

## Entorno de Desarrollo

El proyecto utiliza Docker y Docker Compose para crear un entorno de desarrollo consistente y fácil de gestionar.

### Prerrequisitos

- Docker
- Docker Compose

### Puesta en Marcha

1.  **Clonar el repositorio:**
    ```sh
    git clone https://github.com/TonyGLL/gofetch
    cd gofetch
    ```

2.  **Levantar el entorno:**
    El siguiente comando construirá la imagen de la aplicación Go, iniciará los contenedores de la aplicación y de MongoDB, y montará el código fuente para habilitar la recarga en vivo (`hot-reloading`).

    ```sh
    docker-compose up --build
    ```

3.  **Acceder a los servicios:**
    -   **API Server (próximamente):** `http://localhost:8080`
    -   **MongoDB:** `mongodb://admin:password@localhost:27017`

4.  **Detener el entorno:**
    ```sh
    docker-compose down
    ```

## Uso del Makefile

El `Makefile` proporciona comandos útiles para las tareas comunes de desarrollo:

-   `make lint`: Ejecuta el linter para analizar el código.
-   `make test`: Ejecuta todas las pruebas del proyecto.
-   `make test-coverage`: Genera un informe de cobertura de pruebas.
-   `make build-indexer`: Compila el binario del indexador.
-   `make run-indexer`: Ejecuta el indexador sobre el directorio `data/`.
-   `make watch`: Inicia la aplicación en modo de desarrollo con recarga en vivo (`air`).