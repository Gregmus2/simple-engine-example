package main

import (
	"github.com/ByteArena/box2d"
	"github.com/Gregmus2/nnga"
	physics "github.com/Gregmus2/simple-engine-physics-module"
	"github.com/Gregmus2/simple-engine/common"
	"github.com/Gregmus2/simple-engine/graphics"
	"github.com/go-gl/gl/v4.6-core/gl"
	"math"
)

type Agent struct {
	phys      *physics.Object
	prog      *graphics.Program
	shape     *graphics.ShapeHelper
	color     graphics.Color
	lineColor graphics.Color
	person    *nnga.Person
	cursor    *common.Pos

	targetPos box2d.B2Vec2
	distance  float64
	Radius    float32
}

func (f *ObjectFactory) NewAgent(x, y float64, p *nnga.Person) *Agent {
	body := f.phys.SimpleBody(x, y, box2d.B2BodyType.B2_dynamicBody)
	body.SetFixedRotation(false)
	body.SetLinearDamping(10.0)
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(5 / f.pCfg.Scale)
	object := f.pFactory.NewObject(body, &shape, 1.0)

	angle := body.GetAngle()
	pos := body.GetPosition()
	radius := float32(5)
	x1, y1 := float32(pos.X)*f.gCfg.Graphics.Scale, float32(pos.Y)*f.gCfg.Graphics.Scale
	x2 := x1 + (radius * float32(math.Cos(angle)))
	y2 := y1 + (radius * float32(math.Sin(angle)))

	return &Agent{
		phys:      object,
		prog:      graphics.NewProgram(),
		shape:     f.shape,
		color:     graphics.Blue(),
		lineColor: graphics.White(),
		person:    p,
		cursor:    &common.Pos{X: x2, Y: y2},
		Radius:    radius,
	}
}

func (a *Agent) Draw(scale, offsetX, offsetY float32) error {
	pos := a.phys.Body.GetPosition()
	a.prog.ApplyProgram(a.color)
	a.shape.Circle((float32(pos.X)+offsetX)*scale, (float32(pos.Y)+offsetY)*scale, a.Radius)
	gl.UseProgram(0)

	angle := a.phys.Body.GetAngle()
	x, y := (float32(pos.X)+offsetX)*scale, (float32(pos.Y)+offsetY)*scale
	x2 := x + (a.Radius * float32(math.Cos(angle)))
	y2 := y + (a.Radius * float32(math.Sin(angle)))
	a.prog.ApplyProgram(a.lineColor)
	a.shape.Line(x, y, x2, y2)
	gl.UseProgram(0)

	return nil
}
