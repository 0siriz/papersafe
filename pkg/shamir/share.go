package shamir

import "crypto/subtle"

const (
	ShareXSize = 1
)

type Share struct {
	X      byte
	Y      []byte
	sealed bool
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

func (s *Share) Sealed() bool {
	return s.sealed
}

func (s *Share) Seal(key []byte) error {
	if s.sealed {
		return ErrInvalidSealCall
	}

	if err := s.seal(key); err != nil {
		return err
	}

	s.sealed = true

	return nil
}

func (s *Share) Unseal(key []byte) error {
	if !s.sealed {
		return ErrInvalidUnsealCall
	}

	if err := s.seal(key); err != nil {
		return err
	}

	s.sealed = false

	return nil
}

func (s *Share) seal(key []byte) error {
	if len(key) != len(s.Y) {
		return ErrInvalidSealKeyLength
	}

	_ = subtle.XORBytes(s.Y, s.Y, key)

	return nil
}
