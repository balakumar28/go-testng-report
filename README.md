# go-testng-report
A tool to convert go test results to testng reports xml.
<pre>
<b>Usage of testng-report:</b>
	-json-report string
		Golang json test report (default "report.json")
	-out string
		TestNG Report XML (default "testng-results.xml")
	-testng-groups string
		TestNG group mapping (default "groups.properties")
</pre>

## Basic Use

1. Generate json test report for your go module
<pre>
Example:
      # cd <i>your_module</i>
      # go test -json ./... > report.json
</pre>
2. Run go-testng-report
<pre>
Example:
      # cd <i>your_module</i>
      # go-testng-report
</pre>

## TestNG Groups
By default, all tests are tagged with group "UT". This can be customized by placing groups.properties in the project or by passing -testng-groups argument to go-testng-report.<br/>
If the test has a match in the groups.properties, then test is tagged with the groups mentioned, otherwise it is tagged with "UT". The match may be by package.testName or by package, package.testName takes higher precedence.<br/>
For reporting convenience on Jenkins, module name is used as package name and all subpackage names are hyphenated, so that all packages are grouped by module name.<br/>

**groups.properties convention:**<br/>
*module.package=G1,G2,...*<br/>
*module.package.testName=Group1,Group2,...*<br/>

*Note:* Generate a testng-reports.xml with default groups to better understand the naming convention followed and update groups.properties accordingly.<br/>
