# SST Resource Binding for Go

SST has a nice way to bind resources to functions which both adds the right permissions and also injects environment variables containing information about that resouce, for example it's ARN.

This package is a port of the node package for reading those values inside an SST app.

### Example

If you have an S3 bucket that is bound to a function:

```ts
// Create an S3 bucket
const bucket = new Bucket(stack, "myFiles");

new Function(stack, "myFunction", {
  handler: "main.go",
  // Bind to function
  bind: [bucket],
});
```

Then in `main.go` you can access the bucket name by doing the following:

```go
import (
    "fmt"
    "github.com/teamkeel/sst-go"
)

func main() {
    // load the resource
    bucket := sst.Bucket("myFiles")

    // access properties
    fmt.Println("bucket name:", bucket.BucketName)
}
```
