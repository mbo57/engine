import os
import time
import google.generativeai as genai

# 環境変数からAPIキーを取得
API_KEY = os.environ["GENAI_API_KEY"]
genai.configure(api_key=API_KEY)
model = genai.GenerativeModel('gemini-1.5-flash')

# 会話の履歴がhistoryに追加される．
chat = model.start_chat(history=[])

chat.send_message(
    "10単語程度の英文を生成してください。毎回50文作っていください。同じ文が作成されないようにしてください。説明をなくして直接例文で始まるようにしてください。また、文の先頭の文番号は必要ありません。")

# geminiクエリを送り10単語程度の英文を生成させる
# sample2.txtに作られた文章のみを追加していく
for i in range(200):
    res = chat.send_message(
        "10単語程度の英文を生成してください。毎回50文作っていください。同じ文が作成されないようにしてください。説明をなくして直接例文で始まるようにしてください。また、文の先頭の文番号は必要ありません。")
    # res = chat.receive()
    print(res)
    with open("sample2.txt", "a") as f:
        f.write(res.text)
    time.sleep(5)
