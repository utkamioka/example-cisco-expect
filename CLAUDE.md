# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

Cisco Catalystスイッチに SSH 接続してコマンドを自動実行するGoツールです。`goexpect`ライブラリを使用してexpectパターンによるネットワーク機器自動化を実装しています。

## 開発コマンド

### ビルドとテスト
```bash
# ビルド
go build

# 依存関係確認
go mod tidy

# 構文チェック
go vet

# フォーマット
go fmt
```

### 実行とテスト
```bash
# コンパイル確認のみ（実機接続なし）
go build

# 実際の実行（Ciscoデバイス必要）
./example-cisco-expect -H 192.168.1.1 -u admin -p password
./example-cisco-expect --host IP --username USER --password PASS --enable-secret SECRET
```

## アーキテクチャ概要

### 単一ファイル構成
- **main.go**: 全ての機能を含む単一ファイル実装
  - CLI引数処理（Cobra）
  - SSH接続とgoexpect初期化  
  - Cisco expectパターン処理
  - ファイル出力機能

### 主要コンポーネント
1. **runCiscoExpect()**: メイン処理フロー（SSH接続→コマンド実行）
2. **executeAndSaveCommand()**: コマンド実行とファイル保存の汎用関数
3. **enterPrivilegedMode()**: enable特権モード移行処理

### 依存ライブラリ
- `github.com/google/goexpect`: SSH expectパターン実装
- `github.com/spf13/cobra`: CLI引数処理
- `golang.org/x/crypto/ssh`: SSH接続

### セキュリティ考慮事項
- **`ssh.InsecureIgnoreHostKey()`**: ホスト鍵検証を無効化（プロトタイプ用）
- パスワードはコマンドライン引数で渡すため、プロセス一覧で見える可能性

### Expectパターンの実装
- プロンプト待機: `regexp.MustCompile(\`[>#]$\`)`
- タイムアウト: 基本30秒、プロンプト確認は5-10秒
- 特権モード判定: `#`プロンプトの検出