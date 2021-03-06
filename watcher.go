package bitcoinwalletwatcher

import (
	"context"
	"log"
	"time"
)

// Watcher struct
type Watcher struct {
	info      *InfoFile
	loopInSec time.Duration
	addresses []string
	callbacks []WatcherCallback
}

// WatcherCallback method
type WatcherCallback func(string, uint64)

// NewWatcher returns new bitcoin wallet watcher
func NewWatcher(cfg *Config) (*Watcher, error) {
	info, err := NewInfoStorage(cfg.InfoFile)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		info:      info,
		loopInSec: time.Duration(cfg.DefaultLoopSec),
		addresses: cfg.Adresses,
	}, nil
}

// SetCallback sets a new callback
func (w *Watcher) SetCallback(callback WatcherCallback) {
	w.callbacks = append(w.callbacks, callback)
}

// Run runs the watcher
func (w *Watcher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.Stop()
			return
		default:
			w.readBlock()
			time.Sleep(time.Second * w.loopInSec)
		}
	}

}

func (w *Watcher) readBlock() {
	block, err := getBlockInformation(w.info.CurrentBlock)
	if err != nil {
		log.Fatal(err)
		return
	}

	// block not found yet
	if block == nil {
		return
	}

	for _, tx := range block.Tx {
		w.checkTx(tx)
	}

	w.info.Update(block.Height + 1)
}

func (w *Watcher) checkTx(tx transactionDetail) {
	for _, out := range tx.Out {
		if contains(w.addresses, out.Addr) {
			for _, c := range w.callbacks {
				c(out.Addr, out.Amount)
			}
		}
	}
}

// Stop runs when watcher stopping
func (w *Watcher) Stop() {
	if err := w.info.Save(); err != nil {
		log.Fatalf("error saving file %v", err)
	}
}
