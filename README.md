# kube-utils
This repository offers utility functions for easy integrating go applications in kubernetes

Works for both client side and in-cluster testing.

Lets say you'd like to spin up your test environment, run integration tests and cleanup everything, this is all you have to do:

```go
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

		g.After(func() {
			// tear down
			lazy.Cleanup()
		})

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
		})
	})
}
```
