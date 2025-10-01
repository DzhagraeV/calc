package LoggerListTimer

import (
	"bufio"
	"container/list"
	"errors"
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

// фрагмент — ключевые части (исправленный вариант предыдущего unbounded Logger)
func (l *FileLogger) Log(message string) error {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return ErrLoggerClosed
	}
	l.queue.PushBack(l.format(message))
	l.cond.Signal() // быстро разбудить воркер
	l.mu.Unlock()
	return nil
}

func (l *FileLogger) backgroundWriter() {
	defer l.wg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		// забираем очередь в pending и обрабатываем вне lock
		l.mu.Lock()
		for l.queue.Len() == 0 && !l.closed {
			l.cond.Wait()
		}
		// выход если закрыт и пусто
		if l.queue.Len() == 0 && l.closed {
			l.mu.Unlock()
			break
		}
		// swap: берем текущую очередь как pending и ставим новую пустую
		pending := l.queue
		l.queue = list.New()
		l.mu.Unlock()

		// обработать pending вне мьютекса
		for e := pending.Front(); e != nil; e = e.Next() {
			msg := e.Value.(string)
			_, _ = l.writer.WriteString(msg + "\n")
		}

		// flush по таймеру
		select {
		case <-ticker.C:
			_ = l.writer.Flush()
		default:
		}
	}

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
