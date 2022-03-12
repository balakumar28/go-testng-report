package main

import (
	"encoding/xml"
	"fmt"
	"github.com/balakumar28/go-testng-report/json"
	"github.com/balakumar28/go-testng-report/properties"
	"github.com/balakumar28/go-testng-report/testng"
	"github.com/balakumar28/go-testng-report/utils"
	"io/ioutil"
	"time"
)

func GenerateReport(jsonFile string, groupsFile string, reportFile string) error {
	if !utils.FileExists(jsonFile) {
		return fmt.Errorf("gotest report file %s does not exist", jsonFile)
	}
	goReport, err := json.ParseReportJson(jsonFile)
	if err != nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("failed to parse json file %s, err: %w", jsonFile, err)
	}
	groups := properties.NewProperties()
	if utils.FileExists(groupsFile) {
		err := groups.LoadFile(groupsFile)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("testng groups file", groupsFile, "does not exist, using defaults")
	}

	report := ToTestNGReport(goReport, groups)
	bytes, err := xml.MarshalIndent(report, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(reportFile, bytes, 0755)
}

func ToTestNGReport(goReport []json.Report, groups *properties.Properties) testng.Report {
	var report testng.Report

	for i, r := range goReport {
		if i == 0 {
			report.Suite.StartedAt = r.Time.Format(time.RFC3339Nano)
			report.Suite.Name = "golang suite"
			report.Suite.Test.StartedAt = report.Suite.StartedAt
			report.Suite.Test.Name = "golang test"
		}
		if i == len(goReport) - 1 {
			started, err := time.Parse(time.RFC3339Nano, report.Suite.StartedAt)
			if err != nil {
				fmt.Println("failed to parse started time", err)
				continue
			}
			finished := r.Time.Add(time.Duration(int(r.Elapsed * 1000)))
			report.Suite.FinishedAt = finished.Format(time.RFC3339Nano)
			report.Suite.DurationMills = int(finished.Sub(started) / time.Millisecond)
			report.Suite.Test.FinishedAt = report.Suite.FinishedAt
			report.Suite.Test.DurationMills = report.Suite.DurationMills
		}
		if len(r.Test) > 0 && (r.Action == json.Pass || r.Action == json.Fail || r.Action == json.Skip) {
			var testMethod testng.TestMethod
			testMethod.Status = r.Action.String()
			testMethod.Class = r.Package.String()
			testMethod.Name = r.Test
			testMethod.Signature = r.Package.String() + "." + r.Test
			testMethod.DurationMills = int(r.Elapsed * 1000)
			testMethod.StartedAt = r.Time.Format(time.RFC3339Nano)
			testMethod.FinishedAt = r.Time.Add(time.Duration(testMethod.DurationMills)).Format(time.RFC3339Nano)
			report.CountTest(r.Action)
			report.Suite.AddGroups(groups, testMethod)
			report.Suite.AddTest(testMethod)
		}
	}
	return report
}