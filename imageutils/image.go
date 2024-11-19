package imageutils

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"

	// "image/jpeg"
	// "filepath"
	// "image/gif"  // 导入 gif 解码器
	"image/jpeg" // 导入 jpeg 解码器
	"image/png"  // 导入 png 解码器
	"math"
	"os"
	"sort"
)

const u = 3.62 // 混沌系统常数

// produceLogisticArray 生成混沌序列
func produceLogisticArray(x float64, n int) []float64 {
	arr := make([]float64, n)
	arr[0] = x
	for i := 1; i < n; i++ {
		arr[i] = u * arr[i-1] * (1 - arr[i-1])
	}
	return arr
}

// getPos 修正版
func getPos(array []float64, npMap map[float64]int, pnMap map[int]float64) {
	indices := make([]int, len(array))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return array[indices[i]] < array[indices[j]]
	})

	for i, idx := range indices {
		npMap[array[idx]] = i
		pnMap[i] = array[idx]
	}
}

func dePos(array []float64, nnmap map[float64]int, ppmap map[float64]float64) {
	indices := make([]int, len(array))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return array[indices[i]] < array[indices[j]]
	})
	for i, idx := range indices {
		nnmap[array[i]] = i
		ppmap[array[i]] = array[idx]
	}
}

// decrypt 修正版
func decrypt(pixel [][]color.Color, x1 float64) [][]color.Color {
	M := len(pixel)
	N := len(pixel[0])
	ans := make([][]color.Color, M)
	for i := range ans {
		ans[i] = make([]color.Color, N)
	}

	// 生成m个数
	array := produceLogisticArray(x1, M)
	nnMap := map[float64]int{}     // 值:位置
	ppMap := map[float64]float64{} //位置:值
	dePos(array, nnMap, ppMap)

	// 恢复行的位置
	for i := 0; i < M; i++ {
		// originalRow := invNpMap[i]
		row := nnMap[ppMap[array[i]]]
		for j := 0; j < N; j++ {
			ans[i][j] = pixel[row][j]
		}
	}
	fmt.Println("行恢复成功")
	return ans
}

// encrypt 图像置乱
func encrypt(pixel [][]color.Color, x1 float64) [][]color.Color {
	M := len(pixel)
	N := len(pixel[0])
	ans := [][]color.Color{}
	// 生成m个数
	array := produceLogisticArray(x1, M)
	npMap := map[float64]int{} // 值:位置
	temp := map[float64]int{}
	pnMap := map[int]float64{} //位置:值
	for i, v := range array {
		npMap[v] = i
		temp[v] = i
		pnMap[i] = v
	}

	getPos(array, npMap, pnMap) // 对混沌序列进行排序,按排序后的顺序交换每一行的元素
	fmt.Println("位置序列生成成功")
	// for i := 0; i < M; i++ {
	// 	x1 = rowEncrypt(pixel[i], x1)
	// }
	for i := 0; i < M; i++ {
		row := npMap[array[i]]
		rps := pixel[row]
		ans = append(ans, make([]color.Color, N))
		for j, v := range rps {
			ans[i][j] = v
		}

	}
	fmt.Println("行置乱成功")
	return ans
}

// 将一个 float64 数值转换为小于 1 的数
func convertToLessThanOne(num float64) float64 {
	if num < 1 {
		return num
	}
	// 计算 num 的位数
	digits := math.Floor(math.Log10(math.Abs(num))) + 1
	// 计算需要除以的因子
	factor := math.Pow(10, digits)
	// 转换数值
	return num / factor
}

func ProcessImage(filePath string, passwd float64, isEncrypt bool, outpath string) error {
	fmt.Println(filePath, passwd, isEncrypt, outpath)
	passwd = convertToLessThanOne(passwd)
	file, err := os.Open(filePath)
	ext := filepath.Ext(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 重置文件指针到文件开头
	_, err = file.Seek(0, 0)

	fmt.Println("打开文件成功")
	fmt.Println("图片格式为:", ext)
	var img image.Image
	if ext == ".jpeg" || ext == ".jpg" {
		fmt.Println("图片格式为:", ext)
		img, err = jpeg.Decode(file)
		if err != nil {
			return err
		}
	} else if ext == ".png" {
		img, err = png.Decode(file)
		if err != nil {
			return err
		}
	}
	// img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	fmt.Println("图片加载成功")
	width, height := bounds.Max.X, bounds.Max.Y
	pixels := make([][]color.Color, height)
	for i := range pixels {
		pixels[i] = make([]color.Color, width)
		for j := 0; j < width; j++ {
			pixels[i][j] = img.At(j, i)
		}
	}

	x1 := passwd // 初始混沌值
	ans := [][]color.Color{}

	if isEncrypt {
		fmt.Println("开始加密")
		// outpath = "output.jpeg"
		ans = encrypt(pixels, x1)
	} else {
		fmt.Println("开始解密")
		// outpath = "en_output.jpeg"
		ans = decrypt(pixels, x1)
	}

	// 创建新图像
	newImg := image.NewRGBA(bounds)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			newImg.Set(j, i, ans[i][j])
		}
	}

	// 保存图像
	outputFile, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	if ext == ".jpeg" || ext == ".jpg" {
		err = jpeg.Encode(outputFile, newImg, nil)
		if err != nil {
			return err
		}
	} else if ext == ".png" {
		err = png.Encode(outputFile, newImg)
		if err != nil {
			return err
		}
		fmt.Println("png over")
	}
	// err = jpeg.Encode(outputFile, newImg, nil)
	// if err != nil {
	// 	return err
	// }

	return nil
}
