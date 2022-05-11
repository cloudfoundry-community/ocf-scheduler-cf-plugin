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

var _ = Describe("Scheduler Jobs", func() {
	var (
		jobName string
	)

	BeforeEach(func() {
		jobName = random_name.CATSRandomName("JOB")
	})

	AfterEach(func() {
		Expect(cf.Cf("delete-job", jobName).Wait()).To(Exit(0))
	})

	Describe("create-job", func() {
		It("test correct job creation", func() {
			Expect(cf.Cf("create-job", appName, jobName, `pwd`).
				Wait(time.Second * 10)).To(Exit(0))
		})
	})

	Describe("schedule-job", func() {
		It("test correct job scheduling", func() {
			Expect(cf.Cf("create-job", appName, jobName, `pwd`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("schedule-job", jobName, `15 * * * *`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("job-schedules").
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring(jobName))
		})
	})

	Describe("run-job", func() {
		It("test correct job manual execution", func() {
			Expect(cf.Cf("create-job", appName, jobName, `pwd`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("run-job", jobName).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))
		})
	})

	Describe("delete-job", func() {
		It("test correct job deletion", func() {
			Expect(cf.Cf("create-job", appName, jobName, `pwd`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("delete-job", jobName).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("jobs").
				Wait(time.Second * 10).Out.Contents()).
				ShouldNot(ContainSubstring(jobName))

			Expect(cf.Cf("job-history", jobName).
				Wait(time.Second * 10).Out.Contents()).
				ShouldNot(ContainSubstring("OK"))
		})
	})

	Describe("delete-job-schedule", func() {
		It("test correct job schedule deletion", func() {
			Expect(cf.Cf("create-job", appName, jobName, `pwd`).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("schedule-job", jobName, `15 * * * *`).
				Wait(time.Second * 10)).To(Exit(0))

			schedules := cf.Cf("job-schedules").
				Wait(time.Second * 10)

			Expect(schedules).To(Exit(0))

			var schedule string
			re := regexp.MustCompile(`^(.*?)[\s]+(.*?)[\s]+([a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12})`)
			for _, line := range strings.Split(string(schedules.Out.Contents()), "\n") {
				for _, i := range re.FindAllStringSubmatch(line, -1) {
					if i[1] == jobName {
						schedule = i[3]
					}
				}
			}

			Expect(schedule).NotTo(BeEmpty())

			Expect(cf.Cf("delete-job-schedule", jobName, schedule).
				Wait(time.Second * 10)).To(Exit(0))

			Expect(cf.Cf("job-schedules").
				Wait(time.Second * 10).Out.Contents()).
				ShouldNot(ContainSubstring(jobName))
		})
	})
})
