# image-convert
goのバージョンは1.19です。

こういった画像を、

![gopher.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/504414/d5719da2-d560-aaf6-b33f-fa4f907aa2ee.png)
by [Renée French](http://qiita.com).


こうする。切り抜く位置は変更可能です。正方形ではなく長方形も可能です。

![gopher2.jpg](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/504414/119d3b04-5658-cc16-7a59-fc8ff0b30b47.jpeg)

画像フォーマットも`.png`から`.jpg`に変更される。
元画像のファイル形式はpng、jpg、jpeg、gif、webpに対応しています。
変換後のファイル形式はpng、jpg、jpeg、gifに対応しています。
webp形式は、標準パッケージでは読み込めないので対応しないつもりでしたが、友人からの要望で対応することにしました。
こちら(https://pkg.go.dev/golang.org/x/image/webp) のpackageをお借りしました。
![かに.jpg](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/504414/0f139819-d7a6-7682-f075-41f41196c0a3.jpeg)

対話式CLIツールにしました。
トリミング後の画像サイズや、切り抜き位置、フォーマットを指定できます。

```bash
image-convert $ go run main.go 
トリミング後画像の幅をpxで入力してください(省略すると元画像と同じ大きさになります) >>75
トリミング後画像の高さをpxで入力してください(省略すると元画像と同じ大きさになります) >>75
画像ファイルのパスを拡張子付きで入力してください >>/[適当なパス]/gopher.png
デフォルトのトリミング範囲は画像の中央です
トリミング範囲を右に50px変更したい場は「50」、左に50px変更したい場は「-50」と入力してください(省略可) >>-30
トリミング範囲を上に50px変更したい場は「50」、下に50px変更したい場は「-50」と入力してください(省略可) >>70
出力ファイルのパスを拡張子付きで入力してください >>/[適当なパス]/gopher2.jpg
出力完了！
```

# ISSUE
- コマンドライン引数に非対応
- 円形など、長方形以外のトリミングに非対応
- エラーハンドリングが(おそらく)不十分
- 複数画像の一括変換に非対応
- 画像サイズを小さくする機能がない
