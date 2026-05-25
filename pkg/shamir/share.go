package shamir

const (
	// ShareXSize is the byte size of the Share.X field in binary encoding.
	ShareXSize = 1
)

// Share represents a single Shamir secret share with an X coordinate (the
// share index) and Y values (the evaluated polynomial values for each byte
// of the secret).
type Share struct {
	X byte
	Y []byte
}

// AppendBinary appends the binary encoding of the share to b.
func (s *Share) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, s.X)
	b = append(b, s.Y...)

	return b, nil
}

// MarshalBinary returns the binary encoding of the share.
func (s *Share) MarshalBinary() ([]byte, error) {
	b, err := s.AppendBinary(make([]byte, 0, ShareXSize+len(s.Y)))
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UnmarshalBinary decodes the binary encoding into the share.
func (s *Share) UnmarshalBinary(b []byte) error {
	buf := b

	s.X = buf[0]
	buf = buf[ShareXSize:]

	s.Y = make([]byte, len(buf))
	copy(s.Y, buf)

	return nil
}
