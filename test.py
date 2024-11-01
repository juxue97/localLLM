import ollama

stream = ollama.chat(
    model='llama3.1:8b',
    messages=[{'role': 'user', 'content': 'tell me story of yua mikami with 2000 words'}],
    stream=True,
)

for chunk in stream:
  print(chunk['message']['content'], end='', flush=True)