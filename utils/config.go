package utils

// Config holds all game configuration values
type Config struct {
	// Display settings
	ScreenWidth  int
	ScreenHeight int
	FullScreen   bool
	VSync        bool

	// World settings
	WorldWidth  int
	WorldHeight int

	// Game settings
	TicksPerSecond int
	MaxCreatures   int
	StartingNorns  int

	// Graphics settings
	EnableParticles bool
	EnableShadows   bool
	ParticleLimit   int

	// Audio settings
	MasterVolume  float64
	MusicVolume   float64
	EffectsVolume float64

	// Debug settings
	DebugMode    bool
	ShowFPS      bool
	ShowHitboxes bool

	// Gameplay settings
	DifficultyLevel int
	AutoSave        bool
	AutoSaveMinutes int
}

// LoadConfig loads the game configuration
func LoadConfig() *Config {
	// In a full implementation, this would load from a file
	// For now, return default values
	return &Config{
		// Display
		ScreenWidth:  1280,
		ScreenHeight: 720,
		FullScreen:   false,
		VSync:        true,

		// World - Made larger
		WorldWidth:  4000, // Doubled from 2000
		WorldHeight: 2000, // Doubled from 1000

		// Game
		TicksPerSecond: 60,
		MaxCreatures:   50, // Increased from 20
		StartingNorns:  5,  // Increased from 3

		// Graphics
		EnableParticles: true,
		EnableShadows:   true,
		ParticleLimit:   1000,

		// Audio
		MasterVolume:  0.8,
		MusicVolume:   0.6,
		EffectsVolume: 0.7,

		// Debug
		DebugMode:    false,
		ShowFPS:      true,
		ShowHitboxes: false,

		// Gameplay
		DifficultyLevel: 1, // 0=Easy, 1=Normal, 2=Hard
		AutoSave:        true,
		AutoSaveMinutes: 5,
	}
}

// SaveConfig saves the configuration to file
func (c *Config) SaveConfig() error {
	// In a full implementation, this would save to a file
	// For now, just return nil
	return nil
}

// Validate ensures configuration values are valid
func (c *Config) Validate() {
	// Ensure reasonable bounds
	c.ScreenWidth = ClampInt(c.ScreenWidth, 640, 3840)
	c.ScreenHeight = ClampInt(c.ScreenHeight, 480, 2160)

	c.WorldWidth = ClampInt(c.WorldWidth, 1000, 5000)
	c.WorldHeight = ClampInt(c.WorldHeight, 500, 3000)

	c.TicksPerSecond = ClampInt(c.TicksPerSecond, 30, 120)
	c.MaxCreatures = ClampInt(c.MaxCreatures, 1, 100)
	c.StartingNorns = ClampInt(c.StartingNorns, 1, 10)

	c.ParticleLimit = ClampInt(c.ParticleLimit, 100, 5000)

	c.MasterVolume = Clamp(c.MasterVolume, 0, 1)
	c.MusicVolume = Clamp(c.MusicVolume, 0, 1)
	c.EffectsVolume = Clamp(c.EffectsVolume, 0, 1)

	c.DifficultyLevel = ClampInt(c.DifficultyLevel, 0, 2)
	c.AutoSaveMinutes = ClampInt(c.AutoSaveMinutes, 1, 60)
}
