# exam-digest

exam-viewer用の集計プログラム

特定ディレクトリ配下にある複数のテスト結果を纏めたファイルを生成する

## Requirements

- git 2.34.1
- go 1.19

## Build

```bash
git clone https://github.com/okapon1210/exam-digest.git
cd exam-digest
go build -o exd main.go
```

x64のwindows向けにビルドする場合は最後の行を以下に変更する
```bash
GOOS=windows GOARCH=amd64 go build -o exd.exe main.go
```

## Usage

```bash
exd [-r regexp] [-i file] [-d] [-o file] filepath
```

## Option

### -r

分析対象のファイルを表す正規表現を指定する

### -i

分析が終了したファイルや分析対象から除外するファイル名を各行ごとに記述したファイルのパスを指定する

### -d

分析が終了したファイルのファイル名を分析対象から除外しないようにするフラグ

### -o

分析結果を出力するためのファイルのパスを指定する
