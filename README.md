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

## � Docker

La plataforma utiliza una única imagen multi-etapa que contiene tanto el binario de la API como el del Worker.

### 1. Construir la imagen
```bash
docker build -t project-job-platform .
```

### 2. Ejecutar la API (Puerto 8080)
Como ahora usamos una base de datos en la nube (MongoDB Atlas), la conexión funciona directamente desde el contenedor:
```bash
docker run -p 8080:8080 --env-file .env project-job-platform
```

### 3. Ejecutar el Worker
```bash
docker run --env-file .env project-job-platform ./worker-bin
```

> [!TIP]
> Asegúrate de haber actualizado el campo `<db_password>` en tu archivo `.env` por la contraseña real de Atlas antes de correr los comandos. Como la URI apunta a la nube, ya no es necesario usar `host.docker.internal`.

## �🚀 API Endpoints

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


## Siguientes pasos
- En nustro ejemplo, queremos usar SNS para encolar las peticiones del zip. Entiendo que la seccion worker, estaria subcripto al SNS esperando consumir de a 1 mensaje a la vez para evitar el colapso. Una vez que hace el zip, este podria guardarlo en S3, y tener una lamba que gatille cuando se crea un zip nuevo en el bucket de zip, y mandar un mail al cliente dueño de ese file... pero como hago para q la lambda sepa el mail del zip que se sube al s3? ademas, aws tiene un sistema de mail que pueda usar para hacer las pruebas?

SNS   
contract-package-requested
arn:aws:sns:us-east-1:432162757798:contract-package-requested

SQS
contract-package-requested-queue
arn:aws:sqs:us-east-1:432162757798:contract-package-requested-queue