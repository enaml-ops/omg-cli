package main

import (
	"bytes"

	"github.com/enaml-ops/pluginlib/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("omg cli", func() {
	Context("when listing products multiple times", func() {
		It("should generate the same output each time", func() {
			products := map[string]registry.Record{
				"foo": registry.Record{
					Name: "foo",
					Path: "foo-plugin-linux",
					Properties: map[string]interface{}{
						"prop1": "value1",
						"prop2": "value2",
					},
				},
				"bar": registry.Record{
					Name: "bar",
					Path: "bar-plugin-linux",
					Properties: map[string]interface{}{
						"prop1": "bar1",
						"prop2": "bar2",
					},
				},
			}

			buf := &bytes.Buffer{}
			ListProducts(buf, products)
			orig := buf.String()

			// re-list products several times and make sure the output
			// is identical to the first run
			for i := 0; i < 10; i++ {
				buf := &bytes.Buffer{}
				ListProducts(buf, products)
				Ω(orig).Should(Equal(buf.String()))
			}
		})
	})
})
