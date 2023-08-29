# SST Resource Binding for Go

SST lets you [bind resources](https://docs.sst.dev/resource-binding) to Lambda functions which both adds the right permissions and also injects environment variables containing information about that resource, for example it's ARN.

This package is a port of the [node package](https://docs.sst.dev/clients/) for using resource bindng in Lambda functions written in Go.

### Simple Example

If you have an S3 bucket that is bound to a function:

```ts
// Create an S3 bucket
const bucket = new Bucket(stack, "MyFiles");

new Function(stack, "MyFunction", {
  handler: "main.go",
  // Bind to function
  bind: [bucket],
});
```

Then in your Go function you can access the bucket name by doing the following:

```go
package main

import (
    "fmt"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/teamkeel/sst-go"
)

func main() {
    lambda.Start(handler)
}

func handler(ctx context.Context) error {
    // load the resource
    bucket, err := sst.Bucket("MyFiles")
    if err != nil {
      return err
    }

    bucket.BucketName // name of the bucket
}
```

### Supported Bindings

- `Secret` (including fallback values)
- `Parameter`
- `Function`
- `EventBus`
- `Topic`
- `Queue`
- `RDS`
- `Table`
