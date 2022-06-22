# Registro de cambios

Todos los cambios notables en este proyecto se documentarán en este archivo.

El formato se basa en [Mantener un registro de cambios](https://keepachangelog.com/en/1.0.0/),
y este proyecto se adhiere a [Control de versiones semántico](https://semver.org/spec/v2.0.0.html).

## Inédito

## \[3.2.0] - 2022-06-05

### Fija

*   Refactorizar el ciclo de vida de cron (https://github.com/distribworks/dkron/pull/1119)
*   Transferencia de liderazgo (https://github.com/distribworks/dkron/pull/1109)

### Cambios

*   Agregar integración de telemetría de prewebhook y cronitor (https://github.com/distribworks/dkron/pull/1099)
*   Implementar el ejecutor GRPC (https://github.com/distribworks/dkron/pull/1049)
*   Usar golang/cross para crear docker multiarch (https://github.com/distribworks/dkron/pull/1105)
*   Implementar algunas sugerencias para la interfaz de usuario (https://github.com/distribworks/dkron/pull/1120)
*   Nuevo sitio web (https://github.com/distribworks/dkron/pull/1072)
*   Bump deps

## \[3.1.11] - 2022-04-07

### Fija

*   Intente arreglar el programador que no se inicia [#1053](https://github.com/distribworks/dkron/pull/1053)
*   Busque complementos en la ruta de configuración utilizada (como se documenta) en lugar de /etc/dkron codificado [#1024](https://github.com/distribworks/dkron/pull/1024)
*   No salir por un liderazgo fallido [#1082](https://github.com/distribworks/dkron/pull/1082)

### Cambios

*   Permitir múltiples corredores en kafka executor [#1037](https://github.com/distribworks/dkron/pull/1037)
*   Permitir el paso a través de una clave de mensaje para el ejecutor de Kafka [#1021](https://github.com/distribworks/dkron/pull/1021)
*   Dependencias bump
*   Crear codeql-analysis.yml
*   Actualizar la página actual después de utilizar el trabajo de alternancia/ejecución [#1026](https://github.com/distribworks/dkron/pull/1026)
*   enlace de proveedor de terraform compatible [#1029](https://github.com/distribworks/dkron/pull/1029)
*   Actualizar documentos de métricas para el archivo de configuración de Prometheus [#1058](https://github.com/distribworks/dkron/pull/1058)
*   Truncar la salida de ejecución en la vista predeterminada [#1025](https://github.com/distribworks/dkron/pull/1025)
*   Agregar extremo de ejecución [#1085](https://github.com/distribworks/dkron/pull/1085)

## \[3.1.10] - 2021-10-01

### Fija

*   El ejecutor de Nats consume todas las conexiones disponibles [#1020](https://github.com/distribworks/dkron/pull/1020)

### Cambios

*   Actualizar activos con los últimos cambios

## \[3.1.9] - 2021-09-29

### Funciones

*   Actualizar el texto de estado del nodo [#1012](https://github.com/distribworks/dkron/pull/1012)

### Fija

*   Agregar el campo 'siguiente' a la definición de swagger [#991](https://github.com/distribworks/dkron/pull/991)
*   Corregir error tipográfico de comentarios [#1000](https://github.com/distribworks/dkron/pull/1000)
*   Corregir errores tipográficos [#1005](https://github.com/distribworks/dkron/pull/1005)
*   Corregir error de notificación [#993](https://github.com/distribworks/dkron/pull/993)
*   fd-leak cuando falla la comunicación del cliente RPC [#1009](https://github.com/distribworks/dkron/pull/1009)
*   Pánico o errores de red [#1008](https://github.com/distribworks/dkron/pull/1008)
*   Espere todos los trabajos realizados [#1010](https://github.com/distribworks/dkron/pull/1010)
*   Espere a que el programador se detenga primero [#1016](https://github.com/distribworks/dkron/pull/1016)

### Cambios

*   Activos de Swagger de cdn [#997](https://github.com/distribworks/dkron/pull/997)
*   Usar API_URL relativas [#1006](https://github.com/distribworks/dkron/pull/1006)

## \[3.1.8] - 2021-06-16

### Fija

*   Las direcciones URL no funcionaban bien en la documentación [#979](https://github.com/distribworks/dkron/pull/979)
*   Mejorar el documento de actualización [#980](https://github.com/distribworks/dkron/pull/980)
*   Actualizar las etiquetas del agente al recargar la configuración [#983](https://github.com/distribworks/dkron/pull/983)
*   Corregir falsos negativos en la prueba intermitente [#982](https://github.com/distribworks/dkron/pull/982)
*   Corregir TLSRaftLayer init [#987](https://github.com/distribworks/dkron/pull/987)

### Cambios

*   Usar Buildkit estable [#977](https://github.com/distribworks/dkron/pull/977)
*   Controlar errores en el inicio del programador [#978](https://github.com/distribworks/dkron/pull/978)

## \[3.1.7] - 2021-05-29

### Funciones

*   Implementar efímero y caduca en la función [#972](https://github.com/distribworks/dkron/pull/972)
*   Agregar @minutely programación personalizada de nuevo [#970](https://github.com/distribworks/dkron/pull/970)

### Fija

*   Correcciones de la interfaz de usuario de DataGrid para campos largos [#965](https://github.com/distribworks/dkron/pull/965)
*   Arreglar las condiciones de la carrera [#967](https://github.com/distribworks/dkron/pull/967)
*   Corregir el ejecutor de shell que se bloquea en el comando que falta [#948](https://github.com/distribworks/dkron/pull/948)

### Cambios

*   Revisar el registro para evitar el var a nivel de paquete [#963](https://github.com/distribworks/dkron/pull/963)
*   Mejorar las pruebas para el ejecutor http [#936](https://github.com/distribworks/dkron/pull/936)
*   Proceso de refactorizaciónFilteredNodes para pruebas [#968](https://github.com/distribworks/dkron/pull/968)

## \[3.1.6] - 2021-03-23

### Funciones

*   Agregar un filtro en el estado deshabilitado [#923](https://github.com/distribworks/dkron/pull/923) @educlos
*   Proporcionar consulta de filtro por trabajo displayName, agregar informes de trabajos prístinos [#897](https://github.com/distribworks/dkron/pull/897) @MGSousa

### Fija

*   Corregir la vulnerabilidad de XSS [#922](https://github.com/distribworks/dkron/pull/922) @yvanoers
*   Título correcto de la página del ejecutor de NATS [#929](https://github.com/distribworks/dkron/pull/929) @yvanoers
*   Correcciones de la interfaz de usuario [#926](https://github.com/distribworks/dkron/pull/926) @yvanoers

### Cambios

*   Usar ir incrustar para incrustar activos [#931](https://github.com/distribworks/dkron/pull/931)

## \[3.1.5] - 2021-03-08

### Funciones

*   Tiempo de espera de trabajo configurable [#906](https://github.com/distribworks/dkron/pull/906)
*   Agregar ejecutor de kafka y nats [#854](https://github.com/distribworks/dkron/pull/854)
*   Agregar estadísticas de uso de informes [#910](https://github.com/distribworks/dkron/pull/910)

### Cambios

*   Golpear algunos deps

### Fija

*   Agregar la interfaz de usuario/dir público [#919](https://github.com/distribworks/dkron/pull/919)

## \[3.1.4] - 2021-01-25

### Cambios

*   Golpear algunos deps
*   Correcciones de la interfaz de usuario
    *   Volver a agregar la zona horaria a la interfaz de usuario
    *   Mostrar estado de ejecución

## \[3.1.3] - 2021-01-17

### Cambios

*   Varias mejoras en la interfaz de usuario [#891](https://github.com/distribworks/dkron/pull/891)
    *   Estado visual de los trabajos
    *   Acciones masivas para alternar y ejecutar
    *   Diseño de cuadrícula de datos de trabajos flexible
    *   Fijar estilo de reloj
    *   Representación de fechas fijas en ejecuciones no finalizadas

## \[3.1.2] - 2021-01-08

### Funciones

*   Se han agregado algunos filtros de consulta en las ejecuciones de trabajos [#878](https://github.com/distribworks/dkron/pull/878) @MGSousa
*   Ordenación de ejecuciones en la interfaz de usuario [#885](https://github.com/distribworks/dkron/pull/885)

### Fija

*   Solucionar el pánico en la recuperación del clúster (con peers.json) [#882](https://github.com/distribworks/dkron/pull/882) @fopina
*   Utilice el ajuste correcto y la fuente monoespacio para la salida de la ejecución [#879](https://github.com/distribworks/dkron/pull/879) @sc0rp10

## \[3.1.1] - 2020-12-21

### Funciones

*   Algunas tarjetas de información en el panel de control [#873](https://github.com/distribworks/dkron/pull/873)
*   Agregar estado a JobOption para filtrar desde la interfaz de usuario y la API [#872](https://github.com/distribworks/dkron/pull/872)

## \[3.1.0] - 2020-12-18

### Funciones

*   Interfaz de usuario web de administración de React [#864](https://github.com/distribworks/dkron/pull/864)

### Cambios

*   Usar la versión más reciente de la biblioteca gRPC [#855](https://github.com/distribworks/dkron/pull/855)
*   Bump deps

### Fija

*   Limpiar mensaje de registro [#860](https://github.com/distribworks/dkron/pull/860) @yvanoers
*   Quitar líneas duplicadas de la documentación de recuperación [#861](https://github.com/distribworks/dkron/pull/861) @vishalsngl
*   Cómo se tratan los errores en la llamada AgentRun [#858](https://github.com/distribworks/dkron/pull/858)
*   Arreglar el equilibrio desigual [#865](https://github.com/distribworks/dkron/pull/865) @yvanoers

## \[3.0.8] - 2020-11-20

### Cambios

*   Limpiar y pelar un poco de código [#853](https://github.com/distribworks/dkron/pull/853)
*   Mejores métricas [#852](https://github.com/distribworks/dkron/pull/852)

### Fija

*   Representar los tiempos de ejecución en la zona horaria del trabajo [#615](https://github.com/distribworks/dkron/pull/615) @yvanoers
*   Evite el puntero nulo si el trabajo se eliminó en ExecutionDone [#851](https://github.com/distribworks/dkron/pull/851)

## \[3.0.7] - 2020-11-10

### Cambios

*   Bump deps

### Fija

*   Bloqueo del servidor al agregar un nuevo trabajo [#840](https://github.com/distribworks/dkron/pull/840)
*   Método Fix de punto final ocupado en swagger.yaml [#843](https://github.com/distribworks/dkron/pull/843) @yvanoers
*   Corregir error de cardinalidad multietiqueta [#842](https://github.com/distribworks/dkron/pull/842) @yvanoers

## \[3.0.6] - 2020-10-15

### Cambios

*   Revertir "feat: Incluir entradas del programador en la API de estado" [#829](https://github.com/distribworks/dkron/pull/829)
*   Golpear algunos deps

## \[3.0.5] - 2020-09-04

### Cambios

*   Acción de Github para la versión (binario y docker) [#770](https://github.com/distribworks/dkron/pull/770)
*   Incluir entradas del programador en la API de estado [#813](https://github.com/distribworks/dkron/pull/813)
*   Bump deps [#814](https://github.com/distribworks/dkron/pull/814)

### Fija

*   s.Cron gratis en la oportunidad adecuada para evitar un bloqueo inesperado [#779](https://github.com/distribworks/dkron/pull/779)

## \[3.0.4] - 2020-06-12

### Fija

*   processFilteredNodes no devuelve nodos cuando la etiqueta especificada no tiene nodos [#785](https://github.com/distribworks/dkron/pull/785)

## \[3.0.3] - 2020-06-10

### Fija

*   Trabajo de inicio de registro y prefijo de registro en el agente grpc [#776](https://github.com/distribworks/dkron/pull/776)

### Funciones

*   clasificación de ejecución de busyhandler [#772](https://github.com/distribworks/dkron/pull/772)
*   Analizar direcciones de reintento-unir con plantillas sockaddr [#783](https://github.com/distribworks/dkron/pull/783)

### Cambios

*   Bump varias dependencias

## \[3.0.2] - 2020-05-15

### Fija

*   Corregir el enlace a la dirección de anuncio [#763](https://github.com/distribworks/dkron/pull/763)

## \[3.0.1] - 2020-05-12

### Funciones

*   Nuevo procesador que envía la salida del trabajo a un destino fluido [#759](https://github.com/distribworks/dkron/pull/759) @andreygolev
*   Tiempo de espera de reconexión del siervo configurable [#757](https://github.com/distribworks/dkron/pull/757) @andreygolev

### Fija

*   Corregir alertas JS [#762](https://github.com/distribworks/dkron/pull/762)

## \[3.0.0] - 2020-05-09

### Funciones

*   Agregar punto de enlace prometheus para métricas [#740](https://github.com/distribworks/dkron/pull/740)

### Fija

*   Parámetro de configuración de RaftMultiplier ignorado \[#753] (https://github.com/distribworks/dkron/pull/753)
*   Aumentar el tamaño del búfer de eventos de serf [#732](https://github.com/distribworks/dkron/pull/732)
*   Restablecimiento del estado y los siguientes parámetros [#730](https://github.com/distribworks/dkron/pull/730)

### Cambios

*   Actualizar deps y agregar nombre en clave a la versión [#751](https://github.com/distribworks/dkron/pull/751)
*   Mejores alertas con notificaciones de gruñidos [#750](https://github.com/distribworks/dkron/pull/750)
*   Refactorizar ejecutar trabajos [#749](https://github.com/distribworks/dkron/pull/749)
*   Agregar etiquetas de nombre de trabajo a eventos de registro para mejorar la depuración [#739](https://github.com/distribworks/dkron/pull/739)

## \[2.2.2] - 2020-04-22

### Fija

*   Aumentar el tamaño del búfer de eventos de serf [#732](https://github.com/distribworks/dkron/pull/732)
*   Restablecimiento del estado y los siguientes parámetros [#730](https://github.com/distribworks/dkron/pull/730)

### Cambios

*   Bump protobuf a 1.4.0 [#729](https://github.com/distribworks/dkron/pull/729)

## \[2.2.1] - 2020-04-15

### Cambios

*   Restaurar trabajos con archivo [#654](https://github.com/distribworks/dkron/pull/654) @vision9527
*   Actualizar deps [#724](https://github.com/distribworks/dkron/pull/724) [#725](https://github.com/distribworks/dkron/pull/725) [#726](https://github.com/distribworks/dkron/pull/726)

## \[2.2.0] - 2020-04-11

### Cambios

*   Dependencias bump
*   Cambiar el tipo de salida de ejecución de \[]byte -> cadena, esto funciona como necesitamos con JSON Marshal de Go

### Cambios de última hora

*   Ejecuciones de streaming: Implemente conexiones gRPC persistentes de agentes a servidor durante las ejecuciones, interfaz de complementos refactorizados para proporcionar la capacidad de transmitir la salida al servidor e implementar el nuevo `/busy` endpoint para mostrar las ejecuciones en ejecución. También se refactorizó el proceso de estado del trabajo, para simplificarlo eliminando el `running` estado, esto podría ser calculado por el usuario utilizando el `/busy` Extremo. (#716, #719, #720, #721, #723)

## \[2.1.1] - 2020-03-20

### Fija

*   Apagado elegante [#690](https://github.com/distribworks/dkron/pull/690) @andreygolev
*   Corrige el bloqueo cuando la configuración del complemento no está definida en un trabajo [#689](https://github.com/distribworks/dkron/pull/689) @andreygolev
*   Aplazar la corrección de pánico en ExecutionDone gRPC call [#691](https://github.com/distribworks/dkron/pull/691) @andreygolev

### Cambios

*   La configuración predeterminada iniciará y arrancará un servidor
*   Se ha agregado el controlador isLeader [#695](https://github.com/distribworks/dkron/pull/695)
*   Compilar con go 1.14
*   Equilibrio de carga de ejecución [#692](https://github.com/distribworks/dkron/pull/692) @andreygolev
*   Actualizar Bootstrap y JQuery [#700](https://github.com/distribworks/dkron/pull/700)
*   Actualizar todas las dependencias [#703](https://github.com/distribworks/dkron/pull/703)

### Cambios de última hora

*   Al disminuir el tamaño del plugin en un 75%, la interfaz de codificación de plugins refactorizada podría afectar el desarrollo de nuevos plugins y requerir adaptaciones para el plugin existente. [#696](https://github.com/distribworks/dkron/pull/696)
*   Usar BuntDB para el almacenamiento local, correcciones [#687](https://github.com/distribworks/dkron/issues/687), requieren una actualización continua. [#702](https://github.com/distribworks/dkron/pull/702) @andreygolev

## \[2.0.6] - 2020-02-14

### Fija

*   Consumo de memoria en el arranque [#682](https://github.com/distribworks/dkron/pull/682)

## \[2.0.5] - 2020-02-12

### Fija

*   Establecer el agente en ejecución dependiente [#675](https://github.com/distribworks/dkron/pull/675)
*   Devolver el código de estado correcto en el extremo principal [#671](https://github.com/distribworks/dkron/pull/671)

### Cambios

*   Comprobar si falta un agente [#675](https://github.com/distribworks/dkron/pull/675)
*   Agregar comentario de código [#675](https://github.com/distribworks/dkron/pull/675)

## \[2.0.4] - 2020-01-31

*   Elimine la dependencia del agente en la tienda y reduzca el uso en Job [#669](https://github.com/distribworks/dkron/pull/669)
*   Actualizar ginebra [#669](https://github.com/distribworks/dkron/pull/669)
*   Agregar métodos auxiliares [#669](https://github.com/distribworks/dkron/pull/669)
*   Mover la creación de directorios a la creación de instancias de la Tienda [#669](https://github.com/distribworks/dkron/pull/669)
*   Aceptar middlewares para rutas de API [#669](https://github.com/distribworks/dkron/pull/669)
*   Documentos de ACL

## \[2.0.3] - 2020-01-04

### Fija

*   Corregir la indexación modal en la interfaz de usuario [#666](https://github.com/distribworks/dkron/pull/666)

### Cambios

*   Bump BadgerDB a 2.0.1 [#667](https://github.com/distribworks/dkron/pull/667)
*   Permitir plantillas de formato de tiempo en notificaciones webhook [#665](https://github.com/distribworks/dkron/pull/665)

## \[2.0.2] - 2019-12-11

### Funciones

*   Buscar todos los trabajos en paneles con el cuadro de búsqueda [#653](https://github.com/distribworks/dkron/pull/653)

### Fija

*   Validar nombres de trabajo vacíos [#659](https://github.com/distribworks/dkron/pull/659)
*   Error de comunicación de die on plugin [#658](https://github.com/distribworks/dkron/pull/658)
*   Revertir GetStatus con simultaneidad prohibida [#655](https://github.com/distribworks/dkron/pull/655)

### Cambios

*   Actualice Angular a la versión más reciente [#641](https://github.com/distribworks/dkron/pull/641)

## \[2.0.1] - 2019-12-03

### Fijo

*   Trabajos secundarios que no se ejecutan debido a la ejecución del estado del trabajo [#651](https://github.com/distribworks/dkron/pull/651)

### Refactorizado

*   Eliminar cron lib copiado y agregar como dependencia [#646](https://github.com/distribworks/dkron/pull/646)

## \[2.0.0] - 2019-11-27

### Cambiado

*   Este archivo se basará en las directrices de Keep a Changelog

### Añadido

*   TBD

## 1.2.5

*   Corrección: actualizaciones de trabajo dependientes (@yvanoers)
*   Corrección: Reprogramación de trabajos en cada actualización de página
*   Corrección: espacio del paginador y agregar el primer y último botón
*   Corrección: los nuevos trabajos cuentan como fallidos en el panel (@yvanoers)
*   Característica: soporte TLS backend (@fopina)
*   Característica: soporte para autenticación con backend etcd (@fopina)
*   Corrección: ejemplo dkron.yml para la configuración de slack (bifurcación @kakakikikeke)

## 1.2.4

*   Actualizar la especificación de swagger: corregir executor_config, agregar estado
*   Nuevo diseño del sitio
*   Pruebas: Parametrizar el conjunto de pruebas para usar cualquier backend
*   Refactorizar: GetLastExecutionGroup para simplificar el código

## 1.2.3

*   Corrección: Bump valkeyrie con redis watches fix
*   Implementar la selección del servidor mediante hash coherente
*   Actualizar siervo a 0.8.2
*   refactorizar: La tienda debe implementar la interfaz

## 1.2.2

*   Corrección: Usar valkeyrie ramificada arreglando DynamoDB devuelve todos los elementos

## 1.2.1

*   Corrección: scheduler_started solución de expvar

## 1.2.0

*   Corrección: Actualizar el registro de errores de soporte del ejecutor en lugar de solo fallar (@tengattack)
*   Característica: Devuelve el siguiente campo de ejecución en la API y el panel.
*   Característica: Agregar marca de contraseña redis de backend (@lisuo3389)
*   Característica: Agregar token de cónsul backend
*   Mejora: Gráfico principal que muestra los trabajos en ejecución

**NOTA: El cambio de ruptura para los complementos de 3rd party, los complementos de ejecutores de la firma de la interfaz cambiados, deben ser recompilados y adaptados.**

## 1.1.1

*   Comando Agregar hoja RPC
*   Corrección: Falta tzdata en la imagen de lanzamiento

## 1.1.0

*   Agregar compatibilidad con DynamoDB
*   Interrupción: eliminar el mensaje de soporte y obsolescencia para los parámetros antiguos de comandos, environment_variables y shell
*   Cambiar scheduler_started expvar a int para que sea analizable
*   Varias mejoras en los documentos
*   Se ha corregido el error de sintaxis swagger.yaml
*   Agregar nueva línea antes de END OUTPUT para el procesador de registros

## 1.0.2

*   Permitir el envío de correo sin credenciales
*   Corregir el etiquetado de Docker
*   Corrección y mejoras del plugin de registro
*   Registro de uso de plugins de procesador más específico

## 1.0.1

*   No dockerignore dist carpeta, es necesario para gorelease docker builder
*   fa323c2 Omitir módulos de nodo
*   2475b37 Mover activos estáticos a su propio directorio dentro de la carpeta estática
*   987dd5d Eliminar hash de la url en el cierre modal
*   455495c Eliminar node_modules
*   0c02ce0 Regeneración de activos en subpaquetes

## 1.0.0

*   c91852b Consentimiento de cookies en el sitio web
*   9865012 No instalar herramientas de compilación
*   e280d31 Inicio de sesión de Docker
*   3229DC2 Garantizar la creación de binarios estáticos
*   01e62b6 Error en la prueba
*   69380f5 Ignorar archivos de sistema
*   A02a1ab Release con docker
*   e795210 Eliminar entradas antiguas de dockerignore
*   c9c692c Quitar unmarshalTags de dkron.go
*   C5F5DE0 Errores de informe en la configuración sin mayúsculas
*   62e1e15 Sumas para la liberación
*   1cf235a UnmarshalTags pertenece al agente y debe ser público
*   36f9318 Archivo Léame de actualización
*   80b2ab1 bandera de puerto de correo es uint

## 1.0.0-rc4

*   913ee87 Estructura de mapa de protuberancias
*   5bd120f Eliminar la carga de configuración heredada
*   f20fbe5 Actualizar mods

## 1.0.0-rc3

*   4811e48 Corregir interfaz de usuario ejecutar y eliminar
*   8695242 Redirigir al panel de control

## 1,0.0-rc2

*   d6dbb1a Agregar palanca a swagger
*   ffa4feb Vinculación profunda a vistas de trabajo
*   fdc5344 No te furiones en la marca
*   236b5f4 No consultar trabajos en intervalos en el panel
*   ea5e60b Corrige la reprogramación en la tienda boltdb
*   f55e2e3 Jardinería y enlaces de anclaje modal abierto
*   b22b362 Generación
*   6887c36 Información de registro
*   d21cf16 Información de registro y uso de la tienda. Tipo de back-end en lugar de cadenas en la configuración
*   28c130b Modalidad abierta con enlaces de anclaje y jardinería
*   1afb3df Se introdujeron varias correcciones de interfaz de usuario al migrar a glifos

## 1.0.0-rc1

*   ef86e13 Bump go-plugin
*   d09b942 Bump varias dependencias
*   f96d622 GRPC
*   8e3b4b9 Ignorar carpeta dist
*   1b7d4bc Plantilla de problema
*   caf4711 Logrus
*   5821c8c Principalmente etcd
*   33a12c5 Revertir "Bump varias dependencias"
*   fb9460d Actualizar cron-spec.md
*   706e65d Actualizar pflag

## 0.11.3

*   723326f Agregar registro para la respuesta de ejecuciones pendientes
*   df76e9c Agregar ejemplos reales a la especificación de swagger
*   d1318a1 Agregar etiquetas param a la especificación swagger
*   4da0b3b Refactorización de big docs
*   2d91a5e Interrupción de errores
*   Comando 8fac831 para generar documentos cli a partir de cobra
*   bdcd09c No uses swagger2markup
*   253fe57 ECS y correo electrónico pro docs
*   e89b353 Expvar dep
*   187190e Fijar sangría
*   c8320b5 Pruebas de corrección
*   9037d65 Corregir error tipográfico en los primeros pasos
*   Formato 9c60fe8
*   f11ed84 Formato
*   20be8e5 Integre swagger-ui para una visualización de API un poco mejor
*   2cede00 Fusionar rama 'master' en boltdb
*   53D8464 Sólo consulta para ejecuciones pendientes cuando hay alguna
*   712be35 Eliminar el bloqueo inútil adicional introducido en 88c072c
*   dacb379 Esto debería ser TrimSuffix
*   dec6701 Actualizar contactos
*   c21e565 Actualizar getting-started.md
*   3fdba5f Usar boltdb como almacenamiento predeterminado
*   70D9229 Guión incorrecto en el archivo de configuración de ejemplo
*   9653bbc expvars están de vuelta y son un punto final de salud simple

## 0.11.2

*   7d88742 Añadir código de conducta
*   aed2f44 Registro de depuración de serf adecuado
*   1226c93 Publicar docker
*   a0b6f59 Publicar docker
*   f1aaecc Reorg importaciones
*   8758bac Las pruebas deben usar etcdv3
*   Fa3aaa4 Las pruebas deben utilizar el cliente v3
*   5Bcea4c Actualizar crear o actualizar punto de conexión de api de trabajo
*   39728d0 refactor: nombre methond
*   Refactor 1c64da4: registro y modo de ginebra adecuados

## 0.11.1 (2018-10-07)

*   Agregar compatibilidad para pasar la carga útil al comando STDIN (@gustavosbarreto)
*   agregar soporte para etcdv3 (@kevynhale)
*   Usar etcdv3 de forma predeterminada
*   Trabajos de URL estáticas fijas

## 0.11.0 (2018-09-24)

*   1.11 estable aún no está en docker hub
*   Añadir plugin http incorporado
*   Add executor shell su option (@tengattack)
*   Mejor dockerfile para pruebas
*   Mejor ayuda para marcar
*   No dependas de michellh/cli
*   Filtrar trabajos por etiquetas (@digitalcrab)
*   Corregir el error de pánico del clúster (@tengattack)
*   Lanzamiento con goreleaser
*   Usa cobra para banderas
*   Usar módulos go
*   agregar funciones de creación y actualización de trabajos (@wd1900)

## 0.10.4 (2018-07-30)

*   Reemplazar RPC por gRPC
*   Reparar archivos de redacción (@kevynhale)

## 0.10.3 (2018-06-20)

*   Reemplazar goxc con makefile
*   Documentos profesionales

## 0.10.2 (2018-05-23)

### Correcciones

*   Comprobación de estado de corrección
*   Eliminar actualizaciones innecesarias de los tiempos de finalización del trabajo (@sysadmind)
*   Reflejar el estado de la tienda en la API
*   Arreglar plugins de Windows (@sysadmind)
*   Detener la actualización del trabajo en el error de análisis JSON en la API (@gromo)

## 0.10.1 (2018-05-17)

### Correcciones

*   Corregir modales de vista/eliminación de trabajos de panel

## 0.10.0 (2018-05-14)

### Correcciones

*   Corregir la dirección de falta de consulta RPCconfig (#336 y relacionada)
*   Solucionar el problema de simultaneidad debido a la condición de carrera en los trabajos de bloqueo [#299](https://github.com/distribworks/dkron/pull/299)
*   Corregir la ejecución que falta al reiniciar el bloqueo de la simultaneidad prohíbe los trabajos [#349](https://github.com/distribworks/dkron/pull/349)
*   Arreglar rutas de carga del plugin [#275](https://github.com/distribworks/dkron/pull/275)
*   Corregir la dirección RPC perdida al recargar la configuración [#262](https://github.com/distribworks/dkron/pull/262)

### Características y mejoras en el código

*   Mejorar ligeramente el procesamiento del último grupo de ejecución (@sysadmind)
*   Mejorar el manejo de dependencias laborales (@sysadmind)
*   Mover el comando dkron a su propio paquete
*   Creación o actualización de trabajos de API de rango de milisegundos
*   Reinicio del programador de refactorización
*   Reemplace bower con npm
*   Plugins de ejecutor basados en GRPC
*   Alternar trabajo desde la interfaz de usuario
*   Buscar trabajo por nombre y paginación en la interfaz de usuario
*   Agregar redis como back-end de almacenamiento (@digitalcrab)
*   Refactorización de la interfaz de usuario con la nueva versión de bootstrap y reemplazo de fontawesome con glifos
*   Calcular el estado del trabajo y devolver el valor de la API que proporciona al usuario más información
*   Programación consciente de la zona horaria (@didiercrunch)

## 0.9.8 (2018-04-27)

*   Arreglar la versión rota 0.9.7

## 0.9.7 (2018-02-12)

*   Registro de plugins menos detallado
*   Actualizar dep osext roto (@ti)
*   Cambiar de libkv a valkeyrie
*   Refactorizar para el código principal utilizable
*   Corregir grupos de ejecución no ordenados (@firstway)
*   Corregir GetLastExecutionGroup (@firstway)

## 0.9.6 (2017-11-14)

*   Migrar de glide a dep
*   Fijar precedencia de parámetros, parámetros cli en la parte superior
*   Conjunto de pruebas más robusto
*   Registro de ginebra a registrador común
*   Mejor script systemd
*   No entre en pánico ni fatal al enviar notificaciones
*   Actualización de siervos
*   Corregir el cambio de ruptura de plantillas en la actualización de Go 1.9

## 0.9.5 (2017-09-12)

Funciones

*   Nuevo sitio web de documentos usando hugo

Correcciones:

*   Limpiar clientes con una señal de salida (@danielhan)
*   Solución #280 (@didiecrunch)
*   Actualizar varias dependencias
*   Corregir la ruta relativa de los activos estáticos

## 0.9.4 (2017-06-07)

*   Corregir errores en documentos de API
*   El uso de "trabajos", "1 am" o "1 pm" en el nombre del trabajo conduce a un error en el tablero
*   Corregir el bloqueo en el nombre del complemento inexistente
*   Incrustar todos los activos en binario, eliminado -ui-dir config param

### Notas de actualización

Este es un cambio rompedor; `ui-dir` se ha eliminado el parámetro de configuración, todas las secuencias de comandos que utilizan este parámetro deben actualizarse.

## 0.9.3 (2017-02-12)

*   Corregir la dirección de escucha del servidor RPC (@firstway)
*   Implementación básica de la infraestructura de pruebas utilizando swarm
*   Implementación de telemetría básica para el envío de métricas a statsd y datadog
*   Corregir un bloqueo en una falla de backend
*   Ejecuciones de ordenación inversa en la interfaz de usuario (@Eyjafjallajokull)

## 0.9.2 (2016-12-28)

Funciones:

*   Implementar la directiva de simultaneidad
*   Interfaz de usuario mejorada: permitir eliminar trabajos de la interfaz de usuario, resaltar JSON
*   Plugins del procesador de ejecución, permite el enrutamiento flexible de los resultados de ejecución
*   Variables de plantilla para la personalización de correos electrónicos de notificación (@oldmantaiter)
*   Ir 1.7
*   Pruebe con docker-compose, esto permitirá probar múltiples tiendas fácilmente

Correcciones:

*   Corregir pruebas que fallan aleatoriamente (@oldmantaiter)
*   Devolver lista vacía cuando no hay trabajos en lugar de null
*   Permitir el uso de POST en el método /leave, dejar en desuso GET

## 0.9.1 (2016-09-24)

Correcciones:

*   Corregir estadísticas de trabajo que no se actualizan #180
*   Fix zookeeper obtener lista de ejecuciones #184
*   Corregir un bloqueo al eliminar un trabajo que no existe #182
*   Arreglar Travis en bifurcaciones

## 0.9.0 (2016-08-24)

Funciones:

*   Admite trabajos de cualquier tamaño
*   Apoyar trabajos encadenados
*   Validación de programación y otras propiedades de trabajo
*   Nuevo diseño de sitio, logotipo y tablero

Correcciones:

*   Corregir reintentos de ejecución
*   Corregir ejecuciones combinadas por el mismo prefijo
*   Corregir el estado HTTP correcto en la creación/actualización

## 0.7.3 (2016-07-12)

*   Trabajos únicos
*   Se ha añadido una especificación cron a los documentos
*   Reintento de ejecución en caso de error
*   Cambie la especificación del esquema JSON para su especificación de API abierta correspondiente
*   Recargar configuración
*   Corregir error de programación
*   El estado del nuevo trabajo proporciona más información sobre la ejecución del trabajo

## 0.7.2 (2016-06-01)

*   Agregar algunos ayudantes y correcciones de errores
*   Agregar la propiedad shell al trabajo, reintrodujo el método de ejecución del shell, pero ahora es una opción
*   Agregar nodo de informes a informes de ejecución
*   Reemplace la etiqueta del servidor por dkron_server y agregue dkron_version

### Notas de actualización

Debido al cambio en las etiquetas internas `server` Para `dkron_server`, deberá ajustar las etiquetas de trabajo si estaban usando esa etiqueta.

## 0.7.1 (2016-05-03)

*   No use la llamada de shell al ejecutar comandos, explotando la línea de comandos.
*   Añadir publicidad, añadir `advertise` opción que resuelve la unión entre hosts al ejecutar Docker
*   Validar el tamaño del trabajo, limitar el tamaño máximo del siervo
*   Sobrescribir trabajos, ahora enviar trabajos existentes no sobrescribe campos no existentes en la solicitud
*   Solución para el bloqueo del panel en un líder inexistente

## 0.7.0 (2016-04-10)

*   Refactorizar la elección del líder, el método antiguo podría conducir a casos en los que 2 o más nodos podrían tener el programador en ejecución sin darse cuenta del otro maestro.
*   Deshágase de `keys`, en un clúster de siervos, los nombres de nodo son únicos, por lo que ahora se usan para claves líderes.
*   Arreglar [#85](https://github.com/distribworks/dkron/issues/85) Reiniciar el programador en la eliminación de trabajos
*   Refactorizar el registro, reemplazar `debug` con `log-level`
*   Ordenar nodos en la interfaz de usuario [#81](https://github.com/distribworks/dkron/issues/81) (felicitaciones @whizz)
*   Agregue vars expuestos a una fácil depuración
*   Ir 1.6
*   Agregar @minutely como programación predefinida (felicitaciones @mlafeldt)

### Actualizar desde 0.6.x

Para actualizar una instalación existente, primero debe eliminar la clave directriz de salida previa del almacén. La clave del líder está en forma de: `[keyspace]/leader`

## 0.6.4 (2016-02-18)

*   Usar expvars para exponer métricas
*   arreglar https://github.com/distribworks/dkron/issues/71
*   Mejor configuración de ejemplo en paquetes y documentos

## 0.6.3 (2015-12-28)

*   INTERFAZ DE USUARIO: Mejor vista del trabajo
*   Lógica para almacenar solo las últimas 100 ejecuciones

## 0.6.2 (2015-12-22)

*   Fijo [#62](https://github.com/distribworks/dkron/issues/55)

## 0.6.1 (2015-12-21)

*   Errores corregidos [#55](https://github.com/distribworks/dkron/issues/55), [#52](https://github.com/distribworks/dkron/issues/52)etc.
*   Compilación para linux arm

## 0.6.0 (2015-12-11)

*   Algunas otras mejoras y corrección de errores
*   Vendoring ahora usando go vendor experiment + glide
*   Corrección: Eliminar ejecuciones en la eliminación de trabajos
*   Mostrar la salida de ejecución completa en el modal de la interfaz de usuario
*   Nuevas ejecuciones resultados internos usando RPC
*   Registro estandarizado
*   Mostrar información sobre herramientas de trabajo con información
*   Acepte solo "bonito" para formatear solicitudes de API
*   Cambie el funcionamiento de los grupos de ejecución para no utilizar el concepto de directorio.

## 0.5.5 (2015-11-19)

*   Más compatibilidad con backend
*   Acepte solo bonito para formatear solicitudes de API
*   Mostrar ejecuciones agrupadas en la interfaz de usuario web
*   Mostrar información sobre herramientas de trabajo con toda la información JSON del trabajo en la interfaz de usuario web
*   Mejores alertas

## 0.5.4 (2015-11-17)

*   Corrección de rutas de acceso de interfaz de usuario web

## 0.5.3 (2015-11-16)

*   La interfaz de usuario web funciona detrás del proxy http

## 0.5.2 (2015-11-09)

*   Corrige un error en el parámetro join config que lo hacía inutilizable desde el archivo de configuración.

## 0.5.1 (2015-11-06)

*   Paquete Deb
*   Libkv actualizado a la última
*   Nuevas opciones de configuración (nivel de registro, ruta de interfaz de usuario web)

## 0.5.0 (2015-09-27)

*   Notificaciones configurables por correo electrónico y Webhook para ejecuciones de trabajos.
*   Capacidad para cifrar el tráfico de red de siervos entre nodos.
*   Respuestas de API de formato bonito
*   La interfaz de usuario ahora muestra el estado de ejecución con codificación de colores y ejecución parcial.
*   Más estabilidad y previsibilidad de la API
*   Esquema JSON de API proporcionado, documentos de API generados basados en el esquema
*   Probado en Travis
*   El uso de Libkv permite utilizar diferentes backends de almacenamiento (etcd, cónsul, zookeeper)
*   Agregar el control de versiones v1 a las rutas de la API

## 0.0.4 (2015-08-23)

*   Compilado con Go 1.5
*   Incluye la vista de nodos de clúster en la interfaz de usuario

## 0.0.3 (2015-08-20)

*   Versión inicial
