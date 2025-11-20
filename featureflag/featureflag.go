package featureflag

import (
	"github.com/aidapedia/gdk/featureflag/module"
	"github.com/aidapedia/gdk/featureflag/module/consul"
	"github.com/aidapedia/gdk/featureflag/module/file"
)

type Option struct {
	Address string
	Module  module.Module
	Prefix  string
}

// New creates a new feature flag module.
func New(opt Option) module.Interface {
	switch opt.Module {
	case module.FileModule:
		return file.New(opt.Address, opt.Prefix)
	case module.ConsulModule:
		return consul.New(opt.Address, opt.Prefix)
	default:
		return nil
	}
}
