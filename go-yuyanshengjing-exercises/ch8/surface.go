package ch8

import (
	"fmt"
	"io"
	"math"
	"runtime"
	"sync"
	"time"
)

const (
	width, height = 600, 320            // 画布大小（像素）
	cells         = 100                 // 网格单元数
	xyrange       = 30.0                // 轴范围（-xyrange..+xyrange）
	xyscale       = width / 2 / xyrange // 每x或y单位像素数
	zscale        = height * 0.4        // 每z单位像素数
	angle         = math.Pi / 6         // x、y轴的角度（=30°）
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func surface(f io.Writer) {
	fmt.Fprintf(f, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	var wg sync.WaitGroup
	fmt.Println(runtime.NumCPU())       // 返回当前CPU内核数
	fmt.Println(runtime.GOMAXPROCS(2))  // 设置运行时最大可执行CPU数
	fmt.Println(runtime.NumGoroutine()) // 当前正在运行的goroutine 数
	ch := make(chan struct{}, 8)
	start := time.Now()
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ch <- struct{}{}
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				handlePolygon(f, i, j)
				<-ch
			}(i, j)
		}
	}
	go func() {
		wg.Wait()
		fmt.Println("done")
	}()
	fmt.Printf("time: %s\n", time.Now().Sub(start).String())
	fmt.Fprintf(f, "</svg>")
}

func handlePolygon(f io.Writer, i, j int) {
	ax, ay := corner(i+1, j)
	bx, by := corner(i, j)
	cx, cy := corner(i, j+1)
	dx, dy := corner(i+1, j+1)
	fmt.Fprintf(f, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
		ax, ay, bx, by, cx, cy, dx, dy)
}

func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// 将（x，y，z）等轴投影到二维SVG画布（sx，sy）上。
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}
