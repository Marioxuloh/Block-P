# Block-P
main.go: Este es el punto de entrada de tu aplicación. Aquí se inicializan el servidor y el cliente de gRPC y se inicia la ejecución.

go.mod y go.sum: Estos archivos son parte del sistema de módulos de Go y se utilizan para gestionar las dependencias de tu proyecto.

api/: En este directorio, puedes mantener tus archivos relacionados con los protocolos de gRPC.

    servicio.proto: Este archivo define los mensajes y servicios gRPC que se utilizarán en tu aplicación. Debes definir aquí los protocolos de comunicación entre el maestro y los esclavos.

    servicio.pb.go y servicio_grpc.pb.go: Estos archivos se generan automáticamente a partir del archivo .proto y contienen las definiciones de mensajes y servicios generados por el compilador de protoc.

cmd/: Este directorio contiene los programas principales de tu aplicación, como el servidor y el cliente, aqui iria el codigo del dashboard.

    cli/: codigo relacionado con la linea de comandos en mi aplicacion

    dashboard/: codigo relacionado con servicio web como puede ser el dashboard y dentro de esta iran la vista y el controlador.

            vista/: En el dashboard, implementarás la lógica para establecer una conexión WebSocket con el servidor. Cuando el servidor envíe un mensaje a través del WebSocket (por ejemplo, una notificación sobre el libro recién agregado), el dashboard escuchará esos mensajes y actualizará la vista en consecuencia, todo en tiempo real sin necesidad de recargar la página.

            controller/: En el controlador del servidor, implementarás la lógica para manejar los WebSockets. Esto incluirá la gestión de conexiones WebSocket entrantes y el envío de mensajes a los clientes (en este caso, el dashboard). Cuando ocurra un evento que requiera una actualización en el dashboard (por ejemplo, la adición de un nuevo libro a la biblioteca), el controlador enviará un mensaje a través del WebSocket a todos los clientes conectados.

pkg/: Este directorio es adecuado para colocar código compartido entre el servidor y el cliente.

    server/: El código relacionado con el servidor, que puede incluir funciones de inicialización, lógica de manejo de solicitudes y lógica de conmutación de roles.

    client/: El código relacionado con el cliente, que puede incluir funciones de conexión, comunicación y lógica específica del cliente.

    models/: El codigo relaciuonado con los modelos del sistema el cual se comunicara con el controlador para actualizar la vista,con la base de datos. en resumen
    cada uno representa una entidad o concepto específico dentro de tu aplicación.

    dao/: patron utilizado para separar la logica de negocio del acceso a datos, encapsula la base de datos

config/: Un directorio que almacena archivos de configuración. En este caso, config.json podría ser un archivo que contiene información sobre el rol del nodo (maestro o esclavo) y otras configuraciones necesarias.

    config.json: Un archivo de configuración que define el rol del nodo, entre otros parámetros.


logica del codigo: en el caso de que en el dashboard yo presente una colección de libros que hay en una base de datos, cuando a mi servidor le llegue una petición de añadir un libro, mi código de recepción del mensaje recibe el mensaje y llamara a una función del modelo(libro.go) que sea añadir un libro, este se agregara a la base de datos y se notificara al controlador para posteriormente actualizar el dashboard.