package warped

import "io"

const (
	// TODO: Use vet-style tristate flag.
	// Leaving as constants for the moment since I'm playing
	// with the stdlib and flag imports too much.

	zeno = true
	// zeno = false
	stall = true
	// stall = false
	corrupt = true
	// corrupt = false
)

func Reader(r io.Reader) io.Reader {
	if r == nil {
		return r
	}
	if zeno {
		r = &zenoreader{r: r}
	}
	if stall {
		r = &stallreader{r: r}
	}
	if corrupt {
		r = &corruptreader{r: r}
	}
	return r
}

// stallreader introduces (0, nil) responses.
type stallreader struct {
	r     io.Reader
	stall int
}

func (r *stallreader) Read(b []byte) (n int, err error) {
	r.stall++
	if r.stall%10 == 1 {
		return 0, nil
	}
	return r.r.Read(b)
}

// zenoreader responds with partial data.
type zenoreader struct {
	r   io.Reader
	buf []byte
}

func (r *zenoreader) Read(b []byte) (int, error) {
	if len(r.buf) == 0 {
		r.buf = make([]byte, len(b))
		n, err := r.r.Read(r.buf)
		if err != nil || n == 0 {
			copy(b, r.buf)
			r.buf = nil
			return n, err
		}
		r.buf = r.buf[:n]
	}
	n := len(r.buf)/2 + 1
	copy(b, r.buf[:n])
	r.buf = r.buf[n:]
	return n, nil
}

// corruptreader corrupts the extra (unused) read buffer.
type corruptreader struct {
	r io.Reader
}

func (r *corruptreader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	// TODO: Be way more efficient here (use copy, doubling each time)
	// TODO: Scribble on the whole buffer
	// TODO: scribble in a pattern like 5ADFACE
	for i := n; 0 <= i && i < len(b) && i < n+128; i++ {
		b[i] = byte(i)
	}
	return n, err
}
