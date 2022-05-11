package acceptance_test

import (
	"regexp"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry/cf-acceptance-tests/helpers/random_name"
	"github.com/cloudfoundry/cf-test-helpers/cf"
)

var _ = Describe("Scheduler Calls", func() {
	var (
		callName string
	)

	BeforeEach(func() {
		callName = random_name.CATSRandomName("CALL")
	})

	AfterEach(func() {
		Expect(cf.Cf("delete-call", callName).Wait()).To(Exit(0))
	})

	Describe("create-call", func() {
		It("test correct call creation", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10)).To(Exit(0))
		})
	})

	Describe("schedule-call", func() {
		It("test correct call scheduling", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("schedule-call", callName, `15 * * * *`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("call-schedules").
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring(callName))
		})
	})

	Describe("run-call", func() {
		It("test correct call manual execution", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("run-call", callName).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))
		})
	})

	Describe("delete-call", func() {
		It("test correct call deletion", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("delete-call", callName).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("calls").
				Wait(time.Second * 10).Out.Contents()).
				ShouldNot(ContainSubstring(callName))

			Expect(cf.Cf("call-history", callName).
				Wait(time.Second * 10).Out.Contents()).
				ShouldNot(ContainSubstring("OK"))
		})
	})

	Describe("delete-call-schedule", func() {
		It("test correct call schedule deletion", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("schedule-call", callName, `15 * * * *`).
				Wait(time.Second * 10)).To(Exit(0))

			schedules := cf.Cf("call-schedules").
				Wait(time.Second * 10)

			Expect(schedules).To(Exit(0))

			var schedule string
			re := regexp.MustCompile(`^(.*?)[\s]+(.*?)[\s]+([a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12})`)
			for _, line := range strings.Split(string(schedules.Out.Contents()), "\n") {
				for _, i := range re.FindAllStringSubmatch(line, -1) {
					if i[1] == callName {
						schedule = i[3]
					}
				}
			}

			Expect(schedule).NotTo(BeEmpty())

			Expect(cf.Cf("delete-call-schedule", callName, schedule).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("call-schedules").
				Wait(time.Second * 10).Out.Contents()).
				ShouldNot(ContainSubstring(callName))
		})
	})
})
