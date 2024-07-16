package worklist

type Entry struct {
	Path string
}

type Worklist struct {
	Jobs chan Entry
}

func (w *Worklist) Add(entry Entry) {
	w.Jobs <- entry
}

func (w *Worklist) Next() Entry {
	next := <-w.Jobs
	return next
}

func New(bufSize int) Worklist {
	worklist := Worklist{make(chan Entry, bufSize)}
	return worklist
}

func NewJob(name string) Entry {
	job := Entry{name}
	return job
}

func (worklist *Worklist) Finalize(worklistSize int) {
	for i := 0; i <= worklistSize; i++ {
		worklist.Add(Entry{""})
	}
}
