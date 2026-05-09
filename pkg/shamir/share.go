package shamir

const (
	ShareXSize = 1
)

type Share struct {
	X byte
	Y []byte
}

func (s *Share) AppendBinary(b []byte) ([]byte, error) {
	b = append(b, s.X)
	b = append(b, s.Y...)

	return b, nil
}

func (s *Share) MarshalBinary() ([]byte, error) {
	b, err := s.AppendBinary(make([]byte, 0, ShareXSize+len(s.Y)))
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Share) UnmarshalBinary(b []byte) error {
	buf := b

	s.X = buf[0]
	buf = buf[ShareXSize:]

	s.Y = make([]byte, len(buf))
	copy(s.Y, buf)

	return nil
}
