package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/xpwu/go-log/log"
)

// DecompressResponse 根据 Content-Encoding 头部解压缩响应数据
func DecompressResponse(data []byte, header http.Header, logger *log.Logger) (string, error) {
	contentEncoding := strings.ToLower(header.Get("Content-Encoding"))

	logger.Debug(fmt.Sprintf("response: %d bytes, encoding: %s", len(data), contentEncoding))

	switch contentEncoding {
	case "gzip":
		return decompressGzip(data)

	case "deflate":
		return decompressDeflate(data)

	case "":
		// 无压缩
		return string(data), nil

	default:
		// 不支持的编码格式，返回错误
		return "", fmt.Errorf("unsupported Content-Encoding: %s", contentEncoding)
	}
}

// decompressGzip 解压缩 GZIP 数据
func decompressGzip(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("gzip reader creation failed: %v", err)
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("gzip decompression failed: %v", err)
	}

	return string(decompressed), nil
}

// decompressDeflate 解压缩 Deflate 数据
func decompressDeflate(data []byte) (string, error) {
	reader := flate.NewReader(bytes.NewReader(data))
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("deflate decompression failed: %v", err)
	}

	return string(decompressed), nil
}
