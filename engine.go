package kraken

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"strings"

	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

// Engine holding all Graphs.
type Engine struct {
	ID      uuid.UUID
	Started time.Time
	Graphs  map[*Graph]bool
}

// Inspect inspects this Engine.
func (e *Engine) Inspect() {
	fmt.Printf("ID:\t\t%s\n", e.ID)
	fmt.Printf("Type:\t\tEngine\n")
	fmt.Printf("Started:\t%s\n", e.Started.Format(TimeFormat))
	fmt.Printf("Graphs:\t\t%d\n", e.CountGraphs())
	fmt.Printf("\n")
}

// GetGraph tries to find a Graph based on an ID.
func (e *Engine) GetGraph(id string) (g *Graph, err error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	for elem := range e.Graphs {
		if elem.ID == uid {
			return elem, nil
		}
	}
	return nil, errors.New("graph not found")
}

// FindGraph tries to find a Graph by its name.
func (e *Engine) FindGraph(name string) (g *Graph, err error) {
	for elem := range e.Graphs {
		if elem.Name == name {
			return elem, nil
		}
	}
	return nil, errors.New("graph not found")
}

// ToYaml transforms the content of this Engine to yaml.
func (e *Engine) ToYaml() (y string, er error) {
	yam, err := yaml.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(yam), nil
}

// AddGraph to Engine.
func (e *Engine) AddGraph(g *Graph) {
	e.Graphs[g] = true
}

// DropGraph from Engine.
func (e *Engine) DropGraph(g *Graph) {
	delete(e.Graphs, g)
}

// CountGraphs counts the total number of all Graphs in this Engine.
func (e *Engine) CountGraphs() int {
	return len(e.Graphs)
}

// LoadDirectory loads all .kraken files in the given directory.
func (e *Engine) LoadDirectory(path string) error {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), FileSuffix) {
			name := strings.TrimSuffix(f.Name(), FileSuffix)
			g, err := e.ReadFromDisk(name)
			if err != nil {
				return nil
			}
			e.AddGraph(g)
		}
	}
	return nil
}

// WriteToDisk writes the content of this graph to disk.
func (e *Engine) WriteToDisk(g *Graph) error {
	g.Saved = time.Now()
	fileName := g.Name + FileSuffix

	y, err := g.ToYaml()
	if err != nil {
		return err
	}
	data := []byte(y)

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadFromDisk loads the graph from the disk.
// Needs the name of the graph to load.
func (e *Engine) ReadFromDisk(name string) (g *Graph, er error) {
	fileName := name + FileSuffix

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	gra, err := FromYaml(string(data))
	if err != nil {
		return nil, err
	}

	return gra, nil
}

// NewEngine creates a brand new Engine.
func NewEngine() *Engine {
	return &Engine{
		ID:      uuid.NewV4(),
		Graphs:  make(map[*Graph]bool),
		Started: time.Now(),
	}
}