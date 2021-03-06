package egress_test

import (
	"errors"
	v2 "plumbing/v2"
	"rlp/internal/egress"

	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryServer", func() {
	var (
		spy    *spyContainerMetricFetcher
		server *egress.QueryServer
	)

	BeforeEach(func() {
		spy = newSpyContainerMetricFetcher()
		server = egress.NewQueryServer(spy)
	})

	It("requests the correct appID", func() {
		ctx := context.TODO()
		server.ContainerMetrics(ctx, &v2.ContainerMetricRequest{
			SourceId: "some-app",
		})
		Expect(spy.appID).To(Equal("some-app"))
		Expect(spy.ctx).To(Equal(ctx))
	})

	It("converts v1 container envelopes to v2 envelopes", func() {
		e := &v2.Envelope{
			SourceId: "some-app",
		}
		spy.results = []*v2.Envelope{e}

		results, err := server.ContainerMetrics(context.TODO(), &v2.ContainerMetricRequest{
			SourceId: "some-app",
		})

		Expect(err).ToNot(HaveOccurred())
		Expect(results.Envelopes).To(HaveLen(1))
		Expect(results.Envelopes[0]).To(Equal(e))
	})

	It("returns an error if the source_id is empty", func() {
		_, err := server.ContainerMetrics(context.TODO(), &v2.ContainerMetricRequest{})
		Expect(err).To(HaveOccurred())
	})

	It("returns an error if fetcher fails", func() {
		spy.err = errors.New("some-error")

		_, err := server.ContainerMetrics(context.TODO(), &v2.ContainerMetricRequest{
			SourceId: "some-app",
		})
		Expect(err).To(HaveOccurred())
	})
})

type spyContainerMetricFetcher struct {
	results []*v2.Envelope
	appID   string
	err     error
	ctx     context.Context
}

func newSpyContainerMetricFetcher() *spyContainerMetricFetcher {
	return &spyContainerMetricFetcher{}
}

func (s *spyContainerMetricFetcher) ContainerMetrics(ctx context.Context, appID string) ([]*v2.Envelope, error) {
	s.ctx = ctx
	s.appID = appID
	return s.results, s.err
}
