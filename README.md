# Creatures Clone - A Life Simulation Game

A Go implementation of a Creatures-inspired artificial life simulation game using the Ebiten game engine.

## Features

- **Neural Network-based AI**: Each creature has its own neural network for learning and decision making
- **Genetic System**: Creatures can breed and pass on genetic traits
- **Metabolism System**: Creatures need food, can get hungry, tired, and age
- **Learning System**: Creatures learn from experiences and can be taught behaviors
- **Emotion System**: Creatures express happiness, fear, anger, and curiosity
- **Language Learning**: Creatures can learn simple words and associate them with objects
- **Interactive Environment**: Multiple objects and food sources to interact with

## Project Structure

```
creatures-clone/
├── main.go                 # Entry point
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies
├── README.md              # This file
├── assets/                # Asset descriptions
│   └── assets.md         # Detailed asset specifications
├── game/                  # Core game logic
│   ├── game.go           # Main game struct and loop
│   ├── world.go          # World management
│   └── camera.go         # Camera system
├── creature/              # Creature implementation
│   ├── creature.go       # Main creature struct
│   ├── brain.go          # Neural network implementation
│   ├── genetics.go       # Genetic system
│   ├── metabolism.go     # Hunger, energy, health
│   ├── emotions.go       # Emotion system
│   ├── movement.go       # Movement and physics
│   ├── learning.go       # Learning system
│   └── language.go       # Language learning
├── objects/               # Game objects
│   ├── object.go         # Base object interface
│   ├── food.go           # Food items
│   ├── toy.go            # Interactive toys
│   └── plant.go          # Growing plants
├── ui/                    # User interface
│   ├── hud.go            # HUD display
│   ├── menu.go           # Game menus
│   └── debug.go          # Debug overlay
├── utils/                 # Utilities
│   ├── vector.go         # 2D vector math
│   ├── random.go         # Random number generation
│   └── config.go         # Game configuration
└── renderer/              # Rendering system
    ├── renderer.go       # Main renderer
    ├── sprite.go         # Sprite management
    └── animation.go      # Animation system
```

## Building and Running

### Prerequisites

- Go 1.19 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/olivierh59500/creatures-clone.git
cd creatures-clone
```

2. Install dependencies:
```bash
go mod download
```

3. Run the game:
```bash
go run main.go
```

### Building

To build an executable:
```bash
go build -o creatures-clone main.go
```

## Controls

- **Left Click**: Select creature or object
- **Right Click**: Place food
- **WASD/Arrow Keys**: Move camera
- **Mouse Wheel**: Zoom in/out
- **Space**: Pause/Resume
- **Tab**: Toggle debug overlay
- **ESC**: Open menu

## Gameplay

### Getting Started

1. The game starts with 2-3 adult Norns in the world
2. Click on a Norn to select it and see its status
3. Right-click to place food when Norns are hungry
4. Teach Norns by clicking objects and typing words
5. Watch as Norns learn, grow, and interact

### Creature Care

- **Feeding**: Norns need regular food to survive
- **Teaching**: Click objects and type words to teach language
- **Playing**: Use toys to keep Norns happy
- **Breeding**: Happy, healthy adult Norns may breed

### Creature Stats

- **Health**: Overall health (0-100)
- **Hunger**: Hunger level (0-100)
- **Energy**: Energy/tiredness (0-100)
- **Age**: Age in minutes
- **Happiness**: Emotional state (-100 to 100)

## Architecture Details

### Neural Network System

Each creature has a simple feedforward neural network with:
- Input layer: Sensory inputs (vision, hunger, pain, etc.)
- Hidden layers: 2 layers with 20 neurons each
- Output layer: Actions (move, eat, speak, etc.)

The network uses:
- Sigmoid activation functions
- Backpropagation for learning
- Genetic inheritance of weights

### Genetic System

Creatures have genes that control:
- Appearance (color variations)
- Metabolism rates
- Learning speed
- Personality traits
- Initial neural network weights

### Learning System

Creatures learn through:
- **Association**: Linking words with objects
- **Reinforcement**: Positive/negative feedback
- **Imitation**: Watching other creatures
- **Experience**: Trial and error

## Asset Specifications

All assets are described programmatically using basic shapes and colors. See `assets/assets.md` for detailed specifications.

### Creature Sprites
- Basic round body with simple animations
- Color variations based on genetics
- Expression changes based on emotions

### Environment Objects
- Food items (fruits, seeds)
- Interactive toys (ball, music box)
- Plants that grow over time
- Terrain features

## Configuration

Game settings can be modified in `utils/config.go`:

```go
const (
    ScreenWidth  = 1280
    ScreenHeight = 720
    TPS          = 60  // Ticks per second
    
    // Creature settings
    MaxCreatures = 20
    StartingNorns = 3
    
    // World settings
    WorldWidth  = 2000
    WorldHeight = 1000
)
```

## Development

### Adding New Objects

1. Create a new file in `objects/`
2. Implement the `Object` interface
3. Register in `world.go`

### Modifying AI Behavior

1. Edit `creature/brain.go` for neural network
2. Modify `creature/learning.go` for learning rules
3. Adjust reward values in `creature/emotions.go`

### Creating New Genes

1. Add gene type in `creature/genetics.go`
2. Implement expression in relevant systems
3. Add to breeding/mutation logic

## Troubleshooting

### Performance Issues
- Reduce `MaxCreatures` in config
- Lower `WorldWidth/WorldHeight`
- Disable debug overlay with Tab

### Creatures Not Learning
- Check learning rate in `brain.go`
- Ensure proper reinforcement signals
- Verify neural network convergence

## Credits

Inspired by the original Creatures series by Steve Grand and Creature Labs.

## License

MIT License - See LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## Future Features

- [ ] Save/Load system
- [ ] More complex genetics
- [ ] Weather system
- [ ] Day/night cycle
- [ ] More creature species (Grendels, Ettins)
- [ ] Modding support
- [ ] Network multiplayer