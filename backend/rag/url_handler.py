from chunker import document_chunker, extract_contacts_with_context
from embedder import embed
from config import supabase, RAG_SOURCES_TABLE_NAME, RAG_CHUNKS_TABLE_NAME, RAG_CONTACT_TABLE_NAME
from transformers import AutoTokenizer, AutoModel
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options

model_name = "BAAI/bge-small-en-v1.5"
model = AutoModel.from_pretrained(model_name)
tokenizer = AutoTokenizer.from_pretrained(model_name)

async def upload_url(url):
    chrome_options = Options()
    chrome_options.add_argument("--headless")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")
    
    driver = webdriver.Chrome(options=chrome_options)
    
    try:
        driver.get(url)
        text = driver.find_element(By.TAG_NAME, "body").text
    finally:
        driver.quit()

    # add entry for url in sources table
    response = supabase.table(RAG_SOURCES_TABLE_NAME).insert([
            {"url": url, "type": "WEBSITE"}
        ]).execute()
    url_id = response.data[0].get('id')

    # chunk document
    chunks = document_chunker(text, tokenizer)

    # process chunks
    index = 1
    for chunk in chunks:
        supabase.table(RAG_CHUNKS_TABLE_NAME).insert([{
            "text": chunk,
            "rag_source_id": url_id,
            "pos_in_source": index,
            "embedding": embed(chunk, tokenizer, model).tolist()
        }]).execute()
        index += 1

    # extract and process contacts
    contacts = extract_contacts_with_context(text)
    contact_index = 1
    
    # process each contact type
    for contact_list, contact_type in zip(contacts, ["email", "phone number", "url"]):
        for value, context in contact_list:
            supabase.table(RAG_CONTACT_TABLE_NAME).insert([{
                "context": context,
                "rag_source_id": url_id,
                "pos_in_source": contact_index,
                "contact_type": contact_type,
                "embedding": embed(context, tokenizer, model).tolist(),
                "contact": value
            }]).execute()
            contact_index += 1

def delete_url(url):
    try:
        response = supabase.table(RAG_SOURCES_TABLE_NAME).delete().eq("url", url).execute()
    except Exception as e:
        print(f"An error occurred deleting the url: {e}")
