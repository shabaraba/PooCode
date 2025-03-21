# uncode使用マニュアル

## 1. はじめに

uncodeへようこそ！このマニュアルでは、uncode言語のインストール方法から基本的な使い方、サンプルプログラムの解説まで、uncodeを使いこなすために必要な情報を提供します。

## 2. インストール方法

### 2.1 前提条件

uncodeを実行するには、以下のソフトウェアが必要です：

- Go言語（バージョン1.18以上）

### 2.2 インストール手順

1. リポジトリをクローンします：
   ```
   git clone https://github.com/yourusername/uncode.git
   cd uncode
   ```

2. ビルドスクリプトを実行します：
   ```
   ./build.sh
   ```

3. インストールが成功すると、`bin/uncode`実行ファイルが生成されます。

## 3. 基本的な使い方

### 3.1 プログラムの実行

uncodeプログラムを実行するには、以下のコマンドを使用します：

```
./bin/uncode プログラムファイル名.poo
```

または、絵文字の拡張子を使用する場合：

```
./bin/uncode プログラムファイル名.💩
```

### 3.2 最初のプログラム

以下は、uncodeでの「Hello, World!」プログラムの例です：

```
"Hello, World!" |> print
```

このプログラムを`hello.poo`として保存し、以下のコマンドで実行します：

```
./bin/uncode hello.poo
```

## 4. 言語の基本

### 4.1 変数と代入

uncodeでは、変数に値を代入するには`>>`演算子を使用します：

```
42 >> answer
"Hello" >> greeting

answer |> print  // 42を出力
greeting |> print  // Helloを出力
```

### 4.2 特殊変数

uncodeには2つの特殊変数があります：

- `🍕` (ピザ): 関数の入力パラメータを表します
- `💩` (プー): 関数の戻り値を表します

例：
```
def double(): int -> int {
    🍕 * 2 >> 💩
}

5 |> double |> print  // 10を出力
```

### 4.3 パイプライン記法

uncodeの特徴的な構文はパイプライン記法（`|>`）です：

```
5 |> add 3 |> mul 2 |> print  // (5 + 3) * 2 = 16を出力
```

これは以下の従来の記法と同等です：

```
print(mul(add(5, 3), 2))
```

### 4.4 関数定義

関数は`def`キーワードを使用して定義します：

```
def add(a, b): int -> int {
    🍕 + b >> 💩
}

def greet(name): str -> str {
    "Hello, " + 🍕 + "!" >> 💩
}
```

### 4.5 条件付き関数

uncodeでは、条件によって異なる実装を持つ関数を定義できます：

```
def process() if 🍕 > 0: int -> str {
    "正の数です" >> 💩
}

def process() if 🍕 < 0: int -> str {
    "負の数です" >> 💩
}

def process(): int -> str {
    "ゼロです" >> 💩
}

5 |> process |> print   // "正の数です"を出力
-3 |> process |> print  // "負の数です"を出力
0 |> process |> print   // "ゼロです"を出力
```

### 4.6 クラスとオブジェクト

クラスを定義し、インスタンスを作成する例：

```
class Person
    public name
    public age
    
    def init(name, age): null -> null {
        name >> 🍕's name
        age >> 🍕's age
        null >> 💩
    }
    
    def greet(): null -> str {
        "こんにちは、" + 🍕's name + "さん！" >> 💩
    }

Person("田中", 30) |> init >> person
person |> .greet |> print  // "こんにちは、田中さん！"を出力
```

## 5. サンプルプログラム解説

### 5.1 FizzBuzz

以下は、uncodeでのFizzBuzzプログラムです：

```
def Int.is_multiple_of(num): int -> bool {
    🍕 % num == 0 >> 💩
}

def fizzbuzz(num): int -> str {
    "" >> result
    
    num |> is_multiple_of 3 >> is_multiple_of_3
    num |> is_multiple_of 5 >> is_multiple_of_5
    
    is_multiple_of_3 && is_multiple_of_5 >> is_multiple_of_both
    
    is_multiple_of_both |> eq true >> is_fizzbuzz
    is_multiple_of_3 && not is_multiple_of_both |> eq true >> is_fizz
    is_multiple_of_5 && not is_multiple_of_both |> eq true >> is_buzz
    
    is_fizzbuzz |> eq true |> add "FizzBuzz" >> result
    is_fizz |> eq true |> add "Fizz" >> result
    is_buzz |> eq true |> add "Buzz" >> result
    
    not is_fizzbuzz && not is_fizz && not is_buzz |> eq true |> add num |> to_string >> result
    
    result >> 💩
}

1 >> i
100 >> max

{
    i |> fizzbuzz |> print
    i |> add 1 >> i
    i |> le max
} |> eq true |> print
```

このプログラムの特徴：
- `Int.is_multiple_of`メソッドを定義して、数値が特定の数の倍数かどうかを判定
- パイプライン記法を使用して、処理の流れを明確に表現
- 条件判定の結果を変数に格納し、それに基づいて出力を決定
- ループ処理を条件式と再帰で表現

### 5.2 TODOリスト

以下は、uncodeでのTODOリストプログラムです：

```
class Todo
    public title
    public expired_at
    public status
    
    def Todo.STATUS_TODO(): null -> str { "TODO" >> 💩 }
    def Todo.STATUS_DOING(): null -> str { "DOING" >> 💩 }
    def Todo.STATUS_DONE(): null -> str { "DONE" >> 💩 }
    
    def is_same_title(target): Todo -> bool {
        🍕's title == target's title >> 💩
    }
    
    def to_string(): null -> str {
        "[" + status + "] " + title + " (期限: " + expired_at + ")" >> 💩
    }

class TodoApp
    public todo_list
    
    def init(): null -> null {
        [] >> todo_list
        null >> 💩
    }
    
    def add(item): Todo -> TodoApp {
        todo_list |> add item >> todo_list
        🍕 >> 💩
    }
    
    def find_index(title): str -> int {
        0 >> i
        -1 >> result
        
        {
            i |> lt todo_list's length |> eq true >> continue
            
            todo_list |> get i >> current_item
            current_item's title |> eq title |> eq true >> found
            
            found |> eq true |> add i >> result
            found |> eq true |> add false >> continue
            
            i |> add 1 >> i
            continue
        } |> eq true
        
        result >> 💩
    }
    
    def change_status(title, new_status): str -> bool {
        title |> find_index >> index
        
        index |> ge 0 |> eq true >> found
        
        found |> eq true |> add {
            todo_list |> get index >> item
            new_status >> item's status
            true
        } >> 💩
        
        found |> eq false |> add false >> 💩
    }
    
    def show_all(): null -> null {
        "===== TODOリスト =====" |> print
        
        todo_list |> each {
            🍕 |> .to_string |> print
        }
        
        "======================" |> print
        
        null >> 💩
    }

TodoApp() |> init >> app

Todo(
    title: "牛乳を買う",
    expired_at: "2023-01-01",
    status: Todo.STATUS_TODO()
) |> app's add

Todo(
    title: "レポートを書く",
    expired_at: "2023-01-15",
    status: Todo.STATUS_DOING()
) |> app's add

app |> .show_all

"牛乳を買う" |> app's change_status Todo.STATUS_DONE() |> print

app |> .show_all
```

このプログラムの特徴：
- クラス定義を使用して、`Todo`と`TodoApp`の2つのクラスを作成
- クラスメソッドとインスタンスメソッドの両方を実装
- プロパティアクセスに`'s`記法を使用
- 配列操作とループ処理を組み合わせて、TODOリストの管理機能を実装

### 5.3 計算機

以下は、uncodeでの計算機プログラムです：

```
def add(a, b): int -> int {
    🍕 + b >> 💩
}

def sub(a, b): int -> int {
    🍕 - b >> 💩
}

def mul(a, b): int -> int {
    🍕 * b >> 💩
}

def div(a, b): int -> float {
    🍕 / b >> 💩
}

def show_result(expr, result): str -> null {
    expr + " = " + result |> to_string |> print
    null >> 💩
}

5 |> add 3 >> result1
"5 + 3" |> show_result result1

10 |> sub 4 >> result2
"10 - 4" |> show_result result2

6 |> mul 7 >> result3
"6 * 7" |> show_result result3

20 |> div 4 >> result4
"20 / 4" |> show_result result4

5 |> add 3 |> mul 2 |> sub 1 >> result5
"(5 + 3) * 2 - 1" |> show_result result5

def is_even() if 🍕 % 2 == 0: int -> str { "偶数です" >> 💩 }
def is_even(): int -> str { "奇数です" >> 💩 }

2 |> is_even |> print
3 |> is_even |> print
```

このプログラムの特徴：
- 基本的な四則演算関数を定義
- パイプライン記法を使用して、計算の流れを明確に表現
- 条件付き関数を使用して、偶数と奇数を判定する機能を実装

## 6. 高度な機能

### 6.1 並列パイプ

並列パイプ（`|`）を使用すると、条件分岐を表現できます：

```
condition | do_if_true | do_if_false
```

### 6.2 メソッドチェーン

メソッドチェーンを使用して、オブジェクトのメソッドを連続して呼び出すことができます：

```
object |> .method1 |> .method2 |> .method3
```

### 6.3 クロージャ

uncodeでは、関数内で定義された変数にアクセスするクロージャを作成できます：

```
def make_counter(): null -> function {
    0 >> count
    
    def increment(): null -> int {
        count |> add 1 >> count
        count >> 💩
    }
    
    increment >> 💩
}

make_counter() >> counter
counter() |> print  // 1を出力
counter() |> print  // 2を出力
counter() |> print  // 3を出力
```

## 7. デバッグとトラブルシューティング

### 7.1 一般的なエラー

- 構文エラー: 括弧の不一致や不正な演算子の使用など
- 型エラー: 互換性のない型の操作
- 名前エラー: 未定義の変数や関数の参照

### 7.2 デバッグのヒント

- `print`関数を使用して、変数の値を確認する
- 複雑な式を小さな部分に分割して、段階的にテストする
- エラーメッセージを注意深く読み、問題の箇所を特定する

## 8. ベストプラクティス

### 8.1 コーディング規約

- 変数名と関数名には、意味のある名前を使用する
- コメントを適切に使用して、コードの意図を明確にする
- 複雑な処理は小さな関数に分割する

### 8.2 パフォーマンスの最適化

- 大きな配列やハッシュマップの操作には注意する
- 再帰呼び出しの深さに注意する
- 不必要な計算を避ける

## 9. 参考資料

- [uncode言語仕様書](language_specification.md): 言語の詳細な仕様
- [サンプルプログラム集](../examples/): 様々なサンプルプログラム
- [GitHub リポジトリ](https://github.com/yourusername/uncode): ソースコードとイシュートラッカー

## 10. おわりに

uncodeは学習と実験のために設計された言語です。この言語を通じて、プログラミングの楽しさと創造性を体験していただければ幸いです。質問やフィードバックがあれば、GitHubリポジトリのイシュートラッカーでお知らせください。

Happy coding with uncode! 💩
