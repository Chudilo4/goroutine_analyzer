# Анализатор горутин

Сама установка библиотеки проходит по старинке
```bash
go get github.com/Chudilo4/goroutine_analyzer@latest
```
## Пример использования
```go
package main
import (
	"github.com/Chudilo4/goroutine_analyzer"
	"time"
	"log"
)

func main() {
	stats := goroutine_analyzer.NewGoroutineStats("my_service")
	// Запуск HTTP сервера для экспорта метрик
	// Маршрут для того что бы забирать информацию о статистике
	// маршрут и адрес можете поменять по вашему усмотрению
	err := stats.RunExportMetricPoint("/metrics", ":8070")
	if err != nil {
		log.Panic(err)
    }

	// Отмечаем запуск горутины
	stats.Add("func")
	go func() {
        // Отмечаем остановку горутины
		defer stats.Done("func")
		// Логика горутины
	}()

    // Отмечаем запуск горутины
    stats.Add("worker.func")
	go func() {
		// Отмечаем остановку горутины
		defer stats.Done("worker.func")
		// Логика горутины
	}()

	// Обновление метрик каждую секунду
	go func() {
		for {
            stats.UpdateMetrics()
			<-time.After(time.Second)
		}
	}()

    // Отображение метрик в логах каждую секунду
	go func() {
		for {
            log.Println(stats.GetMapCount())
            <-time.After(time.Second)
        }   
    }()
	
    // Для корректного завершения программы ожидаем остановку горутин
    stats.Wait()
}
```

В docker-compose указан пример как интегрировать показания в grafana и prometheus

## Новые версии
Для того что бы вести новую версию нужно влить изменение в main ветку после чего проставить тег следующей версии
### Пример
Создаём новую ветку
```bash
git branch feature/new_feature
```
Добавляем изменения
```bash
git add -A
````
Комитим
```bash
git commit -m "Новые изменения"
```
Пушим
```bash
git push
```
Переходим на ветку Develop
```bash
git checkout develop
```
Сливаем ветку Develop
```bash
git merge feature/new_feature
```
Переходим на ветку Main
```bash
git checkout main
```
Сливаем в ветку Main
```bash
git merge develop
```
Ставим Тег
```bash
git tag v0.0.2
```
Пушим Тег
```bash
git push origin v0.0.2
```
