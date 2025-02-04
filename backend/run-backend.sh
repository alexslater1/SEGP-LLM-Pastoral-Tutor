#! /bin/sh

# Installing Go
wget -q -O - https://git.io/vQhTU | sh -s -- --version 1.23.5
source /app/.bashrc

# Installing Python Dependencies
cd rag
python -m venv venv
source venv/bin/activate
pip install -r ../requirements.txt

# Starting services
python main.py > /segp/rag-log.txt 2>&1 &
cd ..
./run-agents-server.sh > /segp/agent-log.txt
