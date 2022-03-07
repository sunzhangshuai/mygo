package unitconv

import "fmt"

// Pound 磅
type Pound float64

// Kg 公斤
type Kg float64

// PtoKConversionRate 折算率
const PtoKConversionRate = 0.4535924

func (p Pound) String() string { return fmt.Sprintf("%g lb", p) }
func (k Kg) String() string    { return fmt.Sprintf("%g kg", k) }

// PtoK converts a Pound temperature Kg Celsius.
func PtoK(p Pound) Kg { return Kg(p * PtoKConversionRate) }

// KtoP converts a Kg temperature Pound Celsius.
func KtoP(k Kg) Pound { return Pound(k / PtoKConversionRate) }
