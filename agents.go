package main

import (
	"github.com/ByteArena/box2d"
	"github.com/Gregmus2/nnga"
	"github.com/Gregmus2/simple-engine/common"
	"github.com/Gregmus2/simple-engine/graphics"
	"github.com/Gregmus2/simple-engine/scenes"
	"github.com/patrikeh/go-deep"
	"math"
	"math/rand"
	"time"
)

type Agents struct {
	scenes.Base
	factory *ObjectFactory
	agents  []*Agent
	food    []*Food
	ga      *nnga.GA

	fixturesToMove []*box2d.B2Fixture
}

func NewAgents(base scenes.Base, f *ObjectFactory, w *box2d.B2World) *Agents {
	agents := &Agents{
		Base:           base,
		factory:        f,
		fixturesToMove: make([]*box2d.B2Fixture, 0),
	}
	w.SetContactListener(agents)

	return agents
}

func (d *Agents) Init() {
	// world set contact callback

	rand.Seed(time.Now().UTC().UnixNano())

	d.walls()

	d.ga = nnga.NewGA(200, &deep.Config{
		/* Input dimensionality */
		Inputs: 2,
		/* Two hidden layers consisting of two neurons each, and a single output */
		Layout: []int{2, 4, 2},
		/* Activation functions: Sigmoid, Tanh, ReLU, Linear */
		Activation: deep.ActivationSigmoid,
		/* Determines output layer activation & loss function:
		ModeRegression: linear outputs with MSE loss
		ModeMultiClass: softmax output with Cross Entropy loss
		ModeMultiLabel: sigmoid output with Cross Entropy loss
		ModeBinary: sigmoid output with binary CE loss */
		Mode: deep.ModeBinary,
		/* Weight initializers: {deep.NewNormal(μ, σ), deep.NewUniform(μ, σ)} */
		Weight: deep.NewNormal(1.0, 0.0),
		/* Apply bias */
		Bias: true,
	}, &nnga.Coefficients{
		Scale:                   1,
		Selection:               0.2,
		MutationClassic:         0.1,
		MutationGrowth:          2,
		MutationGenesMaxPercent: 0.2,
		MutationOffset:          0.1,
	})

	for i := 0; i < 50; i++ {
		food := d.factory.NewFood(float64(rand.Intn(common.Config.Window.W)), float64(rand.Intn(common.Config.Window.H)))
		d.DrawObjects.Put(food)
		d.food = append(d.food, food)
	}

	for _, person := range d.ga.Persons {
		agent := d.factory.NewAgent(float64(rand.Intn(common.Config.Window.W)), float64(rand.Intn(common.Config.Window.H)), person)
		d.DrawObjects.Put(agent)
		d.agents = append(d.agents, agent)
	}
}

func (d *Agents) PreUpdate() {
	d.Base.PreUpdate()
	for _, agent := range d.agents {
		pos := agent.phys.Body.GetPosition()
		targetPos := d.food[0].phys.Body.GetPosition()
		min := math.Sqrt(math.Pow(pos.X-targetPos.X, 2) + math.Pow(pos.Y-targetPos.Y, 2))
		for _, piece := range d.food {
			foodPos := piece.phys.Body.GetPosition()
			distance := math.Sqrt(math.Pow(pos.X-foodPos.X, 2) + math.Pow(pos.Y-foodPos.Y, 2))
			if distance < min {
				min = distance
				targetPos = foodPos
			}
		}

		// todo Math folding plugin
		a := min
		b := math.Sqrt(math.Pow(pos.X-float64(agent.cursor.X), 2) + math.Pow(pos.Y-float64(agent.cursor.Y), 2))
		c := math.Sqrt(math.Pow(targetPos.X-float64(agent.cursor.X), 2) + math.Pow(targetPos.Y-float64(agent.cursor.Y), 2))
		angle := math.Acos((math.Pow(a, 2) + math.Pow(b, 2) - math.Pow(c, 2)) / (2 * a * b))
		output := agent.person.Predict([]float64{math.Cos(angle), min})
		dAngle, force := output[0], (output[1]-0.5)*0.1

		newAngle := agent.phys.Body.GetAngle() + dAngle
		xForce := force * math.Cos(newAngle)
		yForce := force * math.Sin(newAngle)
		agent.phys.Body.ApplyForceToCenter(box2d.MakeB2Vec2(xForce, yForce), true)

		agent.targetPos = targetPos
		agent.distance = min
	}

	for _, fixture := range d.fixturesToMove {
		fixture.GetBody().SetTransform(box2d.MakeB2Vec2(float64(rand.Intn(common.Config.Window.W))/100, float64(rand.Intn(common.Config.Window.H))/100), 0)
	}
	if len(d.fixturesToMove) > 0 {
		d.fixturesToMove = make([]*box2d.B2Fixture, 0)
	}
}

func (d *Agents) Update() {
	for _, agent := range d.agents {
		pos := agent.phys.Body.GetPosition()
		distance := math.Sqrt(math.Pow(pos.X-agent.targetPos.X, 2) + math.Pow(pos.Y-agent.targetPos.Y, 2))
		agent.person.Score(agent.distance - distance)
	}

	d.ga.Evolve()
	for i, person := range d.ga.Persons {
		d.agents[i].person = person
	}
}

func (d *Agents) walls() {
	b := BoxModel{
		X:       graphics.PercentToPosX(50),
		Y:       0,
		W:       graphics.PercentToPosX(100),
		H:       1,
		T:       box2d.B2BodyType.B2_staticBody,
		Color:   graphics.White(),
		Density: 0,
	}
	d.DrawObjects.Put(d.factory.NewBox(b))

	b.Y = graphics.PercentToPosY(100)
	d.DrawObjects.Put(d.factory.NewBox(b))

	b.X, b.Y, b.W, b.H = 0, graphics.PercentToPosY(50), 1, graphics.PercentToPosY(100)
	d.DrawObjects.Put(d.factory.NewBox(b))

	b.X = graphics.PercentToPosX(100)
	d.DrawObjects.Put(d.factory.NewBox(b))
}

func (d *Agents) BeginContact(contact box2d.B2ContactInterface) {
	d.handleContact(contact.GetFixtureA())
	d.handleContact(contact.GetFixtureB())
}

func (d *Agents) handleContact(fixture *box2d.B2Fixture) {
	userData := fixture.GetUserData()
	switch userData.(type) {
	case []Tag:
		tags := userData.([]Tag)
		for _, tag := range tags {
			if tag == FoodTag {
				d.fixturesToMove = append(d.fixturesToMove, fixture)
			}
		}
	}
}

func (d *Agents) EndContact(contact box2d.B2ContactInterface)                                 {}
func (d *Agents) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold)     {}
func (d *Agents) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {}
