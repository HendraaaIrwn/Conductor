package entity

import (
	"sort"
	"time"

	"github.com/example/conductor/internal/render"
)

// Manager owns all entities in a scene. It supports adding and removing
// entities, updating them each frame, culling expired/off-screen entities, and
// rendering them in layer order.
//
// To avoid mutating the entity slice while iterating, additions and removals
// are deferred to the next flush. This means an entity spawned during Update
// will not be updated until the next frame, and a removed entity will not
// disappear until after the current frame's cleanup pass.
type Manager struct {
	entities      []*Entity
	pendingAdd    []*Entity
	pendingRemove map[string]bool

	nextID int
}

// NewManager creates an empty Manager.
func NewManager() *Manager {
	return &Manager{
		pendingRemove: make(map[string]bool),
	}
}

// Add queues an entity for addition. The entity is assigned a unique ID if it
// does not already have one. The entity will be inserted into the active list
// during the next flush.
func (m *Manager) Add(e *Entity) {
	if e.ID == "" {
		m.nextID++
		e.ID = autoID(m.nextID, e.Type)
	}
	m.pendingAdd = append(m.pendingAdd, e)
}

// Remove queues an entity for removal by ID. The entity will be removed during
// the next flush.
func (m *Manager) Remove(id string) {
	m.pendingRemove[id] = true
}

// Count returns the number of active entities (excluding pending additions).
func (m *Manager) Count() int { return len(m.entities) }

// Flush integrates pending adds and removes into the active entity list. This
// is called automatically during Update, but can be called manually when
// entities need to be queryable immediately after Add (e.g. in a constructor).
func (m *Manager) Flush() {
	m.flush()
}

// ByType returns all active entities matching the given type. The returned
// slice is a copy and can be safely iterated.
func (m *Manager) ByType(t Type) []*Entity {
	var result []*Entity
	for _, e := range m.entities {
		if e.Type == t {
			result = append(result, e)
		}
	}
	return result
}

// FirstByType returns the first active entity matching the given type, or nil.
func (m *Manager) FirstByType(t Type) *Entity {
	for _, e := range m.entities {
		if e.Type == t {
			return e
		}
	}
	return nil
}

// Update advances all entities by the given delta. It:
//  1. Flushes pending adds/removes.
//  2. Calls each entity's Behavior.
//  3. Increments Age.
//  4. Marks expired and off-screen entities for removal.
//  5. Flushes pending adds/removes again (behaviors may have spawned/removed).
func (m *Manager) Update(width, height int, delta time.Duration) {
	// Always flush pending adds/removes so entities are queryable even
	// when delta is zero (paused).
	m.flush()

	if delta <= 0 {
		return
	}

	ctx := UpdateContext{
		Width:  width,
		Height: height,
		Spawn:  m.Add,
	}

	for _, e := range m.entities {
		if e.Dead || !e.Visible {
			continue
		}
		if e.Behavior != nil {
			e.Behavior.Update(e, ctx, delta)
		}
		e.Age += delta

		if e.Lifetime > 0 && e.Age >= e.Lifetime {
			e.Dead = true
			continue
		}
		if e.RemoveOffscreen && isOffscreen(e, width, height) {
			e.Dead = true
		}
	}

	m.flush()
}

// Render draws all visible entities onto the canvas in ascending layer order.
// Entities on the same layer are drawn in insertion order.
func (m *Manager) Render(canvas *render.Canvas) {
	sorted := m.sortedByLayer()
	for _, e := range sorted {
		if !e.Visible || e.Dead {
			continue
		}
		if e.RenderFunc != nil {
			e.RenderFunc(canvas, e)
		} else if e.Sprite != nil {
			e.Sprite.Draw(canvas, int(e.X), int(e.Y), e.Frame)
		}
	}
}

// sortedByLayer returns a copy of the entity slice sorted by layer ascending.
func (m *Manager) sortedByLayer() []*Entity {
	out := make([]*Entity, len(m.entities))
	copy(out, m.entities)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Layer < out[j].Layer
	})
	return out
}

// flush integrates pending adds and removes into the active entity list.
func (m *Manager) flush() {
	// Remove dead and explicitly removed entities.
	if len(m.pendingRemove) > 0 || hasDead(m.entities) {
		kept := m.entities[:0]
		for _, e := range m.entities {
			if e.Dead || m.pendingRemove[e.ID] {
				continue
			}
			kept = append(kept, e)
		}
		m.entities = kept
		m.pendingRemove = make(map[string]bool)
	}

	// Append pending additions.
	if len(m.pendingAdd) > 0 {
		m.entities = append(m.entities, m.pendingAdd...)
		m.pendingAdd = m.pendingAdd[:0]
	}
}

// hasDead reports whether any entity is marked Dead.
func hasDead(entities []*Entity) bool {
	for _, e := range entities {
		if e.Dead {
			return true
		}
	}
	return false
}

// isOffscreen reports whether an entity is entirely outside the viewport.
func isOffscreen(e *Entity, width, height int) bool {
	x, y := int(e.X), int(e.Y)
	w, h := 1, 1
	if e.Sprite != nil {
		w = e.Sprite.Width
		h = e.Sprite.Height
	}
	return x+w <= 0 || x >= width || y+h <= 0 || y >= height
}

// autoID generates a deterministic entity ID from a counter and type.
func autoID(n int, t Type) string {
	return string(t) + "-" + itoa(n)
}

// itoa is a minimal int-to-string.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
