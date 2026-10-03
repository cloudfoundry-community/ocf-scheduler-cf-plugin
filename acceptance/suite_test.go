package acceptance_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry/cf-test-helpers/v2/cf"
	"github.com/cloudfoundry/cf-test-helpers/v2/generator"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Scheduler")
}

var appName string

// memory sizes the test app and the tasks its jobs run; a hello-world
// app needs far less than the 1024M task default.
const memory = "64M"

var _ = BeforeSuite(func() {
	appName = generator.PrefixedRandomName("CATS", "APP")

	// This command is expensive, lets do it only once.
	Expect(cf.Cf("push", appName,
		"-m", memory,
		"-p", "assets/golang",
		"-f", "assets/golang/manifest.yml",
	).Wait(time.Second * 600)).To(Exit(0))
})

var _ = AfterSuite(func() {
	Expect(cf.Cf("delete", appName, "-f", "-r").Wait(time.Second * 20)).To(Exit(0))
})
