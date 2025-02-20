import asyncio
from embedder import embed
import numpy as np
from config import supabase, RAG_CHUNKS_TABLE_NAME
import time


#string, (text, embedding)
def search(query, tokenizer, model, chunks, k=5):
    embedded_query = embed(query, tokenizer, model)
    normalised_embedded_query = np.linalg.norm(embedded_query)

    chunk_scores = []
    for (chunk, embedded_chunk) in chunks:
        normalised_embedded_chunk = np.linalg.norm(embedded_chunk)

        score = np.dot(embedded_query, embedded_chunk) / (normalised_embedded_query * normalised_embedded_chunk)
        chunk_scores.append((score, chunk))
    
    chunk_scores.sort(reverse=True, key=lambda x: x[0])

    return chunk_scores[:k]

async def match_rag_chunks(embedded_query, k, similarity_threshold):
    start_time = time.time()
    # Create a coroutine by wrapping the synchronous call in asyncio.to_thread
    chunks = await asyncio.to_thread(
        lambda: supabase.rpc(
            'match_rag_chunks',
            {
                'query_embedding': embedded_query.tolist(),
                'match_count': k,
                'match_threshold': similarity_threshold
            }
        ).execute()
    )
    end_time = time.time()
    print(f"Time taken to match chunks: {end_time - start_time} seconds")
    return chunks

async def match_rag_contacts(embedded_query, k, similarity_threshold):
    start_time = time.time()
    # Create a coroutine by wrapping the synchronous call in asyncio.to_thread
    contacts = await asyncio.to_thread(
        lambda: supabase.rpc(
            'match_rag_contacts',
            {
                'query_embedding': embedded_query.tolist(),
                'match_count': k,
                'match_threshold': similarity_threshold
            }
        ).execute()
    )
    end_time = time.time()
    print(f"Time taken to match contacts: {end_time - start_time} seconds")
    return contacts

async def get_supabase_rag_chunks(query, tokenizer, model, k=5, similarity_threshold=0.5):
    embedded_query = embed(query, tokenizer, model)
    
    try:
        # Measure total time including embedding
        total_start_time = time.time()
        
        # Create tasks but don't await them yet
        chunks_task = match_rag_chunks(embedded_query, k, similarity_threshold)
        contacts_task = match_rag_contacts(embedded_query, k, similarity_threshold)
        
        # Run both tasks concurrently
        chunks, contacts = await asyncio.gather(chunks_task, contacts_task)
        
        total_time = time.time() - total_start_time
        print(f"Total time including parallel operations: {total_time} seconds")
        return {"chunks": chunks, "contacts": contacts}
    except Exception as e:
        print(f"Error in search_database: {e}")
        return []

def get_doc_ids_from_chunks(chunks):
    doc_set = set()
    for entry in chunks:
        doc_set.add(entry['doc_id'])
    return doc_set

def get_text_from_chunks(chunks, k=5):
    blocks = []
    for entry in chunks:
        doc_id = entry['doc_id']
        pos_in_doc = int(entry['pos_in_doc'])
        response = supabase.table(RAG_CHUNKS_TABLE_NAME).select("text").eq('doc_id', doc_id).gte('pos_in_doc', pos_in_doc - k).lte('pos_in_doc', pos_in_doc + k).execute()
        data = response.data
        blocks.append("".join([chunk['text'] for chunk in data]))
    return blocks