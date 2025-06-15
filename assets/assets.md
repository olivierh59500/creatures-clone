# Asset Specifications for Creatures Clone

This document describes all visual assets in the game. Since we're using programmatic rendering, all assets are created using basic shapes and colors.

## Creatures (Norns)

### Base Norn Design
- **Body**: Oval shape (40x50 pixels)
- **Head**: Circle (30x30 pixels) attached to top of body
- **Eyes**: Two circles (8x8 pixels each) with pupils (4x4 pixels)
- **Arms**: Two small ovals (15x8 pixels) on sides
- **Legs**: Two ovals (10x15 pixels) at bottom
- **Tail**: Small triangle or oval at back

### Color Variations (Genetic)
1. **Forest Norn**: Green tones (#228B22 base, #006400 dark)
2. **Desert Norn**: Brown/tan tones (#D2691E base, #8B4513 dark)
3. **Ocean Norn**: Blue tones (#4682B4 base, #191970 dark)
4. **Mountain Norn**: Purple tones (#9370DB base, #4B0082 dark)

### Animations
- **Idle**: Gentle breathing (scale body 98-102%)
- **Walk**: Leg movement, slight body bounce
- **Eat**: Mouth opening, chewing motion
- **Sleep**: Eyes closed, slower breathing
- **Happy**: Bounce animation, smile
- **Sad**: Drooped posture, tears
- **Sick**: Green tint, wobbly movement

### Age Stages
- **Baby**: 70% size, larger head ratio
- **Child**: 85% size, normal proportions
- **Adult**: 100% size
- **Elder**: 95% size, slightly hunched

## Food Items

### Apple
- **Shape**: Circle (20x20 pixels)
- **Color**: Red (#FF0000) with green (#00FF00) stem
- **Animation**: Slight bounce when dropped

### Carrot
- **Shape**: Triangle/cone (10x25 pixels)
- **Color**: Orange (#FFA500) with green top
- **Animation**: None

### Honey
- **Shape**: Hexagon (15x15 pixels)
- **Color**: Golden (#FFD700)
- **Animation**: Slight glow effect

### Seed
- **Shape**: Small oval (8x6 pixels)
- **Color**: Brown (#8B4513)
- **Animation**: None

### Berry
- **Shape**: Small circle cluster (12x12 pixels)
- **Color**: Purple (#800080) or Red (#DC143C)
- **Animation**: Slight shine

## Interactive Objects

### Ball (Toy)
- **Shape**: Circle (25x25 pixels)
- **Color**: Multicolor stripes or solid bright color
- **Animation**: Bounce physics, rotation when rolling

### Music Box
- **Shape**: Rectangle (30x20 pixels) with handle
- **Color**: Wood brown (#8B4513) with gold (#FFD700) details
- **Animation**: Handle rotation when playing, musical notes floating up

### Learning Computer
- **Shape**: Rectangle (40x30 pixels) with screen
- **Color**: Gray (#808080) with blue (#0000FF) screen
- **Animation**: Blinking cursor, text display

### Bed
- **Shape**: Rectangle (60x30 pixels) with pillow
- **Color**: Blue (#4169E1) with white (#FFFFFF) pillow
- **Animation**: None

## Plants

### Tree
- **Shape**: Brown rectangle trunk (20x40) with green circle top (60x60)
- **Color**: Brown (#8B4513) trunk, green (#228B22) leaves
- **Animation**: Gentle leaf sway
- **Growth**: Starts small, grows over time

### Flower
- **Shape**: Stem (2x20) with petal circle (20x20)
- **Color**: Various (red, yellow, blue, pink)
- **Animation**: Slight sway
- **Growth**: Bud -> bloom -> wilt cycle

### Grass
- **Shape**: Multiple thin rectangles (2x10)
- **Color**: Green (#90EE90)
- **Animation**: Wind sway

## Environment

### Ground
- **Texture**: Gradient from grass green (#90EE90) to dirt brown (#8B7355)
- **Pattern**: Some scattered pixels for texture

### Sky
- **Color**: Gradient from light blue (#87CEEB) at horizon to deeper blue (#4682B4) at top
- **Clouds**: White (#FFFFFF) soft circles with transparency

### Water
- **Color**: Blue (#4682B4) with lighter ripples
- **Animation**: Wave motion using sine waves

## UI Elements

### Health Bar
- **Shape**: Rectangle (60x8 pixels)
- **Colors**: 
  - Background: Dark gray (#333333)
  - Health: Green (#00FF00) to Red (#FF0000) gradient
  - Border: Black (#000000)

### Emotion Icons
- **Happy**: Yellow circle with smile curve
- **Sad**: Blue circle with frown curve
- **Angry**: Red circle with angry eyebrows
- **Scared**: Purple circle with wide eyes
- **Curious**: Green circle with raised eyebrow

### Speech Bubble
- **Shape**: Rounded rectangle with tail
- **Color**: White (#FFFFFF) with black (#000000) border
- **Text**: Black pixel font

### Selection Indicator
- **Shape**: Animated dashed circle around selected creature
- **Color**: White (#FFFFFF) or Yellow (#FFFF00)
- **Animation**: Rotating dashes

## Particle Effects

### Food Particles
- **When**: Creature eating
- **Visual**: Small colored squares matching food color
- **Motion**: Arc outward and fade

### Happy Particles
- **When**: Creature very happy
- **Visual**: Small stars or hearts
- **Color**: Yellow (#FFFF00) or Pink (#FFC0CB)
- **Motion**: Float upward and fade

### Sleep Particles
- **When**: Creature sleeping
- **Visual**: "Z" letters
- **Color**: Light blue (#ADD8E6)
- **Motion**: Float upward in wavy pattern

### Learning Particles
- **When**: Creature learns something
- **Visual**: Small "!" or lightbulb
- **Color**: Yellow (#FFFF00)
- **Motion**: Pop above head and fade

## Visual Effects

### Glow Effect
- Used for: Selected objects, important items
- Implementation: Render object multiple times with slight offset and transparency

### Shadow
- Simple dark oval under creatures and objects
- Color: Semi-transparent black
- Size: 80% of object width

### Weather (Future)
- **Rain**: Blue lines falling
- **Snow**: White dots falling with sway
- **Wind**: Affects particle and plant animations

## Color Palette

### Primary Colors
- **Norn Skin Tones**: Various based on genetics
- **UI Blue**: #4169E1
- **UI Green**: #32CD32
- **UI Red**: #DC143C
- **Ground Brown**: #8B7355
- **Sky Blue**: #87CEEB

### Status Colors
- **Healthy**: #00FF00
- **Hungry**: #FFA500
- **Tired**: #9370DB
- **Sick**: #ADFF2F
- **Dying**: #FF0000

## Font Specifications

### Main Font
- Pixel font style
- Size: 8px for normal text
- Size: 12px for headers
- Color: Black (#000000) or White (#FFFFFF) depending on background

### Debug Font
- Monospace style
- Size: 10px
- Color: Green (#00FF00) on black background

## Animation Frame Rates

- **Creature animations**: 8-12 FPS
- **Environmental animations**: 4-6 FPS
- **UI animations**: 15-30 FPS
- **Particle effects**: 30 FPS