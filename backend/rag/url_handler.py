from chunker import document_chunker, extract_contacts_with_context
from embedder import embed
from config import supabase, RAG_SOURCES_TABLE_NAME, RAG_CHUNKS_TABLE_NAME, RAG_CONTACT_TABLE_NAME
from transformers import AutoTokenizer, AutoModel
from selenium import webdriver
from selenium.webdriver.common.by import By

def upload_url(url):

    model_name = "BAAI/bge-small-en-v1.5"

    model = AutoModel.from_pretrained(model_name)
    tokenizer = AutoTokenizer.from_pretrained(model_name)

    #add entry for url in sources table
    response = supabase.table(RAG_SOURCES_TABLE_NAME).insert([
            {"url": url, "type": "URL"}
        ]).execute()
    url_id = response.data[0].get('id')

    driver = webdriver.Chrome()
    driver.get(url)
    text = driver.find_element(By.TAG_NAME, "body").text
    driver.quit()

    #chunk document
    chunks = document_chunker(text, tokenizer)

    index = 1
    for t in chunks:
        supabase.table(RAG_CHUNKS_TABLE_NAME).insert([
            {"text": t, "rag_source_id": url_id, "pos_in_source": index, "embedding": embed(t, tokenizer, model).tolist()}
        ]).execute()
        index += 1
    print("4")
    #extract contacts
    contacts = extract_contacts_with_context(text)
    contact_index = 1
    for (email, email_context) in contacts[0]:
        supabase.table(RAG_CONTACT_TABLE_NAME).insert([
            {"context": email_context,
            "rag_source_id": url_id, 
            "pos_in_source": contact_index,
            "contact_type": "email", 
            "embedding": embed(email_context, tokenizer, model).tolist(),
            "contact": email}
        ]).execute()

        print(f" - {contact_index}")
        contact_index += 1
    
    for (phone, phone_context) in contacts[1]:
        supabase.table(RAG_CONTACT_TABLE_NAME).insert([
            {"context": phone_context,
             "rag_source_id": url_id,
             "pos_in_source": contact_index,
             "contact_type": "phone number", 
             "embedding": embed(phone_context, tokenizer, model).tolist(), 
             "contact": phone}
        ]).execute()

        print(f" - {contact_index}")
        contact_index += 1
        
    for (url, url_context) in contacts[2]:
        supabase.table(RAG_CONTACT_TABLE_NAME).insert([
            {
            "context": url_context, 
            "rag_source_id": url_id, 
            "pos_in_source": contact_index, 
            "contact_type": "url", 
            "embedding": embed(url_context, tokenizer, model).tolist(),
            "contact": url
            }
        ]).execute()

        print(f" - {contact_index}")
        contact_index += 1
    print("5")

def delete_url(url):
    return ""