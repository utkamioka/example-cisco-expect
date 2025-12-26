# example-cisco-expect

Cisco Catalyst スイッチに SSH 接続し、基本コマンドを自動実行してファイルに保存するコマンドラインツールです。

## 概要

このツールは以下の処理を自動実行します：

1. SSH でCisco Catalystにログイン
2. `show version` コマンドを実行し、結果を `show_version_output.txt` に保存
3. 特権モードに移行（必要に応じて）
4. `show logging` コマンドを実行し、結果を `show_logging_output.txt` に保存
5. ログアウト

## 前提条件

- Go 1.19 以上
- SSH 接続可能なCisco Catalyst スイッチ
- ログイン可能なユーザーアカウント

## インストール

```bash
git clone <repository>
cd example-cisco-expect
go build
```

## 使用方法

### 基本的な使用方法

```bash
./example-cisco-expect --host 192.168.1.1 --username admin --password yourpassword
```

### 特権モードのパスワードが必要な場合

```bash
./example-cisco-expect --host 192.168.1.1 --username admin --password yourpassword --enable-secret enablepassword
```

### 短縮オプションを使用

```bash
./example-cisco-expect -H 192.168.1.1 -u admin -p yourpassword -e enablepassword
```

### ポート番号を指定

```bash
./example-cisco-expect -H 192.168.1.1 -P 2222 -u admin -p yourpassword
```

## コマンドラインオプション

| オプション | 短縮形 | 説明 | 必須 | デフォルト値 |
|-----------|--------|------|------|-------------|
| `--host` | `-H` | 接続先IPアドレス | ✓ | - |
| `--port` | `-P` | SSH接続ポート番号 | | 22 |
| `--username` | `-u` | ログインユーザー名 | ✓ | - |
| `--password` | `-p` | ログインパスワード | ✓ | - |
| `--enable-secret` | `-e` | 特権モードパスワード | | - |
| `--help` | `-h` | ヘルプ表示 | | - |

## 出力ファイル

実行後、以下のファイルが生成されます：

- `show_version_output.txt` - `show version` コマンドの実行結果
- `show_logging_output.txt` - `show logging` コマンドの実行結果

## 使用ライブラリ

- [google/goexpect](https://github.com/google/goexpect) - SSH接続とexpectパターン処理
- [spf13/cobra](https://github.com/spf13/cobra) - コマンドライン引数処理  
- [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) - SSH接続

## 動作例

```bash
$ ./example-cisco-expect -H 192.168.1.100 -u admin -p cisco123
Cisco Catalyst接続開始: 192.168.1.100:22
SSH接続成功
show versionコマンド実行中...
show version完了
特権モードに移行中...
特権モード移行完了
show loggingコマンド実行中...
show logging完了
全ての処理が完了しました
```

## 注意事項

- **セキュリティ**: パスワードがコマンドライン引数として表示される可能性があります
- **ホスト鍵検証**: SSH接続時のホスト鍵検証を無効化しています
- **プロトタイプ**: 本ツールは学習・検証目的のプロトタイプです
- **エラー処理**: 基本的なエラー処理のみ実装されています

## トラブルシューティング

### SSH接続エラー
- ホスト名・IPアドレスが正しいか確認
- ポート番号が正しいか確認（デフォルト22）
- ユーザー名・パスワードが正しいか確認

### 特権モード移行エラー  
- `--enable-secret` オプションでenable パスワードを指定
- ユーザーアカウントに適切な権限があるか確認

### タイムアウトエラー
- ネットワーク接続が安定しているか確認
- デバイスの応答が遅い場合は、コード内のタイムアウト値を調整

## ライセンス

このプロジェクトは学習・検証目的のサンプルコードです。