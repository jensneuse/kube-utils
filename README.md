# kube-utils
This repository offers utility functions for easy integrating go applications in kubernetes

Works for both client side and in-cluster testing.

Lets say you'd like to spin up your test environment, run integration tests and cleanup everything, this is all you have to do:

```go
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
```

Your environment might be a bit more complex. Working with multiple pods is as easy as:

````go
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
````

## Running tests

````bash
go test ./pkg/... -namespace=default
````

## Contributions

Feel free to submit additions (e.g. more pod templates) via pull requests.
Don't forget to add tests.