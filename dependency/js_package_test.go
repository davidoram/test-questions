package dependency

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSPackageUnmarshal(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("data", "app-basic", "package.json"))
	require.NoError(t, err)

	var pkg JSPackage
	require.NoError(t, json.Unmarshal(data, &pkg))

	require.Equal(t, "app-basic", pkg.Name())
	require.Equal(t, []string{"alpha", "beta"}, pkg.DependencyNames())
}

func TestJSPackageDependencyNamesIncludesDevDependencies(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("data", "app-dev", "package.json"))
	require.NoError(t, err)

	var pkg JSPackage
	require.NoError(t, json.Unmarshal(data, &pkg))

	require.Equal(t, []string{"alpha", "tools"}, pkg.DependencyNames())
}

func TestJSPackageUnmarshalAllFixtures(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("data", "*", "package.json"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			var pkg JSPackage
			require.NoError(t, json.Unmarshal(data, &pkg))
			require.NotEmpty(t, pkg.Name())
		})
	}
}
