package solid

type Options struct {
	Center Vector
	Radius float64
	Size   Vector
	Slices int
	Stacks int
	Start  Vector
	End    Vector
}

func OptionsFrom(options []Option) *Options {
	o := &Options{
		Center: Vector{},
		Radius: 1,
		Size:   Vector{X: 2, Y: 2, Z: 2},
		Slices: 16,
		Stacks: 8,
		Start:  Vector{Y: -1},
		End:    Vector{Y: 1},
	}
	for _, option := range options {
		option(o)
	}
	return o
}

type Option func(*Options)

func Center(x, y, z float64) Option {
	option := func(config *Options) {
		config.Center = Vector{X: x, Y: y, Z: z}
	}
	return option
}

func Radius(radius float64) Option {
	option := func(config *Options) {
		config.Radius = radius
	}
	return option
}

func Size(x, y, z float64) Option {
	option := func(config *Options) {
		config.Size = Vector{X: x, Y: y, Z: z}
	}
	return option
}

func Slices(slices int) Option {
	option := func(config *Options) {
		config.Slices = slices
	}
	return option
}

func Stacks(stacks int) Option {
	option := func(config *Options) {
		config.Stacks = stacks
	}
	return option
}

func Start(x, y, z float64) Option {
	option := func(config *Options) {
		config.Start = Vector{X: x, Y: y, Z: z}
	}
	return option
}

func End(x, y, z float64) Option {
	option := func(config *Options) {
		config.End = Vector{X: x, Y: y, Z: z}
	}
	return option
}
