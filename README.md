# Distributed Job Processing Platform (MVP Profesional)

Plataforma de procesamiento distribuido de contratos diseñada para ser escalable, eficiente y robusta, utilizando las mejores prácticas de Go y servicios modernos de infraestructura.

## 🛠️ Stack Tecnológico

- **Lenguaje:** Go (Golang) 1.22+
- **Base de Datos:** MongoDB (Utilizando ObjectIDs nativos para máxima eficiencia).
- **Almacenamiento de Archivos:** AWS S3 (Con generación de URLs pre-firmadas).
- **Enrutador:** Chi Router (Ligero y compatible con `net/http`).

## 🏗️ Arquitectura y Decisiones Técnicas

- **ID Management:** Se utilizan `primitive.ObjectID` de MongoDB. Esto permite inserciones ordenadas naturalmente y ahorro de espacio en memoria para los índices.
- **Concurrencia:** Implementación del patrón **Fan-Out / Fan-In** mediante Goroutines y Channels para el zipeado concurrente de archivos, optimizando el uso de la red y el CPU.
- **Persistencia Asíncrona:** Los metadatos se guardan en MongoDB mientras que los documentos físicos se almacenan en S3, manteniendo una referencia cruzada segura.

## 🚀 API Endpoints

Todos los endpoints que requieren identificar a un cliente utilizan el header `client_id` (Hexadecimal de 24 caracteres).

### 👥 Clientes
- `POST /api/v1/clients`: Crea un nuevo cliente.
- `GET /api/v1/clients/{id}`: Obtiene el detalle de un cliente por su ID.

### 📄 Contratos
- `POST /api/v1/contracts`: Sube un contrato (Multipart Form). Requiere `client_id` en el Header. Guarda el archivo en S3 y los metadatos en Mongo.
- `GET /api/v1/contracts/{id}`: Obtiene un contrato específico con una **URL pre-firmada** temporal para descarga.
- `GET /api/v1/contracts/client`: Lista todos los contratos asociados al `client_id` enviado por Header.
- `GET /api/v1/contracts/client/zip`: **(Procesamiento Asíncrono)** Descarga todos los contratos del cliente y los entrega en un archivo `.zip` generado de forma concurrente "al vuelo".

## 📦 Infraestructura
- **MongoDB:** Colecciones de `clients` y `contracts`.
- **AWS S3:** Almacenamiento organizado por carpetas según el ID del cliente: `bucket/client_id/filename.pdf`.
