// ---

package catalog

import (
	"github.com/simbiont-runtime/graphengine/internal/logutil"
	"github.com/simbiont-runtime/graphengine/parser/model"
)

// PatchType represents the type of patch.
type PatchType byte

const (
	PatchTypeCreateGraph PatchType = iota
	PatchTypeCreateLabel
	PatchTypeCreateIndex
	PatchTypeDropGraph
	PatchTypeDropLabel
	PatchTypeDropIndex
	PatchTypeCreateProperties
)

type (
	// Patch represents patch which contains a DDL change.
	Patch struct {
		Type PatchType
		Data interface{}
	}

	// PatchLabel represents the payload of patching create/drop label DDL.
	PatchLabel struct {
		GraphID   int64
		LabelInfo *model.LabelInfo
	}

	// PatchProperties represents the payload of patching create properties
	PatchProperties struct {
		MaxPropID  uint16
		GraphID    int64
		Properties []*model.PropertyInfo
	}
)

// Apply applies the patch to catalog.
// Note: we need to ensure the DDL changes have applied to persistent storage first.
func (c *Catalog) Apply(patch *Patch) {
	switch patch.Type {
	case PatchTypeCreateGraph:
		data := patch.Data.(*model.GraphInfo)
		graph := NewGraph(data)
		c.mu.Lock()
		c.byName[data.Name.L] = graph
		c.byID[data.ID] = graph
		c.mu.Unlock()

	case PatchTypeDropGraph:
		data := patch.Data.(*model.GraphInfo)
		c.mu.Lock()
		delete(c.byName, data.Name.L)
		delete(c.byID, data.ID)
		c.mu.Unlock()

	case PatchTypeCreateLabel:
		data := patch.Data.(*PatchLabel)
		graph := c.GraphByID(data.GraphID)
		if graph == nil {
			logutil.Errorf("Create label on not exists graph. GraphID: %d", data.GraphID)
			return
		}
		graph.CreateLabel(data.LabelInfo)

	case PatchTypeDropLabel:
		data := patch.Data.(*PatchLabel)
		graph := c.GraphByID(data.GraphID)
		if graph == nil {
			logutil.Errorf("Drop label on not exists graph. GraphID: %d", data.GraphID)
			return
		}
		graph.DropLabel(data.LabelInfo)

	case PatchTypeCreateProperties:
		data := patch.Data.(*PatchProperties)
		graph := c.GraphByID(data.GraphID)
		graph.SetNextPropID(data.MaxPropID)
		for _, p := range data.Properties {
			graph.CreateProperty(p)
		}
	}
}
