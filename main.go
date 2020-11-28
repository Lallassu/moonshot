package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	screenWidth  = 1280
	screenHeight = 1024
)

type ObjectType int

const (
	ObjectRocket ObjectType = iota + 1
	ObjectMoon
	ObjectSatellite
	ObjectAlien
	ObjectShot
	ObjectBackground
	ObjectWorld
	ObjectDebris
)

// Obj is an interface for game objects
type Obj interface {
	Update(dt float64)
	Hit(x, y int, objType ObjectType)
	Draw(dt float64)
	GetBorderPixels() []float64
	GetX() float64
	GetY() float64
	IsActive(x, y int) bool
	GetObjType() ObjectType
	IsRemoved() bool
	Clear()
	Remove()
}

var particles = &ParticlePool{}
var world = &World{}
var background = &Background{}
var shader = &Shader{}
var stats = &Stats{}
var objects = []Obj{}
var moon = &Moon{}
var rocket = &Rocket{}
var conf = Config{}
var gameMap = Map{}
var menu = Menu{}
var sound = Sound{}

var shaders = map[ObjectType]*Shader{}

func init() {
	runtime.LockOSThread()
}

var view = mgl32.Translate3D(-screenWidth/2, -screenHeight/2, -screenWidth+43)

//var view = mgl32.Translate3D(-20, -20, -180) //-screenWidth+43)

var projection = mgl32.Mat4{}

func main() {
	conf = LoadConfiguration("./gameconf.json")

	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	window, err := glfw.CreateWindow(screenWidth, screenHeight, "MoonShot", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(frameBufferSizeCallback)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	// Initiate shaders
	shaders = make(map[ObjectType]*Shader)
	shader, err := NewShader(conf.Shaders["regular_vs"], conf.Shaders["regular_fs"])
	if err != nil {
		panic(err)
	}
	bgShader, err := NewShader(conf.Shaders["texture_vs"], conf.Shaders["texture_fs"])
	if err != nil {
		panic(err)
	}

	shaders[ObjectMoon] = shader
	shaders[ObjectSatellite] = shader
	shaders[ObjectRocket] = shader
	shaders[ObjectAlien] = shader
	shaders[ObjectShot] = shader
	shaders[ObjectBackground] = bgShader
	shaders[ObjectWorld] = shader
	shaders[ObjectDebris] = shader

	// Initiate misc stuff
	stats.Init()
	particles.Init()
	menu.Init()
	gameMap.Init()
	sound.Init()

	// Load sounds
	for k, v := range conf.Sounds {
		sound.Load(v, k)
	}

	sound.Play("bgmusic", 0.3)

	// Key handler
	keyHandler := &KeyHandler{}
	window.SetCursorPosCallback(keyHandler.MousePos)
	window.SetMouseButtonCallback(keyHandler.MouseDown)

	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(screenWidth)/float32(screenHeight), 0.1, 2000.0)

	// Start level 1
	gameMap.StartLevel(1)

	// render loop
	dt := float64(0)
	lastTS := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)

		keyHandler.Process(window)

		stats.Update()

		if keyHandler.WireFrame {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}

		background.Draw(dt)
		// Detect collision between objects
		DetectCollisions(dt)

		if rand.Intn(10) > 8 {
			Star()
		}

		for i := range objects {
			if objects[i].IsRemoved() {
				continue
			}
			objects[i].Update(dt)
		}

		// Remove old objects
		RemoveObjects()

		world.Draw(dt)
		particles.Draw(dt)

		if keyHandler.Debug {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			stats.Draw()
		}

		gameMap.Draw(dt)

		if menu.showMenu {
			menu.Draw(dt)
		}
		if menu.shouldQuit {
			window.SetShouldClose(true)
		}

		window.SwapBuffers()
		glfw.PollEvents()

		dt = float64(time.Since(lastTS).Milliseconds())
		lastTS = time.Now()
	}

	world.Clear()
	gameMap.ClearCurrent()
}

func RemoveObjects() {
	removed := 0
	for i := range objects {
		n := i - removed
		if objects[n].IsRemoved() {
			objects[n].Clear()
			objects[n] = objects[len(objects)-1]
			objects[len(objects)-1] = nil
			objects = objects[:len(objects)-1]
			removed++
		}
	}
}

func frameBufferSizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func DetectCollisions(dt float64) {
	for n := range objects {
		if objects[n].IsRemoved() {
			continue
		}
		p1 := objects[n]
		for i := range objects {
			// Skip own check
			if i == n || objects[i].IsRemoved() {
				continue
			}
			p := objects[i]

			pbp := p.GetBorderPixels()
			if len(pbp) > 0 {
				for i := 0; i < len(pbp); i += 2 {
					x := p.GetX() + pbp[i]
					y := p.GetY() + pbp[i+1]

					if p1.IsActive(int(x), int(y)) {
						p.Hit(int(x), int(y), p1.GetObjType())
						p1.Hit(int(x), int(y), p.GetObjType())
						//world.Add(int(x), int(y), 0xFFFFF, 0xFFFFFF, 0, 0xFFFFFF)
						break
					}
				}
			}
		}
	}
}
