from pdf_to_text import pdf_to_text2
from chunker import document_chunker2, extract_contacts_with_context
from embedder import embed
from supabase import create_client, Client
import os
from config import supabase, DOCUMENTS_BUCKET_NAME, DOCUMENTS_TABLE_NAME, RAG_TABLE_NAME, CONTACT_TABLE_NAME

url = "./pdfs/Student_Code_of_Conduct_2023_24.pdf"
url2 = "./pdfs/Computing-UG-Handbook-2425-v1b.pdf"
url3 = "./pdfs/Orientation-for-Visiting-(non-degree)-Students-(start-date---beginning-of-the-2024-25-academic-year).pdf"



def upload(url):
    file_name = os.path.basename(url)
    pdf_text = pdf_to_text2(url)
    print("1")

    #upload file to bucket storage
    try:
        with open(url, "rb") as file:
            response = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).upload(file_name, file)
    except Exception as e:
        print(f"An error uploading the file occurred: {str(e)}")
        return None
    public_url = supabase.storage.from_(DOCUMENTS_BUCKET_NAME).get_public_url(file_name)
    print("2")

    #add entry for file in documents table
    response = supabase.table(DOCUMENTS_TABLE_NAME).insert([
            {"name": file_name, "url": public_url}
        ]).execute()
    print(response)
    doc_id = response.data[0].get('doc_id')
    print(doc_id)
    print("3")

    #chunk document
    chunks = document_chunker2(pdf_text, "BAAI/bge-small-en-v1.5")
    # for chunk in chunks:
    #     print("-----------------------------------------------------------------")
    #     print(chunk)
    #     print('-----------------------------------------------------------------')

    print("4")
    index = 1
    for text in chunks:
        supabase.table(RAG_TABLE_NAME).insert([
            {"text": text, "doc_id": doc_id, "pos_in_doc": index, "embedding": embed(text).tolist()}
        ]).execute()
        print(f" - {index}")
        index += 1
    print("5")

    #extract contacts
    contacts = extract_contacts_with_context(pdf_text)
    print("6")
    contact_index = 1
    for email_context in contacts[0]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {"context": email_context, "doc_id": doc_id, "pos_in_contacts": contact_index, "contact_type": "email", 
             "embedding": embed(email_context).tolist()}
        ]).execute()
        print(f" - {contact_index}")
        contact_index += 1
    for phone_context in contacts[1]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {"context": phone_context, "doc_id": doc_id, "pos_in_contacts": contact_index, "contact_type": "phone number", 
             "embedding": embed(phone_context).tolist()}
        ]).execute()
        print(f" - {contact_index}")
        contact_index += 1
    for url_context in contacts[2]:
        supabase.table(CONTACT_TABLE_NAME).insert([
            {"context": url_context, "doc_id": doc_id, "pos_in_contacts": contact_index, "contact_type": "url", 
             "embedding": embed(url_context).tolist()}
        ]).execute()
        print(f" - {contact_index}")
        contact_index += 1
    print("7")

upload(url)
upload(url2)
upload(url3)
    