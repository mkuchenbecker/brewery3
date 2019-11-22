package thermometer

import "github.com/mkuchenbecker/brewery3/data/datasink"

type Thermometer interface {
	Read() error
	Sink() datasink.DataSink
	WithWink(datasink.DataSink) Thermometer
}
