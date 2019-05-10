// raffle is a GopherJS application that generates Twitter-based entries into a simple raffle
//
// For more details see github.com/go-london-user-group/raffle
package main

import (
	"myitcv.io/react"

	"honnef.co/go/js/dom"
)

//go:generate reactGen

var document = dom.GetWindow().Document()

func main() {
	domTarget := document.GetElementByID("app")

	react.Render(App(), domTarget)
}
