package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2oracle"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	dlog "oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
)

func main() {
	ctx := dlog.WithDefault(context.Background())

	// Start with a new, empty graph
	_, graph, _ := d2lib.Compile(ctx, "", nil, nil)

	// Initialize a ruler to measure glyphs of text
	ruler, _ := textmeasure.NewRuler()

	f, _ := os.ReadFile(filepath.Join("plan.sql"))

	queries := parseSQL(string(f))

	for i, q := range queries {
		graph = q.transformGraph(graph)

		// Turn the graph into a script
		script := d2format.Format(graph.AST)

		layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
			return d2dagrelayout.DefaultLayout, nil
		}
		// Compile the script with given theme and layout
		diagram, _, _ := d2lib.Compile(ctx, script, &d2lib.CompileOptions{
			LayoutResolver: layoutResolver,
			Ruler:          ruler,
		}, nil)

		// Render to SVG
		padding := int64(d2svg.DEFAULT_PADDING)
		out, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
			Pad: &padding,
		})

		// Write to disk
		_ = os.WriteFile(filepath.Join("svgs", fmt.Sprintf("step%d.svg", i)), out, 0600)
	}

	_ = os.WriteFile("out.d2", []byte(d2format.Format(graph.AST)), 0600)
}

type Query struct {
	Command string
	Table   string
	Column  string
	Type    string

	ForeignTable  string
	ForeignColumn string
}

func (q Query) transformGraph(g *d2graph.Graph) *d2graph.Graph {
	switch q.Command {
	case "create_table":
		// Create an object with the ID set to the table name
		newG, newKey, _ := d2oracle.Create(g, nil, q.Table)
		// Set the shape of the newly created object to be D2 shape type "sql_table"
		shape := "sql_table"
		newG, _ = d2oracle.Set(g, nil, fmt.Sprintf("%s.shape", newKey), nil, &shape)
		return newG
	case "add_column":
		newG, _ := d2oracle.Set(g, nil, fmt.Sprintf("%s.%s", q.Table, q.Column), nil, &q.Type)
		return newG
	case "add_foreign_key":
		newG, _, _ := d2oracle.Create(g, nil, fmt.Sprintf("%s.%s -> %s.%s", q.Table, q.Column, q.ForeignTable, q.ForeignColumn))
		return newG
	}
	return nil
}

func parseSQL(plan string) (out []Query) {
	lines := strings.Split(plan, "\n")

	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		out = append(out, parseSQLCommand(strings.Trim(line, ";")))
	}
	return out
}

func parseSQLCommand(command string) Query {
	q := Query{}

	words := strings.Split(command, " ")
	if strings.HasPrefix(command, "CREATE") {
		q.Command = "create_table"
		q.Table = words[2]
	} else if strings.Contains(command, "ADD COLUMN") {
		q.Command = "add_column"
		q.Table = words[2]
		q.Column = words[5]
		q.Type = words[6]
	} else if strings.Contains(command, "ADD CONSTRAINT") {
		q.Command = "add_foreign_key"
		q.Table = words[2]
		q.Column = strings.Trim(strings.Trim(words[8], "("), ")")
		q.ForeignTable = words[10]
		q.ForeignColumn = strings.Trim(strings.Trim(words[11], "("), ")")
	}

	return q
}
