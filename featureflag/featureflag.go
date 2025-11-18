package featureflag

import (
	"github.com/aidapedia/gdk/featureflag/module"
	"github.com/aidapedia/gdk/featureflag/module/file"
)

type Option struct {
	address string
	module  module.Module
	prefix  string
}

// New creates a new feature flag module.
func New(opt Option) module.Interface {
	switch opt.module {
	case module.FileModule:
		return file.New(opt.address, opt.prefix)
	default:
		return nil
	}
}
