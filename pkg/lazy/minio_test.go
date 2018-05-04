package lazy

import (
	. "github.com/franela/goblin"
	"github.com/jensneuse/kube-utils/pkg/config"
	"github.com/minio/minio-go"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMinioTemplate(t *testing.T) {

	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	podName := "minio-simple-test"
	endpoint := "localhost:9000"
	accessKey := "AKIAIOSFODNN7EXAMPLE"
	secretKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	l := New()

	g.Describe("minio", func() {

		g.After(func() {
			// tear down
			l.Cleanup()
		})

		g.It("should be ready to connect", func() {

			// request whatever we need
			l.CreateAndForwardMinio(config.NAMESPACE, podName, accessKey, secretKey)

			minioClient, err := minio.New(endpoint, accessKey, secretKey, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err = minioClient.ListBuckets()
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
