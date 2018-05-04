package lazy

import (
	"database/sql"
	. "github.com/franela/goblin"
	"github.com/jensneuse/kube-utils/pkg/config"
	_ "github.com/lib/pq"
	. "github.com/onsi/gomega"
	"testing"
)

func TestPostgresTemplate(t *testing.T) {

	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	dsn := "postgresql://postgres@localhost:15432/postgres?sslmode=disable"

	l := New()

	g.Describe("postgres", func() {

		g.After(func() {
			// tear down
			l.Cleanup()
		})

		g.It("should be ready to connect", func() {

			// request whatever we need
			l.CreateAndForwardPostgres(config.NAMESPACE, "postgres-test")

			db, err := sql.Open("postgres", dsn)
			if err != nil {
				t.Fatal(err)
			}

			defer db.Close()

			Expect(db.Ping()).To(BeNil())
		})
	})
}
