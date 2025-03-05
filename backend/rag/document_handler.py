from pdf_to_text import file_to_text
from chunker import document_chunker, extract_contacts_with_context
from embedder import embed
import os
from config import supabase, DOCUMENTS_BUCKET_NAME, RAG_SOURCES_TABLE_NAME, RAG_CHUNKS_TABLE_NAME, RAG_CONTACT_TABLE_NAME
from transformers import AutoTokenizer, AutoModel

model_name = "BAAI/bge-small-en-v1.5"
model = AutoModel.from_pretrained(model_name)
tokenizer = AutoTokenizer.from_pretrained(model_name)

async def upload_doc(file):
    print("Uploading document:", file.filename)
    file_name = file.filename
    file_contents = await file.read()
    
    # upload file to bucket storage
    response = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).upload(file_name, file_contents)
    public_url = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).get_public_url(file_name)
    size = len(file_contents)
    print("Size:", size)
    
    #add entry for file in documents table
    response = supabase.table(RAG_SOURCES_TABLE_NAME).insert([
        {"url": public_url,
         "name": file_name,
         "type": "DOCUMENT",
         "document_size": size,
         "document_type": file_name.split(".")[-1].upper()}
    ]).execute()
    
    doc_id = response.data[0].get('id')
    
    #chunk document
    doc = await file_to_text(file_contents, file_name)
    chunks = document_chunker(doc, tokenizer)

    index = 1
    for text in chunks:
        supabase.table(RAG_CHUNKS_TABLE_NAME).insert([
            {"text": text, "rag_source_id": doc_id, "pos_in_source": index, "embedding": embed(text, tokenizer, model).tolist()}
        ]).execute()
        index += 1
    
    #extract contacts
    contacts = extract_contacts_with_context(doc)
    contact_index = 1
    for (email, email_context) in contacts[0]:
        supabase.table(RAG_CONTACT_TABLE_NAME).insert([
            {"context": email_context,
            "rag_source_id": doc_id, 
            "pos_in_source": contact_index,
            "contact_type": "email", 
            "embedding": embed(email_context, tokenizer, model).tolist(),
            "contact": email}
        ]).execute()
        contact_index += 1
    
    for (phone, phone_context) in contacts[1]:
        supabase.table(RAG_CONTACT_TABLE_NAME).insert([
            {"context": phone_context,
             "rag_source_id": doc_id,
             "pos_in_source": contact_index,
             "contact_type": "phone number", 
             "embedding": embed(phone_context, tokenizer, model).tolist(), 
             "contact": phone}
        ]).execute()
        contact_index += 1
        
    for (url, url_context) in contacts[2]:
        supabase.table(RAG_CONTACT_TABLE_NAME).insert([
            {
            "context": url_context, 
            "rag_source_id": doc_id, 
            "pos_in_source": contact_index, 
            "contact_type": "url", 
            "embedding": embed(url_context, tokenizer, model).tolist(),
            "contact": url
            }
        ]).execute()
        contact_index += 1

def fetch_docs():
    response = supabase.table(RAG_SOURCES_TABLE_NAME).select("*").execute()
    return response.data

def delete_doc(name):
    try:
        response = supabase.table(RAG_SOURCES_TABLE_NAME).delete().eq("name", name).execute()
        response2 = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).remove([name])

    except Exception as e:
        print(f"An error occurred deleting the document: {e}")

def download_doc(name):
    try:
        response = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).download(name)
        return response
    except Exception as e:
        print(f"An error occurred downloading the document: {e}")

def delete_id(id):
    try:
        response = supabase.table(RAG_SOURCES_TABLE_NAME).delete().eq("id", id).execute()

    except Exception as e:
        print(f"An error occurred deleting the document: {e}")
