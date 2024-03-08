# vigor

A framework for building 2d games. Built on top of `ebitengine`, batteries included but keeping a simple interface.
Inspired by ``HaxeFlixel``.

## Possible topics for implementation

- a resouce management system including json serialization,
- spritesheet and animation utilities,
- effects (WIP),

### higher priority

- TODO: camera
- TODO: particles
- TODO: localization
- TODO: cutscenes
- TODO: state management
- TODO: debug mode showing wireframes and internal infos in running game

### lower priority

- TODO: improve collisions for moving objects
- TODO: input is currently `ebitengine-input`: we might internalize or replace with own system
- TODO: tweening is currently `ganema/tween`: we might internalize or replace with own system
- TODO: thread safety

## Open questions

- How to best handle `dt` in update functions. Pass in or handle implicitly inside of lib? How?

## Architecture

### Entities

- Game: a wrapper for ebiten Game
- State
- Stage
- Scene
- Camera
- Sprite
- Group
