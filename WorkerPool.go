package main

import (
	"sync"
)

// WorkerPool - пул воркеров.
// tasks - очередь задач (неограниченная, хранится в памяти).
type WorkerPool struct {
	tasks      []func()
	workersQty int

	mu        sync.Mutex
	cond      *sync.Cond
	running   bool           // принимает задачи, пока true
	activeWg  sync.WaitGroup // считает выполняющиеся задачи
	workersWg sync.WaitGroup // ждём завершения горутин-воркеров
}

// NewWorkerPool создаёт пул и запускает воркеры.
func NewWorkerPool(numberOfWorkers int) *WorkerPool {
	wp := &WorkerPool{
		tasks:      make([]func(), 0),
		workersQty: numberOfWorkers,
		running:    true,
	}
	wp.cond = sync.NewCond(&wp.mu)

	// стартуем воркеры
	wp.workersWg.Add(numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		go wp.worker()
	}

	return wp
}

// внутренный цикл воркера
func (wp *WorkerPool) worker() {
	defer wp.workersWg.Done()

	for {
		wp.mu.Lock()
		// ждем появления задач, пока пул запущен
		for len(wp.tasks) == 0 && wp.running {
			wp.cond.Wait()
		}

		// если есть задача — забираем и выполним
		if len(wp.tasks) > 0 {
			task := wp.tasks[0]
			// удалить первый элемент (эффективно для интервью-решения)
			wp.tasks = wp.tasks[1:]
			// пометить как активную
			wp.activeWg.Add(1)
			wp.mu.Unlock()

			// выполнить вне мьютекса
			task()
			wp.activeWg.Done()
			// и продолжить цикл, не завершая воркера
			continue
		}

		// иначе: задач нет и пул не принимает новых -> выходим
		wp.mu.Unlock()
		return
	}
}

// Submit добавляет задачу в пул.
// Если пул уже остановлен (после Stop/StopWait), задача выполняется синхронно caller'ом.
func (wp *WorkerPool) Submit(task func()) {
	wp.mu.Lock()
	if !wp.running {
		// пул не принимает задачи — выполнить синхронно
		wp.mu.Unlock()
		task()
		return
	}
	// добавить в очередь и разбудить один воркер
	wp.tasks = append(wp.tasks, task)
	wp.cond.Signal()
	wp.mu.Unlock()
}

// SubmitWait добавляет задачу и блокирует вызывающий поток до её завершения.
func (wp *WorkerPool) SubmitWait(task func()) {
	done := make(chan struct{})
	wrapped := func() {
		defer close(done)
		task()
	}

	wp.Submit(wrapped)
	<-done
}

// Stop — остановить пул: перестать принимать новые задачи, очистить очередь (не выполнять ожидающие),
// дождаться завершения только тех задач, которые выполняются в данный момент.
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	// перестаём принимать задачи
	wp.running = false
	// очищаем очередь — эти задачи не будут выполнены
	wp.tasks = nil
	// разбудить всех воркеров, чтобы они могли проверить состояние и выйти
	wp.cond.Broadcast()
	wp.mu.Unlock()

	// дождаться выполнения уже запущенных задач
	wp.activeWg.Wait()
	// дождаться выхода воркеров
	wp.workersWg.Wait()
}

// StopWait — остановить пул: перестать принимать новые задачи, но дождаться выполнения
// всех задач, включая те, что были в очереди.
func (wp *WorkerPool) StopWait() {
	wp.mu.Lock()
	wp.running = false
	// разбудить всех воркеров: они не будут ждать, увидят running==false, но продолжат обрабатывать
	// оставшиеся в очереди задачи (логика worker'а позволяет этим воркерам забрать задачи и выполнить их).
	wp.cond.Broadcast()
	wp.mu.Unlock()

	// дождёмся, пока все воркеры полностью завершат (значит очередь обработана и активные задачи закончены)
	wp.workersWg.Wait()
	// на всякий случай: если какие-то задачи всё ещё отмечены как активные, дождёмся их
	wp.activeWg.Wait()
}
