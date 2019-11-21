package tachymeter_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamiealquiza/tachymeter"
)

func TestTimelineEmptyMetrics(t *testing.T) {
	tmp, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmp)

	tline := &tachymeter.Timeline{}

	ta := tachymeter.New(&tachymeter.Config{Size: 30})
	m := ta.Calc()

	tline.AddEvent(m)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("WriteHtml panicked: %v", r)
		}
	}()

	tline.WriteHTML(tmp)
}
