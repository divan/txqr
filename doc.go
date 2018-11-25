/*Package txqr implements the transmission protocol over QR codes.

Intro

The protocol allows to send a relatively small (fits into the memory fast)
data of a known size. Stream data is not supported by design.

QR codes are supposed to be sent and received by means of optical displays
and sensors with unknown properties. Sender might be a 85 inch OLED TV with
240Hz rate, while receiver could be an old Android phone with 2MP camera and
bound by CPU allowing only 5FPS. Or vice versa. Protocol must adapt to all
cases.

The basic idea is to split the data into chunks, suitable for encoding as a
single QR frame, add frame header/footer information and run it in the loop.

- splitting into frame is crucial to adapt to desired QR code size/error recovery level

- header and footer contain enough information to uniquely identify frame and be able to restore the whole data even if all frames received out of order.

- loop is needed to make sure slow receiver has enough opportunity to restore from missed frames

All data should be within alphanumeric space.
No error correction is implemented, as QR code layer already has one.

Header

    offset/total|<data>

	both offset and total are represents byte position
	and printed in HEX

For, example:

 First chunk:

    0/11|hello

 Second chunk:

    5/11|world!


*/
package txqr
