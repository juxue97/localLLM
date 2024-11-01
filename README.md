# localLLM Project

Hello World!!

# To run llm in local:

1. First you must run the ollama docker image with gpu enabled.

"""docker run -d --gpus=all -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama"""

2. Second, pull your desired model

for ex;
"""docker exec -it ollama ollama run llama3.2:3b"""

3. In addition, you can check your downloaded model list.

"""docker exec -it ollama ollama list"""

# Now, you are a step away to go! Next thing is to install ollama library

"""pip install ollama"""

# Lastly, refer to its official documentation to run the code

here are the link for it:
"""https://github.com/ollama/ollama-python"""
