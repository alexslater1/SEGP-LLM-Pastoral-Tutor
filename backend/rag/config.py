from supabase import create_client, Client

SUPABASE_URL = "https://pwgnrkdcbldzgtzecfjj.supabase.co"
SUPABASE_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InB3Z25ya2RjYmxkemd0emVjZmpqIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzY4NjY1NTcsImV4cCI6MjA1MjQ0MjU1N30.wf2TYNcEo8A_RzJr0ZS4y6fPOmglkTgCI5dEOuDOkiw"
supabase = create_client(SUPABASE_URL, SUPABASE_KEY)
DOCUMENTS_BUCKET_NAME = "document-files"
DOCUMENTS_TABLE_NAME = "rag_sources"
RAG_TABLE_NAME = "rag_chunks"
CONTACT_TABLE_NAME = "rag_contacts"
