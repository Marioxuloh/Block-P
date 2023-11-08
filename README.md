# Block-P

Este repositorio contiene la aplicación Block-P, diseñada para gestionar la comunicación y el monitoreo de nodos maestros y esclavos. A continuación, se presenta una guía de la estructura y los componentes clave del proyecto:

## Archivos Principales

- **main.go**: Punto de entrada de la aplicación. Inicializa el servidor y el cliente de gRPC y comienza la ejecución.

- **go.mod y go.sum**: Parte del sistema de módulos de Go, utilizados para gestionar las dependencias del proyecto.

## Estructura de mensajes gRPC

El directorio `api/` contiene archivos relacionados con los protocolos de gRPC:

- **servicio.proto**: Define los mensajes y servicios gRPC que se utilizan en la aplicación. Aquí se definen los protocolos de comunicación entre el maestro y los esclavos junto a otros protocolos, cada uno de ellos tendra su propio .proto para la definicion de la estructura de sus mensajes y servicios.

- **servicio.pb.go y servicio_grpc.pb.go**: Generados automáticamente a partir del archivo `.proto`, contienen las definiciones de mensajes y servicios generados por el compilador de protoc.

## Componentes de la Aplicación

### Directorio `cmd/`

Este directorio contiene los programas principales de la aplicación, como el control por linea de comandos. También es el lugar donde se encuentra el código de la interfaz grafica de control:

- **cli/**: Contiene el código relacionado con la línea de comandos en la aplicación.

- **dashboard/**: Aquí se encuentra el código relacionado con el servicio web, como el dashboard. En esta carpeta, encontrarás:

  - **view/**: Implementa la lógica para establecer una conexión WebSocket con el servidor y despliega un dashboard de control a nivel local web. El dashboard escucha los mensajes enviados por el controlador a través del WebSocket y actualiza la vista en tiempo real sin necesidad de recargar la página.

  - **controller/**: El controlador del servidor implementa la lógica para manejar los WebSockets. Esto incluye la gestión de conexiones WebSocket entrantes y el envío de mensajes a los clientes, como el dashboard. Cuando ocurre un evento que requiere una actualización en el dashboard, como la adición de un nuevo libro a la biblioteca, el controlador envía un mensaje a través del WebSocket a todos los clientes conectados.

### Directorio `pkg/`

Este directorio es adecuado para colocar código compartido entre el servidor y el cliente:

- **server/**: Contiene el código relacionado con el servidor, que puede incluir funciones de inicialización, lógica de manejo de solicitudes y lógica de conmutación de roles(manejo de peticiones de entrada).

- **client/**: Aquí se encuentra el código relacionado con el cliente, que puede incluir funciones de conexión, comunicación y lógica específica del cliente(manejo de peticiones salientes).

- **models/**: En este directorio se encuentra el código relacionado con los modelos del sistema. Los modelos se comunican con el controlador para actualizar la vista y con la base de datos. Cada modelo representa una entidad o concepto específico dentro de tu aplicación, no solo es modelo del dashboard si no que aqui se agruparan todos los modelos del sistema.

- **dao/**: El patrón DAO se utiliza para separar la lógica de negocio del acceso a datos y encapsular la interacción con la base de datos.

## Directorio `config/`

Este directorio almacena archivos de configuración. El archivo `config.json` podría contener información sobre el rol del nodo (maestro o esclavo) y otras configuraciones necesarias.

- **config.json**: Un archivo de configuración que define el rol del nodo y otros parámetros.

## Lógica del Código

En el caso de que en el dashboard presentes una colección de libros almacenados en una base de datos, cuando el servidor recibe una solicitud para añadir un libro, el código de recepción del mensaje llama a una función del modelo `libro.go` para añadir el libro. Luego, el libro se agrega a la base de datos y se notifica al controlador para actualizar el dashboard en tiempo real.
