// Two semaphore implementations
//
// Two semaphore implementations, in support of exercises of "The Little Book
// of Semaphores" (LBS) by Allen B. Downey.
package sem

import "sync"

// ChanSem implements a semaphore using a buffered channel of empty structs.
//
// Method implementation is very simple.
//
// Construct with NewChanSem.
type ChanSem chan struct{}

// NewChanSem constructs a new ChanSem with given initial and maximum counts.
//
// A small variation from LBS is that the maxCount must be given.  If the
// Signal method is called when the underlying channel is full, it will block,
// a behavior not described in LBS.  In most applications a maxCount will be
// known and this is not an issue.  Note also that there is no memory cost to
// asking for an arbitrarily high maxCount.
func NewChanSem(initCount, maxCount int) ChanSem {
	if maxCount < 1 {
		panic("maxCount must be > 0")
	}
	if initCount < 0 || initCount > maxCount {
		panic("initCount must >= 0 and <= maxCount")
	}
	s := make(ChanSem, maxCount)
	for i := 0; i < initCount; i++ {
		s.Signal()
	}
	return s
}

// Signal implements the semaphore "signal" operation.
//
// In implementation, it sends an item to the underlying channel.  If another
// goroutine was blocked on receive from the empty channel, it (well one
// goroutine) would be "signaled" to proceed.
func (s ChanSem) Signal() { s <- struct{}{} }

// SignalN signals n times.
//
// Introduced in section 3.7.6 Preloaded turnstile, with the note that "the
// multiple signals are not atomic; that is, the signaling thread might be
// interrupted in the loop.
func (s ChanSem) SignalN(n int) {
	for i := 0; i < n; i++ {
		s.Signal()
	}
}

// Wait implements the semaphore "wait" operation.
//
// In implementation, it receives an item from the underlying channel.  If
// the channel is empty, it "waits".
func (s ChanSem) Wait() { <-s }

// CountSem implements a semaphore with an integer count and a sync.Cond.
//
// In contrast to ChanSem, method implementation for CountSem is more complex.
//
// An advantage though is that no maxCount must be specified.
//
// Construct with NewCountSem.
type CountSem struct {
	Count int
	Cond  sync.Cond
}

// NewCountSem constructs a new CountSem with given initial count.
//
// The given initial count must be >= 0
func NewCountSem(initCount int) *CountSem {
	if initCount < 0 {
		panic("initCount must be >= 0")
	}
	return &CountSem{initCount, sync.Cond{L: &sync.Mutex{}}}
}

// Signal implements the semaphore "signal" operation.
func (cs *CountSem) Signal() {
	cs.Cond.L.Lock()
	cs.Count++
	cs.Cond.L.Unlock()
	cs.Cond.Broadcast()
}

// Wait implements the semaphore "wait" operation.
func (cs *CountSem) Wait() {
	cs.Cond.L.Lock()
	cs.Count--
	for cs.Count < 0 {
		cs.Cond.Wait()
	}
	cs.Cond.L.Unlock()
}

// Barrier introduced in section 3.7.7
type Barrier struct {
	n, count                     int
	mutex, turnstile, turnstile2 ChanSem
}

func NewBarrier(n int) *Barrier {
	return &Barrier{
		n:          n,
		count:      0,
		mutex:      NewChanSem(1, 1),
		turnstile:  NewChanSem(0, 1),
		turnstile2: NewChanSem(0, 1),
	}
}

func (b *Barrier) Phase1() {
	b.mutex.Wait()
	b.count++
	if b.count == b.n {
		b.turnstile.SignalN(b.n)
	}
	b.mutex.Signal()
	b.turnstile.Wait()
}

func (b *Barrier) Phase2() {
	b.mutex.Wait()
	b.count--
	if b.count == 0 {
		b.turnstile2.SignalN(b.n)
	}
	b.mutex.Signal()
	b.turnstile2.Wait()
}

func (b *Barrier) Wait() {
	b.Phase1()
	b.Phase2()
}

// Lightswitch introduced in section 4.2.2
type Lightswitch struct {
	counter int
	mutex   ChanSem
}

func NewLightswitch() *Lightswitch {
	return &Lightswitch{mutex: NewChanSem(1, 1)}
}

func (l *Lightswitch) Lock(semaphore ChanSem) {
	l.mutex.Wait()
	l.counter++
	if l.counter == 1 {
		semaphore.Wait()
	}
	l.mutex.Signal()
}

func (l *Lightswitch) Unlock(semaphore ChanSem) {
	l.mutex.Wait()
	l.counter--
	if l.counter == 0 {
		semaphore.Signal()
	}
	l.mutex.Signal()
}
