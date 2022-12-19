package main

import (
	"github.com/ByteArena/box2d"
	physics "github.com/Gregmus2/simple-engine-physics-module"
	"github.com/Gregmus2/simple-engine/graphics"
)

type Food struct {
	phys   *physics.Object
	color  graphics.Color
	Radius float32
}

func (f *ObjectFactory) NewFood(x, y float64) *Food {
	radius := float64(2)
	body := f.phys.SimpleBody(x, y, box2d.B2BodyType.B2_staticBody)
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(radius / f.pCfg.Scale)
	object := f.pFactory.NewObject(body, &shape, 1.0)
	object.Fixture.SetUserData([]Tag{FoodTag})

	return &Food{
		phys:   object,
		Radius: float32(radius),
		color:  graphics.Yellow(),
	}
}

func (o *Food) Draw(scale, offsetX, offsetY float32) error {
	pos := o.phys.Body.GetPosition()
	graphics.Circle((float32(pos.X)+offsetX)*scale, (float32(pos.Y)+offsetY)*scale, o.Radius, o.color)

	return nil
}
