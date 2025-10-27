### **Sprint 0: La Fundación (Duración: 2 días)**

**Sprint Goal:** **Crear un entorno de desarrollo profesional, repetible y automatizado.**

Este Sprint no produce ninguna funcionalidad visible para el usuario, pero es el más importante para nuestra velocidad y calidad a largo plazo. Se trata de afilar el hacha antes de cortar el árbol.

#### **Tareas del Sprint 0:**

---

**Task ID: F0.T1 - Configuración del Repositorio y Entorno Go**

*   **Objetivo:** Establecer las bases del control de versiones y la gestión de dependencias del proyecto.
*   **Actividades:**
    1.  Ejecuta `git init`.
    2.  Crea un `.gitignore` robusto para Go y artefactos del sistema operativo (puedes encontrar plantillas estándar online).
    3.  Crea el `README.md` inicial con el nombre del proyecto y una breve descripción.
    4.  Ejecuta `go mod init <nombre_del_modulo>`, por ejemplo, `go.search.engine/v1`.
    5.  Crea la estructura de directorios inicial: `/cmd`, `/internal`. No crees todos los subdirectorios aún, solo los que necesitemos para empezar.
*   **Criterios de Aceptación:**
    *   [ ] El repositorio está en GitHub/GitLab.
    *   [ ] El comando `go build ./...` se ejecuta sin errores.
    *   [ ] Archivos como binarios compilados o `*.local` son ignorados por Git.

---

**Task ID: F0.T2 - Contenerización del Entorno de Desarrollo**

*   **Objetivo:** Garantizar que cualquier desarrollador (o el sistema de CI) pueda levantar el entorno completo con un solo comando, eliminando el "funciona en mi máquina".
*   **Actividades:**
    1.  Crea un archivo `docker-compose.yml` en la raíz del proyecto.
    2.  Define dos servicios: `app` (para nuestra aplicación Go) y `db` (para MongoDB).
    3.  Configura el servicio `db` para que use la imagen oficial de `mongo`.
    4.  Configura el servicio `app` para que se construya a partir de un `Dockerfile` que crearás.
    5.  Usa volúmenes para persistir los datos de MongoDB (`/data/db`) fuera del contenedor.
    6.  Usa volúmenes para montar tu código fuente local dentro del contenedor `app` para no tener que reconstruir la imagen en cada cambio.
*   **Criterios de Aceptación:**
    *   [ ] El comando `docker-compose up` levanta ambos contenedores sin errores.
    *   [ ] El contenedor de la aplicación Go se inicia y se mantiene corriendo (aunque todavía no haga nada).
    *   [ ] Puedo ejecutar `docker-compose down` y luego `docker-compose up` de nuevo, y los datos en MongoDB persisten.
    *   [ ] Los servicios pueden comunicarse entre sí por la red de Docker.

> **Nota de PM:** La comunicación entre contenedores en Docker Compose es simple. Si tus servicios se llaman `app` y `db`, desde `app` puedes conectarte a MongoDB usando `mongodb://db:27017`.

---

**Task ID: F0.T3 - Pipeline de Integración Continua (CI) Básico**

*   **Objetivo:** Automatizar la validación de la calidad del código en cada cambio, detectando problemas de forma temprana.
*   **Actividades:**
    1.  Configura `golangci-lint` localmente. Elige un conjunto de linters razonable para empezar.
    2.  Crea un `Makefile` con objetivos (`lint`, `test`, `build`, `run-compose`). Esto estandariza los comandos.
    3.  Crea el directorio `.github/workflows/` y dentro un archivo `ci.yml`.
    4.  Define un workflow de GitHub Actions que se dispare en cada `push` a las ramas `main` y `develop`.
    5.  El workflow debe tener dos pasos: uno que ejecute `make lint` y otro que ejecute `make test`.
*   **Criterios de Aceptación:**
    *   [ ] Al hacer `push` de código con errores de formato, el pipeline de CI falla en el paso de linting.
    *   [ ] Al hacer `push` de código con tests que fallan, el pipeline de CI falla en el paso de testing.
    *   [ ] Un `push` con código limpio y tests que pasan resulta en un pipeline exitoso (verde).

---
### **Sprint 1: El Corazón del Análisis y la Indexación (Duración: 1 semana)**

**Sprint Goal:** **Ser capaces de leer un directorio de archivos de texto, analizar su contenido de forma inteligente y almacenar un índice invertido funcional en MongoDB.**

Este es el núcleo de nuestro motor. Al final de este sprint, tendremos el "cerebro" del sistema, aunque todavía no tengamos cómo "hablar" con él (la API).

#### **Tareas del Sprint 1 (Priorizadas):**

---

**Task ID: F1.T1 - Módulo de Análisis de Texto (`/internal/analysis`)**

*   **User Story:** Como Indexador, necesito analizar un bloque de texto para convertirlo en una lista de tokens limpios y relevantes, para poder construir un índice preciso.
*   **Prioridad:** **Máxima.** Esta es la dependencia fundamental para la indexación.
*   **Actividades:**
    1.  Crea el paquete `internal/analysis` con el archivo `analyzer.go`.
    2.  Define la `struct Analyzer`. Piensa en qué estado necesita mantener. (Pista: la lista de stopwords pre-procesada).
    3.  Implementa la lógica del **Tokenizer**. Debe ser una función privada. Decide tu estrategia: ¿dividir por espacios o usar expresiones regulares para una mayor precisión?
    4.  Implementa los **Filtros**. Serán funciones que toman `[]string` y devuelven `[]string`.
        *   Filtro de minúsculas.
        *   Filtro de puntuación. Considera los bordes: ¿qué hacer con palabras como "e-mail" o números? Por ahora, manténlo simple: elimina todo lo que no sea una letra.
    5.  Implementa el **Filtro de Stopwords**. Este debe ser un **método** del `Analyzer` ya que depende de su configuración.
    6.  Crea un constructor `New(stopwords []string) *Analyzer`. Dentro, convierte el slice de stopwords a una estructura de datos más eficiente para búsquedas (un mapa).
    7.  Ensambla todo en un método público `Analyze(text string) []string` que ejecute el pipeline completo en orden.
    8.  Crea el archivo `stopwords.go` con listas predefinidas para español e inglés.
    9.  Escribe tests unitarios exhaustivos en `analyzer_test.go` usando el patrón de "Table-Driven Tests".
*   **Criterios de Aceptación:**
    *   [ ] El método `Analyze("Este es un TEXTO de prueba, ¡genial!")` con stopwords en español devuelve `["texto", "prueba", "genial"]`.
    *   [ ] El método `Analyze(...)` con una lista de stopwords vacía no filtra ninguna palabra (excepto puntuación y mayúsculas).
    *   [ ] Los tests unitarios cubren casos borde: string vacío, texto solo con stopwords, texto solo con puntuación.
    *   [ ] La cobertura de tests para este paquete es superior al 80%.

> **Pista Técnica:** Para el filtro de stopwords, un `map[string]struct{}` es la forma más performante en Go de representar un "conjunto" para comprobaciones de existencia rápidas.

---

**Task ID: F1.T2 - Abstracción de Almacenamiento (`/internal/storage`)**

*   **User Story:** Como Indexador, necesito una forma simple y fiable de interactuar con MongoDB sin tener que conocer los detalles del driver, para poder persistir y recuperar datos del índice.
*   **Prioridad:** **Alta.** Necesitamos esto antes de poder guardar el trabajo del Analizador. Puede trabajarse en paralelo con F1.T1.
*   **Actividades:**
    1.  Crea el paquete `internal/storage`.
    2.  Define las `structs` de Go que mapean tus modelos de datos de MongoDB (`Document`, `InvertedIndexEntry`, `Posting`, `IndexStats`). Sé fiel al plan.
    3.  Crea un cliente de almacenamiento (`MongoStore` o similar) que contenga la conexión a la base de datos.
    4.  Implementa un método `Connect()` que lea la connection string de una variable de entorno y establezca la conexión.
    5.  Implementa un método `Disconnect()`.
    6.  Implementa los métodos de la interfaz que necesitará el indexador. No implementes todo aún, solo lo esencial para la indexación:
        *   `AddDocument(doc Document) (string, error)` (devuelve el ID del nuevo documento).
        *   `UpsertTerm(term string, posting Posting) error` (Esta es la operación clave: busca un término. Si existe, añade el `Posting` a su lista. Si no existe, crea el documento del término con el `Posting`).
        *   `UpdateIndexStats(totalDocs int, ...) error`
*   **Criterios de Aceptación:**
    *   [ ] El código se conecta exitosamente a la instancia de MongoDB del `docker-compose`.
    *   [ ] Puedo llamar a `AddDocument` y ver un nuevo documento en la colección `documents` usando un cliente de BBDD como Compass.
    *   [ ] La lógica de `UpsertTerm` funciona correctamente tanto para términos nuevos como para términos existentes.
    *   [ ] Los métodos manejan errores de base de datos de forma adecuada (devuelven `error`).

> **Nota de PM:** La operación "upsert" y la actualización de un array dentro de un documento son operaciones atómicas muy potentes en MongoDB. Investiga los operadores `$push` y la opción `upsert: true` en las operaciones de actualización del driver de Go. Usar operaciones `bulk` (en lotes) será clave para el rendimiento, tenlo en mente para la siguiente tarea.

---

**Task ID: F1.T3 - Indexador de Archivos Locales (`/internal/indexer`)**

*   **User Story:** Como Usuario, quiero procesar todos los archivos de texto de un directorio para que su contenido sea analizable y consultable por el motor de búsqueda.
*   **Prioridad:** **Media.** Esta tarea une F1.T1 y F1.T2.
*   **Actividades:**
    1.  Crea el paquete `internal/indexer`.
    2.  Crea una `struct Indexer` que reciba como dependencias una instancia del `Analyzer` y del `MongoStore`. Esto se llama Inyección de Dependencias y es crucial para el testing.
    3.  Implementa un método público `IndexDirectory(path string) error`.
    4.  Dentro de este método, implementa la lógica para caminar recursivamente por el directorio (el paquete `io/fs` es excelente para esto).
    5.  Para cada archivo soportado (`.txt`, `.md`), crea un pool de workers (goroutines) para procesarlos en paralelo. Un canal puede servir para distribuir el trabajo (rutas de archivos) a los workers.
    6.  Cada worker debe:
        a. Leer el contenido del archivo.
        b. Pasarlo al `Analyzer` para obtener los tokens.
        c. Calcular las frecuencias de término (`term_frequency`) y las posiciones para ese archivo.
        d. Preparar los datos para ser escritos en MongoDB.
    7.  Implementa una estrategia de escritura eficiente. En lugar de que cada worker escriba en la BBDD individualmente (lo que crearía contención), los workers pueden enviar sus resultados a través de otro canal a una única goroutine "escritora".
    8.  Esta goroutine escritora agrupará las actualizaciones y usará las operaciones `BulkWrite` de MongoDB para minimizar las idas y venidas a la base de datos.
*   **Criterios de Aceptación:**
    *   [ ] Al ejecutar la indexación sobre un directorio con 5 archivos `.txt`, las colecciones `documents`, `inverted_index` y `index_stats` en MongoDB se pueblan correctamente.
    *   [ ] La frecuencia de término (`tf`) y la frecuencia de documento (`df`) son correctas.
    *   [ ] El proceso utiliza múltiples núcleos de CPU gracias a las goroutines.
    *   [ ] El indexador informa de errores si no puede leer un archivo o escribir en la BBDD.

---

**Task ID: F1.T4 - Interfaz de Línea de Comandos (CLI) (`/cmd/indexer`)**

*   **User Story:** Como Administrador del sistema, quiero ejecutar un comando desde mi terminal para iniciar el proceso de indexación de un directorio específico.
*   **Prioridad:** **Baja.** Es la capa final que nos permite probar todo el flujo del Sprint.
*   **Actividades:**
    1.  Crea el directorio `/cmd/indexer` y dentro un `main.go`.
    2.  Usa el paquete estándar `flag` para aceptar un argumento de línea de comandos que especifique la ruta del directorio a indexar.
    3.  En la función `main`:
        a. Inicializa la conexión a la BBDD (`storage.Connect`).
        b. Usa `defer storage.Disconnect()` para asegurar que la conexión se cierre.
        c. Inicializa el `Analyzer` con una lista de stopwords (por ejemplo, español).
        d. Inicializa el `Indexer` inyectando el analizador y el cliente de almacenamiento.
        e. Llama al método `indexer.IndexDirectory()` con la ruta obtenida del flag.
        f. Imprime un mensaje de éxito o de error al terminar.
*   **Criterios de Aceptación:**
    *   [ ] Puedo ejecutar `go run ./cmd/indexer -path=/ruta/a/mis/archivos` y el proceso de indexación se inicia y completa.
    *   [ ] Si no proporciono el flag `-path`, el programa muestra un error y las instrucciones de uso.