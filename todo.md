# uncode 実装 TODO リスト

## 言語設計
- [x] 仕様書の読み込みと分析
- [x] 言語構造の設計
  - [x] 型システムの詳細化
  - [x] エラーハンドリングの仕組み設計
  - [x] スコープ管理の設計
  - [x] メモリ管理の設計
- [x] 言語機能の設計
  - [x] 標準ライブラリの設計
  - [x] 入出力機能の設計
  - [x] 不足している仕様の補完

## 実装
- [x] レキサー（字句解析器）の実装
- [x] パーサー（構文解析器）の実装
- [x] AST（抽象構文木）の設計と実装
- [x] インタプリタの実装
  - [x] 評価器の実装
  - [x] 環境（変数スコープ）の実装
  - [x] 組み込み関数の実装

## リファクタリング
- [ ] ソースコードのモジュール分割
  - [ ] evaluator パッケージの分割
    - [ ] `evaluator/evaluator.go` - コア評価ロジック
    - [ ] `evaluator/builtins.go` - 組み込み関数の実装
    - [ ] `evaluator/expression_eval.go` - 式の評価ロジック
    - [ ] `evaluator/statements_eval.go` - 文の評価ロジック
    - [ ] `evaluator/function_eval.go` - 関数関連の評価ロジック
    - [ ] `evaluator/pipeline_eval.go` - パイプライン処理の評価ロジック
  - [ ] parser パッケージの分割
    - [ ] `parser/parser.go` - コアパーサーのロジック
    - [ ] `parser/expression_parser.go` - 式のパース処理
    - [ ] `parser/statement_parser.go` - 文のパース処理
    - [ ] `parser/literal_parser.go` - リテラルのパース処理
    - [ ] `parser/function_parser.go` - 関数関連のパース処理
  - [ ] ast パッケージの分割
    - [ ] `ast/ast.go` - コアAST構造とインターフェース
    - [ ] `ast/expressions.go` - 式関連のAST構造
    - [ ] `ast/statements.go` - 文関連のAST構造
    - [ ] `ast/literals.go` - リテラル関連のAST構造
    - [ ] `ast/function.go` - 関数関連のAST構造
  - [ ] object パッケージの分割
    - [ ] `object/object.go` - コアオブジェクト構造とインターフェース
    - [ ] `object/primitive_types.go` - 整数、文字列など基本型
    - [ ] `object/function.go` - 関数関連オブジェクト
    - [ ] `object/container.go` - 配列やハッシュなどのコンテナオブジェクト
    - [ ] `object/environment.go` - 環境関連コード
  - [ ] lexer パッケージの分割
    - [ ] `lexer/lexer.go` - コアレキサーロジック
    - [ ] `lexer/token_readers.go` - 特定のトークンを読むための関数
    - [ ] `lexer/helpers.go` - ヘルパー関数
- [ ] コードのドキュメント化（コメントの追加）
- [ ] テストカバレッジの向上
- [ ] パフォーマンス最適化

## テスト
- [ ] ユニットテストの作成
- [x] サンプルプログラムの作成と実行
- [x] インタプリタのビルドとテスト
- [ ] エラーケースのテスト

## ドキュメント
- [x] 言語仕様書の作成
- [x] 使い方マニュアルの作成
- [x] サンプルプログラムの解説
- [ ] 使い方ガイドの作成
- [ ] サンプルプログラム集の作成

## パッケージング
- [ ] ビルドスクリプトの作成
- [ ] リリースパッケージの作成
