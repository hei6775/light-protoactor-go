package actor

const PH = 1

func (m *PID) ResetEx() {
	m.Address = ""
	m.Id = ""
	m.p = nil
}

func (m PID) Clone() *PID {
	return NewPID(m.Address, m.Id)
}

func ClonePIDSlice(dst []*PID, src []*PID) []*PID {
	dst = []*PID{}

	for _, i := range src {
		dst = append(dst, i)
	}

	return dst
}
