# World

## Introducción

Hola! me he animado a abrir un devlog porque estoy motivado con este proyecto y espero llevarlo lejos!

La idea es hacer un juego (MMORPG ofc) del estilo de `Pocket Pioneers`. Quiero mantener el sistema de casillas, tanto para facilitar el editor de mundo como para poder usar movimiento con touch, enfocado en dispositivos móviles (me ha parecido una forma muy comoda de interactuar con todo, al contrario de los switches digitales a los que estamos acostumbrados), también le estoy dando vueltas a la mejor forma de implementar un sistema de combate.

El primer paso de mi idea es crear los sistemas necesarios para después poder construir el juego en sí sobre ellos, ya que la mayoría de los sistemas serán usables en runtime (editor de escenarios, NPCs, elementos interactivos, etc).

Cuando esté la base lista tengo una idea de las características que quiero que tenga el juego, pero no tengo en mente ninguna temática.

Tenéis todo el código subido a Github https://github.com/kidandcat/godot-world. Licencia MIT (mi propósito es divertirme y si alguien le saca partido a esto, pues yo que me alegro)

## Tecnologías

Godot 3.1 para el cliente. El servidor está desarrollado en Go, la comunicación es mediante WebSockets. La base de datos es BuntDB (https://github.com/tidwall/buntdb), una DB en memoria (persistencia con log) con indexado espacial y transacciones.

## Capturas

De momento no hay mucho que enseñar (me da pereza grabar un video). Por otro lado mi especialidad son redes y sistemas en tiempo real, asi que probablemente no vais a ver grandes avances gráficos (mi idea es encontrar/encargar un pack de assets o conseguir a alguien que se apunte al proyecto para hacer esa parte). 

[img]https://i.imgur.com/qSFXrEa.png[/img]
[img]https://i.imgur.com/YBtkM43.png[/img]
[img]https://i.imgur.com/YBtkM43.png[/img]

## Roadmap

- [x] Server: Conexión e inicialización godot<->servidor
- [x] Mesh: Sistema de casillas
- [x] Mesh: Mundo cargado desde servidor
- [x] Mesh: Mundo editable desde cliente con sincronizacion en tiempo real a otros clientes
- [x] Mesh: Mundo infinito (carga/descarga del mapa según se mueve el jugador)
- [x] Movimiento: Pathfinding del lado del servidor
- [ ] Editor: menú de selección de modelos disponibles
- [ ] Editor: editar caracteristicas de celdas desde cliente (walkable, walkCost, etc)
- [ ] Movimiento: filtrar meshes no walkables al hacer el pathfinding
- [ ] Cliente: Skybox, luces, etc (server side¿?)
- [ ] Editor: preview de mesh
- [ ] Editor: rotacion de meshes
- [ ] Editor: multiples niveles en el eje Y
- [ ] Editor: interfaz móvil (de momento hay que usar el ratón para poder crear meshes)