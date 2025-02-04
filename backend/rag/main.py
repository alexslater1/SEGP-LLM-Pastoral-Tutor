from fastapi import FastAPI, File, HTTPException, UploadFile
from searcher import get_supabase_rag_chunks
from dotenv import load_dotenv
from document_uploader import upload_doc, delete_doc
from transformers import AutoTokenizer, AutoModel

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

        model_name = "BAAI/bge-small-en-v1.5"

        model = AutoModel.from_pretrained(model_name)
        tokenizer = AutoTokenizer.from_pretrained(model_name)
        
        # Search chunks based on query
        retrieved = get_supabase_rag_chunks(query, tokenizer, model, num_chunks, similarity_threshold)
       
        return {"chunks": retrieved["chunks"], "contacts": retrieved["contacts"]}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/rag-doc")
async def upload_rag_document(file: UploadFile = File(...)):
    try:
        # Embed and upload document chunks to database
        await upload_doc(file)
        return {"response": "Document uploaded"}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.delete("/rag-doc")
async def delete_rag_document(name: str):
    name = name.strip()
    # Validate url
    if not name:
        raise HTTPException(
            status_code=400, 
            detail="Url parameter cannot be empty"
        )

    if not name.endswith((".pdf", ".docx", ".txt", ".pptx")):
        raise HTTPException(
            status_code=400, 
            detail="Document must be a .pdf, .docx, .txt or .pptx"
        )
    
    if not name.rsplit('.', 1)[0]:
        raise HTTPException(
            status_code=400, 
            detail="Document name must exist"
        )
    
    try:
        delete_doc(name)
        return {"response": "Document deleted"}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# For local development
if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8000)
