package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

const (
	width  int    = 20
	weight int    = 20
	style  string = "border:1px solid #000000;"
)

type Point struct {
	X int
	Y int
}

var (
	canvas    js.Value
	ctx       js.Value
	bodySlice []*Point
	bodyMap   = make(map[string]*Point)
	position  *Point
	food      *Point
	direction = "ArrowRight"
	stop      bool
)

func getDocumentDom() js.Value {
	return js.Global().Get("document")
}

func jsAlert() js.Value {
	return js.Global().Get("alert")
}

func getElementByID(id string) js.Value {
	return getDocumentDom().Call("getElementById", id)
}

func registerCallbackFunc() {
	keyboard := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		key := args[0].Get("key").String()
		if key == "ArrowUp" {
			direction = key
			stop = false
		} else if key == "ArrowDown" {
			direction = key
			stop = false
		} else if key == "ArrowLeft" {
			direction = key
			stop = false
		} else if key == "ArrowRight" {
			direction = key
			stop = false
		} else if key == " " {
			stop = !stop
		}
		return nil
	})

	getDocumentDom().Call("addEventListener", "keypress", keyboard)
	getDocumentDom().Set("onkeypress", keyboard)
	getDocumentDom().Set("onkeyup", keyboard)

	canvas = getElementByID("canvas")
	ctx = canvas.Call("getContext", "2d")
	canvas.Set("height", width*weight)
	canvas.Set("width", width*weight)
	canvas.Set("style", style)
}

func getRandomNum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	if min >= max || max == 0 {
		return max
	}
	return rand.Intn(max-min+1) + min
}

func initSnake() {
	p1 := Point{X: getRandomNum(width/4, width*3/4), Y: getRandomNum(width/4, width*3/4)}
	point1 := fmt.Sprintf("%d-%d", p1.X, p1.Y)
	bodyMap[point1] = &p1
	bodySlice = append(bodySlice, &p1)

	p2 := Point{X: p1.X + 1, Y: p1.Y}
	point2 := fmt.Sprintf("%d-%d", p2.X, p2.Y)
	bodyMap[point2] = &p2
	bodySlice = append(bodySlice, &p2)
	position = &p2
}

func initFood() {
	for true {
		food = &Point{
			X: getRandomNum(0, width-1),
			Y: getRandomNum(0, width-1),
		}
		point := fmt.Sprintf("%d-%d", food.X, food.Y)
		if _, exist := bodyMap[point]; !exist {
			break
		}
	}
}

func drawCanvas() {
	ctx.Call("clearRect", 0, 0, width*weight, width*weight)

	for i := 0; i < width; i++ {
		for j := 0; j < width; j++ {
			ctx.Call("strokeRect", i*weight, j*weight, weight, weight)
		}
	}
}

func drawSnake() {
	for _, p := range bodySlice {
		ctx.Call("fillRect", p.X*weight, p.Y*weight, weight, weight)
	}
}

func flashFood() {
	ctx.Call("fillRect", food.X*weight, food.Y*weight, weight, weight)
}

func flashPosition() {
	ctx.Call("fillRect", position.X*weight, position.Y*weight, weight, weight)
}

func move() {
	if stop {
		return
	}
	if direction == "ArrowUp" {
		position = &Point{X: position.X, Y: position.Y - 1}
	} else if direction == "ArrowDown" {
		position = &Point{X: position.X, Y: position.Y + 1}
	} else if direction == "ArrowLeft" {
		position = &Point{X: position.X - 1, Y: position.Y}
	} else if direction == "ArrowRight" {
		position = &Point{X: position.X + 1, Y: position.Y}
	}

	point := fmt.Sprintf("%d-%d", position.X, position.Y)

	if _, exist := bodyMap[point]; exist || position.X < 0 || position.Y < 0 || position.X >= width || position.Y >= width {
		stop = true
		jsAlert().Invoke("You lose!!!")
		return
	}

	var (
		start = 0
		end   = len(bodySlice)
	)
	if position.X == food.X && position.Y == food.Y {
		// eat
		initFood()
	} else {
		// move
		start = 1
		delBodySlice := bodySlice[0:1]
		delPoint := fmt.Sprintf("%d-%d", delBodySlice[0].X, delBodySlice[0].Y)
		delete(bodyMap, delPoint)
		ctx.Call("clearRect", delBodySlice[0].X*weight, delBodySlice[0].Y*weight, weight, weight)
		ctx.Call("strokeRect", delBodySlice[0].X*weight, delBodySlice[0].Y*weight, weight, weight)
	}
	end = end + 1
	bodyMap[point] = position
	bodySlice = append(bodySlice, position)
	bodySlice = bodySlice[start:end]
}

func run() {
	for i := 0; true; i++ {
		time.Sleep(500 * time.Millisecond)
		drawCanvas()
		drawSnake()
		if i%2 == 0 {
			flashFood()
		} else {
			flashPosition()
			move()
		}
		drawSnake()
	}
}

func init() {
	initSnake()
	initFood()
	registerCallbackFunc()

	go run()
}

func main() {
	select {}
}

// $ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" static
// GOARCH=wasm GOOS=js go build -o static/snake.wasm snake.go
