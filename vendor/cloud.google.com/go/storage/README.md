## Cloud Storage [![GoDoc](https://godoc.org/cloud.google.com/go/storage?status.svg)](https://godoc.org/cloud.google.com/go/storage)

- [About Cloud Storage](https://cloud.google.com/storage/)
- [API documentation](https://cloud.google.com/storage/docs)
- [Go client documentation](https://godoc.org/cloud.google.com/go/storage)
- [Complete sample programs](https://github.com/GoogleCloudPlatform/golang-samples/tree/master/storage)

### Example Usage

First create a `storage.Client` to use throughout your application:

[snip]:# (storage-1)
```go
client, err := storage.NewClient(ctx)
if err != nil {
	log.Fatal(err)
}
```

[snip]:# (storage-2)
```go
// Read the object1 from bucket.
rc, err := client.Bucket("bucket").Object("object1").NewReader(ctx)
if err != nil {
	log.Fatal(err)
}
defer rc.Close()
body, err := io.ReadAll(rc)
if err != nil {
	log.Fatal(err)
}
```