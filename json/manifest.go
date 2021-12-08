package json

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"
)

type TestAction string

func (ta TestAction) String() string {
	return strings.ToUpper(string(ta))
}

type GoPackage string

func (p *GoPackage) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	str = strings.Replace(str, "/", ".", 1)
	*p = GoPackage(strings.ReplaceAll(str, "/", "-"))
	return nil
}

func (p *GoPackage) String() string {
	return string(*p)
}

const (
	Run    TestAction = "run"
	Output TestAction = "output"
	Pass   TestAction = "pass"
	Skip   TestAction = "skip"
	Fail   TestAction = "fail"
)

type Report struct {
	Time    time.Time
	Action  TestAction
	Package GoPackage
	Test    string
	Output  string
	Elapsed float64
}

func ParseReportJson(filePath string) ([]Report, error) {
	var reports []Report
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return nil, err
		}
		var report Report
		err := json.Unmarshal(scanner.Bytes(), &report)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}
