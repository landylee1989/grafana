// Package dataframe provides the DataFrame type.
// The DataFrame type is used to hold data returned from Grafana Datasources.
// This type is meant to tightly integrated with the DataFrame type in Grafana's Frontend.
package dataframe

import (
	"fmt"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
)

// DataFrame holds Table data.
type DataFrame struct {
	Schema  Schema
	Type    FrameType
	Records []Fields
}

// FrameType indicates the type of data the Dataframe holds
type FrameType int

const (
	// NumericFrame indicates the Dataframe holds numeric values.
	NumericFrame FrameType = iota

	// TimeSeriesFrame indicates the Dataframe holds timeseries data.
	TimeSeriesFrame

	// HistogramFrame indicates the Dataframe holds histograms data.
	HistogramFrame

	// OtherFrame indicates the DataFrame holds mixed or another data type.
	OtherFrame
)

func (ft FrameType) String() string {
	switch ft {
	case NumericFrame:
		return "Number"
	case TimeSeriesFrame:
		return "TimeSeries"
	case HistogramFrame:
		return "Histogram"
	default:
		return "Other"
	}
}

func (ft FrameType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, ft.String())), nil
}

// DataFrames is a collection of DataFrames uniquely identified by key.
type DataFrames []DataFrame

// Fields is a slice fo Field.
type Fields []Field

// Field represents a unique field within a dataframe identified by its column and record position.
type Field interface{}

// ToArrow is an experiment to create an arrow Table from the dataframe
func (d *DataFrame) ToArrow() *array.TableReader {
	arrowFields := make([]arrow.Field, len(d.Schema))
	for i, cs := range d.Schema {
		fmt.Println(cs.ArrowType().ID())
		arrowFields[i] = arrow.Field{Name: cs.GetName(), Type: cs.ArrowType()}
	}
	schema := arrow.NewSchema(arrowFields, nil)

	pool := memory.NewGoAllocator()

	rb := array.NewRecordBuilder(pool, schema)
	defer rb.Release()

	records := []array.Record{}
	for rowIdx, row := range d.Records {
		for fieldIdx, field := range row {
			switch arrowFields[fieldIdx].Type.(type) {
			case *arrow.StringType:
				rb.Field(fieldIdx).(*array.StringBuilder).Append(*(field.(*string)))
			}
			//rb.Field(fieldIdx).(cs[fieldIdx].ArrowBuilderType()).AppendValues(field)
			//rb.Field(fieldIdx).(cs[fieldIdx].ArrowBuilderType()).AppendValues(field)
			//rb.Field(fieldIdx).(*(arrowFields[fieldIdx].Type)).Append(field)
		}
		rec := rb.NewRecord()
		defer rec.Release()
		records[rowIdx] = rb.NewRecord()
	}
	table := array.NewTableFromRecords(schema, records)
	defer table.Release()
	return array.NewTableReader(table, 0)

}
