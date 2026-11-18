// Command icons builds BRQ's app and tray icons from build/brq_logo_source.webp.
//
//	go run ./tools/brq/icons
//
// Writes build/appicon.png (1024), build/trayicon.png (256) and
// build/windows/icon.ico (16-256, PNG-compressed entries). Every icon is a
// round badge: the emblem inside a violet ring, like the BRQ optimizer icon.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

// The emblem (hexagon + bolt) inside the 1024px source; the "#BRQ" caption
// below it is unreadable at icon sizes, so it is cropped away.
var emblem = image.Rect(170, 82, 852, 764)

func main() {
	f, err := os.Open("build/brq_logo_source.webp")
	if err != nil {
		log.Fatal(err)
	}
	src, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	crop := image.NewRGBA(image.Rect(0, 0, emblem.Dx(), emblem.Dy()))
	draw.Draw(crop, crop.Bounds(), src, emblem.Min, draw.Src)

	writePNG("build/appicon.png", icon(crop, 1024))
	writePNG("build/trayicon.png", icon(crop, 256))

	sizes := []int{16, 20, 24, 32, 40, 48, 64, 96, 128, 256}
	var entries [][]byte
	for _, s := range sizes {
		var b bytes.Buffer
		if err := png.Encode(&b, icon(crop, s)); err != nil {
			log.Fatal(err)
		}
		entries = append(entries, b.Bytes())
		// Optional: dump each size for a visual check.
		if dir := os.Getenv("BRQ_ICON_PREVIEW_DIR"); dir != "" {
			_ = os.WriteFile(fmt.Sprintf("%s/icon_%d.png", dir, s), b.Bytes(), 0o644)
		}
	}
	writeICO("build/windows/icon.ico", sizes, entries)
	fmt.Println("icons written")
}

// Badge frame, from the outside in, as fractions of the radius. Matches the
// BRQ optimizer icon: a solid violet ring, a thin lavender line, a sliver of
// violet, a thin dark line, then the emblem.
var (
	ringViolet   = color.RGBA{0x58, 0x30, 0x88, 0xff}
	ringLavender = color.RGBA{0xa6, 0x91, 0xbf, 0xff}
	ringDark     = color.RGBA{0x15, 0x0f, 0x22, 0xff}
)

type band struct {
	inner float64 // band runs from this fraction of the radius out to the previous band
	col   color.RGBA
}

var frame = []band{
	{0.875, ringViolet},
	{0.845, ringLavender},
	{0.835, ringViolet},
	{0.815, ringDark},
}

// emblemScale is the emblem's square size relative to the icon, so the
// hexagon's side points sit just inside the dark line.
const emblemScale = 0.94

// icon renders the badge at size x size: the emblem downscaled cleanly,
// sharpened at small sizes, inside the violet frame.
func icon(src *image.RGBA, size int) *image.RGBA {
	es := int(math.Round(float64(size) * emblemScale))
	img := resize(src, es)
	switch {
	case size <= 24:
		img = sharpen(contrast(img, 1.35), 0.9)
	case size <= 32:
		img = sharpen(contrast(img, 1.2), 0.7)
	case size <= 48:
		img = sharpen(img, 0.6)
	case size <= 96:
		img = sharpen(img, 0.3)
	}
	return badge(img, size)
}

// badge composes the emblem into the framed disc with anti-aliased edges.
func badge(emblem *image.RGBA, n int) *image.RGBA {
	fn := float64(n)
	R := fn / 2
	off := (n - emblem.Bounds().Dx()) / 2
	// Very small icons get relatively thicker lines so they don't vanish.
	minLine := 0.0
	if n <= 32 {
		minLine = 0.9
	}
	out := image.NewRGBA(image.Rect(0, 0, n, n))
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			d := math.Hypot(float64(x)+0.5-R, float64(y)+0.5-R)
			cover := math.Min(1, math.Max(0, R-d+0.5))
			if cover == 0 {
				continue
			}
			// Start from the emblem, then lay each band over it, outermost last.
			ex, ey := x-off, y-off
			var r, g, b float64
			if ex >= 0 && ey >= 0 && ex < emblem.Bounds().Dx() && ey < emblem.Bounds().Dy() {
				c := emblem.RGBAAt(ex, ey)
				r, g, b = float64(c.R), float64(c.G), float64(c.B)
			} else {
				r, g, b = float64(ringDark.R), float64(ringDark.G), float64(ringDark.B)
			}
			for i := len(frame) - 1; i >= 0; i-- {
				inner := frame[i].inner * R
				if i > 0 && minLine > 0 {
					outer := frame[i-1].inner * R
					if outer-inner < minLine {
						inner = outer - minLine
					}
				}
				t := math.Min(1, math.Max(0, d-inner+0.5))
				c := frame[i].col
				r = r*(1-t) + float64(c.R)*t
				g = g*(1-t) + float64(c.G)*t
				b = b*(1-t) + float64(c.B)*t
			}
			out.SetRGBA(x, y, color.RGBA{
				R: uint8(r * cover), G: uint8(g * cover), B: uint8(b * cover), A: uint8(255 * cover),
			})
		}
	}
	return out
}

// resize halves with a box filter until close to the target, then finishes
// with Catmull-Rom. One big Catmull-Rom step from ~680px to 16px samples only a
// few source pixels per output pixel and aliases the thin metal edges.
func resize(src *image.RGBA, size int) *image.RGBA {
	cur := src
	for cur.Bounds().Dx() >= size*2 {
		cur = halve(cur)
	}
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(dst, dst.Bounds(), cur, cur.Bounds(), draw.Src, nil)
	return dst
}

func halve(src *image.RGBA) *image.RGBA {
	w, h := src.Bounds().Dx()/2, src.Bounds().Dy()/2
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var r, g, b, a uint32
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					c := src.RGBAAt(2*x+dx, 2*y+dy)
					r += uint32(c.R)
					g += uint32(c.G)
					b += uint32(c.B)
					a += uint32(c.A)
				}
			}
			dst.SetRGBA(x, y, color.RGBA{uint8(r / 4), uint8(g / 4), uint8(b / 4), uint8(a / 4)})
		}
	}
	return dst
}

// contrast spreads tones around mid-grey so the light bolt separates from the
// darker hexagon once there are only a few pixels to draw them with.
func contrast(src *image.RGBA, k float64) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	f := func(v uint8) uint8 {
		return uint8(math.Max(0, math.Min(255, math.Round((float64(v)-118)*k+118))))
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := src.RGBAAt(x, y)
			dst.SetRGBA(x, y, color.RGBA{f(c.R), f(c.G), f(c.B), c.A})
		}
	}
	return dst
}

// sharpen is a 3x3 unsharp mask: out = in + amount*(in - blur).
func sharpen(src *image.RGBA, amount float64) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	at := func(x, y int) color.RGBA {
		x = max(b.Min.X, min(b.Max.X-1, x))
		y = max(b.Min.Y, min(b.Max.Y-1, y))
		return src.RGBAAt(x, y)
	}
	clamp := func(v float64) uint8 { return uint8(math.Max(0, math.Min(255, math.Round(v)))) }
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			var sr, sg, sb float64
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					c := at(x+dx, y+dy)
					sr += float64(c.R)
					sg += float64(c.G)
					sb += float64(c.B)
				}
			}
			c := at(x, y)
			dst.SetRGBA(x, y, color.RGBA{
				R: clamp(float64(c.R) + amount*(float64(c.R)-sr/9)),
				G: clamp(float64(c.G) + amount*(float64(c.G)-sg/9)),
				B: clamp(float64(c.B) + amount*(float64(c.B)-sb/9)),
				A: c.A,
			})
		}
	}
	return dst
}

func writePNG(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

func writeICO(path string, sizes []int, data [][]byte) {
	var b bytes.Buffer
	le := binary.LittleEndian
	_ = binary.Write(&b, le, [3]uint16{0, 1, uint16(len(sizes))})
	offset := 6 + 16*len(sizes)
	for i, s := range sizes {
		dim := uint8(s)
		if s >= 256 {
			dim = 0
		}
		_ = binary.Write(&b, le, struct {
			W, H, Colors, Reserved uint8
			Planes, BPP            uint16
			Size, Offset           uint32
		}{dim, dim, 0, 0, 1, 32, uint32(len(data[i])), uint32(offset)})
		offset += len(data[i])
	}
	for _, d := range data {
		b.Write(d)
	}
	if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}
