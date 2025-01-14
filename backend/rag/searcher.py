from embedder import embed
import numpy as np


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

