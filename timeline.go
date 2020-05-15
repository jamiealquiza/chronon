package tachymeter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"
)

// Timeline holds a []*timelineEvents,
// which nest *Metrics for analyzing
// multiple collections of measured events.
type Timeline struct {
	timeline []*timelineEvent
}

// timelineEvent holds a *Metrics and
// time that it was added to the Timeline.
type timelineEvent struct {
	Metrics *Metrics
	Created time.Time
}

// AddEvent adds a *Metrics to the *Timeline.
func (t *Timeline) AddEvent(m *Metrics) {
	t.timeline = append(t.timeline, &timelineEvent{
		Metrics: m,
		Created: time.Now(),
	})
}

const (
	tab string = `	`
	nl  string = "\n"
)

// WriteHTML takes an absolute path p and writes an
// html file to 'p/tachymeter-<timestamp>.html' of all
// histograms held by the *Timeline, in series.
func (t *Timeline) WriteHTML(p string) error {
	path, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	fname := fmt.Sprintf("%s/tachymeter-%d.html", path, time.Now().Unix())
	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("can't create %s: %v", fname, err)
	}
	defer f.Close()

	return t.WriteHTMLTo(f)
}

// WriteHTMLTo ...
func (t *Timeline) WriteHTMLTo(w io.Writer) error {
	_, err := io.WriteString(w, head)
	if err != nil {
		return err
	}

	// Append graph + info entry for each timeline
	// event.
	for n := range t.timeline {
		// Graph div.
		fmt.Fprintf(w, `%s<div class="graph">%s`, tab, nl)
		fmt.Fprintf(w, `%s%s<canvas id="canvas-%d"></canvas>%s`, tab, tab, n, nl)
		fmt.Fprintf(w, `%s</div>%s`, tab, nl)
		// Info div.
		fmt.Fprintf(w, `%s<div class="info">%s`, tab, nl)
		fmt.Fprintf(w, `%s<p><h2>Iteration %d</h2>%s`, tab, n+1, nl)
		io.WriteString(w, t.timeline[n].Metrics.String())
		fmt.Fprintf(w, "%s%s</p></div>%s", nl, tab, nl)
	}

	// Write graphs.
	for id, m := range t.timeline {
		err := genGraphHTML(w, m, id)
		if err != nil {
			return fmt.Errorf("can't generate graph.js: %v", err)
		}
	}

	_, err = io.WriteString(w, head)
	return err
}

var graphTmpl = template.Must(template.New("graph").Parse(graph))

// genGraphHTML takes a *timelineEvent and id (used for each graph
// html element ID) and writes a chart.js graph into w.
func genGraphHTML(w io.Writer, te *timelineEvent, id int) error {
	keys := []string{}
	values := []uint64{}

	for _, b := range *te.Metrics.Histogram {
		for k, v := range b {
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	keysj, _ := json.Marshal(keys)
	valuesj, _ := json.Marshal(values)

	err := graphTmpl.Execute(w, struct {
		CanvasID string
		Keys     string
		Values   string
	}{
		CanvasID: strconv.Itoa(id),
		Keys:     string(keysj),
		Values:   string(valuesj),
	})
	if err != nil {
		return fmt.Errorf("can't execute graph template: %v", err)
	}
	return nil
}
