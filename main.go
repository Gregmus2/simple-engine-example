package main

import (
	"github.com/Gregmus2/simple-engine"
	physics "github.com/Gregmus2/simple-engine-physics-module"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	c, err := engine.BuildContainer()
	if err != nil {
		logrus.WithError(err).Fatal("error building DI container")
	}

	err = buildContainer(c)
	if err != nil {
		logrus.WithError(err).Fatal("error building DI container")
	}

	if err := c.Invoke(func(app *engine.App, agents *Agents) {
		app.InitWithScene(agents)
		app.Loop()
	}); err != nil {
		logrus.Fatal(err)
	}
}

func buildContainer(c *dig.Container) error {
	if err := c.Provide(NewObjectFactory); err != nil {
		return err
	}
	if err := c.Provide(NewAgents); err != nil {
		return err
	}

	return physics.BuildContainer(c)
}
