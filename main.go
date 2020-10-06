package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"match/mycompress"
	"time"
)

const (
	compress int8 = iota
	decompress
)

var (
// TODO 压缩参数
)

type compressMethod func(in []byte, do int8) ([]byte, error)

func lzw_(in []byte, do int8) ([]byte, error) {
	if do == compress {
		buf := bytes.NewBuffer(nil)
		w := lzw.NewWriter(buf, lzw.LSB, 8)

		_, err := w.Write(in)
		if err != nil {
			return nil, err
		}
		w.Close()
		return buf.Bytes(), nil

	} else {
		r := lzw.NewReader(bytes.NewReader(in), lzw.LSB, 8)

		out, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		r.Close()
		return out, err
	}
}

func flate_(in []byte, do int8) ([]byte, error) {
	if do == compress {
		buf := bytes.NewBuffer(nil)
		w, err := flate.NewWriter(buf, flate.BestCompression)
		if err != nil {
			return nil, err
		}

		_, err = w.Write(in)
		if err != nil {
			return nil, err
		}
		w.Close()
		return buf.Bytes(), nil

	} else {
		r := flate.NewReader(bytes.NewReader(in))

		out, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		r.Close()
		return out, err
	}
}

func gzip_(in []byte, do int8) ([]byte, error) {
	if do == compress {
		buf := bytes.NewBuffer(nil)
		w, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
		if err != nil {
			return nil, err
		}

		_, err = w.Write(in)
		if err != nil {
			return nil, err
		}
		w.Close()
		return buf.Bytes(), nil

	} else {
		r, err := gzip.NewReader(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}

		out, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		r.Close()
		return out, err
	}
}

func zlib_(in []byte, do int8) ([]byte, error) {
	if do == compress {
		buf := bytes.NewBuffer(nil)
		w, err := zlib.NewWriterLevel(buf, gzip.BestCompression)
		if err != nil {
			return nil, err
		}

		_, err = w.Write(in)
		if err != nil {
			return nil, err
		}
		w.Close()
		return buf.Bytes(), nil

	} else {
		r, err := zlib.NewReader(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}

		out, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		r.Close()
		return out, err
	}
}

func zstd_(in []byte, do int8) ([]byte, error) {
	if do == compress {
		return mycompress.ZstdCompress(in), nil
	} else {
		out, err := mycompress.ZstdDecompress(in)

		if err != nil {
			return nil, err
		}
		return out, err
	}
}

func gozstd_(in []byte, do int8) ([]byte, error) {
	if do == compress {
		return mycompress.GozstdCompress(in), nil
	} else {
		out, err := mycompress.GozstdDecompress(in)

		if err != nil {
			return nil, err
		}
		return out, err
	}
}

func lz4_(in []byte, do int8) ([]byte, error) {
	if do == compress {
		out, err := mycompress.Lz4Compress(in)

		if err != nil {
			return nil, err
		}
		return out, err
	} else {
		out, err := mycompress.Lz4Decompress(in)

		if err != nil {
			return nil, err
		}
		return out, err
	}
}

func test(method compressMethod, inFile, outFile, fileDepcompress string) (int64, int64, int64) {
	start := time.Now()
	in := mycompress.ReadFileBufferio(inFile)
	out, err := method(in, compress)
	if err != nil {
		fmt.Println(err)
	} else {
		mycompress.WriteFileBufferio(outFile, out)
	}
	compressTime := time.Since(start)

	start = time.Now()
	in = mycompress.ReadFileBufferio(outFile)
	out, err = method(in, decompress)

	if err != nil {
		fmt.Println(err)
	} else {
		mycompress.WriteFileBufferio(fileDepcompress, out)
	}
	decompressTime := time.Since(start)

	size := mycompress.GetFileSize(outFile)
	return compressTime.Milliseconds(), decompressTime.Milliseconds(), size
}

func compare(in, out string) bool {
	return mycompress.CompareFile(in, out)
}

func main() {
	file := "xml"

	inFile := "./src/match/" + file
	outFile := "./src/match/" + file + ".compress"
	fileDepcompress := "./src/match/" + file + ".decompress"
	sourceFileSize := float32(mycompress.GetFileSize(inFile))

	// lzw
	compressTime, decompressTime, size := test(lzw_, inFile, outFile, fileDepcompress)
	fmt.Printf("lzw:compressTime:%d,decompressTime:%d,compressionRatio:%f.\n", compressTime, decompressTime, sourceFileSize/float32(size))
	fmt.Println(compare(inFile, fileDepcompress))
	// flate
	compressTime, decompressTime, size = test(flate_, inFile, outFile, fileDepcompress)
	fmt.Printf("flate:compressTime:%d,decompressTime:%d,compressionRatio:%f.\n", compressTime, decompressTime, sourceFileSize/float32(size))
	fmt.Println(compare(inFile, fileDepcompress))
	// gzip
	compressTime, decompressTime, size = test(gzip_, inFile, outFile, fileDepcompress)
	fmt.Printf("gzip:compressTime:%d,decompressTime:%d,compressionRatio:%f.\n", compressTime, decompressTime, sourceFileSize/float32(size))
	fmt.Println(compare(inFile, fileDepcompress))
	// zlib
	compressTime, decompressTime, size = test(zlib_, inFile, outFile, fileDepcompress)
	fmt.Printf("zlib:compressTime:%d,decompressTime:%d,compressionRatio:%f.\n", compressTime, decompressTime, sourceFileSize/float32(size))
	fmt.Println(compare(inFile, fileDepcompress))
	// zstd
	compressTime, decompressTime, size = test(zstd_, inFile, outFile, fileDepcompress)
	fmt.Printf("zstd:compressTime:%d,decompressTime:%d,compressionRatio:%f.\n", compressTime, decompressTime, sourceFileSize/float32(size))
	fmt.Println(compare(inFile, fileDepcompress))
	// gozstd
	compressTime, decompressTime, size = test(gozstd_, inFile, outFile, fileDepcompress)
	fmt.Printf("gozstd:compressTime:%d,decompressTime:%d,compressionRatio:%f.\n", compressTime, decompressTime, sourceFileSize/float32(size))
	fmt.Println(compare(inFile, fileDepcompress))
	// lz4 有问题
	compressTime, decompressTime, size = test(lz4_, inFile, outFile, fileDepcompress)
	fmt.Printf("lz4:compressTime:%d,decompressTime:%d,compressionRatio:%f.\n", compressTime, decompressTime, sourceFileSize/float32(size))
	fmt.Println(compare(inFile, fileDepcompress))
}
