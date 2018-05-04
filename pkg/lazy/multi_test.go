package lazy

import (
	"database/sql"
	. "github.com/franela/goblin"
	"github.com/jensneuse/kube-utils/pkg/config"
	"github.com/jensneuse/kube-utils/pkg/podtemplates"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCreateAndForwardPods(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("CreateAndForwardPods", func() {

		l := New()

		g.After(func() {
			l.Cleanup()
		})

		minioPodName := "minio-multitest"
		endpoint := "localhost:9000"
		accessKey := "AKIAIOSFODNN7EXAMPLE"
		secretKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

		postgresPodName := "postgres-multitest"
		dsn := "postgresql://postgres@localhost:15432/postgres?sslmode=disable"

		g.It("should provide minio + postgres", func() {
			l.CreateAndForwardPods(config.NAMESPACE,
				podtemplates.Minio(minioPodName, accessKey, secretKey),
				podtemplates.Postgresql(postgresPodName),
			)

			minioClient, err := minio.New(endpoint, accessKey, secretKey, false)
			if err != nil {
				t.Fatal(err)
			}

			_, err = minioClient.ListBuckets()
			if err != nil {
				t.Fatal(err)
			}

			db, err := sql.Open("postgres", dsn)
			if err != nil {
				t.Fatal(err)
			}

			defer db.Close()

			Expect(db.Ping()).To(BeNil())
		})
	})
}
