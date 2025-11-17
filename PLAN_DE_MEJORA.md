# Plan de Mejora para el Proyecto gofetch

## Introducción

Este documento detalla una revisión del proyecto `gofetch` y propone un plan de acción para abordar las áreas de mejora identificadas. El objetivo es mejorar la mantenibilidad, robustez y escalabilidad del código.

## Áreas de Mejora Identificadas

### 1. Gestión de la Configuración

*   **Problema:** La configuración está descentralizada y es inconsistente. Algunos valores están hardcodeados (`crawler`), otros se pasan por flags (`indexer`) y otros se mencionan como variables de entorno (`README.md`).
*   **Solución:**
    *   Implementar una librería de gestión de configuración como [Viper](https://github.com/spf13/viper).
    *   Centralizar toda la configuración (MongoDB URI, nombre de la base de datos, puerto del servidor, URLs del crawler, etc.) en un único lugar.
    *   Permitir la configuración a través de archivos (e.g., `config.yaml`) y variables de entorno.

### 2. Manejo de Errores y Logging

*   **Problema:** El uso de `panic(err)` y `log.Fatal` en caso de errores críticos (como la conexión a la base de datos) detiene la aplicación de forma abrupta.
*   **Solución:**
    *   Reemplazar `panic` y `log.Fatal` por un manejo de errores más robusto.
    *   Introducir un logger estructurado (como `logrus` o `zap`) para tener un registro de eventos más detallado y con diferentes niveles de severidad (debug, info, warn, error).

### 3. Estructura del Código y Duplicación

*   **Problema:**
    *   Existe duplicación de código en la inicialización de componentes (`storage`, `analyzer`, `indexer`) en `cmd/crawler/main.go` y `cmd/indexer/main.go`.
    *   Hay dos paquetes `storage`, uno en `internal` y otro en `pkg`, lo que genera confusión.
*   **Solución:**
    *   Refactorizar la creación de dependencias utilizando un contenedor de inyección de dependencias o un patrón de factoría para evitar la duplicación.
    *   Consolidar los paquetes `storage` en uno solo, probablemente en `pkg` si se considera un componente reutilizable.

### 4. Pruebas (Testing)

*   **Problema:** El `README.md` menciona `make test`, pero no está claro el alcance de las pruebas existentes. Un proyecto de esta envergadura necesita una suite de pruebas completa.
*   **Solución:**
    *   Realizar una auditoría de las pruebas existentes.
    *   Añadir pruebas unitarias para los componentes clave: `analysis`, `indexer`, `ranking`, `storage`.
    *   Añadir pruebas de integración para la API, verificando los endpoints y las respuestas.
    *   Configurar un pipeline de CI (GitHub Actions) que ejecute las pruebas automáticamente.

### 5. Crawler

*   **Problema:** El crawler tiene la URL de inicio y la profundidad hardcodeadas. La documentación es inconsistente sobre su estado de implementación.
*   **Solución:**
    *   Hacer que las URLs de inicio y la profundidad del crawler sean configurables.
    *   Actualizar el `README.md` para reflejar que el crawler es una funcionalidad implementada.

### 6. API

*   **Problema:** La API de búsqueda es muy básica.
*   **Solución:**
    *   Añadir paginación a los resultados de búsqueda.
    *   Mejorar las respuestas de error de la API para que sean más descriptivas.
    *   Considerar la adición de `rate limiting` para prevenir abusos.

### 7. Concurrencia

*   **Problema:** El `README.md` menciona una arquitectura concurrente, pero no es evidente cómo se está utilizando.
*   **Solución:**
    *   Revisar la implementación del `indexer` y el `crawler` para asegurarse de que el uso de goroutines es eficiente y seguro (sin race conditions).
    *   Añadir métricas para monitorizar el rendimiento de los procesos concurrentes.

### 8. Seguridad

*   **Problema:** No se han identificado medidas de seguridad explícitas.
*   **Solución:**
    *   Realizar una revisión de seguridad del código de la API para prevenir vulnerabilidades comunes como Inyección de NoSQL.
    *   Añadir validación de entradas en todos los endpoints.

## Plan de Acción

A continuación, se presenta un plan de acción detallado para implementar las mejoras propuestas.

1.  **Fase 1: Refactorización y Configuración**
    *   **Tarea 1.1:** Integrar Viper para la gestión de la configuración. Crear un archivo `config.yaml.example` con la configuración por defecto.
    *   **Tarea 1.2:** Centralizar la lógica de conexión a la base de datos en un solo lugar.
    *   **Tarea 1.3:** Consolidar los paquetes `storage` duplicados.
    *   **Tarea 1.4:** Refactorizar la creación de dependencias para eliminar código duplicado.

2.  **Fase 2: Robustez y Pruebas**
    *   **Tarea 2.1:** Reemplazar `panic` y `log.Fatal` con un logger estructurado.
    *   **Tarea 2.2:** Añadir pruebas unitarias para los paquetes `analysis` e `indexer`.
    *   **Tarea 2.3:** Añadir pruebas de integración para el endpoint de búsqueda de la API.

3.  **Fase 3: Funcionalidades y Mejoras**
    *   **Tarea 3.1:** Hacer configurable el crawler.
    *   **Tarea 3.2:** Añadir paginación a la API de búsqueda.
    *   **Tarea 3.3:** Actualizar la documentación (`README.md`).

4.  **Fase 4: Seguridad y Optimización**
    *   **Tarea 4.1:** Realizar una auditoría de seguridad del código.
    *   **Tarea 4.2:** Analizar y optimizar el uso de concurrencia en el `indexer`.
