package main

type particle struct {
	Phys
	z    float32
	r    float32
	g    float32
	b    float32
	a    float32
	size float64
}

// Update particle
func (p *particle) Update(dt float64) {
	p.Phys.Update(dt)
}

// Stop particle
func (p *particle) Stop() {
	p.life = 0
	if world.IsActive(int(p.x), int(p.y-1)) {
		world.Add(int(p.x), int(p.y), p.r, p.g, p.b, p.a)
	}

}
