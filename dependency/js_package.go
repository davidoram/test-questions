package dependency

import (
	"encoding/json"
	"slices"
)

type JSPackage struct {
	name            string
	version         string
	dependencies    map[string]string
	devDependencies map[string]string
}

type jsPackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (p *JSPackage) UnmarshalJSON(data []byte) error {
	var pkg jsPackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}

	p.name = pkg.Name
	p.version = pkg.Version
	p.dependencies = pkg.Dependencies
	p.devDependencies = pkg.DevDependencies

	return nil
}

func (p JSPackage) Name() string {
	return p.name
}

func (p JSPackage) DependencyNames() []string {
	seen := make(map[string]struct{}, len(p.dependencies)+len(p.devDependencies))
	for name := range p.dependencies {
		seen[name] = struct{}{}
	}
	for name := range p.devDependencies {
		seen[name] = struct{}{}
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
