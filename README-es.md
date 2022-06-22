<p align="center">
<img width="400" src="docs/images/DKRON_STICKER_OK_CMYK_RGB_CONV_300.png" alt="Dkron" title="Dkron" />
</p>

# Dkron - Sistema de programación de trabajos distribuido y tolerante a fallos para entornos nativos de la nube [![GoDoc](https://godoc.org/github.com/distribworks/dkron?status.svg)](https://godoc.org/github.com/distribworks/dkron) [![Actions Status](https://github.com/distribworks/dkron/workflows/Test/badge.svg)](https://github.com/distribworks/dkron/actions) [![Gitter](https://badges.gitter.im/distribworks/dkron.svg)](https://gitter.im/distribworks/dkron)

Sitio web: http://dkron.io/

Dkron es un servicio cron distribuido, fácil de configurar y tolerante a fallos con enfoque en:

*   Fácil: Fácil de usar con una gran interfaz de usuario
*   Fiable: Completamente tolerante a fallos
*   Alta escalabilidad: Capaz de manejar grandes volúmenes de trabajos programados y miles de nodos

Dkron está escrito en Go y aprovecha el poder del protocolo Raft y Serf para proporcionar tolerancia a fallas, confiabilidad y escalabilidad mientras se mantiene simple y fácilmente instalable.

Dkron se inspira en el documento técnico de Google [Cron confiable en todo el planeta](https://queue.acm.org/detail.cfm?id=2745840) y por Airbnb Chronos tomando prestadas las mismas características de él.

Dkron se ejecuta en Linux, OSX y Windows. Se puede utilizar para ejecutar comandos programados en un clúster de servidores utilizando cualquier combinación de servidores para cada trabajo. No tiene puntos únicos de falla debido al uso del protocolo Gossip y bases de datos distribuidas tolerantes a fallas.

Puede usar Dkron para ejecutar la parte más importante de su empresa, los trabajos programados.

## Instalación

[Instrucciones de instalación](https://dkron.io/basics/installation/)

La documentación completa y completa se puede ver en el [Sitio web de Dkron](http://dkron.io)

## Desarrollo Inicio rápido

La mejor manera de probar y desarrollar dkron es usando docker, necesitará [Estibador](https://www.docker.com/) instalado antes de proceder.

Clonar el repositorio.

A continuación, ejecute la configuración de Docker Compose incluida:

`docker-compose up`

Esto iniciará las instancias de Dkron. Para agregar más instancias de Dkron a los clústeres:

    docker-compose up --scale dkron-server=4
    docker-compose up --scale dkron-agent=10

Compruebe la asignación de puertos mediante `docker-compose ps` y utilice el navegador para navegar al panel de control de Dkron utilizando uno de los puertos asignados por componer.

Para agregar trabajos al sistema, lea el botón [Documentos de API](https://dkron.io/api/).

## Desarrollo frontend

El panel de control de Dkron se construye utilizando [React Admin](https://marmelab.com/react-admin/) como una aplicación de una sola página.

Para comenzar a desarrollar el panel, ingrese al `ui` directorio y ejecución `npm install` Para obtener las dependencias front-end y, a continuación, iniciar el servidor local con `npm start` debe iniciar un nuevo servidor web local y abrir una nueva ventana del navegador que sirve de la interfaz de usuario web.

Realice los cambios en el código y, a continuación, ejecute `make ui` para generar archivos de activos. Este es un método para incrustar recursos en aplicaciones Go.

### Recursos

Libro de cocina del chef
https://supermarket.chef.io/cookbooks/dkron

Biblioteca de cliente de Python
https://github.com/oldmantaiter/pydkron

Cliente Ruby
https://github.com/jobandtalent/dkron-rb

Cliente PHP
https://github.com/gromo/dkron-php-adapter

Proveedor de Terraform
https://github.com/bozerkins/terraform-provider-dkron

Gestiona y ejecuta trabajos en Dkron desde tu proyecto django
https://github.com/surface-security/django-dkron

## Ponte en contacto con nosotros

*   Twitter: [@distribworks](https://twitter.com/distribworks)
*   Chat: https://gitter.im/distribworks/dkron
*   Correo electrónico: victor en distrib.works

# Patrocinador

Este proyecto es posible gracias al apoyo de Jobandtalent

![](https://upload.wikimedia.org/wikipedia/en/d/db/Jobandtalent_logo.jpg)
