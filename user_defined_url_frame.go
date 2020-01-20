package id3v2

import "io"

// UserDefinedURLFrame is used to work with WXXX frames.
// There can be many UserDefinedURLFrame but the Desciption fields need to be unique.
type UserDefinedURLFrame struct {
	Encoding    Encoding
	Description string
	Value       string
}

func (uduf UserDefinedURLFrame) Size() int {
	return 1 + encodedSize(uduf.Description, uduf.Encoding) + len(uduf.Encoding.TerminationBytes) + encodedSize(uduf.Value, uduf.Encoding)
}

func (uduf UserDefinedURLFrame) WriteTo(w io.Writer) (n int64, err error) {
	return useBufWriter(w, func(bw *bufWriter) {
		bw.WriteByte(uduf.Encoding.Key)
		bw.EncodeAndWriteText(uduf.Description, uduf.Encoding)
		bw.Write(uduf.Encoding.TerminationBytes)
		bw.EncodeAndWriteText(uduf.Value, uduf.Encoding)
	})
}

func parseUserDefinedURLFrame(br *bufReader) (Framer, error) {
	encoding := getEncoding(br.ReadByte())
	description := br.ReadText(encoding)

	if br.Err() != nil {
		return nil, br.Err()
	}

	value := getBytesBuffer()
	defer putBytesBuffer(value)

	if _, err := value.ReadFrom(br); err != nil {
		return nil, err
	}

	uduf := UserDefinedURLFrame{
		Encoding:    encoding,
		Description: decodeText(description, encoding),
		Value:       decodeText(value.Bytes(), encoding),
	}

	return uduf, nil
}
