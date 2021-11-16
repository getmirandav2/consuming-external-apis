# Consuming external APIs

Este proyecto consta de dos endpoints:

- `/repository`
- `/repositories`

## Estructura general

### app

Es donde se inicializa el proyecto y  las rutas.

### config

Se definen todas las variables de entorno.

### controllers

Son todas las funciones que estan enrutadas con un endpoint, se encargan de validar el request, llamar a un `service`, retornar un error (si es que hay) o retornar la respuesta.

### domain

Se definen los modelos que ocupará el `provider`, incluyen requests, responses, responses de errores y modelos auxiliares. Esta capa tambien se utiliza para realizar consultas a una base de datos tipo SDK.

### providers

Esta capa se encarga de interactuar con APIs externas. No es necesario que esta capa implemente las funciones definidas en una intercace, ya que el unico objeto que se debe mockear es el cliente http.