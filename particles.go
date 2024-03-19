package vigor

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type particle struct {
	*Image

	ttl float32
}

type Emitter struct {
	particles []*particle
	size      int
	ptype     Image
	active    bool
	toSpawn   float32
	id        uint64

	origin   Vec2[float32] // origin is the source point for all emitted particles.
	capacity int           // capacity defines the max amount of particles.
	rate     int           // rate defines the emitted particles per second.
	angle    Vec2[float32] // angle is a range of the min and max angle (rad) defining particle directions.
	lifetime Vec2[float32] // lifetime is a range defining the min and max age of the particles.
	speed    Vec2[float32] // speed is a range defining the min and max velocity of the particles.
}

// TODO: pass config struct to NewParticleEmitter instead of single values
// and assume sane defaults for all not given.
func NewParticleEmitter(img Image, x, y float32, cap, rate int) *Emitter {
	e := &Emitter{
		id:       G.createId(),
		origin:   Vec2[float32]{X: x, Y: y},
		capacity: cap,
		size:     0,
		rate:     rate,
		angle:    Vec2[float32]{X: 0, Y: 2 * math.Pi},
		lifetime: Vec2[float32]{X: 0.5, Y: 1},
		speed:    Vec2[float32]{X: 25, Y: 75},

		ptype:     img,
		particles: make([]*particle, cap),
		active:    false,
		toSpawn:   0,
	}

	for i := 0; i < e.capacity; i++ {
		e.particles[i] = &particle{
			Image: CopyImage(&e.ptype),
			ttl:   0,
		}
	}

	return e
}

func (e *Emitter) Burst() {
	for spawned := e.spawn(); spawned; spawned = e.spawn() {
	}
}

func (e *Emitter) SetOrigin(x, y float32) {
	e.origin.X = x
	e.origin.Y = y
}

func (e *Emitter) SetParticleType(img Image) {
	e.ptype = img
}

func (e *Emitter) Id() uint64 {
	return e.id
}

func (e *Emitter) Visible() bool {
	return e.active
}

func (e *Emitter) Show(v bool) {
	e.active = v
}

func (e *Emitter) Update() {
	if !e.active {
		return
	}

	// Spawn new particles.
	e.toSpawn += float32(e.rate) * G.Dt()
	amount := int(e.toSpawn)
	for i := 0; i < amount; i++ {
		if !e.spawn() {
			break
		}
		e.toSpawn -= 1
	}

	// Remove dead particles by swapping with last "good" one.
	for i := 0; i < e.size; i++ {
		p := e.particles[i]
		p.ttl -= G.Dt()
		if p.ttl <= 0 {
			e.size--
			swp := e.particles[e.size]
			e.particles[e.size] = e.particles[i]
			e.particles[i] = swp
			i--
		} else {
			// This one is still active and receives an update.
			e.particles[i].Update()
		}
	}
}

func (e *Emitter) ActiveParticles() int {
	return e.size
}

func (e *Emitter) draw(target *ebiten.Image) {
	for i := 0; i < e.size; i++ {
		e.particles[i].draw(target)
	}
}

func (e *Emitter) spawn() bool {
	if e.size >= e.capacity {
		return false
	}

	p := e.particles[e.size]
	e.size++
	p.ttl = 0
	p.SetPos(e.origin.X, e.origin.Y)

	// Find random speed within range.
	speed := rand.Float32()*e.speed.Y + e.speed.X

	// Find random angle to emit. No rotation means right (1,0).
	rad := rand.Float32()*e.angle.Y + e.angle.X

	// Rotate default direction with random angle and apply speed.
	dir := Vec2[float32]{X: 1, Y: 0}
	rot := dir.Rotate(float64(rad))
	vec := rot.Multiply(speed)
	// And set the according velocity
	p.SetVel(vec.X, vec.Y)

	// Set random lifetime.
	lifetime := rand.Float32()*e.lifetime.Y + e.lifetime.X
	p.ttl = lifetime

	// TODO: accel
	// TODO: rotation
	return true
}
