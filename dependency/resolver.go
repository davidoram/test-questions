package dependency

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type FileResolver struct {
	path  string
	cache map[string]JSPackage
}

func (r FileResolver) Resolve(name string) (error, JSPackage) {
	if pkg, ok := r.cache[name]; ok {
		return nil, pkg
	}

	data, err := os.ReadFile(filepath.Join(r.path, name, "package.json"))
	if err != nil {
		return err, JSPackage{}
	}

	var pkg JSPackage
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err, JSPackage{}
	}

	r.cache[name] = pkg
	return nil, pkg
}

func (r FileResolver) Packages() []JSPackage {
	packages := make([]JSPackage, 0, len(r.cache))
	for _, pkg := range r.cache {
		packages = append(packages, pkg)
	}

	slices.SortFunc(packages, func(a, b JSPackage) int {
		if a.Name() < b.Name() {
			return -1
		}
		if a.Name() > b.Name() {
			return 1
		}
		return 0
	})

	return packages
}

func NewFileResolver(path string) (error, FileResolver) {
	resolver := FileResolver{
		path:  path,
		cache: map[string]JSPackage{},
	}

	return nil, resolver
}

func MakeGraphFromJS(path string, names ...string) (error, *Graph) {
	err, resolver := NewFileResolver(path)
	if err != nil {
		return err, nil
	}

	graph := NewGraph()
	seen := map[string]bool{}

	var addPackage func(name string) error
	addPackage = func(name string) error {
		if seen[name] {
			return nil
		}
		seen[name] = true

		err, pkg := resolver.Resolve(name)
		if err != nil {
			return err
		}

		graph.AddNode(pkg.Name())
		for _, dependency := range pkg.DependencyNames() {
			if err := addPackage(dependency); err != nil {
				return err
			}
			if err := graph.AddEdge(pkg.Name(), dependency); err != nil {
				return err
			}
		}
		return nil
	}

	for _, name := range names {
		if err := addPackage(name); err != nil {
			return err, nil
		}
	}

	return nil, graph
}
