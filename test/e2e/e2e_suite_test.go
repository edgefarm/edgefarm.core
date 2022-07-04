package e2e

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/edgefarm/edgefarm.core/test/e2e/framework"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestKube(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	gomega.RegisterFailHandler(ginkgo.Fail)
	err := framework.CreateFramework(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to create framework: %v", err)
	}
	ginkgo.RunSpecs(t, "Kube Suite")
}
