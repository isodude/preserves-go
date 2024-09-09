package lib

// Raised whenever [preserves.binary.Decoder][] or [preserves.text.Parser][] detect invalid input.
type DecodeError error

// Raised whenever [preserves.binary.Encoder][] or [preserves.text.Formatter][] are unable to proceed.
type EncodeError error

// Raised whenever [preserves.binary.Decoder][] or [preserves.text.Parser][] discover that
// they want to read beyond the end of the currently-available input buffer in order to
// completely read an encoded value.
type ShortPacket DecodeError

type NotImplementedError error
