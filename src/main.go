package main

import (
	"math"
	"math/rand"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowWidth  = 1280
	WindowHeight = 720
	StarCapacity = 800
	Velocity     = 10.0
)

// this struct represents a star with its x, y, z coordinates and previous z coordinate (pz).
type Star struct {
	x, y, z, pz float64
}

var stars [StarCapacity]Star

// checkError function checks if there is an error and logs it.
func checkError(err error) {
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Failed to create renderer: %s\n", err)
		os.Exit(1)
	}
}

// this function draws a circle on the renderer at position (x, y) with the specified radius.
func renderDrawCircle(renderer *sdl.Renderer, x, y, radius int32) {
	for i := 0; i < 360; i++ {
		for j := 0; j < int(radius); j++ {
			xPos := int32(math.Cos(float64(i))*float64(j)) + x
			yPos := int32(math.Sin(float64(i))*float64(j)) + y
			renderer.DrawPoint(xPos, yPos)
		}
	}
}

func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		checkError(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Goverse", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	checkError(err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkError(err)
	defer renderer.Destroy()

	// Initialize the stars.
	for i := 0; i < StarCapacity; i++ {
		z := float64(rand.Intn(WindowWidth))
		stars[i] = Star{
			x:  float64(WindowWidth - rand.Intn(WindowWidth*2)),
			y:  float64(WindowHeight - rand.Intn(WindowHeight*2)),
			z:  z,
			pz: z,
		}
	}

	// Event handling loop.
	var event sdl.Event
	running := true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Draw the circles.
		renderer.SetDrawColor(255, 255, 255, 255)
		for i := 0; i < StarCapacity; i++ {
			star := &stars[i]

			// Calculate the projected x and y coordinates of the star.
			xx := (star.x / star.z) * float64(WindowWidth)
			yy := (star.y / star.z) * float64(WindowHeight)

			// The previous projected x and y coordinates.
			pxx := (star.x / star.pz) * float64(WindowWidth)
			pyy := (star.y / star.pz) * float64(WindowHeight)

			// The radius of the star based on its z coordinate.
			radius := 5 - ((star.z / float64(WindowWidth)) * 5)

			// The final x and y coordinates of the star on the window.
			destX := int32(xx) + WindowWidth/2
			destY := int32(yy) + WindowHeight/2

			// The final previous x and y coordinates of the star on the window.
			pDestX := int32(pxx) + WindowWidth/2
			pDestY := int32(pyy) + WindowHeight/2

			// Update the previous z coordinate of the star.
			star.pz = star.z

			// Update the z coordinate of the star, moving it towards the viewer.
			star.z += Velocity * -1

			// If the star gets too close to the viewer, reset it.
			if star.z < 1 {
				star.z = float64(WindowWidth)
				star.x = float64(WindowWidth - rand.Intn(WindowWidth*2))
				star.y = float64(WindowHeight - rand.Intn(WindowHeight*2))
				star.pz = star.z
			}

			// If the star is within the window bounds, draw it.
			if destX > 0 && destX < WindowWidth &&
				destY > 0 && destY < WindowHeight {
				renderDrawCircle(renderer, destX, destY, int32(radius))
				renderer.DrawLine(pDestX, pDestY, destX, destY)
			}
		}

		renderer.Present()

		sdl.Delay(16) // 1000ms / 16ms = 62.5 fps
	}
}
