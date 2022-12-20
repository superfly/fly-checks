# Fly checks
Simple health-check framework.

## Usage
```golang
import (
	"context"
	"fmt"

	"github.com/superfly/fly-checks/check"
)

func main() {
	// Set an optional timeout.
	timeout := 200 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	suite := check.NewCheckSuite("my-test-function")

	// Passing check
	suite.AddCheck("passing_check", func() (string, error) {
		return "pass", nil
	})

	// Failing check
	suite.AddCheck("failing_check", func() (string, error) {
		return "failed", fmt.Errorf("random error")
	})

	// Check that will fail due to timeout
	suite.AddCheck("failing_check_due_to_timeout", func() (string, error) {
		time.Sleep(time.Second)
		return "failed", fmt.Errorf("random error")
	})
	
	
	// Optionally run something on check completion.
	// suite.OnCompletion = func() {}

	//
	go func() {
		suite.Process(ctx)
		cancel()
	}()

	<-ctx.Done():
	
        fmt.Println(suite.Result())
}
```

```bash
go run main.go
[✓] passing_check: pass (83ns)
[✗] failing_check: random error (105.38µs)
[✗] failing_check_due_to_timeout: Timed out (201ms)
