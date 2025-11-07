from flask import Flask, request, jsonify
from llama_cpp import Llama
import json

app = Flask(__name__)

llm = Llama(model_path="models/phi3.gguf", n_ctx=2048)

@app.route("/interpret", methods=["POST"])
def interpret():
    user_input = request.json.get("prompt", "")
    prompt = f"""
    You are a file manager assistant. Convert the following instruction
    into a structured JSON command.

    Example:
    User: rename all .txt files to .md
    Output:
    {{"action": "rename", "file_pattern": "*.txt", "replace_extension": "md"}}

    User: {user_input}
    Output:
    """

    output = llm(prompt, max_tokens=200)
    text = output["choices"][0]["text"]

    try:
        command = json.loads(text)
    except Exception:
        command = {"error": "Could not parse model output", "raw_output": text}

    return jsonify(command)

if __name__ == "__main__":
    app.run(port=5001)
