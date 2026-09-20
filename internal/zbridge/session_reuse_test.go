package zbridge

import (
	"context"
	"sync"
	"testing"
	"time"
)

type reuseStub struct {
	mu      sync.Mutex
	next    int
	created []string
	deleted []string
}

func (s *reuseStub) CreateChatSession(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	id := string(rune('a'+s.next)) + "-reuse-test"
	// unique ids: reuse-test-1, ...
	id = "reuse-test-" + itoa(s.next)
	s.created = append(s.created, id)
	return id, nil
}

func (s *reuseStub) DeleteChatSession(ctx context.Context, ids ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deleted = append(s.deleted, ids...)
	return nil
}

func (s *reuseStub) counts() (c, d int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.created), len(s.deleted)
}

func waitReuse(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

func TestPoolReusesUpToMaxUses(t *testing.T) {
	prev := config.SessionReuseCount
	config.SessionReuseCount = 3
	defer func() { config.SessionReuseCount = prev }()

	b := &reuseStub{}
	p := NewSessionPool(b, 1)
	p.SetMaxUses(3)
	p.Start()
	waitReuse(t, "warmup", func() bool { return p.Ready() == 1 })

	// Use the same pooled session 3 times: first 2 releases re-queue, no delete.
	var id string
	for i := 0; i < 2; i++ {
		got, err := p.Acquire(context.Background(), time.Second)
		if err != nil {
			t.Fatalf("acquire %d: %v", i, err)
		}
		if i == 0 {
			id = got
		} else if got != id {
			t.Fatalf("acquire %d = %q, want reused %q", i, got, id)
		}
		p.Release(got)
		waitReuse(t, "re-queue", func() bool { return p.Ready() == 1 })
		if c, d := b.counts(); d != 0 {
			t.Fatalf("after use %d: deleted=%d, want 0 (reused)", i+1, d)
		} else {
			_ = c
		}
		if uses := p.Uses(id); uses != i+1 {
			t.Fatalf("uses(%q) = %d, want %d", id, uses, i+1)
		}
	}

	// 3rd use hits the limit: deleted upstream + refilled with a new ID.
	got, err := p.Acquire(context.Background(), time.Second)
	if err != nil {
		t.Fatalf("acquire 3: %v", err)
	}
	if got != id {
		t.Fatalf("acquire 3 = %q, want reused %q", got, id)
	}
	p.Release(got)
	waitReuse(t, "delete+refill after max", func() bool {
		_, d := b.counts()
		return d == 1 && p.Ready() == 1
	})
	b.mu.Lock()
	del := append([]string{}, b.deleted...)
	b.mu.Unlock()
	if len(del) != 1 || del[0] != id {
		t.Fatalf("deleted = %v, want [%s]", del, id)
	}
	// New generation is stocked (different ID, zero uses).
	p.mu.Lock()
	_, stillTracked := p.uses[id]
	p.mu.Unlock()
	if stillTracked {
		t.Fatalf("uses entry for rotated %q not cleared", id)
	}
	p.Shutdown()
}

func TestPoolMaxUsesOneIsThrowaway(t *testing.T) {
	prev := config.SessionReuseCount
	config.SessionReuseCount = 1
	defer func() { config.SessionReuseCount = prev }()

	b := &reuseStub{}
	p := NewSessionPool(b, 1)
	p.SetMaxUses(1)
	p.Start()
	waitReuse(t, "warmup", func() bool { return p.Ready() == 1 })
	id, err := p.Acquire(context.Background(), time.Second)
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	p.Release(id)
	waitReuse(t, "immediate delete", func() bool {
		_, d := b.counts()
		return d == 1
	})
	p.Shutdown()
}
