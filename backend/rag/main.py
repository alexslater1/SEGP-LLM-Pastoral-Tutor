from chunker import document_chunker, print_chunks
from embedder import embed
from searcher import search
from pdf_to_text import pdf_to_text
from openai import OpenAI

#TODO: get document and convert to text files
# pdf_path = "./pdfs/Student_Code_of_Conduct_2023_24.pdf"
# txt_path = "./documents/Student_Code_of_Conduct_2023_24.txt"
# pdf_to_text(pdf_path, txt_path)

#TODO: split text into chunks
chunks = document_chunker("./documents/", "BAAI/bge-small-en-v1.5")

text_chunks = []
for outer_key, outer_value in chunks.items():
    for inner_key, inner_value in outer_value.items():
        if 'text' in inner_value:
            # print("--------------------")
            # print(inner_value['text'])
            text_chunks.append(inner_value['text'])


#embed chunks as vectors
embedded_chunks = []
for chunk in text_chunks:
   embedded_chunks.append((chunk, embed(chunk)))

#TODO: store vector with document

#TODO: search chunks
query = "What is stalking?"
return_chunks = search(query, embedded_chunks)

retrieved_chunks = []
for (_, chunk) in return_chunks:
    retrieved_chunks.append(chunk)

base_prompt = """You are an AI assistant for RAG. Your task is to understand the user question, and provide an answer using the provided contexts.

Your answers are correct, high-quality, and written by an domain expert. If the provided context does not contain the answer, simply state, "The provided context does not have the answer."

User question: {user_query}

Contexts:
{chunks_information}
"""
prompt = base_prompt.format(user_query=query, chunks_information="\n".join([chunk for chunk in retrieved_chunks]))

print(prompt)

client = OpenAI(
    api_key="idk"
)
response = client.chat.completions.create(
    model="gpt-4o-mini",
    temperature=0,
    messages=[
        {"role": "system", "content": prompt},
    ],
)

print(response.choices[0].message.content)