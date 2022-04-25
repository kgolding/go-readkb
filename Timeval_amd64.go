package readkb

type Timeval [16]byte

func (tv Timeval) Equals(tv2 Timeval) bool {
	for i, b := range tv {
		if b != tv2[i] {
			return false
		}
	}
	return true
}
