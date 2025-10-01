package LoggerList

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// Упрощённый FileLogger:
// - внутренний буфер — container/list (unbounded).
// - один sync.Cond для синхронизации производителей и потребителя.
// - Log() добавляет сообщение и не блокирует (потому что очередь unbounded).
// - фоновой воркер батчит все накопленные сообщения и пишет в файл.
// - Close() помечает logger как закрытый, ждёт фона, делает final flush+sync+close.
var ErrLoggerClosed = errors.New("logger closed")

type FileLogger struct {
	file   *os.File
	writer *bufio.Writer

	mu    sync.Mutex
	cond  *sync.Cond
	queue *list.List

	closed bool
	wg     sync.WaitGroup
}

// NewFileLogger: открывает файл в append и запускает фонового воркера.
func NewFileLogger(fileName string) (*FileLogger, error) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	l := &FileLogger{
		file:   f,
		writer: bufio.NewWriterSize(f, 32*1024),
		queue:  list.New(),
	}
	l.cond = sync.NewCond(&l.mu)

	l.wg.Add(1)
	go l.backgroundWriter()

	return l, nil
}

// простое форматирование: timestamp + message
func (l *FileLogger) format(msg string) string {
	return time.Now().UTC().Format(time.RFC3339Nano) + " " + msg
}

// Log добавляет сообщение в очередь.
// В этой реализации очередь не ограничена, поэтому Log не блокирует (кроме краткой блокировки mutex).
// Если логгер закрыт — вернёт ErrLoggerClosed.
func (l *FileLogger) Log(message string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return ErrLoggerClosed
	}

	l.queue.PushBack(l.format(message))
	// разбудить фонового воркера (если он ждёт)
	l.cond.Signal()
	return nil
}

// backgroundWriter читает все накопленные сообщения пачкой и пишет в файл.
// Делает периодический flush (каждую секунду).
func (l *FileLogger) backgroundWriter() {
	defer l.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		// собрать все имеющиеся сообщения (или ждать, если их нет)
		l.mu.Lock()
		for l.queue.Len() == 0 && !l.closed {
			l.cond.Wait()
		}

		// если очередь пуста и logger закрыт — выйти
		if l.queue.Len() == 0 && l.closed {
			l.mu.Unlock()
			break
		}

		// вынуть все сообщения в срез (батч)
		var msgs []string
		for e := l.queue.Front(); e != nil; e = e.Next() {
			if s, ok := e.Value.(string); ok {
				msgs = append(msgs, s)
			}
		}
		l.queue.Init() // очистить очередь
		l.mu.Unlock()

		// записать батч
		for _, m := range msgs {
			_, _ = l.writer.WriteString(m + "\n")
		}

		// Быстрый ненавязчивый flush: если тикер сработал — flush.
		// (Это даёт баланс между производительностью и потерей при креше.)
		select {
		case <-ticker.C:
			_ = l.writer.Flush()
		default:
		}
	}

	// финальный flush + sync + close
	_ = l.writer.Flush()
	_ = l.file.Sync()
	_ = l.file.Close()
}

// Close помечает логгер как закрытый, пробуждает воркер и ждёт его завершения.
// После Close любые вызовы Log будут возвращать ErrLoggerClosed.
func (l *FileLogger) Close() error {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return ErrLoggerClosed
	}
	l.closed = true
	// разбудить всех (и производителя, и потребителя), чтобы они ушли/вернули ошибку
	l.cond.Broadcast()
	l.mu.Unlock()

	l.wg.Wait()
	return nil
}

// пример использования
func main() {
	logger, err := NewFileLogger("app.log")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		_ = logger.Log(fmt.Sprintf("message %d", i))
	}

	_ = logger.Close()
}
