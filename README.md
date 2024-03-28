# vigor

A framework for building 2d games. Built on top of `ebitengine`, batteries included but keeping a simple interface.
Inspired by `HaxeFlixel`.

## Features

- a resource management system including JSON serialization,
- spritesheet and animation utilities, including tweening
- particle emitter (WIP)
- effects (WIP),
- collision detection for objects
- input management (via `ebitengine-input`)

### higher priority

- TODO: camera
- TODO: localization
- TODO: cutscenes
- TODO: state management
- TODO: debug mode showing wireframes and internal infos in running game

### lower priority

- TODO: improve collisions for moving objects
- TODO: input is currently `ebitengine-input`: we might internalize or replace with own system
- TODO: tweening is currently `ganema/tween`: we might internalize or replace with own system | allow tweening vec2d
- TODO: thread safety

## Architecture

### Entities

