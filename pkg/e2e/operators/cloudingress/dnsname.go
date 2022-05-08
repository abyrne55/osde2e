package cloudingress

import (
	"context"
	"log"
	"time"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift/osde2e/pkg/common/constants"
	"github.com/openshift/osde2e/pkg/common/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// tests
var _ = ginkgo.Describe(constants.SuiteInforming+TestPrefix, func() {
	h := helper.New()
	var dnsnameOriginal string
	// How long to wait for IngressController changes
	pollingDuration := 120 * time.Second
	ginkgo.Context("publishingstrategy-dnsname", func() {
		ginkgo.It("IngressController should be patched when update dnsname", func() {
			ingress1, _ := getingressController(h, "default")
			dnsnameOriginal = string(ingress1.Spec.Domain)
			log.Print(" the Domain name \n", dnsnameOriginal)
			updateDnsName(h, "foo")

			ingress, _ := getingressController(h, "default")
			log.Print(" The new Generation is \n", ingress.Generation)
		}, pollingDuration.Seconds())

		ginkgo.It("IngressController should be patched when return to the original dnsname", func() {
			updateDnsName(h, dnsnameOriginal)

			ingress, _ := getingressController(h, "default")
			Expect(string(ingress.Spec.Domain)).To(Equal(dnsnameOriginal))
		}, pollingDuration.Seconds())
	})
})

func updateDnsName(h *helper.H, newName string) {
	var err error
	PublishingStrategyInstance, ps := getPublishingStrategy(h)
	AppIngress := PublishingStrategyInstance.Spec.ApplicationIngress

	for i, v := range AppIngress {
		if v.Default == true {
			AppIngress[i].DNSName = newName
		}
	}

	PublishingStrategyInstance.Spec.ApplicationIngress = AppIngress

	ps.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&PublishingStrategyInstance)
	Expect(err).NotTo(HaveOccurred())

	// Update the publishingstrategy
	ps, err = h.Dynamic().Resource(schema.GroupVersionResource{Group: "cloudingress.managed.openshift.io", Version: "v1alpha1", Resource: "publishingstrategies"}).Namespace(OperatorNamespace).Update(context.TODO(), ps, metav1.UpdateOptions{})
	Expect(err).NotTo(HaveOccurred())
}