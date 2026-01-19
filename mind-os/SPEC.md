# Mind-OS システム仕様書 (System Specification)

**Version:** 1.4.0  
**Last Updated:** 2026-01-20

## 1. 概要 (Overview)
Mind-OS (Mental Interface & Neural Driver OS) は、生物的な「心」の働きを模倣するAIエージェントコアです。
単なるテキスト応答だけでなく、入力に対する「感情(Emotion)」「記憶(Memory)」「意欲(Motivation)」「理性(Sanity)」「恒常性(Homeostasis)」「概日リズム(Circadian Rhythm)」の内部状態を持ち、それらが相互に作用しながら状態を遷移させます。

## 2. システムアーキテクチャ (Architecture)

システムは生物の脳構造を模した以下のモジュールで構成されています。

```mermaid
graph TD
    Input[Sensory Input] --> Brain
    Brain --> Thalamus[Thalamus<br>(視床)]
    Thalamus --> Hypo[Hypothalamus<br>(視床下部)]
    Thalamus --> Amygdala[Amygdala<br>(扁桃体)]
    Brain --> PFC[Prefrontal Cortex<br>(前頭前皮質)]
    Brain --> Hippo[Hippocampus<br>(海馬)]
    Brain --> Basal[Basal Ganglia<br>(大脳基底核)]
    
    Hypo -->|Hormones & Circadian| PFC
    Hypo -->|Motivation Cap| Basal
    Thalamus -->|Gain Modulation| Amygdala
    Amygdala -->|Raw Emotions| PFC
    PFC -->|Controlled Emotions| Hippo
    Hippo -->|Context/Memory| Brain
    
    Input -->|Physical Signal| Hypo
    Input -->|Physical Signal| Basal
    Basal -->|Motivation| Brain
    
    Brain --> Output[Mind State Response]
    
    Hippo <--> DB[(SQLite Database)]
```

### コアコンポーネント
1.  **Thalamus (視床)**:
    *   **機能**: 感覚入力のフィルタリングと順応（慣れ）。
    *   **ロジック**: 
        *   入力テキストの類似度を判定し、繰り返し入力を検出。
        *   繰り返し回数に応じてゲイン（強度係数）を減衰: `Gain = 1.0 / (1.0 + 0.5 × RepetitionCount)`
        *   スパム対策として、同じ褒め言葉の連打などに対する反応を抑制。

2.  **Hypothalamus (視床下部)**:
    *   **機能**: 生体の恒常性 (Homeostasis) 維持と概日リズム管理。
    *   **ホルモン**:
        *   `Cortisol` (ストレス): 不快刺激で上昇。
        *   `Oxytocin` (愛着): 快感刺激で上昇。
        *   `Melatonin` (睡眠): 夜間（22:00-6:00）に上昇。
        *   `Serotonin` (覚醒): 日中（6:00-22:00）に上昇。
    *   **概日リズム効果**:
        *   夜間: 意欲キャップ（50-100%）、感情感度上昇（1.0-1.2倍）。
        *   日中: ストレス減衰加速（1.0-2.0倍）。

3.  **Amygdala (扁桃体)**:
    *   **機能**: 直感的な感情生成。
    *   **ロジック**: 日本語形態素解析エンジン `Kagome` を使用し、入力テキストをトークン化。内部辞書と照合して原始的な感情値を算出します。

4.  **Prefrontal Cortex (PFC, 前頭前皮質)**:
    *   **機能**: 理性とホルモンバランスによる感情の制御。
    *   **ロジック**:
        *   **基本制御**: `Sanity` (理性値) に基づきネガティブ感情を抑制。
        *   **ストレス影響**: `Cortisol` が高いと理性が弱まり、ネガティブ感情が増幅されます（イライラ状態）。
        *   **愛着影響**: `Oxytocin` が高いと、怒り(`Anger`)が悲嘆・甘え(`Grief`)に変換されます。

5.  **Basal Ganglia (大脳基底核)**:
    *   **機能**: 意欲 (Motivation) と報酬予測誤差 (RPE) の管理。
    *   **ロジック**: 期待する報酬と実際の報酬の差分（RPE）に基づいて意欲を更新します。「飽き」や「期待外れ」による意欲減退をシミュレートします。

6.  **Hippocampus (海馬)**:
    *   **機能**: 記憶の形成、保持、検索。
    *   **ロジック**: 短期記憶 (STM) と長期記憶 (LTM) の2層構造。睡眠処理 (`Sleep`) によりSTMをLTMへ固定化し、SQLiteデータベースに永続化します。

## 3. インターフェース仕様 (API Specification)

すべてのAPIはRESTful JSON形式で提供されます。
ベースURL: `http://localhost:8081/api/v1`

### 3.1 感覚入力 (Sensory Input)
外部からの刺激（会話、物理接触など）を受け取ります。処理パイプラインにおいて、視床フィルタリング、視床下部のホルモン状態更新、感情生成・制御が順次行われます。

*   **Endpoint**: `POST /sensory`
*   **Request Body**:
    ```json
    {
      "type": "chat",          // "chat" (会話) または "physical" (物理刺激)
      "text": "こんにちは",      // 記憶用テキスト記述（必須）
      "signal_value": 0        // 物理刺激の強度 (-100 〜 +100)
    }
    ```

    *   **Type: "chat"**: Amygdalaによる感情解析が行われ、その感情値が二次的にホルモンと意欲に影響します。
    *   **Type: "physical"**: Amygdala解析はスキップされ、`signal_value` が直接的に快感/不快感として処理され、ホルモンと意欲を即座に更新します。

*   **Response Body**: `MindStateResponse`
    ```json
    {
      "current_reaction": [
        {"code": "J", "value": 80}
      ],
      "mood_stability": 0.85,    // 気分の安定度 (0.0-1.0)
      "personality_bias": [...], // 性格傾向
      "motivation": 0.75,        // 現在の意欲 (0.0-1.0) ※概日リズムキャップ適用済み
      "sanity": 0.90,            // 現在の理性 (0.0-1.0)
      "cortisol": 15.5,          // [DEBUG] 現在のストレスレベル (0-100)
      "oxytocin": 60.2,          // [DEBUG] 現在の愛着レベル (0-100)
      "predicted_reward": 55.0   // [DEBUG] 現在の報酬期待値 (0-100)
    }
    ```

### 3.2 睡眠 (Sleep)
記憶の整理・固定化を行います。
*   **Endpoint**: `POST /sleep`

### 3.3 状態取得 (Get State)
現在の脳の内部パラメータを取得します。
*   **Endpoint**: `GET /state`

## 4. データモデル (Data Models)

### 4.1 EmotionCode
| コード | 感情 | 説明 |
| :--- | :--- | :--- |
| `J` | Joy | 喜び、快感 |
| `A` | Anger | 怒り |
| `F` | Fear | 恐れ、痛み |
| `L` | Love | 愛、信頼 |
| `D` | Disgust | 嫌悪、不快 |
| `N` | Neutral | 中立 |
| `S` | Surprise | 驚き |
| `H` | Hope | 希望 |
| `G` | Grief | 悲嘆、悲しみ、甘え |

### 4.2 SensoryInput (Go Struct)
```go
type SignalType string
const (
    SignalChat     SignalType = "chat"
    SignalPhysical SignalType = "physical"
)

type SensoryInput struct {
    Type        SignalType `json:"type"`
    InputText   string     `json:"text"`
    SignalValue int        `json:"signal_value"`
}
```

## 5. 処理パイプライン (Processing Pipeline)

`ProcessInput` における処理順序:
1. **Decay**: ホルモンの時間経過による自然減衰
2. **Circadian Rhythm Update**: 現在時刻に基づく概日リズムホルモン（Melatonin/Serotonin）の更新
3. **Thalamus Filter**: 入力の繰り返し判定とゲイン計算
4. **Sensory Processing**:
   - Physical: 信号値を直接感情・ホルモン・意欲に変換（ゲイン適用）
   - Chat: Amygdala解析 → ゲイン・概日リズム感度適用 → ホルモン・意欲更新
5. **PFC Arbitrate**: ホルモン状態に基づく感情制御
6. **Memory Storage**: Hippocampusへの記憶保存
7. **Response Generation**: 概日リズムキャップを適用した意欲値を含むレスポンス生成

## 6. 永続化とインフラ (Persistence & Infrastructure)

### データベース
*   **Engine**: SQLite (`mind.db`)
*   **スキーマ**: Memoryテーブル (JSONデータ格納)
*   **パス設定**: 環境変数 `DB_PATH` で指定可能。

### Docker
*   **Build**: `docker build -t mind-os .`
*   **Run**: `docker-compose up -d`
*   **Volume**: `./data:/app`

## 7. 開発ガイドライン (Development Guidelines)
*   **言語**: Go 1.25.6
*   **設計原則**:
    *   **OOP**: 各脳部位をオブジェクトとして分離。
    *   **Single Responsibility**: 1ファイル200行以内を遵守。
    *   **Zero Hardcoding**: 設定値や定数は外部化または定数定義を使用。
