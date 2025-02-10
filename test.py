import asyncio
from ollama import AsyncClient
import dotenv
import os
from ollama import Client

dotenv.load_dotenv()
host=os.getenv("localLLM_IP")
# print(host)

async def chat(prompt):
    query = "lol"
    message=[
       {"role":"user","content":f"Instruction: Translate to Malaysia malay language. Do not add any other information.\nInput: {query}. \nOutput:"},
    ]
   #  "http://192.168.50.164:11434"
   #  print(host)
    client = AsyncClient(host=host)
    response = await client.chat(model='llama3.2:3b', 
                                 messages=message, 
                                 stream=True,
                                 options={'temperature':2,'num_predict':-1,'seed':123,'num_predict':512,})
    async for part in response:
      print(part['message']['content'], end='', flush=True)
      # await asyncio.sleep(0.1)
# Run the chat coroutine

alpaca_prompt = """Below is an instruction that describes a task, paired with an input that provides further context. Write a response that appropriately completes the request.

### Instruction:
Translate to Malaysia malay language. Do not add any other information.

### Input:
how are you buddy?

### Response:
"""

instruction = "Translate to Malaysia malay language. Do not add any other information."
input = "how are you buddy?"
response=""

formattedPrompt = alpaca_prompt.format(instruction,input,response)
# print(formattedPrompt)
asyncio.run(chat(formattedPrompt))
