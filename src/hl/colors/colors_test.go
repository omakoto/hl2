package colors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func a(c Color) *Color {
	return &c
}

func TestNewIndexColor(t *testing.T) {
	assert.Equal(t, Color{index: indexOffset + 0}, NewIndexColor(0))
	assert.Equal(t, Color{index: indexOffset + 7}, NewIndexColor(7))
	assert.Panics(t, func() { NewIndexColor(8) })
}

func TestNewRgb888Color(t *testing.T) {
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 0, b: 0}, newRgb888Color(0, 0, 0))
	assert.Equal(t, Color{index: rgbColor, r: 127, g: 0, b: 0}, newRgb888Color(127, 0, 0))
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 127, b: 0}, newRgb888Color(0, 127, 0))
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 0, b: 127}, newRgb888Color(0, 0, 127))
	assert.Equal(t, Color{index: rgbColor, r: 255, g: 0, b: 0}, newRgb888Color(255, 0, 0))
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 255, b: 0}, newRgb888Color(0, 255, 0))
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 0, b: 255}, newRgb888Color(0, 0, 255))
}

func TestNewRgb216Color(t *testing.T) {
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 0, b: 0}, newRgb216Color(0, 0, 0))
	assert.Equal(t, Color{index: rgbColor, r: 153, g: 0, b: 0}, newRgb216Color(3, 0, 0))
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 153, b: 0}, newRgb216Color(0, 3, 0))
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 0, b: 153}, newRgb216Color(0, 0, 3))
	assert.Equal(t, Color{index: rgbColor, r: 255, g: 0, b: 0}, newRgb216Color(5, 0, 0))
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 255, b: 0}, newRgb216Color(0, 5, 0))
	assert.Equal(t, Color{index: rgbColor, r: 0, g: 0, b: 255}, newRgb216Color(0, 0, 5))

	assert.Panics(t, func() { newRgb216Color(6, 0, 0) })
}

func TestColor_IsNone(t *testing.T) {
	assert.Equal(t, true, a(NoColor).IsNone())

	assert.Equal(t, false, a(NewIndexColor(0)).IsNone())
	assert.Equal(t, false, a(NewIndexColor(7)).IsNone())

	assert.Equal(t, false, a(newRgb216Color(0, 0, 0)).IsNone())
	assert.Equal(t, false, a(newRgb888Color(0, 0, 0)).IsNone())
}

func TestColor_IsIndex(t *testing.T) {
	assert.Equal(t, false, a(NoColor).IsIndex())

	assert.Equal(t, true, a(NewIndexColor(0)).IsIndex())
	assert.Equal(t, true, a(NewIndexColor(7)).IsIndex())

	assert.Equal(t, false, a(newRgb216Color(0, 0, 0)).IsIndex())
	assert.Equal(t, false, a(newRgb888Color(0, 0, 0)).IsIndex())
}

func TestColor_IsRgb(t *testing.T) {
	assert.Equal(t, false, a(NoColor).IsRgb())

	assert.Equal(t, false, a(NewIndexColor(0)).IsRgb())
	assert.Equal(t, false, a(NewIndexColor(7)).IsRgb())

	assert.Equal(t, true, a(newRgb216Color(0, 0, 0)).IsRgb())
	assert.Equal(t, true, a(newRgb888Color(0, 0, 0)).IsRgb())
}

func TestColor_GetIndex(t *testing.T) {
	assert.Panics(t, func() { a(NoColor).Index() })

	assert.Equal(t, uint8(0), a(NewIndexColor(0)).Index())
	assert.Equal(t, uint8(7), a(NewIndexColor(7)).Index())

	assert.Panics(t, func() { a(newRgb216Color(0, 0, 0)).Index() })
	assert.Panics(t, func() { a(newRgb888Color(0, 0, 0)).Index() })
}

func TestColor_R(t *testing.T) {
	assert.Panics(t, func() { a(NoColor).R() })
	assert.Panics(t, func() { a(NewIndexColor(0)).R() })

	assert.Equal(t, uint8(0), a(newRgb216Color(0, 0, 0)).R())
	assert.Equal(t, uint8(255), a(newRgb216Color(5, 0, 0)).R())
	assert.Equal(t, uint8(0), a(newRgb888Color(0, 0, 0)).R())
	assert.Equal(t, uint8(255), a(newRgb888Color(255, 0, 0)).R())
}

func TestColor_G(t *testing.T) {
	assert.Panics(t, func() { a(NoColor).G() })
	assert.Panics(t, func() { a(NewIndexColor(0)).G() })

	assert.Equal(t, uint8(0), a(newRgb216Color(0, 0, 0)).G())
	assert.Equal(t, uint8(255), a(newRgb216Color(0, 5, 0)).G())
	assert.Equal(t, uint8(0), a(newRgb888Color(0, 0, 0)).G())
	assert.Equal(t, uint8(255), a(newRgb888Color(0, 255, 0)).G())
}

func TestColor_B(t *testing.T) {
	assert.Panics(t, func() { a(NoColor).B() })
	assert.Panics(t, func() { a(NewIndexColor(0)).B() })

	assert.Equal(t, uint8(0), a(newRgb216Color(0, 0, 0)).B())
	assert.Equal(t, uint8(255), a(newRgb216Color(0, 0, 5)).B())
	assert.Equal(t, uint8(0), a(newRgb888Color(0, 0, 0)).B())
	assert.Equal(t, uint8(255), a(newRgb888Color(0, 0, 255)).B())
}
