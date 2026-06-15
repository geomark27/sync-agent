package watcher

import (
	"context"
	"time"
)

// Debounce agrupa eventos rápidos. Acumula las rutas entrantes y, tras un
// periodo de quietud de 'delay' sin nuevos eventos, emite el conjunto
// acumulado como un único lote (deduplicado). Esto filtra la ráfaga de
// señales de guardado que producen muchos editores.
func Debounce(ctx context.Context, in <-chan string, delay time.Duration) <-chan []string {
	out := make(chan []string)

	go func() {
		defer close(out)

		pending := make(map[string]struct{})
		var timer *time.Timer
		var timerC <-chan time.Time

		flush := func() {
			if len(pending) == 0 {
				return
			}
			batch := make([]string, 0, len(pending))
			for p := range pending {
				batch = append(batch, p)
			}
			pending = make(map[string]struct{})
			select {
			case out <- batch:
			case <-ctx.Done():
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case path, ok := <-in:
				if !ok {
					flush()
					return
				}
				pending[path] = struct{}{}
				if timer != nil {
					timer.Stop()
				}
				timer = time.NewTimer(delay)
				timerC = timer.C
			case <-timerC:
				flush()
				timerC = nil
			}
		}
	}()

	return out
}
