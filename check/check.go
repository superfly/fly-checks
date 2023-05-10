package check

import (
	"fmt"
	"sync"
	"time"
)

type CheckFunction func() (string, error)

type Check struct {
	Name      string
	CheckFunc CheckFunction
	startTime time.Time
	endTime   time.Time
	message   string
	err       error
	m         sync.RWMutex
}

func (c *Check) Process() {
	c.m.Lock()
	c.startTime = time.Now()
	c.m.Unlock()

	msg, err := c.CheckFunc()
	c.m.Lock()
	c.message, c.err = msg, err
	c.endTime = time.Now()
	c.m.Unlock()
}

func (c *Check) Error() string {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.err.Error()
}

func (c *Check) Passed() bool {
	c.m.RLock()
	defer c.m.RUnlock()

	if c.startTime.IsZero() || c.endTime.IsZero() {
		return false
	}
	return c.err == nil
}

func (c *Check) executionTime() time.Duration {
	c.m.RLock()
	defer c.m.RUnlock()

	if !c.endTime.IsZero() {
		return RoundDuration(c.endTime.Sub(c.startTime), 2)
	}
	return RoundDuration(time.Now().Sub(c.startTime), 2)
}

func (c *Check) Result() string {
	c.m.RLock()
	defer c.m.RUnlock()

	if c.startTime.IsZero() {
		return fmt.Sprintf("[-] %s: %s", c.Name, "Not processed")
	}
	if !c.startTime.IsZero() && c.endTime.IsZero() {
		return fmt.Sprintf("[✗] %s: %s (%s)", c.Name, "Timed out", c.executionTime())
	}
	if c.Passed() {
		return fmt.Sprintf("[✓] %s: %s (%s)", c.Name, c.message, c.executionTime())
	} else {
		return fmt.Sprintf("[✗] %s: %s (%s)", c.Name, c.err.Error(), c.executionTime())
	}
}

func (c *Check) RawResult() string {
	c.m.RLock()
	defer c.m.RUnlock()

	if c.Passed() {
		return c.message
	}
	return c.err.Error()
}
