package prefix

import (
	"sort"
	"strings"
)

// printLabelsMultiline prints multiple labels with a user-defined alignment.
func PrintLabelsMultiline(level int, title, indent string, labels map[string]string, w PrefixWriter) {
	inner := indent[len(title)+2:]
	w.Write(level, "%s: %s", title, inner)

	if labels == nil || len(labels) == 0 {
		w.WriteLine("<none>")
		return
	}

	// to print labels in the sorted order
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		w.WriteLine("<none>")
		return
	}
	sort.Strings(keys)

	for i, key := range keys {
		if i != 0 {
			w.Write(level, "%s", indent)
		}
		w.Write(LEVEL_0, "%s=%s\n", key, labels[key])
		i++
	}
}

const maxAnnotationLen int = 140

// outputAnnotationsMultilineWithIndent prints multiple annotations with a user-defined alignment.
// If annotation string is too long, we omit chars more than 140 length.
func PrintAnnotationsMultiline(level int, title, indent string, annotations map[string]string, w PrefixWriter) {
	inner := indent[len(title)+2:]
	w.Write(level, "%s: %s", title, inner)

	// to print labels in the sorted order
	keys := make([]string, 0, len(annotations))
	for key := range annotations {
		keys = append(keys, key)
	}
	if len(annotations) == 0 {
		w.WriteLine("<none>")
		return
	}
	sort.Strings(keys)
	for i, key := range keys {
		if i != 0 {
			w.Write(level, indent)
		}
		value := strings.TrimSuffix(annotations[key], "\n")
		if (len(value)+len(key)+2) > maxAnnotationLen || strings.Contains(value, "\n") {
			w.Write(LEVEL_0, "%s:\n", key)
			for _, s := range strings.Split(value, "\n") {
				w.Write(LEVEL_0, "%s  %s\n", indent, shorten(s, maxAnnotationLen-2))
			}
		} else {
			w.Write(LEVEL_0, "%s: %s\n", key, value)
		}
		i++
	}
}

func shorten(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}
