// ---

package codec

import "errors"

const (
	incomingEdgeSep = 'i'
	outgoingEdgeSep = 'o'
)

const EdgeKeyLen = 1 /*prefix*/ + 8 /*graphID*/ + 8 /*srcVertexID*/ + 1 /*edgeSep*/ + 8 /*dstVertexID*/

// IncomingEdgeKey encodes the incoming edge key.
//
// The key format is: ${Prefix}${GraphID}${DstVertexID}${IncomingEdgeSep}${SrcVertexID}.
func IncomingEdgeKey(graphID, srcVertexID, dstVertexID int64) []byte {
	result := make([]byte, 0, EdgeKeyLen)
	result = append(result, prefix...)
	result = EncodeInt(result, graphID)
	result = EncodeInt(result, dstVertexID)
	result = append(result, incomingEdgeSep)
	result = EncodeInt(result, srcVertexID)
	return result
}

// ParseIncomingEdgeKey parse the incoming edge key.
func ParseIncomingEdgeKey(key []byte) (graphID, srcVertexID, dstVertexID int64, err error) {
	if len(key) < EdgeKeyLen {
		return 0, 0, 0, errors.New("insufficient key length")
	}
	_, graphID, err = DecodeInt(key[len(prefix):])
	if err != nil {
		return
	}
	_, dstVertexID, err = DecodeInt(key[len(prefix)+8:])
	if err != nil {
		return
	}
	_, srcVertexID, err = DecodeInt(key[len(prefix)+8+8+1:])
	return
}

// OutgoingEdgeKey encodes the outgoing edge key.
//
// The key format is: ${Prefix}${GraphID}${SrcVertexID}${outgoingEdgeSep}${DstVertexID}.
func OutgoingEdgeKey(graphID, srcVertexID, dstVertexID int64) []byte {
	result := make([]byte, 0, EdgeKeyLen)
	result = append(result, prefix...)
	result = EncodeInt(result, graphID)
	result = EncodeInt(result, srcVertexID)
	result = append(result, outgoingEdgeSep)
	result = EncodeInt(result, dstVertexID)
	return result
}

// ParseOutgoingEdgeKey parse the outgoing edge key.
func ParseOutgoingEdgeKey(key []byte) (graphID, srcVertexID, dstVertexID int64, err error) {
	if len(key) < EdgeKeyLen {
		return 0, 0, 0, errors.New("insufficient key length")
	}
	_, graphID, err = DecodeInt(key[len(prefix):])
	if err != nil {
		return
	}
	_, srcVertexID, err = DecodeInt(key[len(prefix)+8:])
	if err != nil {
		return
	}
	_, dstVertexID, err = DecodeInt(key[len(prefix)+8+8+1:])
	return
}
