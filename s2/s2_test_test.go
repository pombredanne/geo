package s2

import (
	"reflect"
	"testing"

	"github.com/golang/geo/r1"
	"github.com/golang/geo/s1"
)

func TestKmToAngle(t *testing.T) {
	const earthRadiusKm = 6371.01

	tests := []struct {
		have float64
		want s1.Angle
	}{
		{0.0, 0.0},
		{1.0, 0.00015696098420815537 * s1.Radian},
		{earthRadiusKm, 1.0 * s1.Radian},
		{-1.0, -0.00015696098420815537 * s1.Radian},
		{-10000.0, -1.5696098420815536300 * s1.Radian},
		{1e9, 156960.984208155363007 * s1.Radian},
	}
	for _, test := range tests {
		if got := kmToAngle(test.have); !float64Eq(float64(got), float64(test.want)) {
			t.Errorf("kmToAngle(%f) = %0.20f, want %0.20f", test.have, got, test.want)
		}
	}

}

func TestParsePoint(t *testing.T) {
	tests := []struct {
		have string
		want Point
	}{
		{"0:0", PointFromCoords(1, 0, 0)},
		{"90:0", PointFromCoords(6.123233995736757e-17, 0, 1)},
		{"91:0", PointFromCoords(-0.017452406437283473, -0, 0.9998476951563913)},
		{"179.99:0", PointFromCoords(-0.9999999847691292, -0, 0.00017453292431344843)},
		{"180:0", PointFromCoords(-1, -0, 1.2246467991473515e-16)},
		{"181.0:0", PointFromCoords(-0.9998476951563913, -0, -0.017452406437283637)},
		{"-45:0", PointFromCoords(0.7071067811865476, 0, -0.7071067811865475)},
		{"0:0.01", PointFromCoords(0.9999999847691292, 0.00017453292431333684, 0)},
		{"0:30", PointFromCoords(0.8660254037844387, 0.49999999999999994, 0)},
		{"0:45", PointFromCoords(0.7071067811865476, 0.7071067811865475, 0)},
		{"0:90", PointFromCoords(6.123233995736757e-17, 1, 0)},
		{"30:30", PointFromCoords(0.7500000000000001, 0.4330127018922193, 0.49999999999999994)},
		{"-30:30", PointFromCoords(0.7500000000000001, 0.4330127018922193, -0.49999999999999994)},
		{"180:90", PointFromCoords(-6.123233995736757e-17, -1, 1.2246467991473515e-16)},
		{"37.4210:-122.0866, 37.4231:-122.0819", PointFromCoords(-0.4218751185559026, -0.6728760966593905, 0.6076669670863027)},
	}
	for _, test := range tests {
		if got := parsePoint(test.have); got != test.want {
			t.Errorf("parsePoint(%s) = %v, want %v", test.have, got, test.want)
		}
	}
}

func TestParseRect(t *testing.T) {
	tests := []struct {
		have string
		want Rect
	}{
		{"0:0", Rect{}},
		{
			"1:1",
			Rect{
				r1.Interval{float64(s1.Degree), float64(s1.Degree)},
				s1.Interval{float64(s1.Degree), float64(s1.Degree)},
			},
		},
		{
			"1:1, 2:2, 3:3",
			Rect{
				r1.Interval{float64(s1.Degree), 3 * float64(s1.Degree)},
				s1.Interval{float64(s1.Degree), 3 * float64(s1.Degree)},
			},
		},
		{
			"-90:-180, 90:180",
			Rect{
				r1.Interval{-90 * float64(s1.Degree), 90 * float64(s1.Degree)},
				s1.Interval{180 * float64(s1.Degree), -180 * float64(s1.Degree)},
			},
		},
		{
			"-89.99:0, 89.99:179.99",
			Rect{
				r1.Interval{-89.99 * float64(s1.Degree), 89.99 * float64(s1.Degree)},
				s1.Interval{0, 179.99 * float64(s1.Degree)},
			},
		},
		{
			"-89.99:-179.99, 89.99:179.99",
			Rect{
				r1.Interval{-89.99 * float64(s1.Degree), 89.99 * float64(s1.Degree)},
				s1.Interval{179.99 * float64(s1.Degree), -179.99 * float64(s1.Degree)},
			},
		},
		{
			"37.4210:-122.0866, 37.4231:-122.0819",
			Rect{
				r1.Interval{float64(s1.Degree * 37.4210), float64(s1.Degree * 37.4231)},
				s1.Interval{float64(s1.Degree * -122.0866), float64(s1.Degree * -122.0819)},
			},
		},
		{
			"-876.54:-654.43, 963.84:2468.35",
			Rect{
				r1.Interval{-876.54 * float64(s1.Degree), -876.54 * float64(s1.Degree)},
				s1.Interval{-654.43 * float64(s1.Degree), -654.43 * float64(s1.Degree)},
			},
		},
	}
	for _, test := range tests {
		if got := parseRect(test.have); got != test.want {
			t.Errorf("parseRect(%s) = %v, want %v", test.have, got, test.want)
		}
	}
}

func TestParseLatLngs(t *testing.T) {
	tests := []struct {
		have string
		want []LatLng
	}{
		{"0:0", []LatLng{LatLng{0, 0}}},
		{
			"37.4210:-122.0866, 37.4231:-122.0819",
			[]LatLng{
				LatLng{s1.Degree * 37.4210, s1.Degree * -122.0866},
				LatLng{s1.Degree * 37.4231, s1.Degree * -122.0819},
			},
		},
	}
	for _, test := range tests {
		got := parseLatLngs(test.have)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("parseLatLngs(%s) = %v, want %v", test.have, got, test.want)
		}
	}
}

func TestParsePoints(t *testing.T) {
	tests := []struct {
		have string
		want []Point
	}{
		{"0:0", []Point{PointFromCoords(1, 0, 0)}},
		{"      0:0,    ", []Point{PointFromCoords(1, 0, 0)}},
		{
			"90:0,-90:0",
			[]Point{
				PointFromCoords(6.123233995736757e-17, 0, 1),
				PointFromCoords(6.123233995736757e-17, 0, -1),
			},
		},
		{
			"90:0, 0:90, -90:0, 0:-90",
			[]Point{
				PointFromCoords(6.123233995736757e-17, 0, 1),
				PointFromCoords(6.123233995736757e-17, 1, 0),
				PointFromCoords(6.123233995736757e-17, 0, -1),
				PointFromCoords(6.123233995736757e-17, -1, 0),
			},
		},
	}

	for _, test := range tests {
		got := parsePoints(test.have)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("parsePoints(%s), got %v, want %v",
				test.have, got, test.want)
		}
	}
}