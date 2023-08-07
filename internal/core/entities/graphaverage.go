package entities

import (
	"fmt"
	"github.com/skoona/hubPower/internal/commons"
	"strings"
	"time"
)

type GraphAverage struct {
	seriesName  string
	graphPeriod time.Duration
	size        time.Duration
	dataPoints  []float64
}

//var _ (ports.GraphPointSmoothing) = (*GraphAverage)(nil)

// NewGraphAverage create a data queue of graphPeriod depth
// returns the straight average of the queue,
func NewGraphAverage(hostName, pointName string, graphPeriod time.Duration) *GraphAverage {
	return &GraphAverage{
		seriesName:  pointName,
		graphPeriod: graphPeriod,
		size:        1,
		dataPoints:  []float64{1.0}, // avoids first value being zero
	}
}

// AddValue adds the given float32 value into the queue
// and return the average value of the queue
// value queue's size is limited by graph period config value
func (g *GraphAverage) AddValue(value float64) float64 {
	if g.size >= g.graphPeriod {
		g.dataPoints = commons.ShiftSlice(value, g.dataPoints)
	} else {
		g.dataPoints = append(g.dataPoints, value)
	}

	g.size = time.Duration(len(g.dataPoints))

	return g.computeAverage()
}
func (g *GraphAverage) SeriesName() string {
	return strings.Clone(g.seriesName)
}
func (g *GraphAverage) computeAverage() float64 {
	var sum float64
	for _, fval := range g.dataPoints {
		sum = sum + fval
	}
	return (sum / float64(g.size))
}
func (g *GraphAverage) String() string {
	return fmt.Sprint("series:", g.seriesName, ", graphPeriod:", g.graphPeriod, ", dataPoints:", g.dataPoints)
}
func (g *GraphAverage) IsNil() bool {
	return (g == nil)
}
