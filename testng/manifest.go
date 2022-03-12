package testng

import (
	"encoding/xml"
	"github.com/balakumar28/go-testng-report/json"
	"github.com/balakumar28/go-testng-report/properties"
	"strings"
)

const defaultGroupName = "UT"
const DateFormat = "2006-01-02T15:04:05Z"

type Report struct {
	XMLName        xml.Name `xml:"testng-results"`
	Skipped        int      `xml:"skipped,attr"`
	Failed         int      `xml:"failed,attr"`
	Ignored        int      `xml:"ignored,attr"`
	Total          int      `xml:"total,attr"`
	Passed         int      `xml:"passed,attr"`
	ReporterOutput string   `xml:"reporter-output"`
	Suite          Suite    `xml:"suite"`
}

func (r *Report) CountTest(result json.TestAction) {
	r.Total += 1
	switch result {
	case json.Pass:
		r.Passed += 1
		break
	case json.Fail:
		r.Failed += 1
		break
	case json.Skip:
		r.Skipped += 1
		break
	}
}

type Suite struct {
	Name          string  `xml:"name,attr"`
	DurationMills int     `xml:"duration-ms,attr"`
	StartedAt     string  `xml:"started-at,attr"`
	FinishedAt    string  `xml:"finished-at,attr"`
	Groups        []Group `xml:"groups>group"`
	Test          Test    `xml:"test"`
}

func (s *Suite) AddGroups(groups *properties.Properties, method TestMethod) {
	groupMethod := method.toGroupMethod()
	added := s.addGroups(groups, method.Signature, groupMethod)
	if !added {
		added = s.addGroups(groups, method.Class, groupMethod)
		if !added {
			s.addToDefaultGroup(groupMethod)
		}
	}
}

func (s *Suite) addGroups(groups *properties.Properties, groupKey string, groupMethod Method) bool {
	if g, ok := groups.Get(groupKey); ok {
		grps := strings.Split(g, ",")
		for _, grp := range grps {
			if i := s.indexOfGroup(grp); i != -1 {
				s.Groups[i].Methods = append(s.Groups[i].Methods, groupMethod)
			} else {
				var group Group
				group.Name = grp
				group.Methods = append(group.Methods, groupMethod)
				s.Groups = append(s.Groups, group)
			}
		}
		return true
	}
	return false
}

func (s *Suite) addToDefaultGroup(groupMethod Method) {
	if i := s.indexOfGroup(defaultGroupName); i != -1 {
		s.Groups[i].Methods = append(s.Groups[i].Methods, groupMethod)
	} else {
		var group Group
		group.Name = defaultGroupName
		group.Methods = append(group.Methods, groupMethod)
		s.Groups = append(s.Groups, group)
	}
}

func (s Suite) indexOfGroup(groupName string) int {
	for i, g := range s.Groups {
		if g.Name == groupName {
			return i
		}
	}
	return -1
}

func (s *Suite) AddTest(testMethod TestMethod) {
	index := s.Test.getIndex(testMethod.Class)
	s.Test.Classes[index].TestMethods = append(s.Test.Classes[index].TestMethods, testMethod)
}

type Group struct {
	XMLName xml.Name `xml:"group"`
	Name    string   `xml:"name,attr"`
	Methods []Method `xml:"method"`
}

type Method struct {
	Signature string `xml:"signature,attr"`
	Name      string `xml:"name,attr"`
	Class     string `xml:"class,attr"`
}

type Test struct {
	Name          string  `xml:"name,attr"`
	DurationMills int     `xml:"duration-ms,attr"`
	StartedAt     string  `xml:"started-at,attr"`
	FinishedAt    string  `xml:"finished-at,attr"`
	Classes       []Class `xml:"class"`
}

func (t *Test) getIndex(pName string) int {
	if i := t.indexOfClass(pName); i != -1 {
		return i
	} else {
		var class Class
		class.Name = pName
		t.Classes = append(t.Classes, class)
		return len(t.Classes) - 1
	}
}

func (t *Test) indexOfClass(className string) int {
	for i, c := range t.Classes {
		if c.Name == className {
			return i
		}
	}
	return -1
}

type Class struct {
	Name        string       `xml:"name,attr"`
	TestMethods []TestMethod `xml:"test-method"`
}

type TestMethod struct {
	Status         string `xml:"status,attr"`
	Signature      string `xml:"signature,attr"`
	Name           string `xml:"name,attr"`
	DurationMills  int    `xml:"duration-ms,attr"`
	StartedAt      string `xml:"started-at,attr"`
	FinishedAt     string `xml:"finished-at,attr"`
	ReporterOutput string `xml:"reporter-output"`
	Class          string `xml:"-"`
}

func (m TestMethod) toGroupMethod() Method {
	var method Method
	method.Name = m.Name
	method.Signature = m.Signature
	method.Class = m.Class
	return method
}
