package properties

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadProperties_MissingFile(t *testing.T) {
	propFile := filepath.Join("testdata", "missingProp.properties")
	err := NewProperties().LoadFile(propFile)

	assert.Error(t, err)
}

func TestLoadProperties_CorruptedPropFile(t *testing.T) {
	propFile := filepath.Join("testdata", "corruptProp.properties")
	err := NewProperties().LoadFile(propFile)

	assert.Error(t, err)
}

func TestLoadProperties_ValidFile(t *testing.T) {
	propFile := filepath.Join("testdata", "validProp.properties")
	properties := NewProperties()
	err := properties.LoadFile(propFile)

	assert.Nil(t, err)
	assert.NotNil(t, properties)
	assert.Equal(t, properties.Size(), 3)
	assertPropertiesGet(t, properties, "key1", "value1")
	assertPropertiesGet(t, properties, "key2", "value2")
	assertPropertiesGet(t, properties, "key3", "value3")
}

func assertPropertiesGet(t *testing.T, properties *Properties, key string, expectedValue string) {
	value, _ := properties.Get(key)
	assert.Equal(t, value, expectedValue)
}

func TestLoadProperties_ValidFile_WithReference(t *testing.T) {
	propFile := filepath.Join("testdata", "validPropWithInvalidReference.properties")
	err := NewProperties().LoadFile(propFile)

	assert.Error(t, err)
}

func Test_updateProperty(t *testing.T) {
	p := NewProperties()
	err := p.updateProperty(`MSI_ARGS=TARGETDIR="C:\Program Files (x86)" REINSTALLMODE=vomus REINSTALL=ALL`)
	assert.Nil(t, err)
	assert.Equal(t, `TARGETDIR="C:\Program Files (x86)" REINSTALLMODE=vomus REINSTALL=ALL`, p.properties["MSI_ARGS"])
}
