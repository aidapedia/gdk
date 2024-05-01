package mask

import (
	masker "github.com/ggwhite/go-masker/v2"
)

type Mask struct {
	*masker.MaskerMarshaler
}

func New() *Mask {
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
	return mask
}

// Mask all struct fields
func (m *Mask) MaskStruct(val interface{}) (interface{}, error) {
	return m.Struct(val)
}
