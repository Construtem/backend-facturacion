# Backend - Facturación

Microservicio backend encargado del procesamiento y gestión de documentos tributarios en Construtem.

## 🛠️ Tecnologías
- Go
- PostgreSQL (compartida con ventas)
- JWT
- Docker

## 🚀 Funcionalidades
- Emisión de boletas y facturas.
- Registro de pagos.

## 📦 Instalación
```bash
docker-compose up --build
```

## 📁 Estructura
- `/controllers`
- `/models`
- `/routes`

## 🔗 Integraciones
- Se comunica con `backend-ventas` para emitir documentos a partir de ventas realizadas.
