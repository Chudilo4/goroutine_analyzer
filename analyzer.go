package goroutine_analyzer

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

// GoroutineStats - менеджер горутин
type GoroutineStats struct {
	mu             sync.RWMutex         // Мьютекс на чтение и запись горутин
	wg             sync.WaitGroup       // Сущность для завершения работы горутин
	goroutines     map[string]int       // Карта кол-ва запущенных функций
	serviceName    string               // Имя запущенного сервиса
	goroutineCount *prometheus.GaugeVec // Сборщик статистики
	server         *http.Server
}

// NewGoroutineStats - создать новый менеджер обработки горутин
// serviceName - наименование сервиса
func NewGoroutineStats(serviceName string) *GoroutineStats {
	goroutineCount := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "goroutine_count",
			Help: "Number of goroutines running",
		},
		[]string{"service", "name"},
	)
	prometheus.MustRegister(goroutineCount)
	goroutineStats := GoroutineStats{
		goroutines:     make(map[string]int),
		serviceName:    serviceName,
		goroutineCount: goroutineCount,
	}
	return &goroutineStats
}

// Add - Отметить запуск горутины.
// name - наименование горутины
func (gs *GoroutineStats) Add(name string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.wg.Add(1)
	gs.goroutines[name]++
}

// Done - отметить что горутина завершена.
// name - наименование горутины
func (gs *GoroutineStats) Done(name string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.wg.Done()
	gs.goroutines[name]--
}

// GetMapCount - Получить карту запущенных горутин.
func (gs *GoroutineStats) GetMapCount() map[string]int {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.goroutines
}

// UpdateMetrics - Обновление статистики по горутинам.
func (gs *GoroutineStats) UpdateMetrics() {
	for name, count := range gs.GetMapCount() {
		gs.goroutineCount.WithLabelValues(gs.serviceName, name).Set(float64(count))
	}
}

// Wait - Ожидать завершения горутин.
func (gs *GoroutineStats) Wait() {
	gs.wg.Wait()
}

// RunExportMetricPoint - Запуск снятия метрики для prometheus.
// pattern - точка входа для prometheus
// addr - сетевой адрес для прослушивания tcp соединения
func (gs *GoroutineStats) RunExportMetricPoint(pattern, addr string) error {
	var server http.Server
	server.Addr = addr

	gs.server = &server

	http.Handle(pattern, promhttp.Handler())
	err := gs.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// StopExportMetricPoint - Остановка сервера метрики для prometheus.
// ctx - контекст приложения
func (gs *GoroutineStats) StopExportMetricPoint(ctx context.Context) error {
	err := gs.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
