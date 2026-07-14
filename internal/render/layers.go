package render

// Layer constants define the rendering order. Entities with a lower layer are
// drawn first and may be overdrawn by entities on a higher layer. The values
// are spaced by 10 to allow fine-grained insertion without renumbering.
const (
	LayerSky        = 0  // Sky background fill
	LayerCelestial  = 10 // Sun, moon, stars, clouds
	LayerDistant    = 20 // Distant mountains and buildings
	LayerNearby     = 30 // Trees, utility poles, nearby scenery
	LayerPlatform   = 40 // Platforms, signals, railway infrastructure
	LayerTrackRear  = 50 // Rear track
	LayerTrain      = 60 // Train
	LayerTrackFront = 70 // Foreground track and nearby objects
	LayerParticles  = 80 // Smoke, rain, snow, particles
	LayerOverlay    = 90 // Status and help overlay
)
