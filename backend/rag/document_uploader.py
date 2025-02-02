from pdf_to_text import file_to_text
from chunker import document_chunker, extract_contacts_with_context
from embedder import embed
import os
from pdf_to_text import pdf_to_text
from config import supabase, DOCUMENTS_BUCKET_NAME, DOCUMENTS_TABLE_NAME, RAG_TABLE_NAME, CONTACT_TABLE_NAME
from transformers import AutoTokenizer, AutoModel

def upload_doc(url):

    model_name = "BAAI/bge-small-en-v1.5"

    model = AutoModel.from_pretrained(model_name)
    tokenizer = AutoTokenizer.from_pretrained(model_name)

    file_name = os.path.basename(url)
    print("1")
    mkdwn_pdf = pdf_to_text(url)

    # upload file to bucket storage
    
    with open(url, "rb") as file:
        response = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).upload(file_name, file)

    public_url = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).get_public_url(file_name)
    print("2")
    #add entry for file in documents table
    response = supabase.table(DOCUMENTS_TABLE_NAME).insert([
            {"url": public_url, "name": file_name, "type": "DOCUMENT"}
        ]).execute()
    source_id = response.data[0].get('id')

    text = file_to_text(url)
    print("3")
    #chunk document
    chunks = document_chunker(mkdwn_pdf, tokenizer)

    index = 1
    for text in chunks:
        supabase.table(RAG_TABLE_NAME).insert([
            {"text": text, "rag_source_id": source_id, "pos_in_source": index, "embedding": embed(text, tokenizer, model).tolist()}
        ]).execute()
        index += 1
    print("4")
    #extract contacts
    contacts = extract_contacts_with_context(mkdwn_pdf)
    contact_index = 1
    for (email, email_context) in contacts[0]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {"context": email_context, "rag_source_id": source_id, "pos_in_source": contact_index, "contact_type": "email", 
             "embedding": embed(email_context, tokenizer, model).tolist(), "contact": email}
        ]).execute()

        print(f" - {contact_index}")
        contact_index += 1
    
    for (phone, phone_context) in contacts[1]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {"context": phone_context, "rag_source_id": source_id, "pos_in_source": contact_index, "contact_type": "phone number", 
             "embedding": embed(phone_context, tokenizer, model).tolist(), "contact": phone}
        ]).execute()

        print(f" - {contact_index}")
        contact_index += 1
        
    for (url, url_context) in contacts[2]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {
            "context": url_context, 
            "rag_source_id": source_id, 
            "pos_in_source": contact_index, 
            "contact_type": "website", 
            "embedding": embed(url_context, tokenizer, model).tolist(),
            "contact": url
            }
        ]).execute()

        print(f" - {contact_index}")
        contact_index += 1
    print("5")

def delete_doc(name):
    try:
        response = supabase.table(DOCUMENTS_TABLE_NAME).delete().eq("name", name).execute()
        response2 = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).remove([name])

    except Exception as e:
        print(f"An error occurred deleting the document: {e}")