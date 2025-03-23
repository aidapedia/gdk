package mask

import (
	"reflect"

	masker "github.com/ggwhite/go-masker/v2"
)

type Mask struct {
	*masker.MaskerMarshaler
}

type Maskers struct {
	Type masker.MaskerType
	Mask masker.Masker
}

func New(opts ...Option) *Mask {
	mask := &Mask{}
	mask.MaskerMarshaler = masker.NewMaskerMarshaler()
	list := mask.List()
	for _, v := range list {
		mask.Unregister(v)
	}
	mask.Register(masker.MaskerTypeNone, &masker.NoneMasker{})
	for _, opt := range opts {
		opt.Apply(mask)
	}
	return mask
}

func NewDefault(opts ...Option) *Mask {
	mask := &Mask{}
	// Default maskers are:
	//   - NoneMasker
	//   - PasswordMasker
	//   - NameMasker
	//   - AddressMasker
	//   - EmailMasker
	//   - MobileMasker
	//   - TelephoneMasker
	//   - IDMasker
	//   - CreditMasker
	//   - URLMasker
	mask.MaskerMarshaler = masker.NewMaskerMarshaler()
	for _, opt := range opts {
		opt.Apply(mask)
	}
	return mask
}

// Mask all struct fields
func (m *Mask) MaskStruct(val interface{}) (interface{}, error) {
	return m.Struct(val)
}

// Mask all map values
func (m *Mask) MaskMap(val map[string]interface{}) (map[string]interface{}, error) {
	for k, v := range val {
		// if type struct or map[string]interface{}
		// then call MaskStruct or MaskMap
		switch v.(type) {
		case map[string]interface{}:
			val[k], _ = m.MaskMap(v.(map[string]interface{}))
			continue
		case interface{}:
			typ := reflect.ValueOf(v).Kind()
			if typ == reflect.Struct {
				maskVal, err := m.MaskStruct(v)
				if err != nil {
					return nil, err
				}
				val[k] = maskVal
				continue
			}
		}
		mask, err := m.Get(masker.MaskerType(k))
		if err != nil {
			continue
		}
		val[k] = mask.Marshal("*", v.(string))
	}
	return val, nil
}
