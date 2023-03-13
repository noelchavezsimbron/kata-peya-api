package pet

import (
	"errors"
	"fmt"
	"strings"
)

type (
	Vaccine string

	Vaccines []Vaccine

	Pet struct {
		Id        int
		Name      string
		Vaccines  Vaccines
		AgeMonths int
	}
)

func (p *Pet) GetAgeFormatted() string {
	if p.AgeMonths < 12 {
		return fmt.Sprintf("%d months", p.AgeMonths)
	}

	if p.AgeMonths == 12 {
		return fmt.Sprintf("%d year", 1)
	}

	years := p.AgeMonths / 12

	return fmt.Sprintf("%d years", years)
}

func (v *Vaccines) Scan(src any) error {

	var bs []byte

	switch src.(type) {
	case []uint8:
		bs = src.([]byte)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for Vaccines")
	}

	if len(bs) > 0 {
		values := strings.Split(string(bs), ",")
		vs := make(Vaccines, len(values))
		for i, val := range values {
			vs[i] = Vaccine(val)
		}
		*v = vs
	}
	return nil
}
