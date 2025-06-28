# Sistema de Gestión de E-commerce
## Objetivo del Programa

Este proyecto implementa un sistema backend de gestión de e-commerce desarrollado en **Go (Golang)**. Su objetivo principal es simular las operaciones fundamentales de una tienda en línea, gestionando la información de productos, usuarios y pedidos de forma eficiente y modular, exponiendo sus funcionalidades a través de una **API RESTful**.

## Principales Funcionalidades

El sistema ofrece las siguientes capacidades clave, accesibles a través de sus servicios web:

### Módulo de Productos
* **`POST /products`**: **Creación de Productos.** Permite añadir nuevos productos al inventario con detalles como nombre, descripción, precio, stock y categoría.
* **`GET /products`**: **Listado de Productos.** Obtiene un listado completo de todos los productos disponibles en el inventario.
* **`GET /products/{id}`**: **Consulta de Producto por ID.** Recupera los detalles de un producto específico utilizando su identificador único.
* **`PUT /products/{id}`**: **Actualización de Productos.** Modifica la información de un producto existente.
* **`DELETE /products/{id}`**: **Eliminación de Productos.** Remueve un producto del inventario.

### Módulo de Usuarios
* **`POST /users/register`**: **Registro de Usuarios.** Permite a nuevos usuarios crear una cuenta en el sistema.
* **`POST /users/login`**: **Autenticación de Usuarios.** Valida las credenciales de un usuario (email y contraseña) para permitirle acceder al sistema.

### Módulo de Pedidos
* **`POST /orders`**: **Creación de Pedidos.** Procesa nuevas órdenes de compra, vinculándolas a un usuario, gestionando los ítems seleccionados con sus cantidades, verificando stock y calculando el total.
* **`GET /orders/{userId}`**: **Listado de Pedidos por Usuario.** Obtiene todos los pedidos realizados por un usuario específico.
* **`PUT /orders/{orderId}/status`**: **Actualización de Estado de Pedido.** Modifica el estado de un pedido (ej. de "Pendiente" a "Procesado", "Enviado", "Entregado" o "Cancelado").
* **`GET /orders`**: **Listado de Todos los Pedidos.** Permite consultar todos los pedidos registrados en el sistema (ideal para roles de administración).

## 🛠️ Tecnologías Utilizadas

* **Go (Golang):** Lenguaje de programación principal, elegido por su rendimiento, concurrencia y facilidad para construir APIs.
* **Gorilla Mux:** Librería robusta para el enrutamiento HTTP en Go, facilitando la definición de rutas y métodos para la API.
* **JSON:** Formato estándar para la serialización y deserialización de datos en las comunicaciones de la API, garantizando la interoperabilidad.
* **Almacenamiento en memoria:** Para la persistencia de datos (productos, usuarios, pedidos) durante la ejecución del programa, lo que permite una configuración rápida para demostraciones. (En un entorno de producción, esto se reemplazaría por una base de datos).

## Estructura del Proyecto

El proyecto sigue una estructura modular y limpia para mantener el código organizado y escalable:

* `cmd/api/main.go`: El punto de entrada de la aplicación. Aquí se inicializa el servidor web, se configuran las dependencias de los servicios y se definen todas las rutas de la API utilizando Gorilla Mux.
* `internal/api/handlers.go`: Contiene las funciones que actúan como "manejadores" de las solicitudes HTTP. Son la interfaz entre las peticiones web y la lógica de negocio.
* `internal/products/`: Módulo encapsulado para la gestión de productos.
    * `model.go`: Define las estructuras de datos (structs) para `Product` y `ProductRequest`, incluyendo las etiquetas `json` para la serialización.
    * `repository.go`: Implementación del almacenamiento de datos de productos (actualmente en un mapa en memoria).
    * `service.go`: Contiene la lógica de negocio para las operaciones CRUD de productos y validaciones.
* `internal/users/`: Módulo encapsulado para la gestión de usuarios.
    * `model.go`: Define las estructuras de datos para `User`, `UserRegisterRequest` y `UserLoginRequest`.
    * `repository.go`: Implementación del almacenamiento de datos de usuarios (en memoria).
    * `service.go`: Contiene la lógica de negocio para el registro y autenticación de usuarios.
* `internal/orders/`: Módulo encapsulado para la gestión de pedidos.
    * `model.go`: Define las estructuras de datos para `Order`, `LineItem`, `OrderRequest` y `LineItemRequest`.
    * `repository.go`: Implementación del almacenamiento de datos de pedidos (en memoria).
    * `service.go`: Contiene la lógica de negocio para la creación de pedidos (interactuando con productos y usuarios), consulta y actualización de estados.
      
  * **Fecha de Última Actualización:** 29 de junio de 2025
* **Integrantes del Grupo:**
    * Kevin Daniel Aguilar Baca
