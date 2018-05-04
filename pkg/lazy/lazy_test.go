package lazy

import (
	. "github.com/franela/goblin"
	"github.com/minio/minio-go"
	. "github.com/onsi/gomega"
	"testing"
)

func TestLazy(t *testing.T) {

	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	namespace := "jens-neuse"
	podName := "minio-test-123"
	endpoint := "localhost:9000"
	accessKey := "AKIAIOSFODNN7EXAMPLE"
	secretKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	lazy := New()

	g.Describe("lazy", func() {

		g.It("should have provided a minio pod, ready to connect", func() {

			// request whatever we need
			lazy.CreateAndForwardMinioBlocking(namespace, podName, accessKey, secretKey)

			minioClient, err := minio.New(endpoint, accessKey, secretKey, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err = minioClient.ListBuckets()
			if err != nil {
				t.Fatal(err)
			}

			// tear down
			lazy.Cleanup()
		})
	})
}
