package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Logger interface {
	Log(level, message string) error
	Close() error
}

type FileLogger struct {
	file     *os.File
	writer   *bufio.Writer
	ch       chan string
	closed   int32 // atomic flag
	wg       sync.WaitGroup
	flushDur time.Duration
}

var ErrLoggerClosed = errors.New("logger closed")
var ErrBufferFull = errors.New("log buffer full")

// NewFileLogger opens file in append mode and starts background writer.
// bufferSize - размер канала сообщений (например 10000).
// flushDur - раз в какой период делать flush (например 1s).
func NewFileLogger(fileName string, bufferSize int, flushDur time.Duration) (*FileLogger, error) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	l := &FileLogger{
		file:     f,
		writer:   bufio.NewWriterSize(f, 64*1024), // 64KB буфер
		ch:       make(chan string, bufferSize),
		flushDur: flushDur,
	}

	l.wg.Add(1)
	go l.backgroundWriter()
	return l, nil
}

func (l *FileLogger) format(level, msg string) string {
	// простой формат: timestamp level message
	return fmt.Sprintf("%s [%s] %s", time.Now().UTC().Format(time.RFC3339Nano), level, msg)
}

// Log пишет сообщение в канал. Если канал заполнен, возвращает ErrBufferFull.
func (l *FileLogger) Log(level, message string) error {
	if atomic.LoadInt32(&l.closed) == 1 {
		return ErrLoggerClosed
	}
	formatted := l.format(level, message)

	// non-blocking put: если канал полон — отбросим сообщение и вернём ошибку
	select {
	case l.ch <- formatted:
		return nil
	default:
		// можно пересчитать dropped metric здесь
		return ErrBufferFull
	}
}

func (l *FileLogger) backgroundWriter() {
	defer l.wg.Done()
	ticker := time.NewTicker(l.flushDur)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-l.ch:
			if !ok {
				// канал закрыт — выжать остаток и выйти
				l.flushAll()
				l.writer.Flush()
				l.file.Sync()
				l.file.Close()
				return
			}
			l.writer.WriteString(msg + "\n")
			// если буфер большой — можно Flush, но лучше полагаться на тикер
		case <-ticker.C:
			l.writer.Flush()
			// здесь можно иногда делать file.Sync() для durability по политике
		}
	}
}

func (l *FileLogger) flushAll() {
	for {
		select {
		case msg := <-l.ch:
			l.writer.WriteString(msg + "\n")
		default:
			return
		}
	}
}

// Close закрывает логгер: больше не принимает логи, ждёт фон. горутины.
func (l *FileLogger) Close() error {
	if !atomic.CompareAndSwapInt32(&l.closed, 0, 1) {
		return ErrLoggerClosed
	}
	close(l.ch)
	l.wg.Wait()
	// background writer уже закрыл файл
	return nil
}
