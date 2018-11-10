/* Package protocol defines the transmission protocol over QR codes.

The protocol allows to send a relatively small (fits into the memory fast)
data of a known size. Stream data is not supported by design.

QR codes are supposed to be sent and received by means of optical displays
and sensors with unknown properties. Sender might be a 85 inch OLED TV with
240Hz rate, while receiver could be an old Android phone with 2MP camera and
bound by CPU allowing only 5FPS. Or vice versa. Protocol must adapt to all
cases.

The basic idea is to split the data into chunks, suitable for encoding as a
single QR frame, add frame header/footer information and run it in the loop.
 - splitting into frame is crucial to adapt to desired QR code size/error
   recovery level
 - header and footer contain enough information to uniquely identify frame and
   be able to restore the whole data even if all frames received out of order.
 - loop is needed to make sure slow receiver has enough opportunity to restore
   from missed frames

All data should be within alphanumeric space.
No error correction is implemented, as QR code layer already has one.

Header

    current/total|<data>

	both current and total are represents byte position
	(as in Seek) and printed in HEX

For, example:

 First chunk:

    0/11|hello

 Second chunk:

    5/11|world!


*/
package protocol

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Encoder represents protocol encoder.
type Encoder struct {
	ChunkLen int
}

// NewEncoder creates and inits a new encoder for the given chunk length.
func NewEncoder(n int) *Encoder {
	return &Encoder{
		ChunkLen: n,
	}
}

// EncodeReader encodes data from reader and splits it into chunks to be
// futher converted to QR code frames.
func (e *Encoder) EncodeReader(r io.Reader) ([]string, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read: %v", err)
	}

	return e.Encode(string(data))
}

// Encode encodes data from reader and splits it into chunks to be
// futher converted to QR code frames.
func (e *Encoder) Encode(str string) ([]string, error) {
	if len(str) < e.ChunkLen {
		return []string{str}, nil
	}

	numChunks := len(str)/e.ChunkLen + 1

	// TODO(divan): use sync.Pool as this probably will be used many times
	ret := make([]string, numChunks)
	count := 0
	for start := 0; start < len(str); start += e.ChunkLen {
		end := start + e.ChunkLen
		if end > len(str) {
			end = len(str)
		}

		ret[count] = e.frame(start, len(str), str[start:end])
		count++
	}
	return ret, nil
}

func (e *Encoder) frame(count, total int, str string) string {
	return fmt.Sprintf("%x/%x|%s", count, total, str)
}
