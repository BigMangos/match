package mycompress

import (
	"bytes"
	"io/ioutil"

	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4"
	"github.com/valyala/gozstd"
)

// github.com/klauspost/compress/zstd
func ZstdCompress(in []byte) []byte {
	var encoder, _ = zstd.NewWriter(nil)
	return encoder.EncodeAll(in, make([]byte, 0, len(in)))

}
func ZstdDecompress(in []byte) ([]byte, error) {
	var decoder, _ = zstd.NewReader(nil)
	out, err := decoder.DecodeAll(in, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// github.com/valyala/gozstd
func GozstdCompress(in []byte) []byte {
	return gozstd.CompressLevel(nil, in, gozstd.DefaultCompressionLevel)
}

func GozstdDecompress(in []byte) ([]byte, error) {
	out, err := gozstd.Decompress(nil, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// github.com/pierrec/lz4
func Lz4Compress(in []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zw := lz4.NewWriter(buf)

	_, err := zw.Write(in)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func Lz4Decompress(in []byte) ([]byte, error) {
	zw := lz4.NewReader(bytes.NewBuffer(in))

	out, err := ioutil.ReadAll(zw)
	if err != nil {
		return nil, err
	}
	return out, nil
}
