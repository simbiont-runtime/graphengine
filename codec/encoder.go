// ---

package codec

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/types"
)

// PropertyEncoder is used to encode datums into value bytes.
type PropertyEncoder struct {
	rowBytes

	values []datum.Datum
}

// Encode encodes properties into a value bytes.
func (e *PropertyEncoder) Encode(buf []byte, labelIDs, propertyIDs []uint16, values []datum.Datum) ([]byte, error) {
	e.reform(labelIDs, propertyIDs, values)
	for i, value := range e.values {
		err := e.encodeDatum(value)
		if err != nil {
			return nil, err
		}
		e.offsets[i] = uint16(len(e.data))
	}
	return e.toBytes(buf[:0]), nil
}

func (e *PropertyEncoder) encodeDatum(value datum.Datum) error {
	// Put the type information first.
	e.data = append(e.data, byte(value.Type()))
	switch value.Type() {
	case types.Int:
		e.data = encodeInt(e.data, datum.AsInt(value))
	case types.Float:
		e.data = EncodeFloat(e.data, datum.AsFloat(value))
	case types.String:
		e.data = append(e.data, datum.AsBytes(value)...)
	case types.Date:
		e.data = encodeDate(e.data, datum.AsDate(value))
	default:
		return fmt.Errorf("unsupported encode type %T", value)
	}
	return nil
}

func encodeInt(buf []byte, iVal int64) []byte {
	var tmp [8]byte
	if int64(int8(iVal)) == iVal {
		buf = append(buf, byte(iVal))
	} else if int64(int16(iVal)) == iVal {
		binary.LittleEndian.PutUint16(tmp[:], uint16(iVal))
		buf = append(buf, tmp[:2]...)
	} else if int64(int32(iVal)) == iVal {
		binary.LittleEndian.PutUint32(tmp[:], uint32(iVal))
		buf = append(buf, tmp[:4]...)
	} else {
		binary.LittleEndian.PutUint64(tmp[:], uint64(iVal))
		buf = append(buf, tmp[:8]...)
	}
	return buf
}

func encodeDate(buf []byte, date *datum.Date) []byte {
	return encodeInt(buf, int64(date.UnixEpochDays()))
}

func (e *PropertyEncoder) reform(labelIDs, propertyIDs []uint16, values []datum.Datum) {
	e.labelIDs = append(e.labelIDs[:0], labelIDs...)
	e.propertyIDs = append(e.propertyIDs[:0], propertyIDs...)
	e.offsets = make([]uint16, len(e.propertyIDs))
	e.data = e.data[:0]
	e.values = e.values[:0]

	e.values = append(e.values, values...)

	sort.Slice(e.labelIDs, func(i, j int) bool {
		return e.labelIDs[i] < e.labelIDs[j]
	})
	sort.Sort(&propertySorter{
		propertyIDs: e.propertyIDs,
		values:      e.values,
	})
}

type propertySorter struct {
	propertyIDs []uint16
	values      []datum.Datum
}

// Less implements the Sorter interface.
func (ps *propertySorter) Less(i, j int) bool {
	return ps.propertyIDs[i] < ps.propertyIDs[j]
}

// Len implements the Sorter interface.
func (ps *propertySorter) Len() int {
	return len(ps.propertyIDs)
}

// Swap implements the Sorter interface.
func (ps *propertySorter) Swap(i, j int) {
	ps.propertyIDs[i], ps.propertyIDs[j] = ps.propertyIDs[j], ps.propertyIDs[i]
	ps.values[i], ps.values[j] = ps.values[j], ps.values[i]
}
