package tachymeter_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/jamiealquiza/tachymeter"
)

func newTestTachy() *tachymeter.Tachymeter {
	ta := tachymeter.New(&tachymeter.Config{Size: 4})

	for i := 12; i < 16; i++ {
		ta.AddTime(time.Duration(i) * time.Millisecond)
	}
	return ta
}

func BenchmarkWriteHTMLTo(b *testing.B) {
	ta := newTestTachy()
	tl := tachymeter.Timeline{}
	tl.AddEvent(ta.Calc())

	b.ReportAllocs()
	b.ResetTimer()
	w := ioutil.Discard

	for i := 0; i < b.N; i++ {
		err := tl.WriteHTMLTo(w)
		if err != nil {
			b.Fatalf("can't write html: %s", err)
		}
	}
}
