# コンパイラ理論：体系的な理論解説

**📚 コンパイラ技術の理論的基盤**

## 🎯 コンパイラ理論の全体像

コンパイラ理論は、プログラミング言語を機械語に変換する技術の理論的基盤です。pugプロジェクトを通じて、以下の理論を段階的に実践します：

1. **形式言語理論** - 言語の数学的基盤
2. **字句解析理論** - 正規言語と有限オートマトン
3. **構文解析理論** - 文脈自由文法とプッシュダウンオートマトン
4. **意味解析理論** - 型理論と記号表管理
5. **コード生成理論** - 中間表現と最適化
6. **最適化理論** - データフロー解析と制御フロー解析

---

## 📖 第1章：形式言語理論の基礎

### 言語の数学的定義

プログラミング言語は形式的に以下のように定義されます：

**定義1.1（言語）**
アルファベット Σ 上の言語 L は、Σ* の部分集合である。

```
Σ = {a, b, c, ..., z, 0, 1, ..., 9, +, -, *, /, ...}  // アルファベット
Σ* = すべての有限文字列の集合
L ⊆ Σ*  // 言語は文字列の集合
```

**例：dogプログラミング言語**
```
L_dog = { "let x = 42;", "fn add(a, b) { a + b }", ... }
```

### チョムスキー階層

形式言語は生成能力により4つのクラスに分類されます：

| タイプ | 言語クラス | 文法 | オートマトン | 例 |
|--------|-----------|------|-------------|-----|
| **Type 0** | 句構造言語 | 無制限文法 | チューリングマシン | 任意の計算可能言語 |
| **Type 1** | 文脈依存言語 | 文脈依存文法 | 線形有界オートマトン | {aⁿbⁿcⁿ \| n≥1} |
| **Type 2** | 文脈自由言語 | 文脈自由文法 | プッシュダウンオートマトン | {aⁿbⁿ \| n≥1} |
| **Type 3** | 正規言語 | 正規文法 | 有限オートマトン | {(ab)*} |

**プログラミング言語の位置づけ**
- **字句解析**：Type 3（正規言語）
- **構文解析**：Type 2（文脈自由言語）
- **意味解析**：Type 1またはType 0（文脈依存）

### pugにおける応用

```
字句解析：識別子 = [a-zA-Z][a-zA-Z0-9]*     (正規言語)
構文解析：E → E + E | E * E | (E) | id      (文脈自由言語)
意味解析：変数は使用前に宣言されること          (文脈依存)
```

---

## 🔍 第2章：字句解析理論

### 正規言語と有限オートマトン

字句解析は正規言語理論に基づきます。

**定理2.1（クリーネの定理）**
以下は同等である：
1. 正規表現で表現可能
2. 有限オートマトンで認識可能
3. 正規文法で生成可能

### 有限オートマトンの定義

**定義2.1（決定性有限オートマトン：DFA）**
DFA M = (Q, Σ, δ, q₀, F) ここで：
- Q：状態の有限集合
- Σ：入力アルファベット
- δ：Q × Σ → Q（遷移関数）
- q₀ ∈ Q：初期状態
- F ⊆ Q：受理状態の集合

**例：整数リテラルを認識するDFA**
```
Q = {q₀, q₁, q₂}
Σ = {0, 1, 2, ..., 9, +, -}
δ(q₀, +) = q₁, δ(q₀, -) = q₁, δ(q₀, digit) = q₂
δ(q₁, digit) = q₂
δ(q₂, digit) = q₂
F = {q₂}
```

### 非決定性有限オートマトン（NFA）

**定義2.2（非決定性有限オートマトン：NFA）**
NFA M = (Q, Σ, δ, q₀, F) ここで：
- δ：Q × (Σ ∪ {ε}) → 2^Q（非決定的遷移関数）

**定理2.2（サブセット構成法）**
任意のNFAに対して、同等な言語を認識するDFAが構成可能。

### Thompson構成法

正規表現からNFAを構成する方法：

1. **基本要素**：a → NFA(a)
2. **連接**：AB → NFA(A)とNFA(B)をε遷移で接続
3. **選択**：A|B → ε遷移で分岐・合流
4. **クリーネ閉包**：A* → ε遷移でループ構成

**pugでの実装例**
```go
// 識別子の正規表現: [a-zA-Z][a-zA-Z0-9]*
func (l *Lexer) readIdentifier() string {
    position := l.position
    // 最初の文字は英字のみ
    for isLetter(l.ch) {
        l.readChar()
    }
    // 後続文字は英数字
    for isLetter(l.ch) || isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}
```

---

## 🌳 第3章：構文解析理論

### 文脈自由文法

**定義3.1（文脈自由文法：CFG）**
CFG G = (N, T, P, S) ここで：
- N：非終端記号の有限集合
- T：終端記号の有限集合（N ∩ T = ∅）
- P：N → (N ∪ T)* の生成規則の有限集合
- S ∈ N：開始記号

**例：算術式の文法**
```
G = ({E, T, F}, {+, *, (, ), id}, P, E)
P: E → E + T | T
   T → T * F | F  
   F → (E) | id
```

### 導出と構文木

**定義3.2（左端導出）**
各ステップで最も左の非終端記号を書き換える導出。

**例：id + id * id の左端導出**
```
E ⇒ E + T ⇒ T + T ⇒ F + T ⇒ id + T ⇒ id + T * F 
  ⇒ id + F * F ⇒ id + id * F ⇒ id + id * id
```

### 曖昧性と優先順位

**定義3.3（曖昧な文法）**
ある文字列に対して2つ以上の異なる左端導出（構文木）が存在する文法。

**問題例**：E → E + E | E * E | id
文字列 "id + id * id" に対して複数の構文木が存在。

**解決法1：文法の書き換え**
```
E → E + T | T        // + の優先度を低く
T → T * F | F        // * の優先度を高く  
F → (E) | id
```

**解決法2：演算子優先順位法**
```
優先順位テーブル：
    |  +  *  (  )  id $
----+------------------
 +  |  >  <  <  >  <  >
 *  |  >  >  <  >  <  >
 (  |  <  <  <  =  <  
 )  |  >  >     >     >
id  |  >  >     >     >
 $  |  <  <  <        
```

### LL構文解析

**定義3.4（LL(k)文法）**
左から右に読み、左端導出を行い、k個の先読みで決定的に解析可能な文法。

**FIRST集合とFOLLOW集合**
```
FIRST(α) = { a | α ⇒* aβ, a ∈ T } ∪ { ε | α ⇒* ε }
FOLLOW(A) = { a | S ⇒* αAaβ, a ∈ T } ∪ { $ | S ⇒* αA }
```

**例：E → T E'の構文解析表**
```
     +    *    (   )   id   $
E           E→TE'      E→TE'
E'  E'→+TE'      E'→ε   E'→ε
T           T→FT'      T→FT'
T'  T'→ε  T'→*FT' T'→ε T'→ε
F           F→(E)      F→id
```

### LR構文解析

**定義3.5（LR(k)文法）**
左から右に読み、右端導出の逆を行い、k個の先読みで決定的に解析可能な文法。

**LR構文解析の利点**
- LL文法より広いクラスを扱える
- 左再帰を許可
- エラー検出が早い

**SLR構文解析表の構成**
1. LR(0)項目集合の構成
2. GOTOとCLOSURE関数の計算
3. SHIFTとREDUCEアクションの決定

**pugのPratt Parser**
pugは演算子優先順位を明示的に扱えるPratt Parserを採用：

```go
type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

// 優先順位に基づく解析
func (p *Parser) parseExpression(precedence int) ast.Expression {
    left := p.prefixParseFns[p.curToken.Type]()
    
    for precedence < p.peekPrecedence() {
        infix := p.infixParseFns[p.peekToken.Type]
        left = infix(left)
    }
    
    return left
}
```

---

## 🔤 第4章：意味解析理論

### 型理論の基礎

**定義4.1（型システム）**
型システムは、プログラムの部分に型を割り当て、これらの型が一貫していることを確認する規則の集合。

**型の分類**
```
基本型（Base Types）：int, bool, string
関数型（Function Types）：T₁ → T₂
直積型（Product Types）：T₁ × T₂
直和型（Sum Types）：T₁ + T₂
多態型（Polymorphic Types）：∀α.T
```

### 型判定規則

**型判定の記法**
```
Γ ⊢ e : T
```
環境Γの下で式eは型Tを持つ

**基本的な型判定規則**
```
[Var] Γ(x) = T
      -----------
      Γ ⊢ x : T

[App] Γ ⊢ e₁ : T₁ → T₂    Γ ⊢ e₂ : T₁
      ----------------------------
      Γ ⊢ e₁ e₂ : T₂

[Abs] Γ, x : T₁ ⊢ e : T₂
      --------------------
      Γ ⊢ λx.e : T₁ → T₂
```

### 型推論アルゴリズム

**Algorithm W（Hindley-Milner型推論）**
```
W(Γ, e) = (S, T) where:
  S: 最も一般的な単一化子（substitution）
  T: eの型
```

**例：let多態性**
```
let id = λx.x in (id 3, id true)
```
1. `id`の型を`∀α.α → α`と推論
2. 使用箇所で具体化：`int → int`, `bool → bool`

### pugの型システム実装

```go
type Type interface {
    String() string
}

type IntegerType struct{}
type BooleanType struct{}
type FunctionType struct {
    Parameters []Type
    ReturnType Type
}

func (tc *TypeChecker) checkExpression(expr ast.Expression) Type {
    switch e := expr.(type) {
    case *ast.IntegerLiteral:
        return &IntegerType{}
    case *ast.InfixExpression:
        leftType := tc.checkExpression(e.Left)
        rightType := tc.checkExpression(e.Right)
        return tc.checkInfixOperation(e.Operator, leftType, rightType)
    }
}
```

---

## 🏗️ 第5章：中間表現理論

### 中間表現の必要性

**なぜ中間表現が必要か**
1. **フロントエンドとバックエンドの分離**
2. **最適化の共通基盤**
3. **複数ターゲットへの対応**
4. **解析の簡素化**

### 三番地コード（Three-Address Code）

**定義5.1（三番地コード）**
各命令が最大3つのアドレスを持つ中間表現：
```
x = y op z    // 二項演算
x = op y      // 単項演算
x = y         // コピー
goto L        // 無条件ジャンプ
if x relop y goto L  // 条件分岐
```

**例：(a + b) * (c + d)の三番地コード**
```
t1 = a + b
t2 = c + d  
t3 = t1 * t2
```

### SSA形式（Static Single Assignment）

**定義5.2（SSA形式）**
各変数が正確に一度だけ定義される中間表現。

**φ関数**
制御フローの合流点で値を選択：
```
if (condition) {
    x1 = a + b;
} else {
    x2 = c + d;
}
x3 = φ(x1, x2);  // 制御フローに応じてx1またはx2を選択
```

**SSAの利点**
1. **def-use関係の明確化**
2. **最適化の簡素化**
3. **並行性の解析**
4. **レジスタ割り当ての効率化**

### 制御フローグラフ（CFG）

**定義5.3（制御フローグラフ）**
CFG = (V, E) where:
- V：基本ブロックの集合
- E：制御フローエッジの集合

**基本ブロック**
- 最初の命令から最後の命令まで順次実行
- 分岐のない直線的な命令列
- エントリは先頭のみ、エグジットは末尾のみ

**例：if文のCFG**
```
    [entry]
       |
   [condition]
    /       \
[then]    [else]
    \       /
    [merge]
       |
     [exit]
```

---

## ⚡ 第6章：最適化理論

### データフロー解析

**定義6.1（データフロー解析）**
プログラムの制御フロー全体にわたって情報を伝播させる解析手法。

**格子理論（Lattice Theory）**
```
(L, ⊑)：半順序集合
⊥：最小元（bottom）
⊤：最大元（top）
⊔：上限（join）
⊓：下限（meet）
```

**到達定義解析（Reaching Definitions）**
```
gen[B] = Bで生成される定義
kill[B] = Bで削除される定義
in[B] = ⋃{out[P] | P ∈ pred(B)}
out[B] = gen[B] ⋃ (in[B] - kill[B])
```

**生存変数解析（Liveness Analysis）**
```
use[B] = Bで使用される変数
def[B] = Bで定義される変数
out[B] = ⋃{in[S] | S ∈ succ(B)}
in[B] = use[B] ⋃ (out[B] - def[B])
```

### 最適化変換

**定数畳み込み（Constant Folding）**
```
x = 2 + 3  →  x = 5
```

**定数伝播（Constant Propagation）**
```
x = 5      →  x = 5
y = x + 3      y = 8
```

**デッドコード除去（Dead Code Elimination）**
```
x = a + b    →  x = a + b
y = c + d        z = x * 2
z = x * 2
```

**共通部分式除去（Common Subexpression Elimination）**
```
a = b + c    →  t = b + c
d = b + c        a = t
                 d = t
```

### ループ最適化

**自然ループの識別**
```
定理：逆エッジ(t, h)について、hがtを支配するとき、
     (t, h)は自然ループを形成する。
```

**ループ不変式移動（Loop Invariant Code Motion）**
```
for (i = 0; i < n; i++) {
    x = a + b;      // ループ不変
    array[i] = x;
}

↓

x = a + b;          // ループ外に移動
for (i = 0; i < n; i++) {
    array[i] = x;
}
```

**ループアンローリング（Loop Unrolling）**
```
for (i = 0; i < n; i++) {     for (i = 0; i < n; i += 4) {
    a[i] = b[i] + c[i];   →      a[i] = b[i] + c[i];
}                              a[i+1] = b[i+1] + c[i+1];
                               a[i+2] = b[i+2] + c[i+2];
                               a[i+3] = b[i+3] + c[i+3];
                           }
```

---

## 🎯 第7章：レジスタ割り当て理論

### グラフ彩色法

**定義7.1（干渉グラフ）**
変数間の干渉関係を表現するグラフ G = (V, E) where:
- V：変数の集合
- E：同時に生存する変数対のエッジ

**定理7.1（グラフ彩色とレジスタ割り当て）**
k個のレジスタでのレジスタ割り当て問題は、干渉グラフのk彩色問題と等価。

**Chaitin彩色アルゴリズム**
```
1. 次数がk未満のノードを除去（スタックにプッシュ）
2. グラフが空になるまで繰り返し
3. スタックからポップしながら彩色
4. 彩色不可能な場合はスピル
```

**スピル処理**
```
レジスタ不足の変数をメモリに退避：
load  temp, spill_var    // メモリからロード
...use temp...
store temp, spill_var    // メモリに格納
```

### 線形走査レジスタ割り当て

**Linear Scanアルゴリズム**
```
1. 全ての生存区間を開始点でソート
2. 各区間について：
   - 期限切れの区間を解放
   - 利用可能レジスタがあれば割り当て
   - なければスピル
```

**時間計算量**：O(n log n) vs Chaitinの O(n³)

---

## 🏛️ 第8章：コード生成理論

### 命令選択

**定義8.1（命令選択問題）**
中間表現の各ノードに対して、同等の意味を持つ機械語命令列を選択する問題。

**Tree Pattern Matching**
```
中間表現：    ADD(LOAD(addr1), LOAD(addr2))
命令選択：    mov addr1, %rax
             add addr2, %rax
```

**動的プログラミング法**
```
cost[node] = min{cost[pattern] + Σcost[child]}
```

### 命令スケジューリング

**パイプライン最適化**
```
命令の依存関係を考慮してパイプライン効率を最大化：

Before:             After:
add r1, r2, r3      add r1, r2, r3
add r4, r3, r5  →   add r6, r7, r8    // 並行実行可能
add r6, r7, r8      add r4, r3, r5
```

**List Scheduling**
```
1. 依存グラフを構築
2. 準備完了命令をリストに追加
3. 優先度に基づいて命令を選択・スケジュール
4. 後続命令を準備完了リストに追加
```

### リンク時最適化（LTO）

**全体プログラム最適化**
```
- 関数間インライン展開
- 未使用関数の除去
- 定数伝播の全体適用
- デッドコード除去の全体適用
```

---

## 🚀 第9章：高度な最適化理論

### ポインタ解析

**定義9.1（別名解析）**
プログラム中のポインタが指し示す可能性のある記憶場所を決定する解析。

**Andersen型解析**
```
包含制約の集合を解く：
*p = q  →  points-to(q) ⊆ points-to(*p)
p = &x  →  {x} ⊆ points-to(p)
```

### 並行最適化

**依存関係解析**
```
ループの並行化可能性を判定：
for (i = 0; i < n; i++) {
    a[i] = a[i-1] + b[i];  // 依存関係あり（並行化不可）
}

for (i = 0; i < n; i++) {
    a[i] = b[i] + c[i];    // 依存関係なし（並行化可能）
}
```

**自動ベクトル化**
```
スカラーループをSIMD命令に変換：
for (i = 0; i < n; i++) {
    c[i] = a[i] + b[i];
}
↓
vectorized add operations
```

### プロファイル誘導最適化（PGO）

**実行時情報の活用**
```
1. プロファイル情報収集
2. ホットパスの特定
3. 分岐予測情報の利用
4. 関数のインライン判定
```

---

## 📊 第10章：性能解析理論

### 計算複雑度理論

**字句解析**：O(n) - 線形時間
**構文解析**：
- LL(1), LR(1): O(n)
- CYK: O(n³)
- 一般CFG: O(n³)

**最適化アルゴリズム**：
- データフロー解析：O(n³)
- レジスタ割り当て：NP完全
- 命令スケジューリング：NP完全

### アムダールの法則

**定理10.1（Amdahl's Law）**
```
速度向上比 = 1 / ((1 - P) + P/S)
P: 並行化可能な部分の割合
S: 並行部分の速度向上比
```

**コンパイラ最適化への応用**
```
全体実行時間 = 最適化不可能部分 + 最適化部分/改善率
```

### 性能測定理論

**ベンチマーク設計原則**
1. **代表性**：実際のワークロードを反映
2. **再現性**：同一条件で同一結果
3. **公平性**：比較対象間で同等の条件
4. **統計的有意性**：十分なサンプル数

---

## 🎓 第11章：pugプロジェクトでの理論応用

### Phase 1での理論応用

**字句解析**
- 正規表現 → NFA → DFA変換
- Thompson構成法の実装

**構文解析**  
- Pratt Parserによる演算子優先順位
- 再帰下降法の実装

### Phase 2での理論応用

**型システム**
- 単純型付きλ計算
- 型判定規則の実装

**コード生成**
- 三番地コードからアセンブリへの変換
- レジスタ使用の最適化

### Phase 3での理論応用

**中間表現**
- SSA形式の構築
- φ関数の配置

**データフロー解析**
- 到達定義解析
- 生存変数解析

### Phase 4での理論応用

**LLVM統合**
- LLVM IRの生成
- 高度最適化パスの適用

**並行最適化**
- 依存関係解析
- 自動並行化

---

## 🔬 第12章：実装と理論の対応

### 理論の実装での簡略化

**実用化での妥協点**
```
理論：        実装：
完全な型推論 → 明示的型注釈
完全な最適化 → ヒューリスティック
厳密な解析   → 近似解析
```

### 理論研究と産業実装

**学術研究の成果**
- SSA形式（1991年）
- グラフ彩色レジスタ割り当て（1981年）
- データフロー解析（1970年代）

**産業での採用**
- GCC：2004年にSSA採用
- LLVM：設計当初からSSA
- JIT コンパイラでの活用

---

## 📚 第13章：さらなる学習のために

### 推奨教科書

**基礎理論**
1. **「Introduction to Automata Theory, Languages, and Computation」** - Hopcroft, Motwani, Ullman
2. **「Formal Language Theory」** - Salomaa

**コンパイラ理論**
1. **「Compilers: Principles, Techniques, and Tools」** - Aho, Lam, Sethi, Ullman
2. **「Modern Compiler Implementation」** - Appel
3. **「Engineering a Compiler」** - Cooper, Torczon

**最適化理論**
1. **「Optimizing Compilers for Modern Architectures」** - Allen, Kennedy
2. **「Advanced Compiler Design and Implementation」** - Muchnick

### 研究論文

**重要な論文**
1. **SSA Form** - Cytron et al. (1991)
2. **Graph Coloring Register Allocation** - Chaitin (1982)
3. **Linear Scan Register Allocation** - Poletto & Sarkar (1999)
4. **Profile-Guided Optimization** - Fisher (1981)

### オンライン学習リソース

**講義動画**
- MIT 6.035 Computer Language Engineering
- Stanford CS 143 Compilers  
- Berkeley CS 164 Programming Languages and Compilers

**実装プロジェクト**
- LLVM Tutorial
- TinyLang Compiler
- MiniJava Compiler

---

## 🎯 理解度確認問題

### 基礎理論

1. 正規言語と文脈自由言語の違いを説明せよ
2. LL(1)とLR(1)の利点・欠点を比較せよ
3. SSA形式の利点を3つ挙げよ

### 応用理論

1. データフロー解析の格子理論について説明せよ
2. グラフ彩色法でのスピル処理を説明せよ
3. アムダールの法則をコンパイラ最適化に適用せよ

### pugプロジェクト

1. pugのPratt Parserの理論的基盤は何か
2. Phase 2の型システムは型理論のどの部分を実装しているか
3. 予定されているPhase 3のSSA実装の理論的意義は何か

---

**🎉 コンパイラ理論の学習完了おめでとうございます！**

この理論的基盤を理解することで、pugプロジェクトの各実装がなぜそのように設計されているか、そして産業レベルのコンパイラがどのような理論に基づいているかが明確になりました。理論と実践の両輪で、真に理解の深いコンパイラ技術者を目指しましょう！

**次のステップ：理論を実践で確認し、さらなる高度な最適化技術に挑戦しましょう！**