package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/nullboundary/glfont"
)

type Map struct {
	nextLevel             int
	landedTime            time.Time
	rocketX               float64
	rocketY               float64
	rocketBoostMax        float64
	currentLevel          int
	bgAsset               string
	mapAsset              string
	alienAsset            string
	moonAsset             string
	satelliteAsset        string
	nrOfSatellites        int
	nrOfDebris            int
	Wind                  float64
	nrOfAliens            int
	alienBombs            float64
	alienBombSpeed        float64
	alienBombMaxTime      float64
	alienBombType         ShotType
	alienCanBeHitByDebris bool
	alienBombFreq         int
	alienWaitForRelease   bool
	moonSpeed             float64
	satelliteSpeed        float64
	retries               int
	totalRetries          int
	font                  *glfont.Font
	text                  string
}

func (m *Map) Init() {
	font, err := glfont.LoadFont(conf.Assets["menuFont"], int32(72), screenWidth, screenHeight)
	if err != nil {
		panic(err)
	}
	m.font = font
}

func (m *Map) Draw(dt float64) {
	if time.Since(m.landedTime).Seconds() > 3 && m.nextLevel != 0 {
		m.text = ""
		m.StartLevel(m.nextLevel)
		m.nextLevel = 0
	}
	if m.text != "" {
		m.font.Printf(screenWidth/2-float32(len(m.text)*15), screenHeight/2, 1.0, m.text)
	}

	if !(menu.showMenu || menu.showAbout || menu.showSelectLevel) {
		mc := conf.Colors["fuel"]
		rocket.boostFont.SetColor(mc.R, mc.G, mc.B, mc.A)
		rocket.boostFont.Printf(screenWidth/2-100, screenHeight-60, 1.0, fmt.Sprintf("Fuel: %0.2f %v", (rocket.boost/rocket.maxBoost*100), "%%"))

		if m.Wind != 0 {
			mc := conf.Colors["wind"]
			rocket.boostFont.SetColor(mc.R, mc.G, mc.B, mc.A)
			w := ""
			if m.Wind < 0 {
				w = fmt.Sprintf("Wind: <- %0.2f", math.Abs(m.Wind))
			} else {
				w = fmt.Sprintf("Wind: %0.2f ->", math.Abs(m.Wind))
			}
			rocket.boostFont.Printf(screenWidth/2-100, screenHeight-4, 1.0, w)
		}

		stats.font.SetColor(1.0, 1.0, 1.0, 0.7)
		stats.font.Printf(10, screenHeight-20, 1.1, "<space> - Fuel")
		stats.font.Printf(10, screenHeight-4, 1.1, "<enter> - Launch")
	}
}

func (m *Map) StartLevel(level int) {
	m.font.SetColor(1.0, 1.0, 0.0, 1.0)
	m.text = fmt.Sprintf("Level %d", level)
	go func() {
		time.Sleep(2 * time.Second)
		m.text = ""
	}()

	m.Wind = 0

	switch level {
	case 1:
		m.currentLevel = 1
		m.bgAsset = "bg1"
		m.mapAsset = "map1"
		m.alienAsset = "alien1"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon1"
		m.nrOfSatellites = 0
		m.nrOfAliens = 0
		m.moonSpeed = 0.5
		m.retries = 10
		m.rocketX = screenWidth/2 - 50
		m.rocketY = screenHeight/2 - 250
		m.nrOfDebris = 0
		m.rocketBoostMax = 30
		m.LoadLevel()
	case 2:
		m.currentLevel = 2
		m.bgAsset = "bg2"
		m.mapAsset = "map2"
		m.alienAsset = "alien1"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon2"
		m.nrOfSatellites = 8
		m.nrOfAliens = 5
		m.alienBombs = 0
		m.alienBombSpeed = 0
		m.alienBombType = ShotRandom
		m.alienWaitForRelease = false
		m.alienBombFreq = 0
		m.satelliteSpeed = 1
		m.nrOfDebris = 0
		m.moonSpeed = 1
		m.retries = 10
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight/2 - 250
		m.rocketBoostMax = 30
		m.LoadLevel()
		m.alienCanBeHitByDebris = true
	case 3:
		m.currentLevel = 3
		m.bgAsset = "bg3"
		m.mapAsset = "map3"
		m.alienAsset = "alien1"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon2"
		m.nrOfSatellites = 0
		m.nrOfAliens = 3
		m.alienBombs = 0
		m.alienBombSpeed = 1
		m.alienBombType = ShotRandom
		m.alienBombMaxTime = 5
		m.alienWaitForRelease = false
		m.alienBombFreq = 2
		m.satelliteSpeed = 0
		m.moonSpeed = 1.5
		m.retries = 8
		m.nrOfDebris = 1
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight/2 - 250
		m.rocketBoostMax = 30
		m.LoadLevel()

	case 4:
		m.currentLevel = 4
		m.bgAsset = "bg4"
		m.mapAsset = "map4"
		m.alienAsset = "alien2"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon3"
		m.nrOfSatellites = 5
		m.nrOfAliens = 3
		m.alienBombs = 0
		m.alienBombSpeed = 0
		m.satelliteSpeed = 0
		m.moonSpeed = 1
		m.retries = 10
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight / 2
		m.rocketBoostMax = 30
		m.nrOfDebris = 0

		m.LoadLevel()
	case 5:
		m.currentLevel = 5
		m.bgAsset = "bg5"
		m.mapAsset = "map5"
		m.alienAsset = "alien2"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon3"
		m.nrOfSatellites = 0
		m.nrOfAliens = 2
		m.alienBombs = 3
		m.alienBombSpeed = 2
		m.alienBombType = ShotStraight
		m.alienBombFreq = 5
		m.alienBombMaxTime = 5
		m.alienWaitForRelease = true
		m.satelliteSpeed = 0
		m.moonSpeed = 0.5
		m.retries = 6
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight/2 - 300
		m.rocketBoostMax = 40
		m.nrOfDebris = 1
		m.Wind = 0
		m.alienCanBeHitByDebris = false

		m.LoadLevel()
	case 6:
		m.currentLevel = 6
		m.bgAsset = "bg6"
		m.mapAsset = "map6"
		m.alienAsset = "alien3"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon4"
		m.nrOfSatellites = 0
		m.nrOfAliens = 5
		m.alienBombs = 5
		m.alienBombSpeed = 3
		m.alienBombType = ShotStraight
		m.alienBombFreq = 2
		m.alienBombMaxTime = 2
		m.alienWaitForRelease = true
		m.satelliteSpeed = 0
		m.moonSpeed = 1.5
		m.retries = 5
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight/2 - 300
		m.rocketBoostMax = 50
		m.nrOfDebris = 2
		m.Wind = 3 - rand.Float64()*6
		m.alienCanBeHitByDebris = false

		m.LoadLevel()
	case 7:
		m.currentLevel = 7
		m.bgAsset = "bg7"
		m.mapAsset = "map7"
		m.alienAsset = "alien3"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon4"
		m.nrOfSatellites = 0
		m.nrOfAliens = 2
		m.alienBombs = 2
		m.alienBombSpeed = 3
		m.alienBombType = ShotStraight
		m.alienBombFreq = 4
		m.alienBombMaxTime = 2.5
		m.alienWaitForRelease = false
		m.satelliteSpeed = 0
		m.moonSpeed = 2.5
		m.retries = 5
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight/2 - 300
		m.rocketBoostMax = 40
		m.nrOfDebris = 4
		m.Wind = 3 - rand.Float64()*6
		m.alienCanBeHitByDebris = false

		m.LoadLevel()
	case 8:
		m.currentLevel = 8
		m.bgAsset = "bg8"
		m.mapAsset = "map8"
		m.alienAsset = "alien4"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon5"
		m.nrOfSatellites = 0
		m.nrOfAliens = 6
		m.alienBombs = 5
		m.alienBombSpeed = 4
		m.alienBombType = ShotStraight
		m.alienBombFreq = 2
		m.alienBombMaxTime = 2
		m.alienWaitForRelease = true
		m.satelliteSpeed = 0
		m.moonSpeed = 0.5
		m.retries = 5
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight/2 - 300
		m.rocketBoostMax = 50
		m.nrOfDebris = 6
		m.Wind = 3 - rand.Float64()*6
		m.alienCanBeHitByDebris = true

		m.LoadLevel()
	case 9:
		m.currentLevel = 9
		m.bgAsset = "bg9"
		m.mapAsset = "map9"
		m.alienAsset = "alien4"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon5"
		m.nrOfSatellites = 0
		m.nrOfAliens = 3
		m.alienBombs = 5
		m.alienBombSpeed = 3
		m.alienBombType = ShotSeeking
		m.alienBombFreq = 5
		m.alienBombMaxTime = 3
		m.alienWaitForRelease = true
		m.satelliteSpeed = 0
		m.moonSpeed = 1.5
		m.retries = 5
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight/2 - 200
		m.rocketBoostMax = 50
		m.nrOfDebris = 10
		m.Wind = 3 - rand.Float64()*6
		m.alienCanBeHitByDebris = false

		m.LoadLevel()
	case 10:
		m.currentLevel = 10
		m.bgAsset = "bg10"
		m.mapAsset = "map10"
		m.alienAsset = "alien2"
		m.satelliteAsset = "satellite"
		m.moonAsset = "moon6"
		m.nrOfSatellites = 0
		m.nrOfAliens = 10
		m.alienBombs = 5
		m.alienBombSpeed = 5
		m.alienBombType = ShotSeeking
		m.alienBombFreq = 2
		m.alienBombMaxTime = 1
		m.alienWaitForRelease = false
		m.satelliteSpeed = 0
		m.moonSpeed = 1.5
		m.retries = 10
		m.rocketX = screenWidth / 2
		m.rocketY = screenHeight/2 - 100
		m.rocketBoostMax = 50
		m.nrOfDebris = 20
		m.Wind = 3 - rand.Float64()*6
		m.alienCanBeHitByDebris = false

		m.LoadLevel()
	}
}

func (m *Map) LoadLevel() {
	m.ClearCurrent()

	m.totalRetries = m.retries

	world.shader = shaders[ObjectWorld]
	world.Init(screenWidth, screenHeight)

	// Set background
	background.Init(m.bgAsset)

	// Load the foreground/map
	if err := LoadMap(conf.Assets[m.mapAsset]); err != nil {
		panic(err)
	}

	// Set moon
	moon.Init(600, 800, 1, conf.Assets[m.moonAsset], ObjectMoon)
	moon.speed = m.moonSpeed
	objects = append(objects, moon)

	// Create satellites
	for i := 0; i < m.nrOfSatellites; i++ {
		s := &Satellite{origX: 500 + rand.Float64()*100, origY: -300 * rand.Float64()}
		s.rotation = rand.Float64() * 100
		s.Init(0, 0, 2, conf.Assets[m.satelliteAsset], ObjectSatellite)
		objects = append(objects, s)
	}

	// Debris around moon
	for i := 0; i < m.nrOfDebris; i++ {
		s := &Debris{
			origX: 300 - rand.Float64()*600,
			origY: 300 - rand.Float64()*600,
		}
		dType := fmt.Sprintf("rock%d", rand.Intn(7)+1)
		s.Init(rand.Float64()*screenWidth, screenHeight-screenHeight/4, 2, conf.Assets[dType], ObjectDebris)
		objects = append(objects, s)
	}

	// Create alien ships
	for i := 0; i < m.nrOfAliens; i++ {
		s := &Alien{
			AmmoSpeed:        m.alienBombSpeed,
			AmmoFreq:         m.alienBombFreq,
			AmmoType:         m.alienBombType,
			AmmoMaxTime:      m.alienBombMaxTime,
			WaitForRelease:   m.alienWaitForRelease,
			CanBeHitByDebris: m.alienCanBeHitByDebris,
		}
		s.Init(rand.Float64()*screenWidth, screenHeight-screenHeight/4, 2, conf.Assets[m.alienAsset], ObjectAlien)
		objects = append(objects, s)
	}
	// Create rocket
	m.CreateRocket()
}

func (m *Map) ClearCurrent() {
	background.Clear()
	for i := range objects {
		objects[i].Remove()
	}
	RemoveObjects()

	objects = []Obj{}
	moon = &Moon{}
	rocket = &Rocket{}
	background = &Background{}
	world = &World{}
}

func (m *Map) Reset() {
	m.retries--
	if m.retries == 0 {
		// Project failed!
		m.font.SetColor(1.0, 0.0, 0.0, 1.0)
		m.text = "MoonShot Project: FAILED."
		go func() {
			time.Sleep(3 * time.Second)
			m.text = ""
			menu.showMenu = true
		}()
	} else {
		// Failed attemp, retry!
		m.font.SetColor(1.0, 1.0, 0.0, 1.0)
		m.text = fmt.Sprintf("Attempt: %d/%d", m.retries, m.totalRetries)
		go func() {
			time.Sleep(3 * time.Second)
			m.text = ""
		}()
		m.CreateRocket()
	}
}

func (m *Map) Landed() {
	m.font.SetColor(0.0, 1.0, 0.0, 1.0)
	m.text = "MoonShot project: SUCCESS!"
	sound.Play("success", 1.0)
	go func() {
		time.Sleep(3 * time.Second)
		if gameMap.currentLevel == 10 {
			m.font.SetColor(0.0, 1.0, 0.0, 1.0)
			m.text = "All moons visited, congratz!"
			go func() {
				for {
					if gameMap.currentLevel == 10 {
						time.Sleep(100 * time.Millisecond)
						Explode(rand.Float64()*screenWidth, rand.Float64()*screenHeight, rand.Float64()*50)
					}
				}
			}()
		} else {
			m.nextLevel = m.currentLevel + 1
			m.landedTime = time.Now()
		}
	}()

}

func (m *Map) CreateRocket() {
	rocket.removed = true
	r := Rocket{maxBoost: m.rocketBoostMax, initSleepMs: 2000}
	r.Init(m.rocketX, m.rocketY, 10, conf.Assets["rocket"], ObjectRocket)
	objects = append(objects, &r)
	rocket = &r
}
