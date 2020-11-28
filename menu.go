package main

import (
	"github.com/nullboundary/glfont"
)

type Menu struct {
	font               *glfont.Font
	aboutFont          *glfont.Font
	currentItem        int
	menuItems          []string
	levels             []string
	menuCalls          []func()
	shouldQuit         bool
	showAbout          bool
	showSelectLevel    bool
	currentLevelSelect int
	showMenu           bool
}

func (m *Menu) Init() {
	font, err := glfont.LoadFont(conf.Assets["menuFont"], int32(72), screenWidth, screenHeight)
	if err != nil {
		panic(err)
	}
	m.font = font

	font, err = glfont.LoadFont(conf.Assets["statsFont"], int32(22), screenWidth, screenHeight)
	if err != nil {
		panic(err)
	}
	m.aboutFont = font

	m.currentItem = 0
	m.currentLevelSelect = 0

	m.menuItems = []string{
		"START",
		"SELECT LEVEL",
		"ABOUT",
		"QUIT",
	}
	m.menuCalls = []func(){
		m.Start,
		m.SelectLevel,
		m.About,
		m.Quit,
	}

	m.levels = []string{
		"Level 1 - Easy peazy",
		"Level 2 - Zzzzz",
		"Level 3 - all u got",
		"Level 4 - WTF",
		"Level 5 - Please help me",
		"Level 6 - Cry in a corner",
		"Level 7 - OMG",
		"Level 8 - Please stop",
		"Level 9 - I hate you",
		"Level 10 - yeah right",
	}
}

func (m *Menu) Start() {
	m.showSelectLevel = false
	m.showMenu = false
	gameMap.StartLevel(1)
}

func (m *Menu) SelectLevel() {
	m.showSelectLevel = !m.showSelectLevel
}

func (m *Menu) About() {
	m.showAbout = !m.showAbout
}

func (m *Menu) Quit() {
	m.shouldQuit = true
}

func (m *Menu) Select() {
	if m.showSelectLevel {
		m.showSelectLevel = false
		m.showMenu = false
		gameMap.StartLevel(m.currentLevelSelect + 1)
	} else {
		m.menuCalls[m.currentItem]()
	}
}

func (m *Menu) Up() {
	if m.showSelectLevel {
		if m.currentLevelSelect != 0 {
			m.currentLevelSelect--
		} else {
			m.currentLevelSelect = len(m.levels) - 1
		}
	} else {
		if m.currentItem != 0 {
			m.currentItem--
		} else {
			m.currentItem = len(m.menuItems) - 1
		}
	}
}

func (m *Menu) Down() {
	if m.showSelectLevel {
		if m.currentLevelSelect != len(m.levels)-1 {
			m.currentLevelSelect++
		} else {
			m.currentLevelSelect = 0
		}
	} else {
		if m.currentItem != len(m.menuItems)-1 {
			m.currentItem++
		} else {
			m.currentItem = 0
		}
	}
}

func (m *Menu) Draw(dt float64) {
	// Clear background text if menu is shown.
	gameMap.text = ""

	if m.showAbout {
		lines := []string{
			"     This game was made by Magnus Persson",
			"",
			"           For the #GitHubGameOff 2020",
			"",
			"Source code: github.com/lallassu/moonshot",
			"Twitter: @lallassu",
			"",
			"Credits:",
			"Planets/Ships: https://opengameart.org/content/space-game-art-pack-extended",
			"Map Gfx: ansimuz.com",
			"Sound: soundbible.com",
			"Music: https://dos88.itch.io/dos-88-music-library",
		}
		for i, s := range lines {
			if i < 6 {
				mc := conf.Colors["aboutFirst"]
				m.aboutFont.SetColor(mc.R, mc.G, mc.B, mc.A)
			} else {
				mc := conf.Colors["aboutCredits"]
				m.aboutFont.SetColor(mc.R, mc.G, mc.B, mc.A)
			}
			m.aboutFont.Printf(screenWidth/3, screenHeight/3+float32(i*22), 1.0, s)
		}
	} else if m.showSelectLevel {
		for i, s := range m.levels {
			m.font.SetColor(float32(i)*0.2, 1.0-float32(i)*0.1, 0, 1.0)
			if i == m.currentLevelSelect {
				m.font.Printf(screenWidth/5, screenHeight/5+float32(i*80), 1.3, s)
			} else {
				m.font.Printf(screenWidth/5, screenHeight/5+float32(i*80), 1.0, s)
			}
		}
	} else {
		for i, k := range m.menuItems {
			if i == m.currentItem {
				mc := conf.Colors["menuRegular"]
				m.font.SetColor(mc.R, mc.G, mc.B, mc.A)
				m.font.Printf(screenWidth/3, screenHeight/3+float32(i*90), 1.5, k)
			} else {
				mc := conf.Colors["menuSelected"]
				m.font.SetColor(mc.R, mc.G, mc.B, mc.A)
				m.font.Printf(screenWidth/3, screenHeight/3+float32(i*80), 1.0, k)
			}
		}
	}
}
