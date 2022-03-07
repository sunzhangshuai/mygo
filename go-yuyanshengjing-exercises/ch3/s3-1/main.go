// 如果f函数返回的是无限制的float64值，
// 那么SVG文件可能输出无效的多边形元素（虽然许多SVG渲染器会妥善处理这类问题）。
// 修改程序跳过无效的多边形。
package main

import (
	"fmt"
	"math"
)

const (
	width, height = 600, 320            // 画布大小（像素）
	cells         = 100                 // 网格单元数
	xyrange       = 30.0                // 轴范围（-xyrange..+xyrange）
	xyscale       = width / 2 / xyrange // 每x或y单位像素数
	zscale        = height * 0.4        // 每z单位像素数
	angle         = math.Pi / 6         // x、y轴角度（=30°）
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				fmt.Errorf("corner() 产生非数值")
			} else {
				fmt.Printf("<polygon points='%g,%g,%g,%g,%g,%g,%g,%g'/>\n", ax, ay, bx, by, cx, cy, dx, dy)
			}
		}
	}
	fmt.Println("</svg>")
}

// 角
func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// 将（x，y，z）等距投影到二维SVG画布（sx，sy）上。
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}