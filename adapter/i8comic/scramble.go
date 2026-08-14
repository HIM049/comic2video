package i8comic

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/webp"
)

// UnscrambleImage 还原乱序切片图片,返回 PNG 字节
// (jquery.photo-0.4.js onImageLoaded / jquery.photo-0.6.js onImageLoadedForNew)
//
// data:     乱序图片字节
// aid:      漫画 ID
// filename: 图片文件名，如 "00001.webp"
func UnscrambleImage(data []byte, aid, filename string) ([]byte, error) {
	img, _, err := decodeImage(data)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	reassembled, err := reassemble(img, aid, name)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := png.Encode(&out, reassembled); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}
	return out.Bytes(), nil
}

// decodeImage 解码图片:优先标准库格式,webp 交给 x/image
func decodeImage(data []byte) (image.Image, string, error) {
	if img, format, err := image.Decode(bytes.NewReader(data)); err == nil {
		return img, format, nil
	}
	img, err := webp.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	return img, "webp", nil
}

// reassemble 对应 JS onImageLoaded 的画布重组循环
func reassemble(img image.Image, aid, name string) (*image.RGBA, error) {
	s := scrambleSlices(aid, name)

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w == 0 || h == 0 {
		return nil, fmt.Errorf("image has zero size (%dx%d)", w, h)
	}

	rem := h % s // JS: r = parseInt(H % s)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for g := 0; g < s; g++ {
		sliceH := h / s                // JS: h = Math.floor(H/s)
		dstY := sliceH * g             // JS: c = h*g(目标 y,从顶部往下)
		srcY := h - sliceH*(g+1) - rem // JS: w = H - h*(g+1) - r(源 y,从底部往上)
		if g == 0 {
			sliceH += rem // 首片多取 r 像素
		} else {
			dstY += rem
		}
		draw.Draw(dst, image.Rect(0, dstY, w, dstY+sliceH), img, image.Pt(0, srcY), draw.Src)
	}
	return dst, nil
}

// scrambleSlices 计算图片被切成的片数 (get_num)
func scrambleSlices(aid, name string) int {
	sum := md5.Sum([]byte(aid + name))
	hexStr := hex.EncodeToString(sum[:])
	last := hexStr[len(hexStr)-1] // JS: md5(...).substr(-1)

	// JS: n = 末位字符的 charCodeAt('0'-'9' → 48-57,'a'-'f' → 97-102)。
	// 注意:switch 匹配的是 charCode 值,因此数字末位(48-57)在 aid 不在
	// 特殊区段时不会命中任何 case,返回默认 10。
	n := int(last)

	// JS: aid 落在特殊区段时先取模
	aidNum, _ := strconv.Atoi(aid)
	switch {
	case aidNum >= 268850 && aidNum <= 421925:
		n %= 10
	case aidNum >= 421926:
		n %= 8
	}

	// JS: switch(n) { 0:2, 1:4, 2:6, 3:8, 4:10, 5:12, 6:14, 7:16, 8:18, 9:20 }
	switch n {
	case 0:
		return 2
	case 1:
		return 4
	case 2:
		return 6
	case 3:
		return 8
	case 4:
		return 10
	case 5:
		return 12
	case 6:
		return 14
	case 7:
		return 16
	case 8:
		return 18
	case 9:
		return 20
	}
	return 10 // 无匹配时 JS 返回初始值 10
}
