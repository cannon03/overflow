package assets

import "overflow/internal/render"

// Sprites is the collection of all sprite definitions.
var Sprites = &SpritesCollection{}

type SpritesCollection struct {
	PlayerNormal    *render.Sprite
	PlayerHit       *render.Sprite
	PlayerBullet    *render.Sprite
	PlayerBulletBig *render.Sprite

	EnemyBasic      *render.Sprite
	EnemyFast       *render.Sprite
	EnemyTank       *render.Sprite
	EnemyBoss       *render.Sprite

	EnemyBullet     *render.Sprite
	EnemyBulletBig  *render.Sprite

	Explosion1      *render.Sprite
	Explosion2      *render.Sprite
	Explosion3      *render.Sprite

	ParticleDot     *render.Sprite
	ParticleStar    *render.Sprite
	ParticleDiamond *render.Sprite

	HeartFull       *render.Sprite
	HeartEmpty      *render.Sprite

	Crosshair       *render.Sprite

	// UI elements
	HUDBorderH     *render.Sprite
	HUDBorderV     *render.Sprite
	HUDBorderTL    *render.Sprite
	HUDBorderTR    *render.Sprite
	HUDBorderBL    *render.Sprite
	HUDBorderBR    *render.Sprite
}

func init() {
	// Player sprites
	Sprites.PlayerNormal = render.NewSpriteFromString(
		`╲▲╱
◀■▶
╱▼╲`, render.ColorBrightCyan, render.ColorDarkBlue)

	Sprites.PlayerHit = render.NewSpriteFromString(
		`╲▲╱
◀■▶
╱▼╲`, render.ColorWhite, render.ColorDarkBlue)

	Sprites.PlayerBullet = render.NewSprite(
		1, 1, []rune{'∙'}, render.ColorBrightYellow, render.ColorDarkBlue)

	Sprites.PlayerBulletBig = render.NewSprite(
		1, 1, []rune{'●'}, render.ColorYellow, render.ColorDarkBlue)

	// Enemy sprites
	// Hexagon/diamond silhouettes with eyes — spaces on edges are transparent
	// and show the dark blue game area background, creating a shaped silhouette.
	Sprites.EnemyBasic = render.NewSpriteFromString(
		` ██ 
█◆◆█
 ██ `, render.ColorRed, render.ColorDarkBlue)

	Sprites.EnemyFast = render.NewSpriteFromString(
		` ██ 
█◈◈█
 ██ `, render.ColorMagenta, render.ColorDarkBlue)

	Sprites.EnemyTank = render.NewSpriteFromString(
		` ████ 
█◆◆◆◆█
 ████ `, render.ColorDarkRed, render.ColorDarkBlue)

	Sprites.EnemyBoss = render.NewSpriteFromString(
		`  ████  
 ██████ 
██◆◆◆◆██
██◆◆◆◆██
████████`, render.ColorOrange, render.ColorDarkBlue)

	// Enemy bullet
	Sprites.EnemyBullet = render.NewSprite(
		1, 1, []rune{'•'}, render.ColorRed, render.ColorDarkBlue)

	Sprites.EnemyBulletBig = render.NewSprite(
		1, 1, []rune{'●'}, render.ColorOrange, render.ColorDarkBlue)

	// Explosion sprites
	Sprites.Explosion1 = render.NewSprite(
		3, 3, []rune{
			' ', '✦', ' ',
			'✧', '●', '✧',
			' ', '✦', ' ',
		}, render.ColorYellow, render.ColorDarkBlue)

	Sprites.Explosion2 = render.NewSprite(
		3, 3, []rune{
			'✧', ' ', '✧',
			' ', '★', ' ',
			'✧', ' ', '✧',
		}, render.ColorOrange, render.ColorDarkBlue)

	Sprites.Explosion3 = render.NewSprite(
		5, 5, []rune{
			' ', ' ', '✦', ' ', ' ',
			' ', '✧', '★', '✧', ' ',
			'✦', '★', '●', '★', '✦',
			' ', '✧', '★', '✧', ' ',
			' ', ' ', '✦', ' ', ' ',
		}, render.ColorRed, render.ColorDarkBlue)

	// Particle sprites
	Sprites.ParticleDot = render.NewSprite(
		1, 1, []rune{'·'}, render.ColorWhite, render.ColorDarkBlue)

	Sprites.ParticleStar = render.NewSprite(
		1, 1, []rune{'✦'}, render.ColorOrange, render.ColorDarkBlue)

	Sprites.ParticleDiamond = render.NewSprite(
		1, 1, []rune{'◆'}, render.ColorYellow, render.ColorDarkBlue)

	// Heart
	Sprites.HeartFull = render.NewSprite(
		1, 1, []rune{'♥'}, render.ColorRed, render.ColorBlack)

	Sprites.HeartEmpty = render.NewSprite(
		1, 1, []rune{'♡'}, render.ColorDarkRed, render.ColorBlack)

	// Crosshair
	Sprites.Crosshair = render.NewSprite(
		3, 3, []rune{
			'┌', '┬', '┐',
			'├', '┼', '┤',
			'└', '┴', '┘',
		}, render.ColorDarkGray, render.ColorBlack)

	// HUD border components
	Sprites.HUDBorderTL = render.NewSprite(
		1, 1, []rune{'╔'}, render.ColorHudBorder, render.ColorHudBG)
	Sprites.HUDBorderTR = render.NewSprite(
		1, 1, []rune{'╗'}, render.ColorHudBorder, render.ColorHudBG)
	Sprites.HUDBorderBL = render.NewSprite(
		1, 1, []rune{'╚'}, render.ColorHudBorder, render.ColorHudBG)
	Sprites.HUDBorderBR = render.NewSprite(
		1, 1, []rune{'╝'}, render.ColorHudBorder, render.ColorHudBG)
	Sprites.HUDBorderH = render.NewSprite(
		1, 1, []rune{'═'}, render.ColorHudBorder, render.ColorHudBG)
	Sprites.HUDBorderV = render.NewSprite(
		1, 1, []rune{'║'}, render.ColorHudBorder, render.ColorHudBG)
}
