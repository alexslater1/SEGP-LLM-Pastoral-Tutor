from chunker import document_chunker, extract_contacts_with_context
from embedder import embed
from config import supabase, RAG_SOURCES_TABLE_NAME, RAG_CHUNKS_TABLE_NAME, RAG_CONTACT_TABLE_NAME
from transformers import AutoTokenizer, AutoModel
from selenium import webdriver
from selenium.webdriver.common.by import By

def upload_url(url):
    return ""

def delete_url(url):
    return ""