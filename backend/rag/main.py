from fastapi import FastAPI, HTTPException
from chunker import document_chunker, print_chunks
from embedder import embed
from searcher import search, search_database, get_text_from_chunks
from pdf_to_text import pdf_to_text
from openai import OpenAI
from dotenv import load_dotenv
import os

# Load environment variables from .env file
load_dotenv()

# #TODO: get document and convert to text files
# # pdf_path = "./pdfs/Student_Code_of_Conduct_2023_24.pdf"
# # txt_path = "./documents/Student_Code_of_Conduct_2023_24.txt"
# # pdf_to_text(pdf_path, txt_path)

# #TODO: split text into chunks
# chunks = document_chunker("./documents/", "BAAI/bge-small-en-v1.5")

# text_chunks = []
# for outer_key, outer_value in chunks.items():
#     for inner_key, inner_value in outer_value.items():
#         if 'text' in inner_value:
#             # print("--------------------")
#             # print(inner_value['text'])
#             text_chunks.append(inner_value['text'])


# #embed chunks as vectors
# embedded_chunks = []
# for chunk in text_chunks:
#    embedded_chunks.append((chunk, embed(chunk)))

# #TODO: store vector with document

#TODO: search chunks
query = "What is stalking?"
retrieved_chunks = search_database(query)

base_prompt = """You are an AI assistant for RAG. Your task is to understand the user question, and provide an answer using the provided contexts.

Your answers are correct, high-quality, and written by an domain expert. If the provided context does not contain the answer, simply state, "The provided context does not have the answer."

User question: {user_query}

Contexts:
{chunks_information}
"""
prompt = base_prompt.format(user_query=query, chunks_information="\n".join(get_text_from_chunks(retrieved_chunks, 1)))

print(prompt)

app = FastAPI()
client = OpenAI(api_key=os.getenv('OPENAI_API_KEY'))

@app.get("/rag")
async def get_rag_response(query: str):
    # Validate query
    if not query or query.isspace():
        raise HTTPException(
            status_code=400, 
            detail="Query parameter cannot be empty or only whitespace"
        )
    
    if len(query.strip()) < 3:
        raise HTTPException(
            status_code=400, 
            detail="Query must be at least 3 characters long"
        )
    
    try:
        # Clean query of extra whitespace
        query = query.strip()
        
        # Search chunks based on query
        retrieved_chunks = search_database(query)
        
        # Format prompt with query and retrieved chunks
        prompt = base_prompt.format(
            user_query=query, 
            chunks_information="\n".join(get_text_from_chunks(retrieved_chunks, 1))
        )
        
        # Get response from OpenAI
        response = client.chat.completions.create(
            model="gpt-4o-mini",
            temperature=0,
            messages=[
                {"role": "system", "content": prompt},
            ],
        )
        
        return {"response": response.choices[0].message.content}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# For local development
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)