package obj

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"
)

func Parse(rd io.Reader) iter.Seq[[][]float64] {
	return func(yield func([][]float64) bool) {
		var vertices [][]float64
	out:
		for scanner := bufio.NewScanner(rd); scanner.Scan(); {
			fields := strings.Fields(strings.TrimSpace(scanner.Text()))
			if len(fields) < 2 {
				continue
			}
			command, args := fields[0], fields[1:]
			if strings.HasPrefix(command, "#") {
				continue
			}
			switch command {
			case "v":
				var vertex []float64
				for _, arg := range args {
					if val, err := strconv.ParseFloat(arg, 64); err == nil {
						vertex = append(vertex, val)
					}
				}
				vertices = append(vertices, vertex)
			case "f":
				var points [][]float64
				for _, arg := range args {
					if index, err := strconv.Atoi(arg); err == nil {
						points = append(points, vertices[index-1])
					}
				}
				if !yield(points) {
					break out
				}
			default:
				fmt.Printf("OBJ Parser: Skipping unknown command '%s'\n", command)
			}
		}
	}
}
