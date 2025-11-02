# Hoja de Ruta de Desarrollo para `gofetch`

**Asunto:** Plan de Proyecto Detallado

**De:** Jules, Project Manager

**Para:** Equipo de Desarrollo

Hola equipo,

He preparado una hoja de ruta detallada para las próximas mejoras de `gofetch`. Este documento servirá como nuestra guía, desglosando cada iniciativa en tareas manejables con criterios de aceptación claros. Nuestro objetivo es trabajar de manera incremental, asegurando la calidad en cada paso.

---

### **Fase 1: Mejora de la Calidad del Núcleo y Relevancia de la Búsqueda**

**Objetivo Estratégico:** Aumentar la precisión de los resultados de búsqueda y fortalecer la base técnica del proyecto para futuras expansiones.

---

#### **Épica 1: Mejora del Algoritmo de Análisis de Texto**

*   **Iniciativa:** Implementar *Stemming* para normalizar las palabras y mejorar la relevancia.
*   **Justificación:** Actualmente, "correr" y "corriendo" son tratados como términos distintos. Al normalizarlos a una raíz común, las búsquedas serán más inteligentes y devolverán resultados más relevantes.

    *   **Tarea 1.1: Integración de una Biblioteca de *Stemming***
        *   **Descripción:** Investigar, seleccionar e instalar una biblioteca de Go que soporte *stemming* para múltiples idiomas (inglés y español). El sistema debe ser capaz de utilizar esta biblioteca para procesar texto.
        *   **Criterios de Aceptación:**
            *   [ ] Se ha añadido una nueva dependencia para el *stemming* en el archivo `go.mod`.
            *   [ ] La dependencia se ha descargado y es accesible desde el código.
            *   [ ] El proyecto compila sin errores después de añadir la dependencia.

    *   **Tarea 1.2: Modificación del Servicio de Análisis de Texto**
        *   **Descripción:** El servicio de análisis de texto (`Analyzer`) debe ser modificado para incluir un paso de *stemming* en su pipeline de procesamiento. Debe ser capaz de manejar diferentes idiomas.
        *   **Criterios de Aceptación:**
            *   [ ] El `Analyzer` ahora tiene una configuración para saber en qué idioma operar (ej. "english" o "spanish").
            *   [ ] El proceso de análisis de texto ahora incluye: tokenización, normalización (minúsculas, etc.), filtrado de *stopwords* y, finalmente, *stemming*.
            *   [ ] Si el *stemming* falla para una palabra específica, el sistema debe usar la palabra normalizada (sin *stemming*) como fallback y no debe fallar.
            *   **Archivos afectados (pista):** `internal/analysis/analyzer.go`.

    *   **Tarea 1.3: Actualización y Verificación de Pruebas**
        *   **Descripción:** Las pruebas unitarias del `Analyzer` deben ser actualizadas para reflejar el nuevo comportamiento del *stemming*.
        *   **Criterios de Aceptación:**
            *   [ ] Se ha creado un nuevo caso de prueba que verifica que palabras como "running" y "runner" se reducen a su raíz en inglés.
            *   [ ] Se ha creado un nuevo caso de prueba que verifica el comportamiento equivalente en español (ej. "corriendo", "corredores").
            *   [ ] Todos los casos de prueba existentes han sido revisados y ajustados para que sus resultados esperados coincidan con el nuevo formato de *stemming*.
            *   [ ] El conjunto completo de pruebas (`make test`) se ejecuta con éxito.
            *   **Archivos afectados (pista):** `internal/analysis/analyzer_test.go`.

---

#### **Épica 2: Flexibilidad y Configuración**

*   **Iniciativa:** Permitir que la configuración del idioma del analizador sea dinámica.
*   **Justificación:** Hardcodear el idioma limita la flexibilidad del motor. Permitir la configuración a través de variables de entorno es una práctica estándar que facilita el despliegue y la adaptación a diferentes usos.

    *   **Tarea 2.1: Implementar Configuración por Variable de Entorno**
        *   **Descripción:** Tanto el indexador como el servidor deben poder leer una variable de entorno (`ANALYZER_LANGUAGE`) para determinar qué analizador de idioma instanciar. Se debe establecer un idioma por defecto si la variable no está presente.
        *   **Criterios de Aceptación:**
            *   [ ] Al arrancar, la aplicación busca la variable de entorno `ANALYZER_LANGUAGE`.
            *   [ ] Si la variable es "spanish", se utiliza el analizador en español.
            *   [ ] Si la variable es "english" o no está definida, se utiliza el analizador en inglés por defecto.
            *   [ ] La lógica de selección del analizador está centralizada y no duplica código.
            *   **Archivos afectados (pista):** `internal/server/router.go`, `cmd/indexer/main.go`.

    *   **Tarea 2.2: Documentación de la Nueva Configuración**
        *   **Descripción:** El archivo `README.md` debe ser actualizado para informar a los usuarios y desarrolladores sobre la nueva opción de configuración.
        *   **Criterios de Aceptación:**
            *   [ ] El `README.md` tiene una nueva entrada en la sección de "Configuración" para `ANALYZER_LANGUAGE`.
            *   [ ] La documentación explica claramente los valores aceptados ("english", "spanish") y el comportamiento por defecto.

---

### **Fase 2: Expansión de Funcionalidades y Experiencia de Usuario**

**Objetivo Estratégico:** Mejorar la utilidad del motor de búsqueda añadiendo funcionalidades clave y soportando más tipos de contenido.

---

#### **Épica 3: Mejora de la Interfaz de Búsqueda**

*   **Iniciativa:** Añadir paginación a los resultados de búsqueda.
*   **Justificación:** Devolver todos los resultados a la vez no es escalable y puede sobrecargar tanto el frontend como el backend. La paginación es esencial para una buena experiencia de usuario.

    *   **Tarea 3.1: Modificar la API para Soportar Paginación**
        *   **Descripción:** La API de búsqueda (`/api/v1/search`) debe ser extendida para aceptar parámetros de paginación.
        *   **Criterios de Aceptación:**
            *   [ ] La API acepta dos nuevos parámetros de consulta: `page` (número de página) y `pageSize` (resultados por página).
            *   [ ] Si no se proporcionan, se utilizan valores por defecto razonables (ej. `page=1`, `pageSize=10`).
            *   [ ] La lógica de búsqueda ahora solo devuelve el subconjunto de documentos correspondiente a la página solicitada.
            *   [ ] La respuesta de la API ahora es un objeto JSON que contiene no solo los `results`, sino también metadatos de paginación (`totalResults`, `totalPages`, `currentPage`).
            *   **Archivos afectados (pista):** `internal/server/handler/search.go`, `internal/search/searcher.go`.

    *   **Tarea 3.2: Actualizar la Interfaz de Usuario**
        *   **Descripción:** La página de búsqueda debe ser actualizada para mostrar controles de paginación (ej. botones "Siguiente" y "Anterior") y manejar la nueva estructura de respuesta de la API.
        *   **Criterios de Aceptación:**
            *   [ ] El frontend realiza solicitudes a la API incluyendo los parámetros `page` y `pageSize`.
            *   [ ] Se muestran controles de navegación de página debajo de los resultados.
            *   [ ] Hacer clic en "Siguiente" o "Anterior" carga la página de resultados correspondiente sin recargar la página completa.
            *   **Archivos afectados (pista):** `ui/index.html`, `ui/script.js`.

    *   **Tarea 3.3: Documentar la API Paginada**
        *   **Descripción:** La documentación de la API en el `README.md` debe ser actualizada para reflejar los nuevos parámetros y la nueva estructura de la respuesta.
        *   **Criterios de Aceptación:**
            *   [ ] La sección "API Endpoints" del `README.md` incluye la descripción de los parámetros `page` y `pageSize`.
            *   [ ] Se ha actualizado el ejemplo de la respuesta JSON para mostrar la nueva estructura con metadatos de paginación.

---

#### **Épica 4: Expansión del Soporte de Contenido**

*   **Iniciativa:** Añadir la capacidad de indexar archivos PDF y DOCX.
*   **Justificación:** Limitar el motor a archivos de texto plano (`.txt`, `.md`) restringe su utilidad. Soportar formatos de documentos comunes como PDF y DOCX lo hará mucho más valioso.

    *   **Tarea 4.1: Abstracción de la Extracción de Contenido**
        *   **Descripción:** Crear una nueva abstracción (una interfaz de Go) llamada `Extractor` que defina un método para extraer texto de un archivo. Esto desacoplará el indexador del formato específico del archivo.
        *   **Criterios de Aceptación:**
            *   [ ] Se ha creado una nueva interfaz `Extractor` en un nuevo paquete (ej. `internal/extractor`).
            *   [ ] La interfaz define un método como `Extract(path string) (string, error)`.
            *   [ ] Se ha creado una implementación inicial de `Extractor` para archivos de texto plano que simplemente lee el contenido del archivo.
            *   **Archivos afectados (pista):** `internal/extractor/extractor.go`, `internal/extractor/text_extractor.go`.

    *   **Tarea 4.2: Integración de Bibliotecas para PDF y DOCX**
        *   **Descripción:** Investigar, seleccionar e instalar bibliotecas de Go para leer y extraer texto de archivos PDF y DOCX.
        *   **Criterios de Aceptación:**
            *   [ ] Se han añadido y verificado las dependencias para el procesamiento de PDF.
            *   [ ] Se han añadido y verificado las dependencias para el procesamiento de DOCX.
            *   [ ] Se han creado nuevas implementaciones de la interfaz `Extractor` para cada uno de estos formatos.
            *   **Archivos afectados (pista):** `internal/extractor/pdf_extractor.go`, `internal/extractor/docx_extractor.go`.

    *   **Tarea 4.3: Modificación del Indexador**
        *   **Descripción:** El indexador debe ser actualizado para utilizar el sistema de `Extractor`s.
        *   **Criterios de Aceptación:**
            *   [ ] El indexador ahora reconoce las extensiones `.pdf` y `.docx` además de las existentes.
            *   [ ] Antes de procesar un archivo, el indexador selecciona el `Extractor` apropiado basado en la extensión del archivo.
            *   [ ] El texto extraído (en lugar del contenido crudo del archivo) se pasa al `Analyzer`.
            *   [ ] El sistema maneja con gracia los errores durante la extracción de texto (ej. un PDF corrupto).
            *   **Archivos afectados (pista):** `internal/indexer/indexer.go`.

---

### **Fase 3: Capacidades Avanzadas y Pruebas Robustas**

**Objetivo Estratégico:** Transformar `gofetch` de un motor de búsqueda de archivos a uno con capacidad para indexar la web, y asegurar la estabilidad del sistema con pruebas de integración.

---

#### **Épica 5: Rastreo Web (Web Crawling)**

*   **Iniciativa:** Implementar un rastreador web para indexar contenido de sitios web.
*   **Justificación:** Esta es la evolución natural del proyecto, permitiendo a `gofetch` buscar contenido más allá de los archivos locales.

    *   **Tarea 5.1: Desarrollo del Módulo de Rastreo**
        *   **Descripción:** Crear un nuevo componente, el `Crawler`, responsable de descargar páginas web, extraer su contenido de texto y encontrar nuevos enlaces para visitar.
        *   **Criterios de Aceptación:**
            *   [ ] El `Crawler` puede descargar el HTML de una URL dada.
            *   [ ] Puede extraer el texto visible del HTML, eliminando etiquetas, scripts y estilos.
            *   [ ] Puede identificar y extraer todos los enlaces (`<a>` tags) de una página para añadirlos a una cola de URLs a visitar.
            *   [ ] El `Crawler` gestiona las URLs visitadas para no procesar la misma página dos veces.
            *   **Archivos afectados (pista):** `internal/crawler/crawler.go`.

    *   **Tarea 5.2: Integración del Crawler con el Indexador**
        *   **Descripción:** El texto extraído por el `Crawler` debe ser procesado y almacenado por el indexador existente.
        *   **Criterios de Aceptación:**
            *   [ ] El `Crawler` pasa el contenido de texto de cada página al `Analyzer`.
            *   [ ] El resultado del análisis se almacena en la base de datos de la misma manera que un documento de archivo, pero identificando la URL como fuente.
            *   [ ] Se ha creado un nuevo punto de entrada (`cmd/crawler/main.go`) que permite iniciar el proceso de rastreo desde la línea de comandos con una URL inicial.

    *   **Tarea 5.3 (Avanzado): Rastreo Responsable**
        *   **Descripción:** Implementar funcionalidades para asegurar que el rastreador se comporte de manera ética y eficiente.
        *   **Criterios de Aceptación:**
            *   [ ] El `Crawler` respeta las directivas del archivo `robots.txt` de los sitios web.
            *   [ ] Se ha implementado un retardo configurable entre solicitudes para no sobrecargar los servidores.

---

#### **Épica 6: Pruebas de Integración Completas**

*   **Iniciativa:** Crear un conjunto de pruebas de integración que verifique el sistema de extremo a extremo.
*   **Justificación:** Mientras que las pruebas unitarias son buenas para verificar componentes aislados, las pruebas de integración son cruciales para asegurar que todas las partes del sistema funcionan juntas correctamente.

    *   **Tarea 6.1: Configuración del Entorno de Pruebas de Integración**
        *   **Descripción:** Preparar un entorno de pruebas automatizado que pueda ejecutar la aplicación completa, incluyendo una base de datos.
        *   **Criterios de Aceptación:**
            *   [ ] Se ha configurado `testcontainers-go` o una herramienta similar para iniciar un contenedor de MongoDB para cada ejecución de prueba.
            *   [ ] Las pruebas no dependen de una base de datos externa y son completamente autocontenidas.

    *   **Tarea 6.2: Escritura de Escenarios de Prueba**
        *   **Descripción:** Escribir pruebas que simulen el flujo de trabajo de un usuario.
        *   **Criterios de Aceptación:**
            *   [ ] Se ha creado una prueba que (1) ejecuta el indexador sobre un conjunto de archivos de prueba, (2) inicia el servidor, (3) realiza una solicitud a la API de búsqueda y (4) verifica que los resultados son correctos y están ordenados por relevancia.
            *   [ ] Se ha añadido un nuevo comando `make test-integration` al `Makefile` para ejecutar estas pruebas.

---

Este plan de proyecto debería proporcionar una dirección clara para el futuro de `gofetch`.

Atentamente,

**Jules**
Project Manager
