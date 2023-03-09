# exam-digest

exam-viewer用の集計プログラム

特定ディレクトリ配下にある複数のテスト結果を纏めたファイルを生成する

## usage

```bash
command [-r|-i|-o]... filepath
```

## option

### -r
分析対象のファイルを表す正規表現を指定する

### -i
分析対象から除外するファイル名を各行ごとに記述したファイル
