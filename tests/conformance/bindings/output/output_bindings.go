package bindings

import (
	"fmt"
	"testing"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/tests/conformance/utils"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/sets"
)

type TestConfig struct {
	utils.CommonConfig
}

func NewTestConfig(name string, allOperations bool, operations []string) TestConfig {
	return TestConfig{
		CommonConfig: utils.CommonConfig{
			ComponentType: "output-binding",
			ComponentName: name,
			AllOperations: allOperations,
			Operations:    sets.NewString(operations...),
		},
	}
}

func ConformanceTests(t *testing.T, props map[string]string, binding bindings.OutputBinding, config TestConfig) {
	if config.CommonConfig.HasOperation("init") {
		t.Run(config.GetTestName("init"), func(t *testing.T) {
			err := binding.Init(bindings.Metadata{
				Properties: props,
			})
			assert.NoError(t, err, "expected no error setting up binding")
		})
	}

	if config.CommonConfig.HasOperation("operations") {
		t.Run(config.GetTestName("operations"), func(t *testing.T) {
			ops := binding.Operations()
			for _, op := range ops {
				assert.True(t, config.HasOperation(string(op)), fmt.Sprintf("Operation missing from conformance test config: %v", op))
			}
		})
	}

	// Base operations from bindings definition.
	allOps := []bindings.OperationKind{
		bindings.CreateOperation,
		bindings.GetOperation,
		bindings.ListOperation,
		bindings.DeleteOperation,
	}

	for _, op := range allOps {
		if config.HasOperation(string(op)) {
			t.Run(config.GetTestName(string(op)), func(t *testing.T) {
				req := bindings.InvokeRequest{
					Data:      []byte("Test Data"),
					Operation: op,
					Metadata: map[string]string{
						"key": "test-key",
					},
				}
				_, err := binding.Invoke(&req)
				assert.Nil(t, err, "expected no error invoking output binding")
			})
		}
	}
}
