from io import BytesIO
from fastapi import FastAPI, File, HTTPException, UploadFile, Form
from searcher import get_supabase_rag_chunks
from dotenv import load_dotenv
from document_handler import upload_doc, delete_doc, fetch_docs, download_doc
from transformers import AutoTokenizer, AutoModel
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse, JSONResponse, StreamingResponse
from url_handler import upload_url, delete_url
import validators

# Load environment variables from .env file
load_dotenv()

app = FastAPI()

# Update CORS middleware configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
    expose_headers=["Content-Disposition"]  # Important for file downloads
)

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
        retrieved = await get_supabase_rag_chunks(query, tokenizer, model, num_chunks, similarity_threshold)
       
        return {"chunks": retrieved["chunks"], "contacts": retrieved["contacts"]}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/rag-doc")
async def upload_rag_document(file: UploadFile = File(...)):
    print("Uploading document:", file.filename)
    try:
        # Embed and upload document chunks to database
        await upload_doc(file)
        return {"response": "Document uploaded"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/rag-doc/url")
async def upload_rag_document_url(url: str):
    # Clean url of extra whitespace
    url = url.strip()

    # Validate url
    if not url:
        raise HTTPException(
            status_code=400, 
            detail="Url parameter cannot be empty"
        )

    if not url.endswith((".pdf", ".docx", ".txt", ".pptx")):
        raise HTTPException(
            status_code=400, 
            detail="Document must be a .pdf, .docx, .txt or .pptx"
        )
    
    if not url.rsplit('.', 1)[0]:
        raise HTTPException(
            status_code=400, 
            detail="Document name must exist"
        )
    
    try:
        # Embed and upload document chunks to database
        upload_doc(url)
        return {"response": "Document uploaded"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/rag-doc")
async def fetch_rag_documents():
    print("Fetching documents")
    try:
        docs = fetch_docs()
        print("Documents:", docs)
        return {"documents": docs}


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

@app.get("/rag-doc/download")
async def download_rag_document(name: str):
    print("Downloading document")
    try:
        file = download_doc(name)
        
        # Create BytesIO object that can be streamed
        file_stream = BytesIO(file)
        
        return StreamingResponse(
            file_stream,
            media_type="application/octet-stream",
            headers={
                "Content-Disposition": f"attachment; filename={name}"
            }
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/rag-url")
async def upload_rag_url(url: str):
    # Clean url of extra whitespace
    url = url.strip()
    
    # Validate url
    if not url:
        raise HTTPException(
            status_code=400, 
            detail="Url parameter cannot be empty"
        )

    if not validators.url(url):
        raise HTTPException(
            status_code=400, 
            detail="Input must be a valid url"
        )
    
    try:
        # Embed and upload url chunks to database
        await upload_url(url)
        return {"response": "Url uploaded"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.delete("/rag-url")
async def delete_rag_url(url: str):
    # Clean url of extra whitespace
    url = url.strip()

    # Validate url
    if not url:
        raise HTTPException(
            status_code=400, 
            detail="Url parameter cannot be empty"
        )

    if not validators.url(url):
        raise HTTPException(
            status_code=400, 
            detail="Input must be a valid url"
        )
    
    try:    
        # Embed and upload url chunks to database
        delete_url(url)
        return {"response": "Url deleted"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# For local development
if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8000)
