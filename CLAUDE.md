# example-cisco-expect

## プロジェクト概要

Cisco Catalyst スイッチに SSH 接続し、基本的なコマンドを自動実行するプロトタイプツールです。
Go言語とgoexpectライブラリを使用して、expectパターンによるネットワーク機器自動化の基礎を学習することが目的です。

## 使用技術

- **言語**: Go
- **ライブラリ**:
  - `github.com/google/goexpect` - SSH + expectパターンの実装
  - `github.com/spf13/cobra` - CLIインターフェース
  - `golang.org/x/crypto/ssh` - SSH接続

## 機能

1. SSH接続によるCisco Catalystログイン
2. `show version` コマンド実行と結果保存
3. 特権モード移行（enable）
4. `show logging` コマンド実行と結果保存
5. 適切な接続終了処理

## ビルド・実行

```bash
# ビルド
go build

# 実行例
./example-cisco-expect -H 192.168.1.1 -u admin -p password
./example-cisco-expect --host 192.168.1.1 --username admin --password secret --enable-secret enablepass
```

## 開発状況

- [x] Go module初期化
- [x] 必要ライブラリの追加
- [x] 基本CLI構造の実装
- [x] SSH接続とgoexpectセットアップ
- [x] show versionコマンド実行
- [x] enable特権モード移行処理
- [x] show loggingコマンド実行
- [x] CLIオプション短縮形追加

## テスト方法

実機のCisco Catalystが必要なため、現在はコンパイル確認のみ実施。
実際のテストは以下の環境で実施予定：

```bash
./example-cisco-expect -H <catalyst_ip> -u <username> -p <password>
```

## 今後の拡張可能性

- 他ベンダー機器への対応（Juniper, Arista等）
- 設定変更コマンドの実行
- 複数機器への同時接続
- 実行結果の構造化（JSON出力等）

## 注意事項

- プロトタイプ目的のため、本番環境での使用は推奨しません
- パスワードはコマンドライン引数で指定するため、プロセス一覧で見える可能性があります
- SSH接続時のホスト鍵検証を無効化しています（`InsecureIgnoreHostKey`）