package mask

import masker "github.com/ggwhite/go-masker/v2"

type Option interface {
	Apply(mask *Mask)
}

func WithRegisterMaskers(maskers ...Maskers) Option {
	return &withRegisterMaskers{maskers: maskers}
}

type withRegisterMaskers struct {
	maskers []Maskers
}

func (w *withRegisterMaskers) Apply(mask *Mask) {
	for _, masker := range w.maskers {
		mask.Register(masker.Type, masker.Mask)
	}
}

func WithUnregisterMaskers(maskers ...masker.MaskerType) Option {
	return &withUnregisterMaskers{maskers: maskers}
}

type withUnregisterMaskers struct {
	maskers []masker.MaskerType
}

func (w *withUnregisterMaskers) Apply(mask *Mask) {
	for _, masker := range w.maskers {
		mask.Unregister(masker)
	}
}
