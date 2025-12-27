package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/goexpect"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	host         string
	port         string
	username     string
	password     string
	enableSecret string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cisco-expect",
		Short: "Cisco Catalyst SSH automation tool",
		Long:  "SSH接続でCisco Catalystにログインし、show versionとshow loggingコマンドを実行して結果をファイルに保存します。",
		Run:   runCiscoExpect,
	}

	rootCmd.Flags().StringVarP(&host, "host", "H", "", "接続先IPアドレス (必須)")
	rootCmd.Flags().StringVarP(&port, "port", "P", "22", "SSH接続ポート番号 (デフォルト: 22)")
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "ログインユーザー名 (必須)")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "ログインパスワード (必須)")
	rootCmd.Flags().StringVarP(&enableSecret, "enable-secret", "e", "", "特権モードパスワード (オプション)")

	rootCmd.MarkFlagRequired("host")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCiscoExpect(cmd *cobra.Command, args []string) {
	fmt.Printf("Cisco Catalyst接続開始: %s:%s\n", host, port)

	// SSH設定
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// SSH接続
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		log.Fatalf("SSH接続エラー: %v", err)
	}
	defer sshClient.Close()

	// goexpectでSSH seesionを作成
	expecter, _, err := expect.SpawnSSH(sshClient, 30*time.Second)
	if err != nil {
		log.Fatalf("SSH session作成エラー: %v", err)
	}
	defer expecter.Close()

	fmt.Println("SSH接続成功")

	// ログイン後のプロンプト待機
	_, _, err = expecter.Expect(regexp.MustCompile(`[>#]$`), 10*time.Second)
	if err != nil {
		log.Fatalf("ログイン後プロンプト待機エラー: %v", err)
	}

	// show versionコマンド実行
	fmt.Println("show versionコマンド実行中...")
	if err := executeAndSaveCommand(expecter, "show version", "show_version_output.txt"); err != nil {
		log.Fatalf("show version実行エラー: %v", err)
	}
	fmt.Println("show version完了")

	// 特権モード移行確認と移行
	if err := enterPrivilegedMode(expecter); err != nil {
		log.Fatalf("特権モード移行エラー: %v", err)
	}

	// show loggingコマンド実行
	fmt.Println("show loggingコマンド実行中...")
	if err := executeAndSaveCommand(expecter, "show logging", "show_logging_output.txt"); err != nil {
		log.Fatalf("show logging実行エラー: %v", err)
	}
	fmt.Println("show logging完了")

	fmt.Println("全ての処理が完了しました")
}

// executeAndSaveCommand はコマンドを実行して結果をファイルに保存します
func executeAndSaveCommand(e *expect.GExpect, command, filename string) error {
	// コマンド送信
	if err := e.Send(command + "\r\n"); err != nil {
		return fmt.Errorf("コマンド送信エラー: %w", err)
	}

	// 実行結果の取得（プロンプト待機）
	result, _, err := e.Expect(regexp.MustCompile(`[>#]$`), 30*time.Second)
	if err != nil {
		return fmt.Errorf("コマンド実行結果待機エラー: %w", err)
	}

	// ファイルに保存
	if err := os.WriteFile(filename, []byte(result), 0644); err != nil {
		return fmt.Errorf("ファイル保存エラー: %w", err)
	}

	return nil
}

// enterPrivilegedMode は特権モードに移行します
func enterPrivilegedMode(e *expect.GExpect) error {
	// 現在のプロンプトをチェック
	if err := e.Send("\r\n"); err != nil {
		return fmt.Errorf("プロンプト確認エラー: %w", err)
	}

	result, _, err := e.Expect(regexp.MustCompile(`[>#]$`), 5*time.Second)
	if err != nil {
		return fmt.Errorf("プロンプト待機エラー: %w", err)
	}

	// 既に特権モード（最後の文字が#）の場合はそのまま返す
	if strings.HasSuffix(strings.TrimSpace(result), "#") {
		fmt.Println("既に特権モードです")
		return nil
	}

	// enableコマンドで特権モードに移行
	fmt.Println("特権モードに移行中...")
	if err := e.Send("enable\r\n"); err != nil {
		return fmt.Errorf("enableコマンド送信エラー: %w", err)
	}

	// パスワードプロンプトまたは特権プロンプトを待機
	result, _, err = e.Expect(regexp.MustCompile(`(Password:|[>#]$)`), 10*time.Second)
	if err != nil {
		return fmt.Errorf("enable後の応答待機エラー: %w", err)
	}

	// パスワードが要求された場合
	if strings.Contains(result, "Password:") {
		if enableSecret == "" {
			return fmt.Errorf("enableパスワードが要求されましたが、--enable-secretが指定されていません")
		}

		if err := e.Send(enableSecret + "\r\n"); err != nil {
			return fmt.Errorf("enableパスワード送信エラー: %w", err)
		}

		// 特権プロンプト待機（認証失敗パターンも検出）
		result, _, err := e.Expect(regexp.MustCompile(`(% Bad secrets|#$)`), 10*time.Second)
		if err != nil {
			return fmt.Errorf("特権プロンプト待機エラー: %w", err)
		}

		// デバッグ出力: expect結果を表示
		fmt.Printf("enable後の応答: %q\n", result)

		// 認証失敗の場合はエラーを返す
		if strings.Contains(result, "% Bad secrets") || strings.Contains(result, "Bad secrets") {
			return fmt.Errorf("enable認証に失敗しました: パスワードが正しくありません")
		}
	}

	fmt.Println("特権モード移行完了")
	return nil
}
