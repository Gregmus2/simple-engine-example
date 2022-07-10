package main

import (
	"github.com/ByteArena/box2d"
	physics "github.com/Gregmus2/simple-engine-physics-module"
	"github.com/Gregmus2/simple-engine/graphics"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Box struct {
	w, h  float32
	phys  *physics.Object
	prog  *graphics.Program
	shape *graphics.ShapeHelper
	color graphics.Color
}

type BoxModel struct {
	X, Y, W, H float64
	T          uint8
	Color      graphics.Color
	Density    float64
}

func (f *ObjectFactory) NewBox(model BoxModel) *Box {
	body := f.phys.SimpleBody(model.X, model.Y, model.T)
	body.SetFixedRotation(false)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(model.W/f.pCfg.Scale/2, model.H/f.pCfg.Scale/2)
	object := f.pFactory.NewObject(body, &shape, model.Density)

	return &Box{
		phys:  object,
		w:     float32(model.W),
		h:     float32(model.H),
		prog:  f.program,
		shape: f.shape,
		color: graphics.White(),
	}
}

func (o *Box) Draw(scale, offsetX, offsetY float32) error {
	pos := o.phys.Body.GetPosition()
	o.prog.ApplyProgram(o.color)
	o.shape.Box((float32(pos.X)+offsetX)*scale-o.w/2, (float32(pos.Y)+offsetY)*scale+o.h/2, o.w, o.h)
	gl.UseProgram(0)

	return nil
}
