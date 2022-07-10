package main

import (
	physics "github.com/Gregmus2/simple-engine-physics-module"
	"github.com/Gregmus2/simple-engine/common"
	"github.com/Gregmus2/simple-engine/graphics"
)

type ObjectFactory struct {
	pFactory *physics.ObjectFactory
	phys     *physics.Physics
	pCfg     *physics.Config
	gCfg     *common.Config
	shape    *graphics.ShapeHelper
	program  *graphics.Program
}

func NewObjectFactory(f *physics.ObjectFactory, p *physics.Physics, pCfg *physics.Config, shape *graphics.ShapeHelper, gCfg *common.Config, program *graphics.Program) *ObjectFactory {
	return &ObjectFactory{pFactory: f, phys: p, pCfg: pCfg, shape: shape, gCfg: gCfg, program: program}
}
