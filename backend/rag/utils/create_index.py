import psycopg2

# Define connection parameters (replace with your details)
db_params = {
    'dbname': 'postgres',
    'user': 'postgres.pwgnrkdcbldzgtzecfjj',
    'password': 'SECRET!!!!!',
    'host': 'aws-0-eu-west-2.pooler.supabase.com',
    'port': 5432,
}

# Establish a connection to the database
connection = psycopg2.connect(**db_params)
connection.autocommit = True

# Create a cursor object to interact with the database
cursor = connection.cursor()

# Execute the SQL queries separately
cursor.execute("SET hnsw.ef_search = 75")
cursor.execute("SET statement_timeout = '2000min'")
cursor.execute("""
    CREATE INDEX ON rag_contacts 
    USING hnsw (embedding vector_cosine_ops) 
    WITH (m = 16, ef_construction = 100);
""")

# Close the cursor and connection
cursor.close()
connection.close()

print("HNSW index created successfully!")
