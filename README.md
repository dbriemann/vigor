# vigor

A framework extending ebitengine which more batteries:

## Possible topics for implementation

- a resouce management system including json serialization,
- spritesheet and animation utilities,
- effects (WIP),
- TODO: collision
- TODO: camera
- TODO: input management
- TODO: particles
- TODO: localization
- TODO: cutscenes
- TODO: state management
- TODO: tweening

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
