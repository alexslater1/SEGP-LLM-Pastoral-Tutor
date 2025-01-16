from embedder import embed
import numpy as np
from config import supabase, DOCUMENTS_TABLE_NAME, RAG_TABLE_NAME


#string, (text, embedding)
def search(query, chunks, k=5):
    embedded_query = embed(query)
    normalised_embedded_query = np.linalg.norm(embedded_query)

    chunk_scores = []
    for (chunk, embedded_chunk) in chunks:
        normalised_embedded_chunk = np.linalg.norm(embedded_chunk)

        score = np.dot(embedded_query, embedded_chunk) / (normalised_embedded_query * normalised_embedded_chunk)
        chunk_scores.append((score, chunk))
    
    chunk_scores.sort(reverse=True, key=lambda x: x[0])

    return chunk_scores[:k]

def search_database(query, k=1):
    embedded_query = embed(query)
    normalised_embedded_query = np.linalg.norm(embedded_query)

    response = supabase.table(RAG_TABLE_NAME).select("*").execute()
    data = response.data

    chunk_scores = []
    for entry in data:
        embedded_chunk = np.array(eval(entry['embedding']))
        normalised_embedded_chunk = np.linalg.norm(embedded_chunk)

        score = np.dot(embedded_query, embedded_chunk) / (normalised_embedded_query * normalised_embedded_chunk)
        chunk_scores.append((score, entry))
    
    chunk_scores.sort(reverse=True, key=lambda x: x[0])

    return [entry for (score, entry) in chunk_scores][:k]