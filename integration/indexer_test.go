package integration_test

import (
	"context"
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/kstiehl/index-bouncer/integration/helper"
	"github.com/kstiehl/index-bouncer/pkg/opensearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	oSearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/testcontainers/testcontainers-go"
)

var _ = Describe("opensearch", Ordered, func() {
	var openSearchContainer testcontainers.Container
	var osClient *oSearch.Client

	BeforeAll(func() {
		ctx := context.Background()
		openSearchContainer, osClient = helper.GetOpenSearch(ctx)
	})

	AfterAll(func() {
		openSearchContainer.Terminate(context.Background())
	})

	It("create datastream", func() {
		ctx := logr.NewContext(context.Background(), stdr.New(log.New(os.Stdout, "", log.Lshortfile)))

		existsRequest := opensearchapi.IndicesExistsIndexTemplateRequest{
			Name: "test",
		}

		res, err := existsRequest.Do(ctx, osClient)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.IsError()).To(BeTrue())

		osClient := opensearch.Client{osClient}
		err = opensearch.EnsureIndexTemplate(ctx, osClient, helper.TestStream{
			StreamName: "test",
		})
		Expect(err).ToNot(HaveOccurred())

		existsRequest = opensearchapi.IndicesExistsIndexTemplateRequest{
			Name: "test",
		}
		res, err = existsRequest.Do(ctx, osClient)
		defer res.Body.Close()
		Expect(err).ToNot(HaveOccurred())
		Expect(res.IsError()).To(BeFalse())
	})
})
