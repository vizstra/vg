// +build !goci
package vg

// #cgo linux pkg-config: glfw3 xxf86vm x11 gl glu xrandr xi xcursor gthread-2.0 freetype2
// #cgo linux LDFLAGS: -lm
// #define NANOVG_GL3_IMPLEMENTATION
// #include <stdlib.h>
// #include <math.h>
// #include <GL/gl.h>
// #include "nanovg.h"
// #include "nanovg_gl.h"
// #include "nanovg_gl_utils.h"
// #define FONS_USE_FREETYPE
// #include "fontstash.h"
// #include "stb_image.h"
// #include "stb_truetype.h"
import "C"

import (
	"errors"
	// "fmt"
	"image/color"
	"unsafe"
)

type LineCap C.int

const (
	BUTT LineCap = iota
	ROUND
	SQUARE
	BEVEL
	MITER
)

type Winding C.int

const (
	CCW Winding = 1
	CW  Winding = 2
)

type PatternRepeat C.int

const (
	NOREPEAT = 0
	REPEATX  = 0x01 // Repeat image pattern in X direction
	REPEATY  = 0x02 // Repeat image pattern in Y direction
)

type Align C.int

const (
	// Horizontal align
	ALIGN_LEFT   Align = 1 << 0 // Default, align text horizontally to left.
	ALIGN_CENTER Align = 1 << 1 // Align text horizontally to center.
	ALIGN_RIGHT  Align = 1 << 2 // Align text horizontally to right.

	// Vertical align
	ALIGN_TOP      Align = 1 << 3 // Align text vertically to top.
	ALIGN_MIDDLE   Align = 1 << 4 // Align text vertically to middle.
	ALIGN_BOTTOM   Align = 1 << 5 // Align text vertically to bottom.
	ALIGN_BASELINE Align = 1 << 6 // Default, align text vertically to baseline.
)

type MipmapFlag C.int

const (
	IMAGE_GENERATE_MIPMAPS MipmapFlag = 1 << 0 // Generate mipmaps during creation of the image.
)

func toColor(c color.Color) C.NVGcolor {
	ri, gi, bi, ai := c.RGBA()
	return C.nvgRGBA(C.uchar(ri), C.uchar(gi), C.uchar(bi), C.uchar(ai))
}

type Context struct {
	cbase *C.NVGcontext
	fonts map[string]*Font
}

const FONT_DEFAULT = "DejaVuSans"

// NewContext returns a New Context
func NewContext() Context {
	ctx := Context{
		C.nvgCreateGL3(C.NVG_ANTIALIAS | C.NVG_STENCIL_STROKES),
		make(map[string]*Font),
	}
	// ctx.NewFont(FONT_DEFAULT, "/usr/share/fonts/truetype/freefont/FreeSans.ttf")
	ctx.NewFont(FONT_DEFAULT, "res/Roboto/Roboto-Regular.ttf")
	// ctx.NewFont(FONT_DEFAULT, "/usr/share/fonts/truetype/msttcorefonts/Verdana_Bold.ttf")
	// ctx.NewFont(FONT_DEFAULT, "/usr/share/fonts/truetype/msttcorefonts/Times_New_Roman.ttf")
	return ctx
}

// BeginFrame begins drawing a new frame.
// BeginFrame defines the size of the window to render to in relation currently
// set viewport (i.e. glViewport on GL backends). Device pixel ration allows to
// control the rendering on Hi-DPI devices.
//
// For example, GLFW returns two dimensions for an opened window: window size and
// frame buffer size. In that case you would set the window Width/Height to the
// window size ratio to: frameBufferWidth / windowWidth.
func (self *Context) BeginFrame(w, h int, ratio float64) {
	C.nvgBeginFrame(self.cbase, C.int(w), C.int(h), C.float(ratio))
}

// EndFrame ends drawing, flushing remaining render state.
func (self *Context) EndFrame() {
	C.nvgEndFrame(self.cbase)
}

// Save pushes and saves the current render state into a state stack.
// A matching nvgRestore() must be used to restore the state.
func (self *Context) Save() {
	C.nvgSave(self.cbase)
}

// Restore pops and restores current render state.
func (self *Context) Restore() {
	C.nvgRestore(self.cbase)
}

// Reset will reset current render state to default
// values. Does not affect the render state stack.
func (self *Context) Reset() {
	C.nvgReset(self.cbase)
}

// StrokeColor sets current stroke style to a solid color.
func (self *Context) StrokeColor(color color.Color) {
	C.nvgStrokeColor(self.cbase, toColor(color))
}

// StrokePaint sets current stroke style to a paint, which can
// be a one of the gradients or a pattern.
func (self *Context) StrokePaint(paint Paint) {
	C.nvgStrokePaint(self.cbase, paint.cbase)
}

// FillColor sets current fill cstyle to a solid color.
func (self *Context) FillColor(color color.Color) {
	C.nvgFillColor(self.cbase, toColor(color))
}

// FillPaint sets current fill style to a paint, which can be
// a one of the gradients or a pattern.
func (self *Context) FillPaint(paint Paint) {
	C.nvgFillPaint(self.cbase, C.NVGpaint(paint.cbase))
}

// MiterLimit sets the miter limit of the stroke style.
// Miter limit controls when a sharp corner is beveled.
func (self *Context) MiterLimit(limit float64) {
	C.nvgMiterLimit(self.cbase, C.float(limit))
}

// StrokeWidth sets the stroke witdth of the stroke style.
func (self *Context) StrokeWidth(size float64) {
	C.nvgStrokeWidth(self.cbase, C.float(size))
}

// LineCap sets how the end of the line (cap) is drawn,
// Can be one of: BUTT (default), ROUND, SQUARE, etc.
func (self *Context) LineCap(cap LineCap) {
	C.nvgLineCap(self.cbase, C.int(cap))
}

// LineJoin sets how sharp path corners are drawn.
// Can be one of MITER (default), ROUND, BEVEL.
func (self *Context) LineJoin(join LineCap) {
	C.nvgLineJoin(self.cbase, C.int(join))
}

// Sets the transparency applied to all rendered shapes.
// Alreade transparent paths will get proportionally
// more transparent as well.
func (self *Context) GlobalAlpha(alpha float64) {
	C.nvgGlobalAlpha(self.cbase, C.float(alpha))
}

// ResetTransform resets current transform to a
// identity matrix.
func (self *Context) ResetTransform() {
	C.nvgResetTransform(self.cbase)
}

// Transform premultiplies current coordinate system by specified
// matrix. The parameters are interpreted as matrix as follows:
//   [a c e]
//   [b d f]
//   [0 0 1]
func (self *Context) Transform(a, b, c, d, e, f float64) {
	C.nvgTransform(self.cbase, C.float(a), C.float(b), C.float(c), C.float(d), C.float(e), C.float(f))
}

// Translate the current coordinate system.
func (self *Context) Translate(x, y float64) {
	C.nvgTranslate(self.cbase, C.float(x), C.float(y))
}

// Rotate the current coordinate system.
// The angle is specifid in radians.
func (self *Context) Rotate(angle float64) {
	C.nvgRotate(self.cbase, C.float(angle))
}

// Skew the current coordinate system along X axis.
// The angle is specifid in radians.
func (self *Context) SkewX(angle float64) {
	C.nvgSkewX(self.cbase, C.float(angle))
}

// Skew the current coordinate system along Y axis.
// The angle is specifid in radians.
func (self *Context) SkewY(angle float64) {
	C.nvgSkewY(self.cbase, C.float(angle))
}

// Scale the current coordinate system.
func (self *Context) Scale(x, y float64) {
	C.nvgScale(self.cbase, C.float(x), C.float(y))
}

// NewImage creates image by loading it from the disk from
// specified file name.  Returns the image handle.
func (self *Context) NewImage(filename string, flags int) *Image {
	cs := C.CString(filename)
	defer C.free(unsafe.Pointer(cs))
	return &Image{C.nvgCreateImage(self.cbase, cs, C.int(flags)), self}
}

// Creates image by loading it from the specified chunk of memory.
// Returns handle to the image.
func (self *Context) NewImageFromData(data []byte, flags int) *Image {
	return &Image{
		C.nvgCreateImageMem(self.cbase, C.int(flags), (*C.uchar)(unsafe.Pointer(&data[0])), C.int(len(data))),
		self,
	}
}

// // Creates image from specified image data.
// // Returns handle to the image.
// int nvgCreateImageRGBA(NVGcontext* ctx, int w, int h, int imageFlags, const unsigned char* data);

type Image struct {
	handle  C.int
	context *Context
}

// // Updates image data specified by image handle.
// void nvgUpdateImage(NVGcontext* ctx, int image, const unsigned char* data);

// // Returns the domensions of a created image.
// void nvgImageSize(NVGcontext* ctx, int image, int* w, int* h);
func (self *Image) Size() (float64, float64) {
	var w, h C.int
	C.nvgImageSize(self.context.cbase, self.handle, &w, &h)
	return float64(w), float64(h)
}

// // Deletes created image.
func (self *Image) Delete() {
	C.nvgDeleteImage(self.context.cbase, self.handle)
}

func (self *Image) Draw(x, y float64) {
	w, h := self.Size()
	p := self.context.ImagePattern(x, y, w, h, 0.0/180.0*3.14, *self, 0, 1.0)
	self.context.BeginPath()
	self.context.RoundedRect(x, y, w, h, 10)
	self.context.FillPaint(p)
	self.context.Fill()
	self.context.ClosePath()
}

type Paint struct {
	cbase C.NVGpaint
}

func (self *Paint) Release() {
	// C.free(self.cbase)
}

func (self Context) LinearGradient(x1, y1, x2, y2 float64, a, b color.Color) Paint {
	c, d := toColor(a), toColor(b)
	return Paint{C.nvgLinearGradient(self.cbase, C.float(x1), C.float(y1), C.float(x2), C.float(y2), c, d)}
}

func (self Context) RadialGradient(cx, cy, i, o float64, a, b color.Color) Paint {
	c, d := toColor(a), toColor(b)
	return Paint{C.nvgRadialGradient(self.cbase, C.float(cx), C.float(cy), C.float(i), C.float(o), c, d)}
}

func (self Context) BoxGradient(x, y, w, h, r, f float64, a, b color.Color) Paint {
	c, d := toColor(a), toColor(b)
	return Paint{C.nvgBoxGradient(self.cbase, C.float(x), C.float(y), C.float(w), C.float(h), C.float(r), C.float(f), c, d)}
}

func (self Context) ImagePattern(cx, cy, w, h, angle float64, image Image, repeat int, alpha float64) Paint {
	return Paint{C.nvgImagePattern(self.cbase,
		C.float(cx),
		C.float(cy),
		C.float(w),
		C.float(h),
		C.float(angle),
		C.int(image.handle),
		C.int(repeat),
		C.float(alpha))}
}

// Scissor sets the current scissor rectangle.
// The scissor rectangle is transformed by the current transform.
func (self Context) Scissor(x, y, w, h float64) {
	C.nvgScissor(self.cbase, C.float(x), C.float(y), C.float(w), C.float(h))
}

// IntersectScissor intersects current scissor rectangle with the specified
// rectangle. The scissor rectangle is transformed by the current transform.
// Note: in case the rotation of previous scissor rect differs from
// the current one, the intersection will be done between the specified
// rectangle and the previous scissor rectangle transformed in the current
// transform space. The resulting shape is always rectangle.
func (self Context) IntersectScissor(x, y, w, h float64) {
	C.nvgIntersectScissor(self.cbase, C.float(x), C.float(y), C.float(w), C.float(h))
}

// Reset and disables scissoring.
func (self Context) ResetScissor() {
	C.nvgResetScissor(self.cbase)
}

// Clears the current path and sub-paths.
func (self Context) BeginPath() {
	C.nvgBeginPath(self.cbase)
}

// MoveTo starts new sub-path with specified point as first point.
func (self Context) MoveTo(x, y float64) {
	C.nvgMoveTo(self.cbase, C.float(x), C.float(y))
}

// LineTo adds line segment from the last point in the path to the
// specified point.
func (self Context) LineTo(x, y float64) {
	C.nvgLineTo(self.cbase, C.float(x), C.float(y))
}

// BezierTo adds cubic bezier segment from last point in the path via
// two control points to the specified point.
func (self Context) BezierTo(c1x, c1y, c2x, c2y, x, y float64) {
	C.nvgBezierTo(self.cbase, C.float(c1x), C.float(c1y), C.float(c2x), C.float(c2y), C.float(x), C.float(y))
}

// QuadTo adds quadratic bezier segment from last point in the path
// via a control point to the specified point.
func (self Context) QuadTo(cx, cy, x, y float64) {
	C.nvgQuadTo(self.cbase, C.float(cx), C.float(cy), C.float(x), C.float(y))
}

// ArcTo adds an arc segment at the corner defined by the last path
// point, and two specified points.
func (self Context) ArcTo(x1, y1, x2, y2, radius float64) {
	C.nvgArcTo(self.cbase, C.float(x1), C.float(y1), C.float(x2), C.float(y2), C.float(radius))
}

// ClosePath closes current sub-path with a line segment.
func (self Context) ClosePath() {
	C.nvgClosePath(self.cbase)
}

// PathWinding sets the current sub-path winding, see NVGwinding
// and NVGsolidity.
func (self Context) PathWinding(dir int) {
	C.nvgPathWinding(self.cbase, C.int(dir))
}

// Arc creates new circle arc shaped sub-path. The arc center is
// at cx,cy, the arc radius is r, and the arc is drawn from angle
// a0 to a1, and swept in direction dir (NVG_CCW, or NVG_CW).
// Angles are specified in radians.
func (self Context) Arc(cx, cy, r, a0, a1 float64, dir int) {
	C.nvgArc(self.cbase, C.float(cx), C.float(cy), C.float(r), C.float(a0), C.float(a1), C.int(dir))
}

// Rect creates new rectangle shaped sub-path.
func (self Context) Rect(x, y, w, h float64) {
	C.nvgRect(self.cbase, C.float(x), C.float(y), C.float(w), C.float(h))
}

// RoundedRect creates new rounded rectangle shaped sub-path.
func (self Context) RoundedRect(x, y, w, h, r float64) {
	C.nvgRoundedRect(self.cbase, C.float(x), C.float(y), C.float(w), C.float(h), C.float(r))
}

// Ellipse creates new ellipse shaped sub-path.
func (self Context) Ellipse(cx, cy, rx, ry float64) {
	C.nvgEllipse(self.cbase, C.float(cx), C.float(cy), C.float(rx), C.float(ry))
}

// Circle creates new circle shaped sub-path.
func (self Context) Circle(cx, cy, r float64) {
	C.nvgCircle(self.cbase, C.float(cx), C.float(cy), C.float(r))
}

// Fill the current path with current fill style.
func (self Context) Fill() {
	C.nvgFill(self.cbase)
}

// Stroke the current path with current stroke style.
func (self Context) Stroke() {
	C.nvgStroke(self.cbase)
}

// Creates font by loading it from the disk from specified file name.
// Returns handle to the font.
func (self Context) NewFont(name, filepath string) *Font {
	if _, ok := self.fonts[name]; !ok {
		n := C.CString(name)
		defer C.free(unsafe.Pointer(n))
		p := C.CString(filepath)
		defer C.free(unsafe.Pointer(p))
		self.fonts[name] = &Font{C.nvgCreateFont(self.cbase, n, p), nil}
	}

	return self.fonts[name]
}

// Creates image by loading it from the specified memory chunk, .
// Returns handle to the font.
func (self Context) NewFontMem(name string, data []byte) (*Font, error) {
	if _, ok := self.fonts[name]; !ok {
		n := C.CString(name)
		defer C.free(unsafe.Pointer(n))
		i := C.nvgCreateFontMem(self.cbase, n, (*C.uchar)(unsafe.Pointer(&data[0])), (C.int)(len(data)), 0)
		if i == -1 {
			return nil, errors.New("Font could not be created from data: " + name)
		}
		self.fonts[name] = &Font{i, data}
	}
	return self.fonts[name], nil
}

// Finds a loaded font of specified name, and returns handle to it, or -1 if the font is not found.
func (self Context) FindFont(name string) (*Font, error) {
	if _, ok := self.fonts[name]; !ok {
		n := C.CString(name)
		defer C.free(unsafe.Pointer(n))
		i := C.nvgFindFont(self.cbase, n)
		if i == -1 {
			return nil, errors.New("Font not found: " + name)
		}
		self.fonts[name] = &Font{i, nil}
	}
	return self.fonts[name], nil
}

// SetFontSize sets the font size of current text style in pixels.
func (self Context) SetFontSize(px float64) {
	C.nvgFontSize(self.cbase, C.float(px))
}

// SetFontPtSize sets the font size of current text style in points.
func (self Context) SetFontPtSize(pt float64) {
	C.nvgFontSize(self.cbase, C.float(pt*1.333333333))
}

// FontBlur sets the blur of current text style.
func (self Context) FontBlur(blur float64) {
	C.nvgFontBlur(self.cbase, C.float(blur))
}

// TextLetterSpacing sets the letter spacing of current text style.
func (self Context) TextLetterSpacing(spacing float64) {
	C.nvgTextLetterSpacing(self.cbase, C.float(spacing))
}

// TextLineHeight sets the proportional line height of current text style. The line height is specified as multiple of font size.
func (self Context) TextLineHeight(h float64) {
	C.nvgTextLineHeight(self.cbase, C.float(h))
}

// TextAlign sets the text align of current text style, see NVGaling for options.
func (self Context) TextAlign(align Align) {
	C.nvgTextAlign(self.cbase, C.int(align))
}

// SetFont sets the font to be used.
func (self Context) SetFont(f *Font) {
	C.nvgFontFaceId(self.cbase, f.cbase)
}

// SetFontFace sets the font based on given string name of current text style.
func (self Context) SetFontFace(name string) {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	C.nvgFontFace(self.cbase, n)
}

// TextMetrics returns the vertical metrics based on the current text style.
// Measured values are returned in local coordinate space.
func (self Context) TextMetrics() (ascender, descender, lineh float64) {
	var a, d, lh C.float
	C.nvgTextMetrics(self.cbase, &a, &d, &lh)
	return float64(a), float64(d), float64(lh)
}

// NVGcontext* ctx, float x, float y, const char* string, const char* end, float* bounds
func (self Context) TextBounds(text string, x, y float64) (xmin, ymin, xmax, ymax float64) {
	n := C.CString(text)
	defer C.free(unsafe.Pointer(n))
	var bounds [4]C.float
	C.nvgTextBounds(self.cbase, C.float(x), C.float(y), n, nil, &bounds[0])
	return float64(bounds[0]), float64(bounds[1]), float64(bounds[2]), float64(bounds[3])
}

// RuneBounds calculates the boundry of the specified rune.
func (self Context) RuneBounds(r rune, x, y float64) (X, minX, maxX float64) {
	n := C.CString(string(r))
	defer C.free(unsafe.Pointer(n))
	var p C.NVGglyphPosition
	C.nvgTextGlyphPositions(self.cbase, C.float(x), C.float(y), n, nil, &p, 1)
	return float64(p.x), float64(p.minx), float64(p.maxx)
}

// Text draws text string at specified location.
func (self Context) Text(x, y float64, text string) {
	t := C.CString(text)
	defer C.free(unsafe.Pointer(t))
	C.nvgText(self.cbase, C.float(x), C.float(y), t, nil)
}

func (self Context) Rune(x, y float64, r rune) {
	var t C.char = C.char(r)
	// defer C.free(unsafe.Pointer(t))
	C.nvgText(self.cbase, C.float(x), C.float(y), &t, nil)
}

// WrappedText will wrap text at the specified location at the specified with.
func (self Context) WrappedText(x, y, w float64, text string) float64 {
	t := C.CString(text)
	defer C.free(unsafe.Pointer(t))
	return float64(C.nvgTextBox(self.cbase, C.float(x), C.float(y), C.float(w), t, nil))
}

type GlyphPosition struct {
	cbase *C.NVGglyphPosition
}

type TextRow struct {
	cbase *C.NVGtextRow
}

type Font struct {
	cbase C.int

	// PATCH: Hold reference to avoid GC
	data []byte
}
