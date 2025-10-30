# Plan de Proyecto: Motor de Búsqueda Local con Web Crawler en Go

**Versión:** 1.0
**Fecha:** 2023-10-27
**Autor:** Project Team

## 1. Definición y Visión del Proyecto

### 1.1. Visión
Construir un motor de búsqueda de alto rendimiento, escrito en Go, capaz de indexar contenido de archivos locales y sitios web. El sistema ofrecerá búsquedas por relevancia a través de una API RESTful y una interfaz de usuario web minimalista. El proyecto será autocontenido, fácil de desplegar con Docker y diseñado para ser extensible.

### 1.2. Objetivos Clave
- **Indexación Dual:** Soportar la indexación de directorios locales (archivos `.txt`, `.md`, `.html`) y de sitios web a través de un crawler integrado.
- **Búsqueda por Relevancia:** Implementar algoritmos de ranking estándar como TF-IDF (inicialmente) y BM25 (como mejora).
- **Persistencia Robusta:** Utilizar MongoDB para almacenar el índice invertido, metadatos de documentos y estadísticas del índice.
- **Arquitectura Concurrente:** Diseñar un sistema que aproveche las goroutines de Go para una indexación y crawling eficientes y paralelos.
- **Interfaz Moderna:** Proveer una API REST para la interacción programática y una UI web simple para la búsqueda manual.
- **Operatividad Profesional:** Asegurar que el proyecto sea fácil de configurar, desplegar, monitorear y mantener en diferentes entornos (desarrollo, producción).

### 1.3. Alcance del MVP (Producto Mínimo Viable)
- **Indexer:** Indexación de archivos locales (`.txt`, `.md`).
- **Crawler:** Crawling de un único dominio (sin seguir enlaces externos) respetando `robots.txt`.
- **Análisis de Texto:** Tokenización, conversión a minúsculas y eliminación de stopwords básicas.
- **Índice Invertido:** Almacenamiento del índice en MongoDB.
- **Ranking:** Implementación de TF-IDF.
- **API:** Endpoints para indexar un directorio, iniciar un crawl y buscar (`/index`, `/crawl`, `/search`).
- **UI:** Una única página HTML con un campo de búsqueda y una lista de resultados (título y path/URL).
- **Deployment:** Un `docker-compose.yml` para levantar la aplicación y la base de datos MongoDB en un entorno de desarrollo.

### 1.4. Fuera del Alcance del MVP (Posibles Mejoras Futuras)
- Ranking BM25.
- Stemming y lematización.
- Búsquedas por frases exactas y operadores booleanos (AND, OR, NOT).
- Resaltado de términos (highlighting) en los resultados.
- Paginación avanzada y facetas.
- Indexación de formatos complejos (PDF, DOCX).
- UI avanzada con filtros y previsualizaciones.
- Panel de administración para gestionar el índice.
- Escalado a múltiples nodos.

## 2. Arquitectura y Stack Tecnológico

### 2.1. Diagrama de Flujo de Componentes
```
      +------------------+     +-----------------+      +-----------------+
      | Local File       | --> |                 |      |                 |
      | Ingestor         |     |                 |      |                 |
      +------------------+     |                 |      |   Indexer       | --> | MongoDB |
                               |   Analyzer      |      |   (Writer)      |     | (Índice Invertido,
      +------------------+     |   (Tokenizer,   | -->  |                 |     |  Documentos, Stats)
      | Web Crawler      | --> |    Stopwords)   |      |                 |     +-----------------+
      +------------------+     +-----------------+      +-----------------+              ^
                                                                                         |
+-------------+      +-----------------+      +---------------------+      +-----------------+
|   Usuario   | -->  |   UI Web        | -->  |     API Server      | -->  |    Searcher     |
|             |      |   (HTML/JS)     |      |  (Go, net/http)     |      |    (Ranker)     |
+-------------+      +-----------------+      +---------------------+      +-----------------+
```

### 2.2. Stack Tecnológico
- **Lenguaje:** Go (v1.20+)
- **Base de Datos:** MongoDB (v6.0+)
- **API Framework:** `net/http` (biblioteca estándar de Go, para minimizar dependencias).
- **Driver de MongoDB:** `go.mongodb.org/mongo-driver`
- **Crawler/Parser HTML:** `golang.org/x/net/html` y `github.com/temoto/robotstxt`
- **Contenerización:** Docker & Docker Compose
- **CI/CD:** GitHub Actions
- **Observabilidad:** `prometheus/client_golang` para métricas, `uber-go/zap` para logging estructurado.

### 2.3. Modelado de Datos en MongoDB

#### Colección: `documents`
Almacena metadatos para cada documento indexado.
```json
{
  "_id": "ObjectId('...')", // ID único del documento
  "source_type": "file" | "web",
  "path": "/path/to/document.txt", // Opcional, si es archivo
  "url": "https://example.com/page",  // Opcional, si es web
  "title": "Título del Documento",
  "length": 152, // Número de tokens en el cuerpo
  "indexed_at": "ISODate('...')"
}
```

#### Colección: `inverted_index`
El corazón del motor. Cada documento en esta colección representa un término y contiene su lista de postings.
**Consideración:** Los documentos en MongoDB tienen un límite de 16MB. Para términos muy comunes, la lista de postings podría exceder este límite. Para este proyecto, asumimos que no será un problema. Una solución avanzada sería segmentar la lista de postings.
```json
{
  "_id": "golang", // El término (token) en sí
  "doc_frequency": 42, // Cuántos documentos contienen este término (df)
  "postings": [
    {
      "doc_id": "ObjectId('...')", // Referencia a la colección `documents`
      "term_frequency": 5, // Cuántas veces aparece el término en este doc (tf)
      "positions": [12, 45, 78, 101, 134] // Posiciones para snippets/frases
    },
    // ... más postings
  ]
}
```

#### Colección: `index_stats`
Un único documento para almacenar estadísticas globales del índice.
```json
{
  "_id": "global_stats",
  "total_docs": 15000,
  "avg_doc_length": 250.5
}
```

## 3. Plan Detallado de Implementación (Fases y Tareas)

### Fase 0: Fundación y Configuración del Entorno (Duración: ~2 días)
- **Tarea 0.1:** Inicializar el repositorio Git con `.gitignore` y `README.md`.
- **Tarea 0.2:** Definir la estructura de directorios del proyecto (`/cmd`, `/internal`, `/pkg`, `/api`, `/ui`).
- **Tarea 0.3:** Inicializar el módulo de Go (`go mod init`).
- **Tarea 0.4:** Crear `docker-compose.yml` para levantar el servicio de Go y una instancia de MongoDB.
- **Tarea 0.5:** Configurar un linter (`golangci-lint`) y formateador.
- **Tarea 0.6:** Crear un esqueleto de `Makefile` para tareas comunes (`build`, `run`, `test`, `lint`).
- **Tarea 0.7:** Configurar un pipeline de CI básico en GitHub Actions que ejecute `lint` y `test` en cada push.

### Fase 1: Core de Indexación y Análisis (Duración: ~1 semana)
- **Tarea 1.1: Módulo de Análisis de Texto (`/internal/analysis`)**
    - Sub-tarea 1.1.1: Implementar un `Tokenizer` que convierta texto en un slice de tokens (palabras).
    - Sub-tarea 1.1.2: Añadir filtros: conversión a minúsculas y eliminación de puntuación.
    - Sub-tarea 1.1.3: Implementar un filtro de `stopwords` (con una lista configurable en inglés/español).
    - Sub-tarea 1.1.4: Escribir tests unitarios para el analizador.
- **Tarea 1.2: Conexión y Abstracción de MongoDB (`/internal/storage`)**
    - Sub-tarea 1.2.1: Implementar la conexión a MongoDB y la gestión del ciclo de vida (conectar/desconectar).
    - Sub-tarea 1.2.2: Definir las estructuras de Go que mapean a las colecciones de MongoDB.
    - Sub-tarea 1.2.3: Crear un wrapper con métodos para `AddDocument`, `GetDocument`, `UpdateTermPostings`, etc.
- **Tarea 1.3: Indexador de Archivos Locales (`/internal/indexer`)**
    - Sub-tarea 1.3.1: Implementar la lógica para recorrer un directorio recursivamente y leer archivos (`.txt`, `.md`).
    - Sub-tarea 1.3.2: Crear el pipeline de indexación: `Leer archivo -> Analizar texto -> Generar Postings -> Escribir en MongoDB`.
    - Sub-tarea 1.3.3: Usar goroutines y un pool de workers para paralelizar el análisis de archivos.
    - Sub-tarea 1.3.4: Implementar la lógica de escritura en MongoDB (crear/actualizar documentos y términos). Usar operaciones bulk para eficiencia.
- **Tarea 1.4: CLI (`/cmd/indexer`)**
    - Sub-tarea 1.4.1: Crear un comando de CLI simple (usando `flag` o `cobra`) que tome una ruta de directorio y ejecute el proceso de indexación.

### Fase 2: Búsqueda y API (Duración: ~1 semana)
- **Tarea 2.1: Lógica de Búsqueda (`/internal/searcher`)**
    - Sub-tarea 2.1.1: Implementar la lógica para analizar una query de búsqueda usando el mismo `Analyzer` de la Fase 1.
    - Sub-tarea 2.1.2: Desarrollar la función para recuperar las `postings lists` de MongoDB para cada término de la query.
    - Sub-tarea 2.1.3: Implementar el algoritmo de ranking **TF-IDF**.
        - `TF` ya está en los postings.
        - `IDF` se calcula con `total_docs` (de `index_stats`) y `doc_frequency` (de `inverted_index`).
    - Sub-tarea 2.1.4: Implementar la acumulación de puntajes para los documentos y la ordenación de resultados. Usar un heap para obtener el top-K de manera eficiente.
- **Tarea 2.2: Servidor API (`/cmd/server` y `/internal/api`)**
    - Sub-tarea 2.2.1: Configurar el servidor HTTP utilizando el paquete `net/http` estándar. Esto implica crear un `http.ServeMux` para registrar las rutas y lanzar el servidor con `http.ListenAndServe`.
    - Sub-tarea 2.2.2: Registrar la ruta `GET /api/v1/search` usando `http.HandleFunc`. El handler se encargará de parsear los query parameters (ej. `?q=...`).
    - Sub-tarea 2.2.3: Conectar el handler de la ruta con la lógica del `searcher` implementada en la Tarea 2.1.
    - Sub-tarea 2.2.4: Definir las estructuras de respuesta JSON y utilizar el paquete `encoding/json` para serializar la respuesta.
    - Sub-tarea 2.2.5: Implementar logging de peticiones y manejo de errores centralizado, creando un wrapper de middleware para los `http.Handler` para mantener el código limpio.
- **Tarea 2.3: UI Mínima (`/ui`)**
    - Sub-tarea 2.3.1: Crear un archivo `index.html` con un formulario de búsqueda.
    - Sub-tarea 2.3.2: Escribir un `script.js` de JavaScript vanilla que llame al endpoint `/api/v1/search` y renderice los resultados en una lista.
    - Sub-tarea 2.3.3: Servir los archivos estáticos usando `http.FileServer` desde el servidor de Go.

### Fase 3: Integración del Web Crawler (Duración: ~1 semana)
- **Tarea 3.1: Componente Crawler (`/internal/crawler`)**
    - Sub-tarea 3.1.1: Implementar un fetcher HTTP para descargar el contenido de una URL.
    - Sub-tarea 3.1.2: Integrar una librería para parsear `robots.txt` y respetar sus reglas.
    - Sub-tarea 3.1.3: Implementar un mecanismo de "politeness" (delay entre peticiones al mismo dominio).
    - Sub-tarea 3.1.4: Usar `golang.org/x/net/html` para extraer el título y el texto del cuerpo de las páginas HTML.
    - Sub-tarea 3.1.5: Implementar la lógica para extraer y encolar nuevos enlaces del mismo dominio para el crawling recursivo.
- **Tarea 3.2: Integración del Pipeline**
    - Sub-tarea 3.2.1: Conectar la salida del crawler (texto extraído) al `Analyzer` e `Indexer` existentes.
    - Sub-tarea 3.2.2: Asegurarse de que el modelo de `documents` en MongoDB maneje correctamente las URLs.
    - Sub-tarea 3.2.3: Utilizar un mapa o una estructura en MongoDB para evitar visitar URLs duplicadas (`visited_urls`).
- **Tarea 3.3: API para el Crawler**
    - Sub-tarea 3.3.1: Registrar la ruta `POST /api/v1/crawl` en el servidor `net/http`. El handler leerá el cuerpo de la petición para obtener la URL de inicio (`seed_url`).
    - Sub-tarea 3.3.2: Ejecutar el crawling en una goroutine en segundo plano para no bloquear la respuesta de la API.

### Fase 4: Mejoras y Refinamiento (Duración: ~1.5 semanas)
- **Tarea 4.1: Ranking BM25**
    - Sub-tarea 4.1.1: Añadir `avg_doc_length` a la colección `index_stats`.
    - Sub-tarea 4.1.2: Implementar la fórmula de BM25 en el `searcher`.
    - Sub-tarea 4.1.3: Modificar el handler de búsqueda para permitir la selección del motor de ranking a través de un query param: `GET /search?q=...&ranker=bm25`.
- **Tarea 4.2: Snippets y Highlighting**
    - Sub-tarea 4.2.1: Implementar una función que use las `positions` de los postings para encontrar el mejor fragmento de texto que contenga los términos de la búsqueda.
    - Sub-tarea 4.2.2: Añadir el campo `snippet` a la respuesta de la API.
    - Sub-tarea 4.2.3: Actualizar la UI para mostrar el snippet.
- **Tarea 4.3: Gestión del Índice**
    - Sub-tarea 4.3.1: Implementar la lógica para eliminar un documento del índice.
    - Sub-tarea 4.3.2: Crear el endpoint `DELETE /api/v1/documents/{id}`. La extracción del ID desde la ruta (ej. `/api/v1/documents/12345`) requerirá un parseo manual del `r.URL.Path` en el handler.
    - Sub-tarea 4.3.3: Implementar un endpoint `GET /api/v1/stats` para ver las estadísticas del índice.

### Fase 5: Producción y Operatividad (Duración: ~1 semana)
- **Tarea 5.1: Configuración**
    - Sub-tarea 5.1.1: Externalizar la configuración (puerto del servidor, connection string de MongoDB, etc.) usando variables de entorno o un archivo de configuración (`viper`).
- **Tarea 5.2: Observabilidad**
    - Sub-tarea 5.2.1: Integrar logging estructurado (ej. `zap`) en toda la aplicación.
    - Sub-tarea 5.2.2: Exponer métricas en formato Prometheus en un endpoint `/metrics` (ej. latencia de búsqueda, número de documentos indexados).
- **Tarea 5.3: Docker para Producción**
    - Sub-tarea 5.3.1: Crear un `Dockerfile` multi-etapa para construir un binario estático y ligero de Go en una imagen final mínima (ej. `scratch` o `alpine`).
- **Tarea 5.4: Documentación**
    - Sub-tarea 5.4.1: Mejorar el `README.md` con instrucciones de instalación, configuración y uso de la API.
    - Sub-tarea 5.4.2: Añadir comentarios de código (`godoc`).
- **Tarea 5.5: Refinamiento de CI/CD**
    - Sub-tarea 5.5.1: Extender el pipeline de GitHub Actions para construir y pushear la imagen de Docker a un registro (ej. Docker Hub, GHCR) en los tags de Git.

## 4. Estrategia de Pruebas

- **Tests Unitarios:** Cada paquete (`analysis`, `storage`, `ranker`) debe tener una cobertura de tests unitarios sólida. Se probarán funciones puras y casos borde.
- **Tests de Integración:** Se crearán tests que levanten una instancia de MongoDB en Docker para probar el flujo completo de indexación y búsqueda en un pequeño corpus de documentos.
- **Tests End-to-End (E2E):** (Opcional) Se puede usar un script o framework para lanzar la aplicación con `docker-compose` y hacer peticiones HTTP reales a la API, validando las respuestas.

## 5. Riesgos y Mitigaciones

- **Riesgo 1:** El límite de 16MB de los documentos en MongoDB para las listas de postings de términos muy frecuentes.
  - **Mitigación:** Para el MVP, documentar esta limitación. A futuro, implementar segmentación de postings (ej. `term_golang_1`, `term_golang_2`).
- **Riesgo 2:** El web crawler puede ser bloqueado o considerado abusivo.
  - **Mitigación:** Implementar `User-Agent` descriptivo, respetar `robots.txt` estrictamente y asegurar que el `delay` de "politeness" sea conservador por defecto.
- **Riesgo 3:** El rendimiento de la base de datos se degrada con un índice grande.
  - **Mitigación:** Asegurar la creación de índices correctos en las colecciones de MongoDB. Usar `bulk writes` para la indexación. Optimizar las queries de búsqueda.
- **Riesgo 4:** La complejidad en la gestión de la concurrencia puede introducir race conditions.
  - **Mitigación:** Usar patrones de concurrencia seguros de Go (canales, mutex donde sea necesario) y realizar pruebas de concurrencia. El modelo de un único "writer" al índice puede simplificar esto inicialmente.
## 6. Estrategia de Despliegue (Deployment Strategy)

La estrategia de despliegue se divide en tres entornos para garantizar un ciclo de vida de software robusto y seguro.

### 6.1. Entorno de Desarrollo (Local)
- **Objetivo:** Permitir a los desarrolladores trabajar de forma rápida y aislada.
- **Configuración:** Se utilizará el archivo `docker-compose.yml` definido en la Fase 0.
- **Componentes:**
  - Un servicio para la aplicación Go, con el código fuente montado como un volumen para permitir hot-reloading (usando herramientas como `air` o `CompileDaemon`).
  - Un servicio para la base de datos MongoDB, con sus datos persistidos en un volumen de Docker para no perder el índice entre reinicios.
- **Flujo de trabajo:** El desarrollador clona el repositorio, ejecuta `docker-compose up`, y puede empezar a codificar y probar los cambios en tiempo real.

### 6.2. Entorno de Staging (Pruebas)
- **Objetivo:** Un entorno idéntico a producción para realizar pruebas de integración, E2E, y validación de nuevas características antes del lanzamiento.
- **Configuración:**
  - La aplicación Go se despliega como una imagen de Docker construida por el pipeline de CI/CD. No se monta el código fuente.
  - La base de datos será una instancia de MongoDB separada (puede ser otro contenedor o una instancia gratuita en un servicio como MongoDB Atlas).
  - La configuración (connection strings, niveles de log, etc.) se inyecta a través de variables de entorno.
- **Flujo de trabajo:** Cada vez que una nueva funcionalidad se mezcla en la rama `develop` o `main`, el pipeline de CI/CD construye la imagen y la despliega automáticamente en este entorno. Aquí se realizarían las pruebas de aceptación (UAT).

### 6.3. Entorno de Producción
- **Objetivo:** Servir la aplicación a los usuarios finales de forma fiable, escalable y segura.
- **Configuración:**
  - **Aplicación Go:** Desplegada como un contenedor Docker (usando la imagen optimizada de producción) en un orquestador como **Kubernetes (K8s)**, **AWS ECS**, o una plataforma PaaS como **Google Cloud Run** o **Heroku**. Se recomienda un mínimo de 2 réplicas para alta disponibilidad.
  - **Base de Datos MongoDB:** Se utilizará un servicio gestionado como **MongoDB Atlas**. Esto delega la responsabilidad de las copias de seguridad, el escalado, la seguridad y el mantenimiento a un proveedor especializado.
  - **Networking:** El servicio estará detrás de un balanceador de carga. El acceso a la API se puede proteger con un API Gateway si es necesario.
  - **Configuración:** Todas las configuraciones y secretos (API keys, connection strings) se gestionarán de forma segura a través de los mecanismos del orquestador (ej. K8s Secrets) o un servicio de gestión de secretos.

## 7. Mantenimiento y Operaciones

Un plan para asegurar que el sistema siga funcionando de manera óptima después del despliegue.

### 7.1. Monitorización y Alertas
- **Métricas Clave a Monitorear (vía Prometheus):**
  - **Latencia de la API:** Latencia p95 y p99 del endpoint `/search`.
  - **Tasa de Errores:** Porcentaje de respuestas 5xx en la API.
  - **Rendimiento de Indexación:** Documentos indexados por segundo.
  - **Estado del Crawler:** Número de URLs en cola, tasa de errores HTTP.
  - **Uso de Recursos:** CPU y memoria del contenedor de la aplicación.
  - **Salud de la Base de Datos:** Tiempos de respuesta de queries, conexiones activas (proporcionado por MongoDB Atlas).
- **Alertas (Configuradas en Alertmanager/Grafana):**
  - Alerta si la latencia de búsqueda p99 supera los 500ms.
  - Alerta si la tasa de errores supera el 1% durante 5 minutos.
  - Alerta si el proceso de indexación o crawling falla repetidamente.
  - Alerta si el uso de CPU o memoria alcanza el 85% del límite.

### 7.2. Backups y Recuperación de Desastres
- **Base de Datos:** Al usar MongoDB Atlas, se configurarán copias de seguridad automáticas y recuperación a un punto en el tiempo (PITR). Esto simplifica enormemente la recuperación de desastres.
- **Aplicación:** La aplicación es "stateless" (sin estado), por lo que la recuperación consiste simplemente en volver a desplegar la imagen de Docker funcional desde el registro de contenedores.

### 7.3. Proceso de Actualización
- **Aplicación:** Las actualizaciones se realizarán mediante despliegues "rolling update" en el orquestador de contenedores para garantizar cero tiempo de inactividad.
- **Índice:** En caso de cambios que requieran una re-indexación completa (ej. un nuevo analizador de texto), se seguirá una estrategia de "índice azul/verde":
  1. Se crea un nuevo conjunto de colecciones en MongoDB (ej. `documents_v2`, `inverted_index_v2`).
  2. Se lanza un trabajo de re-indexación que puebla estas nuevas colecciones.
  3. Una vez completado, la aplicación se reconfigura (mediante una variable de entorno) para apuntar a las nuevas colecciones.
  4. Las colecciones antiguas se archivan y eliminan.

## 8. Cronograma Resumido y Milestones

| Fase                                     | Duración Estimada | Entregable Clave                                             |
| ---------------------------------------- | ----------------- | ------------------------------------------------------------ |
| **Fase 0: Fundación y Configuración**    | 2 días            | Repositorio, CI y entorno de desarrollo local funcional.     |
| **Fase 1: Core de Indexación y Análisis**| 1 semana          | CLI para indexar archivos locales en MongoDB.                |
| **Fase 2: Búsqueda y API**               | 1 semana          | API `/search` funcional con ranking TF-IDF y UI web básica.  |
| **Fase 3: Integración del Web Crawler**  | 1 semana          | API `/crawl` para indexar un sitio web.                      |
| **Fase 4: Mejoras y Refinamiento**       | 1.5 semanas       | Ranking BM25, snippets y endpoints de gestión del índice.    |
| **Fase 5: Producción y Operatividad**    | 1 semana          | Dockerfile de producción, métricas, logging y documentación. |
| **Total Estimado**                       | **~ 5-6 semanas** |                                                              |