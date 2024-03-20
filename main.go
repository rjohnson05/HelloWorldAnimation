package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math/rand/v2"
	"time"
)

/*
Serves as the main controller for the "Hello, World!" text animator. This program allows a user to create as many appearances
of "Hello, World!" as they desire. These texts may be animated, giving the appearance that they are bouncing off the edges
of the window. At any time, the animation of these texts may be stopped with the click of a button, and all text may be
cleared from the screen.

Author: Ryan Johnson
*/

type ContainerData struct {
	colorList                 []color.Color
	canvasWidth, canvasHeight int
	animating                 bool
	stopAnimationChan         chan int

	window           fyne.Window
	textObjects      []fyne.CanvasObject
	textContainer    fyne.CanvasObject
	buttonsContainer fyne.CanvasObject
}

// Continually monitors the text objects and creates a new animation whenever they hit the window boundaries, moving the
// text objects to a new location on one of the borders. This animation continues until a stop signal is received, at which
// point the animation ceases.
func (data *ContainerData) animateText() {
	textObjects := data.textObjects

	// Sets the boundaries of the animation space, ensuring that no text moves outside the window. This accounts for the
	// size of the text itself.
	rightBorder := data.canvasWidth - int(textObjects[0].Size().Width)
	bottomBorder := data.canvasHeight - 3*int(textObjects[0].Size().Height)
	xBorderCoords := []int{0, rightBorder}
	yBorderCoords := []int{0, bottomBorder}

	for data.animating {
		for i := 0; i < len(textObjects); i++ {
			text := textObjects[i]
			var move *fyne.Animation

			// Creates a new animation if the text object has hit one of the window boundaries
			if (text.Position().X <= 0) || (text.Position().X >= float32(rightBorder)) ||
				(text.Position().Y <= 0) || (text.Position().Y >= float32(bottomBorder)) {
				vertHorizPick := rand.IntN(2)
				if vertHorizPick == 0 {
					// Create new animation to the top or bottom borders
					newXCoord := rand.IntN(rightBorder)
					newYCoordIndex := rand.IntN(len(yBorderCoords))
					move = canvas.NewPositionAnimation(fyne.NewPos(text.Position().X, text.Position().Y), fyne.NewPos(float32(newXCoord), float32(yBorderCoords[newYCoordIndex])), time.Second, text.Move)
				} else {
					// Create new animation to the left or right borders
					newXCoordIndex := rand.IntN(len(xBorderCoords))
					newYCoord := rand.IntN(bottomBorder)
					move = canvas.NewPositionAnimation(fyne.NewPos(text.Position().X, text.Position().Y), fyne.NewPos(float32(xBorderCoords[newXCoordIndex]), float32(newYCoord)), time.Second, text.Move)
				}
				move.Start()
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Creates a new text object reading "Hello, World!" This object is automatically placed in the upper-left corner of the
// window until animation begins.
func (data *ContainerData) createNewText() {
	text := canvas.NewText("Hello, World!", color.Black)
	text.Resize(fyne.NewSize(90, 20))
	data.textObjects = append(data.textObjects, text)

	data.refreshTextContainer()

	// If the text is being animated, the animation must be restarted to allow the new text object to begin the
	// animation.
	if data.animating {
		data.stopAnimationChannel()
		data.startAnimationChannel()
	}
}

// Creates a goroutine to animate the text, allowing the program to still be responsive. A channel is also created to
// listen for when the animation should be stopped.
func (data *ContainerData) startAnimationChannel() {
	//data.stopAnimationChan = make(chan int)
	data.animating = true
	go data.animateText()
}

// Sets the animating boolean to "false", triggering the animation loop to stop.
func (data *ContainerData) stopAnimationChannel() {
	data.animating = false
}

// Creates a new container containing all the old texts plus a newly created text. This must run every time a new text
// is created in order for the new text to be rendered to the screen.
func (data *ContainerData) refreshTextContainer() {
	data.textContainer = container.NewWithoutLayout(data.textObjects...)

	data.window.SetContent(container.NewVBox(data.buttonsContainer, data.textContainer))
}

// Removes all text objects from the screen, stopping the animation if necessary.
func (data *ContainerData) clearText() {
	data.animating = false
	data.textObjects = nil
	data.refreshTextContainer()
}

// Creates the buttons at the top of the screen for user interaction.
func (data *ContainerData) createButtonsContainer() {
	// Adds a new "Hello, World!" text object
	addHelloButton := widget.NewButton("Add Hello", func() {
		data.createNewText()
	})

	// Removes all text objects from the screen
	clearTextButton := widget.NewButton("Clear Text", func() {
		data.clearText()
	})

	// Begins the animation of all text objects
	startAnimationButton := widget.NewButton("Start Animation", func() {
		if !data.animating {
			data.startAnimationChannel()
		}
	})

	// Stops the animation of all text objects
	stopAnimationButton := widget.NewButton("Stop Animation", func() {
		if data.animating {
			data.stopAnimationChannel()
		}
	})

	data.buttonsContainer = container.NewHBox(addHelloButton, startAnimationButton, stopAnimationButton, clearTextButton)
}

func main() {
	app := app.New()

	ContainerData := ContainerData{}
	ContainerData.canvasWidth = 400
	ContainerData.canvasHeight = 400
	ContainerData.animating = false

	ContainerData.window = app.NewWindow("Hello Go")
	ContainerData.window.Resize(fyne.NewSize(400, 400))

	ContainerData.createButtonsContainer()
	ContainerData.refreshTextContainer()

	ContainerData.window.ShowAndRun()
}
