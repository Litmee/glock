package treaty

import (
	"bufio"
)

// Decode message decoding
func Decode(r *bufio.Reader) ([]byte, error) {

	// Identification parameters used to solve TCP subcontract problems
	var sign byte = 0

	// Container for fetching byte streams
	s := make([]byte, 8)

	// Return the current number of readable bytes in the buffer through the Buffered method
	// Less than 5 to be processed
	if r.Buffered() < 8 {
		// return nil, err
		sign = byte(r.Buffered())
		s = make([]byte, sign)
	}

	// read message entity
	_, err := r.Read(s)
	if err != nil {
		return nil, err
	}

	// sign != 0 indicates that there is subcontracting and needs to be spliced
	if sign != 0 {
		_, _ = r.Peek(1)
		newS := make([]byte, 8-sign)
		_, err = r.Read(newS)
		if err != nil {
			return nil, err
		}
		s = append(s, newS...)
	}

	// Return
	return s, nil
}
