🚀 Proyecto en Go: Autenticación y Persistencia con MySQL utilizando Arquitectura Hexagonal (Ports & Adapters).
Este proyecto es una API desarrollada en Go que implementa autenticación con JWT y persistencia de datos utilizando MySQL con el paquete nativo de sql en Go.

📌 Tecnologías utilizadas
Go: Lenguaje principal
JWT: Para autenticación segura
Variables de entorno: Configuración segura de credenciales
🔧 Configuración
Antes de ejecutar el proyecto, configura tus variables de entorno en un archivo .env:

MYSQL_USER=tu_usuario_mysql
MYSQL_ROOT_PASSWORD=tu_contraseña_mysql
MYSQL_DB_HOST=localhost
MYSQL_DB_PORT=3306
MYSQL_DB_NAME=nombre_de_tu_base_de_datos
MYSQL_TABLE_NAME=nombre_de_tu_tabla

JWT_TOKEN_SECRET=tu_token_secreto

▶️ Ejecución
Instala las dependencias:
go mod tidy
Ejecuta la aplicación:
go run cmd/api/main.go
📌 Funcionalidades
✅ Registro y autenticación de usuarios
✅ Generación y validación de tokens JWT
✅ CRUD de usuarios con data persistente en MySQL
✅ Manejo de configuración con variables de entorno

💡 Contribuciones y feedback son bienvenidos. Si quieres probarlo o mejorarlo, ¡hablemos! 🚀
