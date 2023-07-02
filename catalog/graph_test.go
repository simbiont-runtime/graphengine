// ---

package catalog

import (
	"testing"

	"github.com/simbiont-runtime/graphengine/parser/model"
	"github.com/stretchr/testify/assert"
)

func TestNewGraph(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	assert.Equal(graph.Meta(), meta)
}

func TestGraph_Label(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
		Labels: []*model.LabelInfo{
			{
				ID:   2,
				Name: model.NewCIStr("label1"),
			},
		},
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	assert.Equal(graph.Label("label1").Meta(), meta.Labels[0])
}

func TestGraph_LabelByID(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
		Labels: []*model.LabelInfo{
			{
				ID:   2,
				Name: model.NewCIStr("label1"),
			},
		},
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	assert.Equal(graph.LabelByID(2).Meta(), meta.Labels[0])
}

func TestGraph_CreateLabel(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	label := &model.LabelInfo{
		ID:   2,
		Name: model.NewCIStr("label1"),
	}
	graph.CreateLabel(label)
	assert.Equal(graph.LabelByID(2).Meta(), label)
	assert.Equal(graph.Label("label1").Meta(), label)
}

func TestGraph_DropLabel(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
		Labels: []*model.LabelInfo{
			{
				ID:   2,
				Name: model.NewCIStr("label1"),
			},
		},
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	graph.DropLabel(meta.Labels[0])
	assert.Nil(graph.LabelByID(2))
}

func TestGraph_Property(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
		Properties: []*model.PropertyInfo{
			{
				ID:   2,
				Name: model.NewCIStr("property1"),
			},
		},
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	assert.Equal(graph.Property("property1"), meta.Properties[0])
}

func TestGraph_PropertyByID(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
		Properties: []*model.PropertyInfo{
			{
				ID:   2,
				Name: model.NewCIStr("property1"),
			},
		},
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	assert.Equal(graph.PropertyByID(2), meta.Properties[0])
}

func TestGraph_CreateProperty(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	property := &model.PropertyInfo{
		ID:   2,
		Name: model.NewCIStr("property1"),
	}
	graph.CreateProperty(property)
	assert.Equal(graph.PropertyByID(2), property)
	assert.Equal(graph.Property("property1"), property)
}

func TestGraph_DropProperty(t *testing.T) {
	meta := &model.GraphInfo{
		ID:   1,
		Name: model.NewCIStr("test-graph"),
		Properties: []*model.PropertyInfo{
			{
				ID:   2,
				Name: model.NewCIStr("property1"),
			},
		},
	}

	graph := NewGraph(meta)
	assert := assert.New(t)
	graph.DropProperty(meta.Properties[0])
	assert.Nil(graph.PropertyByID(2))
	assert.Nil(graph.Property("property1"))
}
