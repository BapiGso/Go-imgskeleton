package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"strings"
)

func min(l []int) (min int) {
	min = l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return
}

func dealarr(bianjie *image.RGBA) [][]uint8 {
	bounds := bianjie.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	yslice := [][]uint8{}
	for y := 0; y < dy; y++ {
		xslice := []uint8{}
		for x := 0; x < dx; x++ {
			colorRgb := bianjie.At(x, y)
			_, g, _, _ := colorRgb.RGBA()
			//黑色边界
			if (g >> 8) == 0 {
				xslice = append(xslice, 0)
			}
			//灰色内部
			if (g >> 8) == 128 {
				xslice = append(xslice, 1)
			}
			//白色外部
			if (g >> 8) == 255 {
				xslice = append(xslice, 2)
			}
			//fmt.Printf("%03d", newG)
			//fmt.Print(" ")
		}
		yslice = append(yslice, xslice)
		//fmt.Printf("%v\n", xslice)
	}
	return yslice
}

func mindis(bianjie *image.RGBA, yslice [][]uint8, shendu int) *image.RGBA {
	bounds := bianjie.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	for y := 1; y < dy-1; y++ {
		for x := 1; x < dx-1; x++ {
			//colorRgb := bianjie.At(x, y)
			//_, g, _, _ := colorRgb.RGBA()
			//染成黑色，不过为了不混淆先染成1
			newG := uint8(1)
			//如果是白色外部
			if yslice[y][x] == 2 {
				//newG = uint8(255)
				//也先染成黑色好看
				newG = uint8(1)
			}
			//如果是灰色内部点
			if yslice[y][x] == 1 {
				//也染黑色在加
				newG = uint8(1)
				//算该点到最近的一个边界（不是指边界点）长度，防止溢出
				l := []int{x, y, dx - x, dy - y}
				minl := min(l)
				for i := 0; i < minl; i++ {
					var up, down, right, left, ur, ul, dr, dl = true, true, true, true, true, true, true, true
					//如果往周围走碰到了黑色边界
					if yslice[y+i][x] == 0 {
						up = false
					}
					if yslice[y-i][x] == 0 {
						down = false
					}
					if yslice[y][x+i] == 0 {
						right = false
					}
					if yslice[y][x-i] == 0 {
						left = false
					}
					if yslice[y+i][x+i] == 0 {
						ur = false
					}
					if yslice[y+i][x-i] == 0 {
						ul = false
					}
					if yslice[y-i][x+i] == 0 {
						dr = false
					}
					if yslice[y-i][x-i] == 0 {
						dl = false
					}
					if !up || !down || !left || !right || !ur || !ul || !dr || !dl {
						//tmpslice = append(tmpslice, "b")
						newG = newG + uint8(i*shendu/minl)
						break
					}
				}
			}

			newRgba.SetRGBA(x, y, color.RGBA{R: newG, G: newG, B: newG, A: 255})
		}
	}
	return newRgba
}

func setNeighbourhood(arraylist [][]uint8, erzhi *image.RGBA, dy int, dx int) *image.RGBA {
	bounds := erzhi.Bounds()
	//newslice := [][]string{}
	newRgba := image.NewRGBA(bounds)
	for y := 1; y < dy-1; y++ {
		//tmpslice := []string{}
		for x := 1; x < dx-1; x++ {
			var up, down, right, left = true, true, true, true

			if arraylist[x-1][y] == 0 {
				up = false
			}
			if arraylist[x+1][y] == 0 {
				down = false
			}
			if arraylist[x][y-1] == 0 {
				right = false
			}
			if arraylist[x][y+1] == 0 {
				left = false
			}
			//边界染黑色
			newG := uint8(0)
			//外部染白色
			if !up && !down && !left && !right {
				//tmpslice = append(tmpslice, "b")
				newG = uint8(255)
			}
			//内部染灰色
			if up && down && left && right {
				//tmpslice = append(tmpslice, "i")
				newG = uint8(128)
			}
			newRgba.SetRGBA(y, x, color.RGBA{R: newG, G: newG, B: newG, A: 255})
		}
		//newslice = append(newslice, tmpslice)
		//fmt.Printf("%v\n", tmpslice)
	}
	return newRgba
	//fmt.Printf("%v\n", newslice)
}

func erzhihua(m image.Image, yuzhi uint32) (*image.RGBA, [][]uint8, int, int) { //灰度化图像。
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	yslice := [][]uint8{}
	for y := 0; y < dy; y++ {
		xslice := []uint8{}
		for x := 0; x < dx; x++ {
			colorRgb := m.At(x, y)
			_, g, _, _ := colorRgb.RGBA()
			newG := uint8(0)
			//阈值范围0~65535
			if (g >> 8) > yuzhi {
				newG = uint8(255)
				xslice = append(xslice, 0)
			} else {
				xslice = append(xslice, 1)
			}
			//fmt.Printf("%03d", newG)
			//fmt.Print(" ")
			newRgba.SetRGBA(x, y, color.RGBA{R: newG, G: newG, B: newG, A: 255})
		}
		yslice = append(yslice, xslice)
		//fmt.Printf("%v\n", xslice)
	}
	return newRgba, yslice, dy, dx
}

func grayingImage(m image.Image) *image.RGBA { //灰度化图像。
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			colorRgb := m.At(x, y)
			_, g, _, _ := colorRgb.RGBA()
			newG := uint8(g >> 8)
			newRgba.SetRGBA(x, y, color.RGBA{R: newG, G: newG, B: newG, A: 255})
		}
		//fmt.Printf("%v\n", xslice)
	}
	return newRgba
}

func imageEncode(fileName string, file *os.File, rgba *image.RGBA) error {
	// 将图片和扩展名分离
	stringSlice := strings.Split(fileName, ".")
	// 根据图片的扩展名来运用不同的处理
	switch stringSlice[len(stringSlice)-1] {
	case "jpg":
		return jpeg.Encode(file, rgba, nil)
	case "jpeg":
		return jpeg.Encode(file, rgba, nil)
	case "gif":
		return gif.Encode(file, rgba, nil)
	case "png":
		return png.Encode(file, rgba)
	default:
		panic("不支持的图片类型")
	}
}
func outimg(name string, img *image.RGBA) {
	outFile1, _ := os.Create(name)
	defer outFile1.Close()
	if err := imageEncode(name, outFile1, img); err != nil {
		panic(err)
	}
}

func main() {
	var imagePath string
	var yuzhi uint32
	var shendu int
	fmt.Println("推荐使用白底图片")
	fmt.Println("请输入图片路径:")
	fmt.Scanln(&imagePath)
	fmt.Println("请输入二值化灰度阈值(0-255):")
	fmt.Scanln(&yuzhi)
	fmt.Println("请输入骨架染色深度(我自己弄的一个参数，推荐200±50):")
	fmt.Scanln(&shendu)

	file, _ := os.Open(imagePath)
	defer file.Close() //这个是方式防止忘记关掉。内存溢出。
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	//灰度
	graychange := grayingImage(img)
	//二值化 数组，长宽
	erzhi, arraylist, dx, dy := erzhihua(img, yuzhi)
	//划分边界点内部点
	bianjie := setNeighbourhood(arraylist, erzhi, dy, dx)
	//先变成数组好计算
	slice123 := dealarr(bianjie)
	//欧式距离变换
	oushi := mindis(bianjie, slice123, shendu)
	outimg("欧式.png", oushi)
	outimg("灰度.png", graychange)
	outimg("二值.png", erzhi)
	outimg("边界.png", bianjie)

}
