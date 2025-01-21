from fastapi import FastAPI, HTTPException
from searcher import get_supabase_rag_chunks
from openai import OpenAI
from dotenv import load_dotenv
from document_uploader import upload_doc, delete_doc
import os

# Load environment variables from .env file
load_dotenv()

app = FastAPI()

@app.get("/rag")
async def get_rag_response(query: str, num_chunks: int = 5, similarity_threshold: float = 0.5):
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
    
    # Validate num_chunks
    if num_chunks < 1:
        raise HTTPException(
            status_code=400,
            detail="num_chunks must be at least 1"
        )
    
    # Validate similarity_threshold
    if not 0 <= similarity_threshold <= 1:
        raise HTTPException(
            status_code=400,
            detail="similarity_threshold must be between 0 and 1"
        )
    
    try:
        # Clean query of extra whitespace
        query = query.strip()
        
        # Search chunks based on query
        retrieved = get_supabase_rag_chunks(query, num_chunks, similarity_threshold)
       
        return {"chunks": retrieved["chunks"], "contacts": retrieved["contacts"]}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/rag")
async def upload_rag_document(url: str):
    # Validate url
    if not url or url.isspace():
        raise HTTPException(
            status_code=400, 
            detail="Url parameter cannot be empty or only whitespace"
        )

    if url.strip()[-4:] != ".pdf":
        raise HTTPException(
            status_code=400, 
            detail="Document must be a pdf"
        )
    
    if not url.strip()[:-4]:
        raise HTTPException(
            status_code=400, 
            detail="Document name must exist"
        )
    
    try:
        # Clean url of extra whitespace
        url = url.strip()
        
        # Embed and upload document chunks to database
        upload_doc(url)
       
        return {"response": "Document uploaded"}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.delete("/rag")
async def delete_rag_document(name: str):
    # Validate url
    if not name or name.isspace():
        raise HTTPException(
            status_code=400, 
            detail="Url parameter cannot be empty or only whitespace"
        )

    if name.strip()[-4:] != ".pdf":
        raise HTTPException(
            status_code=400, 
            detail="Document must be a pdf"
        )
    
    if not name.strip()[:-4]:
        raise HTTPException(
            status_code=400, 
            detail="Document name must exist"
        )
    
    try:
        # Clean url of extra whitespace
        name = name.strip()
        
        # Embed and upload document chunks to database
        delete_doc(name)
       
        return {"response": "Document deleted"}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# For local development
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)