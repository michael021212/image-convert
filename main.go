package main

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/webp"
)

// 後述するSubImageメソッドによるトリミングのために用意
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// 出力画像の幅と高さをpxで指定
	fmt.Print("トリミング後画像の幅をpxで入力してください(省略すると元画像と同じ大きさになります) >>")
	scanner.Scan()
	OUTPUT_WIDTH, _ := strconv.Atoi(scanner.Text())
	fmt.Print("トリミング後画像の高さをpxで入力してください(省略すると元画像と同じ大きさになります) >>")
	scanner.Scan()
	OUTPUT_HEIGHT, _ := strconv.Atoi(scanner.Text())

	// 入力画像パス
	fmt.Print("画像ファイルのパスを拡張子付きで入力してください >>")
	scanner.Scan()
	iPath := scanner.Text()
	inputPath := strings.Replace(iPath, `"`, "", -1) // windowsでファイルのパスをコピーするとデフォルトでダブルクォーテーションがつくので除去するように

	// 画像ファイルを開く。戻り値は*os.File型とerror型
	inputFile, err := os.Open(inputPath)
	assert(err, "パスが不正です '"+inputPath+"'")
	defer inputFile.Close() // リソース破棄のためにdeferを使ってファイルをクローズする。deferを使うことで関数終了時に実行される

	// image.Decodeでファイルオブジェクトを画像オブジェクトに変換。戻り値はImage型とstring型とerror型
	// Image型は、形式によって*image.YCbCr型だったり*image.NRGBA型だったりする
	// JPEGの場合はRGBでなくYCbCrで保存し、PNGの場合はNRGBAかRGBAで保存するらしいがこの辺よくわかっていません
	// string型は、"png"や"jpg"だったりするが、今回は不要なのでアンダースコア変数で無視する
	var inputImg image.Image // if文の中で変数の代入ができないので先に宣言する
	if strings.ToLower(filepath.Ext(inputPath)) == ".webp" {
		inputImg, err = webp.Decode(inputFile)
	} else {
		inputImg, _, err = image.Decode(inputFile)
	}
	assert(err, "画像データとして読み込めません")

	// Bounds()の戻り値はRectangle型。Dx()はRectangle型のクラスメソッドで幅を取得する
	inputWidth := inputImg.Bounds().Dx()
	inputHeight := inputImg.Bounds().Dy()

	fmt.Println("デフォルトのトリミング位置は画像の中央です")
	fmt.Print("トリミング位置を右に50px変更したい場は「50」、左に50px変更したい場は「-50」と入力してください(省略可) >>")
	scanner.Scan()
	xSlide, _ := strconv.Atoi(scanner.Text())
	fmt.Print("トリミング位置を上に50px変更したい場は「50」、下に50px変更したい場は「-50」と入力してください(省略可) >>")
	scanner.Scan()
	ySlide, _ := strconv.Atoi(scanner.Text())

	// 切り取りに使う座標を計算する処理
	// if OUTPUT_(WIDTH/HEIGHT) == 0 の条件分岐は、トリミング後画像の幅が未指定の場合は元画像と同じ大きさで出力するための記述
	x1, x2, y1, y2 := 0, 0, 0, 0
	if xSlide != 0 {
		if OUTPUT_WIDTH == 0 {
			x1 = 0 + xSlide
		} else {
			x1 = ((inputWidth - OUTPUT_WIDTH) / 2) + xSlide
		}
	} else {
		if OUTPUT_WIDTH == 0 {
			x1 = 0
		} else {
			x1 = (inputWidth - OUTPUT_WIDTH) / 2
		}
	}

	if OUTPUT_WIDTH == 0 {
		x2 = inputWidth
	} else {
		x2 = x1 + OUTPUT_WIDTH
	}

	if ySlide != 0 {
		if OUTPUT_HEIGHT == 0 {
			y1 = 0 + ySlide
		} else {
			y1 = ((inputHeight - OUTPUT_HEIGHT) / 2) + ySlide
		}
	} else {
		if OUTPUT_HEIGHT == 0 {
			y1 = 0
		} else {
			y1 = (inputHeight - OUTPUT_HEIGHT) / 2
		}
	}

	if OUTPUT_HEIGHT == 0 {
		y2 = inputHeight
	} else {
		y2 = y1 + OUTPUT_HEIGHT
	}

	// Imageインターフェース(image.Image)にはSubImageメソッドがないので、
	// SubImagerインターフェースを作って型アサーションすることでSubImageメソッドを使えるようにしている
	// ※よくわかってない
	trimmedImg := inputImg.(SubImager).SubImage(image.Rect(x1, y1, x2, y2))

	// 出力画像パス
	fmt.Print("出力ファイルのパスを拡張子付きで入力してください >>")
	scanner.Scan()
	oPath := scanner.Text()
	outputPath := strings.Replace(oPath, `"`, "", -1) // windowsでファイルのパスをコピーするとデフォルトでダブルクォーテーションがつくので除去するように

	// 出力ファイルを生成
	outputFile, err := os.Create(outputPath)
	assert(err, "ファイルの作成に失敗しました '"+outputPath+"'")
	defer outputFile.Close()

	// 出力ファイルパスの拡張子によってエンコードする形式を変える
	switch strings.ToLower(filepath.Ext(outputPath)) {
	case ".jpeg", ".jpg":
		jpeg.Encode(outputFile, trimmedImg, nil)
		fmt.Println("出力完了！")
	case ".png":
		png.Encode(outputFile, trimmedImg)
		fmt.Println("出力完了！")
	case ".gif":
		gif.Encode(outputFile, trimmedImg, nil)
		fmt.Println("出力完了！")
	default:
		fmt.Println("拡張子が対応していません (対応拡張子: jpeg/jpg/png/gif)")
		err := os.Remove(outputPath)
		assert(err, "出力ファイルは不完全です")
	}
}

// errorがあれば例外を出す
func assert(err error, msg string) {
	if err != nil {
		panic(err.Error() + ":" + msg)
	}
}
