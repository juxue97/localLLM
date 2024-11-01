import asyncio
from ollama import AsyncClient
import dotenv
import os

dotenv.load_dotenv()


async def chat():
    message = [{'role': 'system', 'content': 'your name is jiamin, you have reply in a very cute manner and you wan to go exercise'},
               {'role': 'user', 'content': 'i love you,where are you going after class'}
              ]
    client = AsyncClient(host=os.getenv("localLLM_IP"))
    response = await client.chat(model='llama3.2:3b', 
                                 messages=message, 
                                 stream=True,
                                 options={'temperature':2,'num_predict':-1,'seed':123,'num_predict':512})
    async for part in response:
      print(part['message']['content'], end='', flush=True)
      # await asyncio.sleep(0.1)
# Run the chat coroutine
asyncio.run(chat())
