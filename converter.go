/*=================================
* Copyright(c)2015-2016 gostores
* From: github.com/nfnt/resize
*=================================*/
package resize

import "image"

// Keep value in [0,255] range.
func clampUint8(in int32) uint8 {
	// casting a negative int to an uint will result in an overflown
	// large uint. this behavior will be exploited here and in other functions
	// to achieve a higher performance.
	if uint32(in) < 256 {
		return uint8(in)
	}
	if in > 255 {
		return 255
	}
	return 0
}

// Keep value in [0,65535] range.
func clampUint16(in int64) uint16 {
	if uint64(in) < 65536 {
		return uint16(in)
	}
	if in > 65535 {
		return 65535
	}
	return 0
}

func resizeGeneric(in image.Image, out *image.NRGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					r, g, b, a := in.At(xi+in.Bounds().Min.X, x+in.Bounds().Min.Y).RGBA()

					// reverse alpha-premultiplication.
					if a != 0 {
						r *= 0xffff
						r /= a
						g *= 0xffff
						g /= a
						b *= 0xffff
						b /= a
					}

					rgba[0] += int64(coeff) * int64(r)
					rgba[1] += int64(coeff) * int64(g)
					rgba[2] += int64(coeff) * int64(b)
					rgba[3] += int64(coeff) * int64(a)
					sum += int64(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[offset+2] = uint8(value >> 8)
			out.Pix[offset+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[offset+4] = uint8(value >> 8)
			out.Pix[offset+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[offset+6] = uint8(value >> 8)
			out.Pix[offset+7] = uint8(value)
		}
	}
}

func resizeRGBA(in *image.RGBA, out *image.NRGBA, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(maxX):
						xi *= 4
					case xi >= maxX:
						xi = 4 * maxX
					default:
						xi = 0
					}

					r := uint32(row[xi+0])
					g := uint32(row[xi+1])
					b := uint32(row[xi+2])
					a := uint32(row[xi+3])

					// reverse alpha-premultiplication.
					if a != 0 {
						r *= 0xff
						r /= a
						g *= 0xff
						g /= a
						b *= 0xff
						b /= a
					}

					rgba[0] += int32(coeff) * int32(r)
					rgba[1] += int32(coeff) * int32(g)
					rgba[2] += int32(coeff) * int32(b)
					rgba[3] += int32(coeff) * int32(a)
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*4
			out.Pix[xo+0] = clampUint8(rgba[0] / sum)
			out.Pix[xo+1] = clampUint8(rgba[1] / sum)
			out.Pix[xo+2] = clampUint8(rgba[2] / sum)
			out.Pix[xo+3] = clampUint8(rgba[3] / sum)
		}
	}
}

func resizeNRGBA(in *image.NRGBA, out *image.NRGBA, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(maxX):
						xi *= 4
					case xi >= maxX:
						xi = 4 * maxX
					default:
						xi = 0
					}
					rgba[0] += int32(coeff) * int32(row[xi+0])
					rgba[1] += int32(coeff) * int32(row[xi+1])
					rgba[2] += int32(coeff) * int32(row[xi+2])
					rgba[3] += int32(coeff) * int32(row[xi+3])
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*4
			out.Pix[xo+0] = clampUint8(rgba[0] / sum)
			out.Pix[xo+1] = clampUint8(rgba[1] / sum)
			out.Pix[xo+2] = clampUint8(rgba[2] / sum)
			out.Pix[xo+3] = clampUint8(rgba[3] / sum)
		}
	}
}

func resizeRGBA64(in *image.RGBA64, out *image.NRGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(maxX):
						xi *= 8
					case xi >= maxX:
						xi = 8 * maxX
					default:
						xi = 0
					}

					r := uint32(uint16(row[xi+0])<<8 | uint16(row[xi+1]))
					g := uint32(uint16(row[xi+2])<<8 | uint16(row[xi+3]))
					b := uint32(uint16(row[xi+4])<<8 | uint16(row[xi+5]))
					a := uint32(uint16(row[xi+6])<<8 | uint16(row[xi+7]))

					// reverse alpha-premultiplication.
					if a != 0 {
						r *= 0xffff
						r /= a
						g *= 0xffff
						g /= a
						b *= 0xffff
						b /= a
					}

					rgba[0] += int64(coeff) * int64(r)
					rgba[1] += int64(coeff) * int64(g)
					rgba[2] += int64(coeff) * int64(b)
					rgba[3] += int64(coeff) * int64(a)
					sum += int64(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[xo+0] = uint8(value >> 8)
			out.Pix[xo+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[xo+2] = uint8(value >> 8)
			out.Pix[xo+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[xo+4] = uint8(value >> 8)
			out.Pix[xo+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[xo+6] = uint8(value >> 8)
			out.Pix[xo+7] = uint8(value)
		}
	}
}

func resizeNRGBA64(in *image.NRGBA64, out *image.NRGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(maxX):
						xi *= 8
					case xi >= maxX:
						xi = 8 * maxX
					default:
						xi = 0
					}
					rgba[0] += int64(coeff) * int64(uint16(row[xi+0])<<8|uint16(row[xi+1]))
					rgba[1] += int64(coeff) * int64(uint16(row[xi+2])<<8|uint16(row[xi+3]))
					rgba[2] += int64(coeff) * int64(uint16(row[xi+4])<<8|uint16(row[xi+5]))
					rgba[3] += int64(coeff) * int64(uint16(row[xi+6])<<8|uint16(row[xi+7]))
					sum += int64(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[xo+0] = uint8(value >> 8)
			out.Pix[xo+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[xo+2] = uint8(value >> 8)
			out.Pix[xo+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[xo+4] = uint8(value >> 8)
			out.Pix[xo+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[xo+6] = uint8(value >> 8)
			out.Pix[xo+7] = uint8(value)
		}
	}
}

func resizeGray(in *image.Gray, out *image.Gray, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-newBounds.Min.X)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					gray += int32(coeff) * int32(row[xi])
					sum += int32(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x - newBounds.Min.X)
			out.Pix[offset] = clampUint8(gray / sum)
		}
	}
}

func resizeGray16(in *image.Gray16, out *image.Gray16, scale float64, coeffs []int32, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray int64
			var sum int64
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(maxX):
						xi *= 2
					case xi >= maxX:
						xi = 2 * maxX
					default:
						xi = 0
					}
					gray += int64(coeff) * int64(uint16(row[xi+0])<<8|uint16(row[xi+1]))
					sum += int64(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*2
			value := clampUint16(gray / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
		}
	}
}

func resizeYCbCr(in *ycc, out *ycc, scale float64, coeffs []int16, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var p [3]int32
			var sum int32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(maxX):
						xi *= 3
					case xi >= maxX:
						xi = 3 * maxX
					default:
						xi = 0
					}
					p[0] += int32(coeff) * int32(row[xi+0])
					p[1] += int32(coeff) * int32(row[xi+1])
					p[2] += int32(coeff) * int32(row[xi+2])
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*3
			out.Pix[xo+0] = clampUint8(p[0] / sum)
			out.Pix[xo+1] = clampUint8(p[1] / sum)
			out.Pix[xo+2] = clampUint8(p[2] / sum)
		}
	}
}

func nearestYCbCr(in *ycc, out *ycc, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var p [3]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case uint(xi) < uint(maxX):
						xi *= 3
					case xi >= maxX:
						xi = 3 * maxX
					default:
						xi = 0
					}
					p[0] += float32(row[xi+0])
					p[1] += float32(row[xi+1])
					p[2] += float32(row[xi+2])
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*3
			out.Pix[xo+0] = floatToUint8(p[0] / sum)
			out.Pix[xo+1] = floatToUint8(p[1] / sum)
			out.Pix[xo+2] = floatToUint8(p[2] / sum)
		}
	}
}
