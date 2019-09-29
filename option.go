package fsm

type Options struct {
	SkipGuards bool
}

type Option func(*Options)

func SkipGuard(value bool) Option {
	return func(args *Options) {
		args.SkipGuards = value
	}
}
