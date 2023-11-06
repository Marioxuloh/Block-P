# Block-P
main.go: Este es el punto de entrada de tu aplicación. Aquí se inicializan el servidor y el cliente de gRPC y se inicia la ejecución.

go.mod y go.sum: Estos archivos son parte del sistema de módulos de Go y se utilizan para gestionar las dependencias de tu proyecto.

api/: En este directorio, puedes mantener tus archivos relacionados con los protocolos de gRPC.

    servicio.proto: Este archivo define los mensajes y servicios gRPC que se utilizarán en tu aplicación. Debes definir aquí los protocolos de comunicación entre el maestro y los esclavos.

    servicio.pb.go y servicio_grpc.pb.go: Estos archivos se generan automáticamente a partir del archivo .proto y contienen las definiciones de mensajes y servicios generados por el compilador de protoc.

cmd/: Este directorio contiene los programas principales de tu aplicación, como el servidor y el cliente.

    servidor/: Aquí puedes mantener el código relacionado con el servidor de tu aplicación.

    cliente/: Este directorio alberga el código del cliente de tu aplicación.

pkg/: Este directorio es adecuado para colocar código compartido entre el servidor y el cliente.

    server/: El código relacionado con el servidor, que puede incluir funciones de inicialización, lógica de manejo de solicitudes y lógica de conmutación de roles.

    client/: El código relacionado con el cliente, que puede incluir funciones de conexión, comunicación y lógica específica del cliente.

config/: Un directorio que almacena archivos de configuración. En este caso, config.json podría ser un archivo que contiene información sobre el rol del nodo (maestro o esclavo) y otras configuraciones necesarias.

    config.json: Un archivo de configuración que define el rol del nodo, entre otros parámetros.