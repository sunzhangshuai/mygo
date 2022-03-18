package ch3

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
)

const (
	width, height = 1800, 960           // 画布大小（像素）
	cells         = 100                 // 网格单元数
	XYRange       = 50.0                // 轴范围（-XYRange..+XYRange）
	XYScale       = width / 2 / XYRange // 每x或y单位像素数
	zScale        = height * 0.4        // 每z单位像素数
	angle         = math.Pi / 6         // x、y轴的角度（=30°）
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

var gShape = "f"

// SurfaceHandler http 访问
func SurfaceHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if err = r.ParseForm(); err != nil {
		return
	}
	shape := r.FormValue("shape")
	w.Header().Set("Content-Type", "image/svg+xml")
	Surface(w, shape)
}

// Surface 画布
func Surface(writer io.Writer, shape string) {
	var err error

	// 生成形状
	if shape != "" {
		gShape = shape
	}

	zMin, zMax := minOrMax()

	// 获取画布的整体大小和样式

	if _, err = fmt.Fprintf(writer, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='fill: white; stroke-width: 2' "+
		"width='%d' height='%d'>\n", width, height); err != nil {
		log.Fatal(err)
		return
	}

	// 每行
	for i := 0; i < cells; i++ {
		// 每列
		for j := 0; j < cells; j++ {
			// 获取四边形的四个点
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				log.Fatal("corner() 产生非数值")
			} else if _, err = fmt.Fprintf(writer, "<polygon style='stroke: %s;' "+
				"points='%g,%g,%g,%g,%g,%g,%g,%g'/>\n",
				makeColor(i, j, zMin, zMax), ax, ay, bx, by, cx, cy, dx, dy); err != nil {
				log.Fatal(err)
				return
			}
		}
	}

	if _, err = fmt.Fprintln(writer, "</svg>"); err != nil {
		return
	}
}

// corner 获取顶点
func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := XYRange * (float64(i)/cells - 0.5)
	y := XYRange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	var z float64
	switch gShape {
	case "f":
		z = f(x, y)
	case "eggBox":
		z = eggBox(x, y)
	case "saddle":
		z = saddle(x, y)
	}

	// 将（x，y，z）等轴投影到二维SVG画布（sx，sy）上。
	sx := width/2 + (x-y)*cos30*XYScale
	sy := height/2 + (x+y)*sin30*XYScale - z*zScale
	return sx, sy
}

// f 普通算法
func f(x, y float64) float64 {
	// 获取两点斜边
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

// minOrMax 返回给定x和y的最小值/最大值并假设为方域的z的最小值和最大值。
func minOrMax() (min, max float64) {
	min = math.NaN()
	max = math.NaN()
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			for xOff := 0; xOff <= 1; xOff++ {
				for yOff := 0; yOff <= 1; yOff++ {
					x := XYRange * (float64(i+xOff)/cells - 0.5)
					y := XYRange * (float64(j+yOff)/cells - 0.5)
					z := f(x, y)
					if math.IsNaN(min) || z < min {
						min = z
					}
					if math.IsNaN(max) || z > max {
						max = z
					}
				}
			}
		}
	}
	return min, max
}

// makeColor 获取颜色
func makeColor(i, j int, zMin, zMax float64) string {
	min := math.NaN()
	max := math.NaN()
	for xOff := 0; xOff <= 1; xOff++ {
		for yOff := 0; yOff <= 1; yOff++ {
			x := XYRange * (float64(i+xOff)/cells - 0.5)
			y := XYRange * (float64(j+yOff)/cells - 0.5)
			z := f(x, y)
			if math.IsNaN(min) || z < min {
				min = z
			}
			if math.IsNaN(max) || z > max {
				max = z
			}
		}
	}

	color := ""
	// 判断水平线点上下。
	if math.Abs(max) > math.Abs(min) {
		// 渐变色
		red := math.Exp(math.Abs(max)) / math.Exp(math.Abs(zMax)) * 255
		if red > 255 {
			red = 255
		}
		color = fmt.Sprintf("#%02x0000", int(red))
	} else {
		green := math.Exp(math.Abs(min)) / math.Exp(math.Abs(zMin)) * 255
		if green > 255 {
			green = 255
		}
		color = fmt.Sprintf("#00%02x00", int(green))
	}
	return color
}

// eggBox 鸡蛋盒
func eggBox(x, y float64) float64 {
	r := 0.2 * (math.Cos(x) + math.Cos(y))
	return r
}

// saddle 马鞍
func saddle(x, y float64) float64 {
	a := 25.0
	b := 17.0
	a2 := a * a
	b2 := b * b
	r := y*y/a2 - x*x/b2
	return r
}
