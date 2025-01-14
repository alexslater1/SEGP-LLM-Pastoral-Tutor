from pdf_to_text import pdf_to_text2
from chunker import document_chunker2
from embedder import embed
from supabase import create_client, Client
import os

SUPABASE_URL = "https://pwgnrkdcbldzgtzecfjj.supabase.co"
SUPABASE_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InB3Z25ya2RjYmxkemd0emVjZmpqIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzY4NjY1NTcsImV4cCI6MjA1MjQ0MjU1N30.wf2TYNcEo8A_RzJr0ZS4y6fPOmglkTgCI5dEOuDOkiw"
supabase = create_client(SUPABASE_URL, SUPABASE_KEY)
DOCUMENTS_BUCKET_NAME = "document-files"
DOCUMENTS_TABLE_NAME = "documents"
RAG_TABLE_NAME = "rag"

url = "./pdfs/Student_Code_of_Conduct_2023_24.pdf"


def upload(url):
    file_name = os.path.basename(url)
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
    chunks = document_chunker2(pdf_to_text2(url), "BAAI/bge-small-en-v1.5")
    for chunk in chunks:
        print("-----------------------------------------------------------------")
        print(chunk)
        print('-----------------------------------------------------------------')

    print("4")
    index = 1
    for text in chunks:
        supabase.table(RAG_TABLE_NAME).insert([
            {"text": text, "doc_id": doc_id, "pos_in_doc": index, "embedding": embed(text).tolist()}
        ]).execute()
        print(f" - {index}")
        index += 1
    print("5")


upload(url)
    