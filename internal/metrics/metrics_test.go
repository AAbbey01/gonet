package metrics

import (
	"testing"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/AAbbey01/gonet/internal/check"

)


func TestRecord(t *testing.T){
	m := New(prometheus.NewRegistry())
	tgt := check.Target{Name: "web", Kind: "http"}
	m.Record(check.Result{Target: tgt, Up: true, Latency: 20*time.Millisecond})
	if got := testutil.ToFloat64(m.up.WithLabelValues("web","http")); got != 1{
		t.Fatalf("Failed to record up correctly: found %v, want 1", got)
	}
	m.Record(check.Result{Target: tgt, Up: false, Latency: 20*time.Millisecond})
	if got := testutil.ToFloat64(m.up.WithLabelValues("web","http")); got != 0{
		t.Fatalf("Failed to change record correctly: found %v, want 0", got)
	}

}