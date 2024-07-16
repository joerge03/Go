package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"mgrep/worker"
	"mgrep/worklist"

	"github.com/alexflint/go-arg"
)

func discoverDirs(w *worklist.Worklist, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			nextpath := filepath.Join(path, entry.Name())
			discoverDirs(w, nextpath)
		} else {
			w.Add(worklist.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

var args struct {
	SearchTerm string `arg:"positional,required"`
	SearchDir  string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)

	var workerWg sync.WaitGroup

	wl := worklist.New(200)

	results := make(chan worker.Result, 200)

	numOfWorkers := 10
	workerWg.Add(1)

	go func() {
		defer workerWg.Done()
		discoverDirs(&wl, args.SearchDir)
		wl.Finalize(numOfWorkers)
	}()

	for i := 0; i < numOfWorkers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for {
				workEntry := wl.Next()
				if workEntry.Path != "" {
					workerResult := worker.FindInFile(workEntry.Path, args.SearchTerm)
					if workerResult != nil {
						for _, r := range workerResult.Inner {
							results <- r
						}
					}
				} else {
					return
				}
			}
		}()
	}

	blockWorkerWg := make(chan struct{})

	go func() {
		workerWg.Wait()
		close(blockWorkerWg)
	}()

	var displayWorkerWg sync.WaitGroup

	displayWorkerWg.Add(1)
	go func() {
		for {
			select {
			case result := <-results:
				fmt.Printf("%v: %v - %v \n", result.Path, result.Line, result.LineNum)
			case <-blockWorkerWg:
				if len(results) == 0 {
					displayWorkerWg.Done()
					return
				}
			}
		}
	}()

	displayWorkerWg.Wait()
}
