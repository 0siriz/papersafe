package payload

// Parse decodes a PAPERSAFE text encoding into its domain-specific payload
// type. The returned value is either *Keyshard or *Document.
func Parse(b []byte) (Payload, error) {
	domain, rest, err := splitHeader(b)
	if err != nil {
		return nil, err
	}

	switch domain {
	case DomainKeyshard:
		var k Keyshard
		if err := k.unmarshalBody(rest); err != nil {
			return nil, err
		}
		return &k, nil
	case DomainDocument:
		var d Document
		if err := d.unmarshalBody(rest); err != nil {
			return nil, err
		}
		return &d, nil
	default:
		return nil, ErrInvalidDomain
	}
}
