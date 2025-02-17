package main

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
)

type config struct {
	AdditionalParams          string `env:"additional_params"`
	TestsPathPattern          string `env:"tests_path_pattern"`
	ProjectLocation           string `env:"project_location,dir"`
	TestResultsDir            string `env:"bitrise_test_result_dir,dir"`
	TestResultsName           string `env:"bitrise_test_result_name"`
	GenerateCodeCoverageFiles bool   `env:"generate_code_coverage_files,opt[yes,no]"`
}

var ir interrupt = realInterrupt{}
var parser configParser = realConfigParser{interrupt: ir}
var builder commandBuilder = realCommandBuilder{interrupt: ir}
var test testExecutor = realTestExecutor{interrupt: ir, commandBuilder: builder, testExporter: realTestExporter{interrupt: ir}}

func main() {
	cfg := parser.parseConfig()

	stepconf.Print(cfg)

	additionalParams := parser.parseAdditionalParams(cfg.AdditionalParams)

	testPaths := parser.expandTestsPathPattern(cfg.ProjectLocation, cfg.TestsPathPattern)

	additionalParams = append(additionalParams, testPaths...)

	fmt.Println()
	log.Infof("Running test")

	outputBuffer, testErr := test.executeTest(cfg, additionalParams)
	test.exportTestResults(cfg, outputBuffer)

	if testErr {
		ir.fail()
	}
}
