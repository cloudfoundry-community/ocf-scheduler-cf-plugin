package acceptance_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry/cf-acceptance-tests/helpers/random_name"
	"github.com/cloudfoundry/cf-test-helpers/cf"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Scheduler")
}

var appName string

var _ = BeforeSuite(func() {
	appName = random_name.CATSRandomName("APP")

	// This command is expensive, lets do it only once.
	Expect(cf.Cf("push", appName,
		"-m", "256M",
		"-p", "assets/golang",
		"-f", "assets/golang/manifest.yml",
	).Wait(time.Second * 120)).To(Exit(0))
})

var _ = AfterSuite(func() {
	Expect(cf.Cf("delete", appName, "-f", "-r").Wait()).To(Exit(0))
})
