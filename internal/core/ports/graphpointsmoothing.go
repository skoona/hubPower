package ports

type GraphPointSmoothing interface {
	AddValue(value float64) float64
	SeriesName() string
	IsNil() bool
	String() string
}
