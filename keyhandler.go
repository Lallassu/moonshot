package main

import (
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
)

type KeyHandler struct {
	mouseX      int
	mouseY      int
	mouseDown   bool
	mouseButton int
	WireFrame   bool
	kDown       bool
	Debug       bool
	lastTime    time.Time
}

func (k *KeyHandler) MousePos(w *glfw.Window, xpos, ypos float64) {
	k.mouseX = int(xpos)
	k.mouseY = screenHeight - int(ypos)
	if k.mouseDown && k.mouseButton == 0 {
		Explode(float64(k.mouseX), float64(k.mouseY), 50)
	}
}

func (k *KeyHandler) MouseDown(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == 1 {
		k.mouseDown = true
	} else {
		k.mouseDown = false
	}
	k.mouseButton = int(button)
}

func (k *KeyHandler) Process(w *glfw.Window) {
	if glfw.Press == w.GetKey(GLKeys[conf.KeyMenu]) {
		if time.Since(k.lastTime).Milliseconds() > 200 {
			k.lastTime = time.Now()
			if menu.showAbout {
				menu.About()
			} else if menu.showSelectLevel {
				menu.SelectLevel()
			} else {
				menu.showMenu = !menu.showMenu
			}
		}
	}

	// If menu is up, only handle keys for the menu
	if menu.showMenu {
		if glfw.Press == w.GetKey(GLKeys[conf.KeyMenuUp]) {
			if time.Since(k.lastTime).Milliseconds() > 150 {
				k.lastTime = time.Now()
				menu.Up()
			}
		}

		if glfw.Press == w.GetKey(GLKeys[conf.KeyMenuDown]) {
			if time.Since(k.lastTime).Milliseconds() > 150 {
				k.lastTime = time.Now()
				menu.Down()
			}
		}

		if glfw.Press == w.GetKey(GLKeys[conf.KeyMenuSelect]) {
			if time.Since(k.lastTime).Milliseconds() > 150 {
				k.lastTime = time.Now()
				menu.Select()
			}
		}
		return
	}

	if glfw.Press == w.GetKey(GLKeys[conf.KeyLoadFuel]) {
		rocket.Boost()
	}

	if glfw.Press == w.GetKey(GLKeys[conf.KeyRelease]) {
		rocket.Release()
	}

	if glfw.Press == w.GetKey(GLKeys[conf.KeyDebugInfo]) {
		if time.Since(k.lastTime).Milliseconds() > 200 {
			k.lastTime = time.Now()
			k.Debug = !k.Debug
		}
	}

	if glfw.Press == w.GetKey(GLKeys[conf.KeyRespawn]) {
		if time.Since(k.lastTime).Milliseconds() > 200 {
			k.lastTime = time.Now()
			gameMap.CreateRocket()
		}
	}

	if glfw.Press == w.GetKey(GLKeys[conf.KeyWireframe]) {
		if time.Since(k.lastTime).Milliseconds() > 200 {
			k.lastTime = time.Now()
			k.WireFrame = !k.WireFrame
		}
	}
}
