![Header](Imagenes-readme/header-tarjetas.jpg)

# BDD Tarjetas

BDD Tarjetas es un proyecto realizado para la materia Gestión y administración de bases de datos del tercer año de la carrera Licenciatura en sistemas de la Universidad Nacional General Sarmiento.

## Tabla de contenidos

- [Introducción](#introducción)
- [Descripción de las tablas principales](#descripción-de-las-tablas-principales)
  - [Tarjeta](#tarjeta)
  - [Cliente](#cliente)
  - [Compra](#compra) 
- [Stored Procedures y triggers](#stored-procedures-y-triggers)
- [JSON y bases de datos NoSQL](#json-y-bases-de-datos-nosql)
- [Logros](#logros)


## Introducción

El proyecto consiste en modelar una base de datos con información de tarjetas de crédito, con sus clientes y sus respectivas compras. Este modelo se usará en la implementación de dos aplicaciones CLI en Go.

## Descripción de las tablas principales

### Tarjeta

El sistema debe contar con un registro de compras realizadas con cada tarjeta. Cada tarjeta le pertenece a un cliente (a excepción de dos clientes que tendrán dos). Los clientes realizan compras en comercios a través de su tarjeta. Las tarjetas de crédito son el medio para realizar una compra. Las tarjetas de crédito no tienen extensiones. Los usuarios no tienen permitido financiar sus compras en cuotas, todo en un solo pago. Una tarjeta de crédito puede ser suspendida si recibe dos rechazos por límite de compra en un mismo día. Además, las tarjetas presentan las siguientes características:

- Número de tarjeta.
- Número de cliente.
- Tiempo inicial y final de validez.
- Código de seguridad.
- Límite de compra.
- Estado de la tarjeta ("vigente", "suspendida", "anulada").

### Cliente

Los clientes realizan compras con la tarjeta. El cliente no puede tener más de una tarjeta. Un cliente es alertado por posibles fraudes tales como, realizar una compra en un lapso de tiempo muy bajo en diferentes ubicaciones (menos de un minuto si es en el mismo código postal, cinco minutos si es en diferentes códigos postales), también un cliente es alertado si recibe dos rechazos de límite de compra en un mismo día, seguido de una suspensión preventiva de la tarjeta. Un cliente tiene las siguientes características:

- Número de cliente.
- Nombre y apellido.
- Domicilio.
- Teléfono.

### Compra

La compra se encarga de guardar los datos que implican dicha acción. Los datos son:

- Número de operación.
- Número de tarjeta.
- Número de comercio.
- Fecha de realización.
- Monto.
- Estado de la compra ("pagado" o "no pagado").

## Stored Procedures y triggers

El sistema fue construido con los lenguajes Go y Pl/pgSQL. El modelo de datos relativa a tarjetas de crédito son almacenadas en la base de datos relacional PostgreSQL. Por último, en BoltDB se guardan datos de clientes, comercios, tarjetas y compras para comparar el modelo relacional con un modelo no relacional (NoSQL).

Las funciones creadas tienen como objetivo vincular al usuario con la base de datos. La forma de interactuar es a través de una interfaz de líneas de comandos.

- autorizar_compra

Recibe los datos de una posible compra (número de tarjeta, código de seguridad, número de comercio, monto a pagar). La función retorna <i>true</i> si se aprueba y <i>false</i> si ocurrió lo contrario.

Para autorizar una compra la función controla que se cumplan los siguientes requisitos: La tarjeta debe existir y debe estar en vigencia. La compra no debe superar el límite de compra (para eso se le suma las anteriores compras). La tarjeta no debe estar suspendida, ni encontrarse anulada, para eso agregamos una condición más que verifica ese estado.

Ésta función utiliza un método auxiliar que verifica la vigencia de las tarjetas y en caso que devuelva <i>false</i> no se autoriza la compra.

- _func_generar_resumen

Recibe los datos del cliente con su período del año y genera un resumen con todas las compras y el total a pagar.

- func_alerta_rechazo

Esta función es ejecutada cuando se genera un rechazo al autorizar la compra. Se encarga de registrar el rechazo en una tabla de alertas. Si un cliente tuvo dos rechazos por superar el límite de compra en un día, la función establece una suspensión de la tarjeta seguido de una alerta. Se implementa usando un trigger.

Esta función es ejecutada por el trigger rechazo_trig.

- func_alerta_compra

Es ejecutada cuando se realiza una compra. Controla que no se realicen dos compras en un lapso menor a un minuto cuando es dentro del mismo código postal y eb un lapso de cinco minutos cuando es fuera del código postal. En caso de que se cumpla, la función registra la alerta. Se implemtenta usando un trigger.

## JSON y bases de datos NoSQL

Para guardar valores en los buckets primero se crea variables de tipo struct con las diferentes entidades y con sus respectivos atributos. Luego esos valores son pasados a formato JSON para que en el próximo paso sean guardados en buckets. En los buckets van el elemento JSON, el nombre del bucket y una clave para identificar.

Una de las diferencias entre los modelos de datos SQL y NoSQL radica en la forma de hacer consultas. En NoSQL puede ser más complejo agrupar elementos porque no hay relación. La forma de hacer consultas es a través de la clave y el valor.

Otra diferencia es la manera de crear estructuras de datos. Por un lado tenemos la base de datos SQL donde se crean tablas que pueden estar relacionadas con otras. Y por el otro lado, la base de datos NoSQL se modela con elementos clave/valor.

## Logros

Este proyecto fue muy importante para profundizar e implementar los conceptos sobre la gestión y administración de una base de datos: la creación de consultas, usar stored procedures y triggers, y la gestión de tiempo y responsabilidades en el equipo.
