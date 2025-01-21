from pdf_to_text import pdf_to_text2
from chunker import document_chunker2, extract_contacts_with_context
from embedder import embed
from supabase import create_client, Client
import os
from config import supabase, DOCUMENTS_BUCKET_NAME, DOCUMENTS_TABLE_NAME, RAG_TABLE_NAME, CONTACT_TABLE_NAME

def upload_doc(url):
    file_name = os.path.basename(url)

    #upload file to bucket storage
    try:
        with open(url, "rb") as file:
            response = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).upload(file_name, file)
    except Exception as e:
        print(f"An error uploading the file occurred: {str(e)}")
        return None
    public_url = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).get_public_url(file_name)

    #add entry for file in documents table
    response = supabase.table(DOCUMENTS_TABLE_NAME).insert([
            {"name": file_name, "url": public_url}
        ]).execute()
    doc_id = response.data[0].get('doc_id')

    #chunk document
    chunks = document_chunker2(pdf_to_text2(url), "BAAI/bge-small-en-v1.5")

    index = 1
    for text in chunks:
        supabase.table(RAG_TABLE_NAME).insert([
            {"text": text, "doc_id": doc_id, "pos_in_doc": index, "embedding": embed(text).tolist()}
        ]).execute()
        index += 1

    #extract contacts
    contacts = extract_contacts_with_context(pdf_text)
    contact_index = 1
    for (email, email_context) in contacts[0]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {"context": email_context, "doc_id": doc_id, "pos_in_contacts": contact_index, "contact_type": "email", 
             "embedding": embed(email_context).tolist(), "contact": email}
        ]).execute()
        print(f" - {contact_index}")
        contact_index += 1
    for (phone, phone_context) in contacts[1]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {"context": phone_context, "doc_id": doc_id, "pos_in_contacts": contact_index, "contact_type": "phone number", 
             "embedding": embed(phone_context).tolist(), "contact": phone}
        ]).execute()
        print(f" - {contact_index}")
        contact_index += 1
    for (url, url_context) in contacts[2]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {"context": url_context, "doc_id": doc_id, "pos_in_contacts": contact_index, "contact_type": "url", 
             "embedding": embed(url_context).tolist(), "contact": url}
        ]).execute()
        print(f" - {contact_index}")
        contact_index += 1


def delete_doc(name):
    try:
        response = supabase.table(DOCUMENTS_TABLE_NAME).delete().eq("name", name).execute()
        response2 = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).remove([name])

        if response.status_code == 200 and response2.status_code == 200:
            print(f"Successfully deleted the file: '{name}'")
        else:
            print(f"Error: {response.json()}")  # Show error details if any
    except Exception as e:
        print(f"An error occurred: {e}")
