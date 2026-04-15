package scanner

// Diff holds the changes between two port snapshots.
type Diff struct {
	Opened []Port
	Closed []Port
}

// HasChanges reports whether there are any differences.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare computes the difference between a previous and current set of ports.
// Ports are identified by their Protocol + Address + Port combination.
func Compare(previous, current []Port) Diff {
	prevSet := toSet(previous)
	currSet := toSet(current)

	var diff Diff

	for key, p := range currSet {
		if _, exists := prevSet[key]; !exists {
			diff.Opened = append(diff.Opened, p)
		}
	}

	for key, p := range prevSet {
		if _, exists := currSet[key]; !exists {
			diff.Closed = append(diff.Closed, p)
		}
	}

	return diff
}

func toSet(ports []Port) map[string]Port {
	set := make(map[string]Port, len(ports))
	for _, p := range ports {
		set[p.String()] = p
	}
	return set
}
