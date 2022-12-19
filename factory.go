package main

import (
	physics "github.com/Gregmus2/simple-engine-physics-module"
)

type ObjectFactory struct {
	pFactory *physics.ObjectFactory
	phys     *physics.Physics
	pCfg     *physics.Config
}

func NewObjectFactory(f *physics.ObjectFactory, p *physics.Physics, pCfg *physics.Config) *ObjectFactory {
	return &ObjectFactory{pFactory: f, phys: p, pCfg: pCfg}
}
