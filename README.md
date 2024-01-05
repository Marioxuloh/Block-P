# Block-P

Este repositorio contiene la aplicación Block-P, diseñada para gestionar la comunicación y el monitoreo de nodos maestros y esclavos. A continuación, se presenta una guía de la estructura y los componentes clave del proyecto:

## Archivos Principales

- **main.go**: Punto de entrada de la aplicación. Inicializa el servidor gRPC y el cliente gRPC y el Dashboard.

- **go.mod y go.sum**: Parte del sistema de módulos de Go, utilizados para gestionar las dependencias del proyecto.

## Componentes de la Aplicación

### Directorio `cmd/`

Este directorio contiene los programas principales de la aplicación, como el control por linea de comandos. También es el lugar donde se encuentra el código de la interfaz grafica de control:

- **`cli/`**: Contiene el código relacionado con la línea de comandos en la aplicación.

- **`dashboard/`**: Aquí se encuentra el código relacionado con el servicio web, como el dashboard. En esta carpeta, encontrarás:

  - **`view/`**: Implementa la lógica para establecer una conexión WebSocket con el servidor y despliega un dashboard de control a nivel local web. El dashboard escucha los mensajes enviados por el controlador a través del WebSocket y actualiza la vista en tiempo real sin necesidad de recargar la página.

  - **`controller/`**: El controlador del servidor implementa la lógica para manejar los WebSockets. Esto incluye la gestión de conexiones WebSocket entrantes y el envío de mensajes a los clientes, como el dashboard. Cuando ocurre un evento que requiere una actualización en el dashboard, como la adición de un nuevo libro a la biblioteca, el controlador envía un mensaje a través del WebSocket a todos los clientes conectados.

### Directorio `pkg/`

Este directorio es adecuado para almacenar todo el codigo fuente de la parte de servidor, cliente y el modelo:

- **`server/`**: Contiene el código relacionado con el servidor, que puede incluir funciones de inicialización, lógica de manejo de solicitudes y lógica de conmutación de roles(manejo de peticiones de entrada), la carpeta common de dentro es lo que tienen en comun con los otros nodos, el codigo inyectado por inyeccion de archivos remotos para su posterior ejecucion los cuales son temporales hasta que ya no se requiera ese balanceo de cargas.

- **`client/`**: Aquí se encuentra el código relacionado con el cliente, que puede incluir funciones de conexión, comunicación y lógica específica del cliente(manejo de peticiones salientes). la carpeta common de dentro es lo que tienen en comun con los otros nodos, el codigo inyectado por inyeccion de archivos remotos para su posterior ejecucion los cuales son temporales hasta que ya no se requiera ese balanceo de cargas.

- **`models/`**: En este directorio se encuentra el código relacionado con los modelos del sistema. Los modelos se comunican con el controlador para actualizar la vista y con la base de datos. Cada modelo representa una entidad o concepto específico dentro de tu aplicación, no solo es modelo del dashboard si no que aqui se agruparan todos los modelos del sistema.
Una de las cosas mas importantes es la inicializacion y creacion de la configuracion global basada en un archivo config.config

