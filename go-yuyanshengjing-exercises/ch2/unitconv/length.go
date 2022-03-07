package unitconv

import "fmt"

// Foot 英尺
type Foot float64

// Meters 米
type Meters float64

// FtoMConversionRate 折算率
const FtoMConversionRate = 0.3048

func (f Foot) String() string   { return fmt.Sprintf("%g ft", f) }
func (m Meters) String() string { return fmt.Sprintf("%g m", m) }

// FtoM converts a Foot temperature Meters Celsius.
func FtoM(f Foot) Meters { return Meters(f * FtoMConversionRate) }

// MtoF converts a Meters temperature Foot Celsius.
func MtoF(m Meters) Foot { return Foot(m / FtoMConversionRate) }
