package model

import (
	"errors"
	"sync"
)

type Funds struct {
	Amount  int64
	Blocked int64
	mu      sync.Mutex
}

func (f *Funds) HasAvailableAmount(amount int64) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.Amount-f.Blocked >= amount
}

func (f *Funds) Add(amount int64) int64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Amount += amount
	return f.Amount
}

func (f *Funds) Subtract(amount int64) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Amount-f.Blocked < amount {
		return f.Amount, errors.New("not enough funds")
	}
	f.Amount -= amount
	return f.Amount, nil
}

func (f *Funds) Block(amount int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Amount-f.Blocked < amount {
		return errors.New("not enough funds")
	}
	f.Blocked += amount
	return nil
}

func (f *Funds) Unblock(amount int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Blocked < amount {
		return errors.New("not enough blocked funds")
	}
	f.Blocked -= amount
	return nil
}
