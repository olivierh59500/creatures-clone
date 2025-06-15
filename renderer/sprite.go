package renderer

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite represents a drawable sprite
type Sprite struct {
	Image      *ebiten.Image
	Width      int
	Height     int
	FrameCount int
	Origin     image.Point // Origin point for rotation
}

// NewSprite creates a new sprite
func NewSprite(width, height, frameCount int) *Sprite {
	return &Sprite{
		Image:      ebiten.NewImage(width*frameCount, height),
		Width:      width,
		Height:     height,
		FrameCount: frameCount,
		Origin:     image.Point{X: width / 2, Y: height / 2},
	}
}

// DrawFrame draws a specific frame of the sprite
func (s *Sprite) DrawFrame(screen *ebiten.Image, frameIndex int, x, y float64, opts *ebiten.DrawImageOptions) {
	if opts == nil {
		opts = &ebiten.DrawImageOptions{}
	}

	// Calculate source rectangle for the frame
	sx := frameIndex * s.Width
	sy := 0

	// Create sub-image for the frame
	subImage := s.Image.SubImage(image.Rect(sx, sy, sx+s.Width, sy+s.Height)).(*ebiten.Image)

	// Apply translation
	opts.GeoM.Translate(x-float64(s.Origin.X), y-float64(s.Origin.Y))

	// Draw the frame
	screen.DrawImage(subImage, opts)
}

// CreateCircleSprite creates a sprite with a circle
func CreateCircleSprite(radius float32, c color.Color) *Sprite {
	size := int(radius * 2)
	sprite := NewSprite(size, size, 1)

	// Draw circle on sprite
	// In a full implementation, would use proper circle drawing
	sprite.Image.Fill(c)

	return sprite
}

// CreateRectSprite creates a sprite with a rectangle
func CreateRectSprite(width, height int, c color.Color) *Sprite {
	sprite := NewSprite(width, height, 1)
	sprite.Image.Fill(c)
	return sprite
}

// SpriteSheet manages multiple sprites
type SpriteSheet struct {
	sprites map[string]*Sprite
}

// NewSpriteSheet creates a new sprite sheet
func NewSpriteSheet() *SpriteSheet {
	return &SpriteSheet{
		sprites: make(map[string]*Sprite),
	}
}

// AddSprite adds a sprite to the sheet
func (ss *SpriteSheet) AddSprite(name string, sprite *Sprite) {
	ss.sprites[name] = sprite
}

// GetSprite retrieves a sprite by name
func (ss *SpriteSheet) GetSprite(name string) *Sprite {
	return ss.sprites[name]
}

// CreateProgrammaticSprites creates all sprites programmatically
func CreateProgrammaticSprites() *SpriteSheet {
	sheet := NewSpriteSheet()

	// Create creature body parts
	sheet.AddSprite("creature_body", CreateOvalSprite(40, 50, color.RGBA{100, 200, 100, 255}))
	sheet.AddSprite("creature_head", CreateCircleSprite(15, color.RGBA{100, 200, 100, 255}))
	sheet.AddSprite("creature_eye", CreateCircleSprite(4, color.White))
	sheet.AddSprite("creature_pupil", CreateCircleSprite(2, color.Black))

	// Create food sprites
	sheet.AddSprite("apple", CreateCircleSprite(10, color.RGBA{255, 0, 0, 255}))
	sheet.AddSprite("carrot", CreateTriangleSprite(10, 25, color.RGBA{255, 165, 0, 255}))
	sheet.AddSprite("honey", CreateHexagonSprite(15, color.RGBA{255, 215, 0, 255}))

	// Create toy sprites
	sheet.AddSprite("ball", CreateCircleSprite(12, color.RGBA{255, 0, 0, 255}))
	sheet.AddSprite("musicbox", CreateRectSprite(30, 20, color.RGBA{139, 69, 19, 255}))

	// Create plant sprites
	sheet.AddSprite("tree_trunk", CreateRectSprite(20, 40, color.RGBA{139, 69, 19, 255}))
	sheet.AddSprite("tree_leaves", CreateCircleSprite(30, color.RGBA{34, 139, 34, 255}))
	sheet.AddSprite("flower_stem", CreateRectSprite(2, 20, color.RGBA{0, 128, 0, 255}))
	sheet.AddSprite("flower_petals", CreateCircleSprite(10, color.RGBA{255, 192, 203, 255}))

	return sheet
}

// CreateOvalSprite creates an oval-shaped sprite
func CreateOvalSprite(width, height float32, c color.Color) *Sprite {
	sprite := NewSprite(int(width), int(height), 1)

	// Fill with color (simplified - would draw actual oval)
	sprite.Image.Fill(c)

	return sprite
}

// CreateTriangleSprite creates a triangle-shaped sprite
func CreateTriangleSprite(width, height float32, c color.Color) *Sprite {
	sprite := NewSprite(int(width), int(height), 1)

	// Fill with color (simplified - would draw actual triangle)
	sprite.Image.Fill(c)

	return sprite
}

// CreateHexagonSprite creates a hexagon-shaped sprite
func CreateHexagonSprite(size float32, c color.Color) *Sprite {
	sprite := NewSprite(int(size*2), int(size*2), 1)

	// Fill with color (simplified - would draw actual hexagon)
	sprite.Image.Fill(c)

	return sprite
}
