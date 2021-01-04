# kanji-numbers

大字表記の漢数字(壱万五千のような表記)とアラビア数字の表記をお互いに変換できるプログラムです。

## 利用法
kanji-number-back.an.r.appspot.com にデプロイしています。

アラビア数字から漢数字に変換する場合は

/v1/number2kanji/{number}

漢数字からアラビア数字に変換する場合は

/v1/kanji2number/{kanji}

のエンドポイントが利用できます。

### 例
kanji-number-back.an.r.appspot.com/v1/number2kanji/123

→ 壱百弐拾参

kanji-number-back.an.r.appspot.com/v1/kanji2number/壱百弐拾参

→ 123
