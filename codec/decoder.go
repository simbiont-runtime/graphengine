// ---

package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/simbiont-runtime/graphengine/datum"
	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/simbiont-runtime/graphengine/types"
)

// PropertyDecoder is used to decode value bytes into datum
type PropertyDecoder struct {
	rowBytes

	labels     []*model.LabelInfo
	properties []*model.PropertyInfo
}

func NewPropertyDecoder(labels []*model.LabelInfo, properties []*model.PropertyInfo) *PropertyDecoder {
	return &PropertyDecoder{
		labels:     labels,
		properties: properties,
	}
}

func (d *PropertyDecoder) Decode(rowData []byte) (map[uint16]struct{}, map[uint16]datum.Datum, error) {
	err := d.fromBytes(rowData)
	if err != nil {
		return nil, nil, err
	}

	labelIDs := make(map[uint16]struct{})
	for _, label := range d.labels {
		if d.hasLabel(uint16(label.ID)) {
			labelIDs[uint16(label.ID)] = struct{}{}
		}
	}

	row := make(map[uint16]datum.Datum)
	for _, property := range d.properties {
		idx := d.findProperty(property.ID)
		if idx >= 0 {
			propData := d.getData(idx)
			d, err := d.decodeColDatum(propData)
			if err != nil {
				return nil, nil, err
			}
			row[property.ID] = d
		}
	}

	return labelIDs, row, nil
}

func (d *PropertyDecoder) decodeColDatum(propData []byte) (datum.Datum, error) {
	var value datum.Datum
	typ := types.T(propData[0])
	switch typ {
	case types.Int:
		value = datum.NewInt(decodeInt(propData[1:]))
	case types.Float:
		_, v, err := DecodeFloat(propData[1:])
		if err != nil {
			return nil, err
		}
		value = datum.NewFloat(v)
	case types.String:
		value = datum.NewString(string(propData[1:]))
	case types.Date:
		value = decodeDate(propData[1:])
	default:
		// TODO: support more types
		return value, fmt.Errorf("unknown type %s", typ)
	}
	return value, nil
}

func decodeInt(val []byte) int64 {
	switch len(val) {
	case 1:
		return int64(int8(val[0]))
	case 2:
		return int64(int16(binary.LittleEndian.Uint16(val)))
	case 4:
		return int64(int32(binary.LittleEndian.Uint32(val)))
	default:
		return int64(binary.LittleEndian.Uint64(val))
	}
}

func decodeDate(val []byte) *datum.Date {
	return datum.NewDateFromUnixEpochDays(int32(decodeInt(val)))
}
