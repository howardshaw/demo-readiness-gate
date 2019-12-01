package options

import (
	"github.com/spf13/pflag"
)

// Options contains everything necessary to create and run controller-manager.
type Options struct {
	Namespace string
}

// AddFlags adds flags to fs and binds them to options.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Namespace, "namespace", "", "The namespace of pod to watch.")
}

func NewOptions() *Options {
	return &Options{
	}
}
