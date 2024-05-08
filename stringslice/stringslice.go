package stringslice

func Split(slice []string, sep string) [][]string { return genSplit(slice, sep, 0, -1) }

// Remove the leading and trailing empty strings from the slice
func Trim(slices [][]string) [][]string {
	if len(slices) == 0 || slices == nil {
		return slices
	}

	// Trim the head
	for len(slices) != 0 {
		if len(slices[0]) == 0 {
			slices = slices[1:]
			continue
		}
		break
	}

	for len(slices) != 0 {
		lastElement := slices[len(slices)-1]
		if len(lastElement) == 0 {
			slices = slices[:len(slices)-1]
			continue
		}
		break
	}
	return slices
}

func isEmpty(slice []string) bool {
	return slice == nil || len(slice) == 0
}

func Flatten(slices [][]string) []string {
	var out []string
	for _, outer := range slices {
		for _, inner := range outer {
			out = append(out, inner)
		}
	}
	return out
}

func genSplit(slice []string, sep string, sepSave, n int) [][]string {
	if n == 0 {
		return nil
	}
	// if sep == "" {
	// 	return explode(s, n)
	// }
	if n < 0 {
		n = Count(slice, sep) + 1
	}

	if n > len(slice)+1 {
		n = len(slice) + 1
	}
	a := make([][]string, n)
	n--
	i := 0
	for i < n {
		m := Index(slice, sep)
		if m < 0 {
			break
		}
		a[i] = slice[:m+sepSave]
		// Move 1 past the sep
		slice = slice[m+1:]
		i++
	}
	a[i] = slice
	return a[:i+1]
}

// Count counts the number of non-overlapping instances of substr in s.
// If substr is an empty string slice, Count returns 1 + the number of lines in s.
func Count(s []string, substr string) int {
	// special case
	if len(substr) == 0 {
		return len(s) + 1
	}
	n := 0
	for {
		i := Index(s, substr)
		if i == -1 {
			return n
		}
		n++
		// Move 1 past the sep
		s = s[i+1:]
	}
}

// Index returns the index of the first instance of substr in slice, or -1 if substr is not present in slice.
func Index(slice []string, substr string) int {
	for i, s := range slice {
		if s == substr {
			return i
		}
	}
	return -1
}
