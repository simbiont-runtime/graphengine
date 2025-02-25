// ---

package model

// GraphInfo provides meta data describing a graph.
type GraphInfo struct {
	ID         int64           `json:"id"`
	Name       CIStr           `json:"name"`
	Indexes    []*IndexInfo    `json:"indexes"`
	NextPropID uint16          `json:"next_prop_id"`
	Query      string          `json:"query"`
	Labels     []*LabelInfo    `json:"-"`
	Properties []*PropertyInfo `json:"-"`
}

// IndexInfo provides meta data describing a index.
type IndexInfo struct {
	ID         int64   `json:"id"`
	Name       CIStr   `json:"name"`
	Properties []CIStr `json:"properties"`
	Query      string  `json:"query"`
}

// LabelInfo provides meta data describing a label.
type LabelInfo struct {
	ID    int64  `json:"id"`
	Name  CIStr  `json:"name"`
	Query string `json:"query"`
}

type PropertyInfo struct {
	ID   uint16 `json:"id"`
	Name CIStr  `json:"name"`
}

func (info *GraphInfo) Clone() *GraphInfo {
	cloned := *info
	cloned.Indexes = make([]*IndexInfo, len(info.Indexes))
	for i := range info.Indexes {
		cloned.Indexes[i] = info.Indexes[i].Clone()
	}
	cloned.Labels = make([]*LabelInfo, len(info.Labels))
	for i := range info.Labels {
		cloned.Labels[i] = info.Labels[i].Clone()
	}
	cloned.Properties = make([]*PropertyInfo, len(info.Properties))
	for i := range info.Properties {
		cloned.Properties[i] = info.Properties[i].Clone()
	}
	return &cloned
}

func (info *IndexInfo) Clone() *IndexInfo {
	cloned := *info
	cloned.Properties = make([]CIStr, len(info.Properties))
	for i := range info.Properties {
		cloned.Properties[i] = info.Properties[i]
	}
	return &cloned
}

func (info *LabelInfo) Clone() *LabelInfo {
	cloned := *info
	return &cloned
}

func (info *PropertyInfo) Clone() *PropertyInfo {
	cloned := *info
	return &cloned
}
