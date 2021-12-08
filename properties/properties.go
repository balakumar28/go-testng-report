package properties

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Properties struct {
	properties map[string]string
}

func (p *Properties) Get(key string) (string, bool) {
	value, found := p.properties[key]
	return value, found
}

func (p *Properties) Put(key, value string) {
	p.properties[key] = value
}

func (p *Properties) Size() int {
	return len(p.properties)
}

func NewProperties() *Properties {
	properties := make(map[string]string)
	return &Properties{properties: properties}
}

func (p *Properties) LoadFile(filePath string) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading property file: %s", err)
	}
	text := trimText(bytes)
	for _, line := range strings.Split(text, "\n") {
		if err = p.updateProperty(line); err != nil {
			return err
		}
	}
	return nil
}

func (p *Properties) updateProperty(line string) error {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == '#' {
		return nil
	}
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return errors.New("invalid property file line " + line)
	}
	key := strings.TrimSpace(parts[0])
	value, err := p.resolveValue(strings.TrimSpace(parts[1]))
	if err != nil {
		return err
	}
	p.properties[key] = value
	return nil
}

func (p *Properties) resolveValue(valueString string) (string, error) {
	const errString = "PROP_FILE_ERROR_KEY_NOT_FOUND"
	resolvedValue := os.Expand(valueString, func(name string) string {
		value, ok := p.properties[name]
		if ok {
			return value
		} else {
			return errString
		}
	})
	if errString == resolvedValue {
		return "", errors.New("invalid key reference " + valueString)
	}

	return resolvedValue, nil
}

func trimText(bytes []byte) string {
	text := string(bytes)
	text = strings.Replace(text, "\r\n", "\n", -1)
	text = strings.Replace(text, "\r", "\n", -1)
	return text
}
