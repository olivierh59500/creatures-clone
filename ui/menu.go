package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// MenuAction represents possible menu actions
type MenuAction int

const (
	MenuActionNone MenuAction = iota
	MenuActionStart
	MenuActionOptions
	MenuActionQuit
)

// Menu represents the game menu
type Menu struct {
	// Menu items
	items []MenuItem

	// Current selection
	selectedIndex int

	// Visual properties
	bgColor       color.RGBA
	textColor     color.RGBA
	selectedColor color.RGBA

	// Layout
	centerX    float32
	centerY    float32
	itemHeight float32

	// Animation
	animationTime float64
}

// MenuItem represents a single menu option
type MenuItem struct {
	Text   string
	Action MenuAction
}

// NewMenu creates a new menu
func NewMenu() *Menu {
	return &Menu{
		items: []MenuItem{
			{Text: "Start Game", Action: MenuActionStart},
			{Text: "Options", Action: MenuActionOptions},
			{Text: "Quit", Action: MenuActionQuit},
		},
		selectedIndex: 0,
		bgColor:       color.RGBA{0, 0, 0, 200},
		textColor:     color.RGBA{200, 200, 200, 255},
		selectedColor: color.RGBA{255, 255, 100, 255},
		itemHeight:    40,
	}
}

// Update processes menu input
func (m *Menu) Update(mouseX, mouseY int, clicked bool) MenuAction {
	m.animationTime += 0.016 // 60 FPS

	// Check mouse hover
	for i, item := range m.items {
		itemY := m.centerY + float32(i-len(m.items)/2)*m.itemHeight

		if float32(mouseY) > itemY-m.itemHeight/2 && float32(mouseY) < itemY+m.itemHeight/2 {
			m.selectedIndex = i

			if clicked {
				return item.Action
			}
		}
	}

	return MenuActionNone
}

// Draw renders the menu
func (m *Menu) Draw(screen *ebiten.Image) {
	bounds := screen.Bounds()
	m.centerX = float32(bounds.Dx()) / 2
	m.centerY = float32(bounds.Dy()) / 2

	// Draw semi-transparent background
	vector.DrawFilledRect(screen, 0, 0, float32(bounds.Dx()), float32(bounds.Dy()), m.bgColor, false)

	// Draw title
	title := "CREATURES CLONE"
	titleX := m.centerX - float32(len(title)*4)
	titleY := m.centerY - 100
	ebitenutil.DebugPrintAt(screen, title, int(titleX), int(titleY))

	// Draw menu items
	for i, item := range m.items {
		itemY := m.centerY + float32(i-len(m.items)/2)*m.itemHeight

		// Determine color
		textColor := m.textColor
		if i == m.selectedIndex {
			textColor = m.selectedColor

			// Add selection indicator
			indicator := ">"
			indicatorX := m.centerX - 100
			ebitenutil.DebugPrintAt(screen, indicator, int(indicatorX), int(itemY))
		}

		// Draw text centered
		textX := m.centerX - float32(len(item.Text)*4)
		m.drawTextWithColor(screen, item.Text, int(textX), int(itemY), textColor)
	}

	// Draw instructions
	instructions := "Click to select or press ESC to return"
	instrX := m.centerX - float32(len(instructions)*3)
	instrY := m.centerY + 150
	ebitenutil.DebugPrintAt(screen, instructions, int(instrX), int(instrY))
}

// drawTextWithColor draws text with a specific color
func (m *Menu) drawTextWithColor(screen *ebiten.Image, text string, x, y int, c color.RGBA) {
	// Since ebitenutil.DebugPrint doesn't support color,
	// we just use it as-is. In a full implementation,
	// you'd use a proper font rendering system
	ebitenutil.DebugPrintAt(screen, text, x, y)
}
