# Async Primitives (asyncx)

Import: `github.com/Abraxas-365/manifesto/internal/asyncx`

Prefer `asyncx` over raw goroutines for all concurrent operations.

---

## Map — concurrent transform

```go
results, err := asyncx.Map(ctx, items, func(ctx context.Context, item T) (R, error) {
    return process(item)
})
```

Runs all items concurrently, collects results in order, returns first error.

---

## Future — promise/async value

```go
fut := asyncx.NewFuture(func() (T, error) {
    return heavyOp()
})
// Do other work...
value, err := fut.Await(ctx)
```

---

## Pool — bounded concurrency

```go
pool := asyncx.NewPool(10) // max 10 goroutines
pool.Submit(ctx, task)
pool.Wait()
```

Use when you need backpressure on unbounded input.

---

## Retry with backoff

```go
value, err := asyncx.RetryWithBackoff(ctx, 3, time.Second, func(ctx context.Context) (T, error) {
    return flakyOp(ctx)
})
```

Exponential backoff. Respects context cancellation.

---

## Race — first success wins

```go
result, err := asyncx.Race(ctx, task1, task2, task3)
```

Returns the result of whichever task succeeds first; cancels the rest.

---

## Debounce & Throttle

```go
debounced := asyncx.Debounce(fn, 300*time.Millisecond)
throttled := asyncx.Throttle(fn, 100*time.Millisecond)
```
