// Package base91 implements base91 encoding
package base91

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

type Encoding struct {
	encode    string
	decodeMap [256]byte
}

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&()*+,./:;<=>?@[]^_`{|}~\""

// NewEncoding returns a new Encoding defined by the given alphabet,
// which must be a 91-byte string.
func NewEncoding(encoder string) *Encoding {
	e := new(Encoding)
	e.encode = encoder
	for i := 0; i < len(e.decodeMap); i++ {
		e.decodeMap[i] = 0xFF
	}
	for i := 0; i < len(encoder); i++ {
		e.decodeMap[encoder[i]] = byte(i)
	}
	return e
}

// StdEncoding is the standard base91 encoding
var StdEncoding = NewEncoding(encodeStd)

var removeNewlinesMapper = func(r rune) rune {
	if r == '\r' || r == '\n' {
		return -1
	}
	return r
}

func (enc *Encoding) Encode(dst, src []byte) {
	if len(src) == 0 {
		return
	}

	var n, b, v uint64 = 0, 0, 0
	var pos int

	b = 0
	pos = 0

	for _, abyte := range src {
		ubyte := uint64(abyte)

		b |= (ubyte & 255) << n

		n += 8
		if n > 13 {
			v = b & 8191
			if v > 88 {
				b >>= 13
				n -= 13
			} else {
				v = b & 16383
				b >>= 14
				n -= 14
			}
			dst[pos] = enc.encode[v%91]
			pos++
			dst[pos] = enc.encode[v/91]
			pos++
		}
	}

	if n > 0 {
		dst[pos] = enc.encode[b%91]
		if n > 7 || b > 90 {
			dst[pos] = enc.encode[b/91]
		}
	}
}

// EncodeToString returns the base91 encoding of src.
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := make([]byte, len(src)*2, (len(src) * 2))
	enc.Encode(buf, src)
	return string(buf)
}

type encoder struct {
	err error
	enc *Encoding
	w   io.Writer
	out []byte // output buffer
}

func (e *encoder) Write(p []byte) (n int, err error) {
	if e.err != nil {
		return 0, e.err
	}

	e.enc.Encode(e.out, p)
	return
}

func (e *encoder) Close() error {
	return e.err
}

// NewEncoder returns a new base91 stream encoder.  Data written to
// the returned writer will be encoded using enc and then written to w.
func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser {
	return &encoder{enc: enc, w: w}
}

/*
 * Decoder
 */

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base91 data at input byte " + strconv.FormatInt(int64(e), 10)
}

// This method assumes that src has been
// stripped of all supported whitespace ('\r' and '\n').
func (enc *Encoding) decode(dst, src []byte) (num int, err error) {

	var v, b int64 = -1, 0
	var c, outpos int = 0, 0
	var n uint64 = 0

	decodemap := make(map[rune]int)

	for pos, char := range encodeStd {
		decodemap[char] = pos
	}

	for _, char := range string(src) {

		if strings.Contains(encodeStd, string(char)) {
			c = decodemap[char]
		}

		if v < 0 {
			v = int64(c)
		} else {
			v += int64(c) * 91
			b |= v << uint64(n)
			if (v & 8191) > 88 {
				n += 13
			} else {
				n = 14
			}
			for {
				dst[outpos] = byte(b & 255)
				outpos++
				b >>= 8
				n -= 8
				if n <= 7 {
					break
				}
			}
			v = -1
		}
	}
	if (v + 1) > 0 {
		dst[outpos] = byte((b | v<<n) & 255)
		outpos++
	}

	return len(dst), nil
}

// Decode decodes src using the encoding enc.  It writes at most
// DecodedLen(len(src)) bytes to dst and returns the number of bytes
// written.
// New line characters (\r and \n) are ignored.
func (enc *Encoding) Decode(dst, src []byte) (n int, err error) {
	src = bytes.Map(removeNewlinesMapper, src)
	n, err = enc.decode(dst, src)
	return
}

// DecodeString returns the bytes represented by the base91 string s.
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	s = strings.Map(removeNewlinesMapper, s)
	dbuf := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dbuf, []byte(s))
	return dbuf[:n], err
}

type decoder struct {
	err    error
	enc    *Encoding
	r      io.Reader
	out    []byte
	outbuf [1024 / 4 * 3]byte
}

func (d *decoder) Read(p []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}

	d.enc.Decode(d.out, p)
	return len(d.out), nil

}

type newlineFilteringReader struct {
	wrapped io.Reader
}

func (r *newlineFilteringReader) Read(p []byte) (int, error) {
	n, err := r.wrapped.Read(p)
	for n > 0 {
		offset := 0
		for i, b := range p[0:n] {
			if b != '\r' && b != '\n' {
				if i != offset {
					p[offset] = b
				}
				offset++
			}
		}
		if offset > 0 {
			return offset, err
		}
		// Previous buffer entirely whitespace, read again
		n, err = r.wrapped.Read(p)
	}
	return n, err
}

// NewDecoder constructs a new base91 stream decoder.
func NewDecoder(enc *Encoding, r io.Reader) io.Reader {
	return &decoder{enc: enc, r: &newlineFilteringReader{r}}
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base91-encoded data.
func (enc *Encoding) DecodedLen(n int) int { return n * 2 }
