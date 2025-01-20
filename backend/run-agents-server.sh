# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed on your system."
    echo "Please install Go from https://golang.org/doc/install"
    exit 1
fi

# If Go is installed, run the program
cd ./agents/agents-main

echo "To test the server send a POST request to http://localhost:8080/completion with body"
echo "{
    \"query\": String
}"

echo ""
echo "FOR EXAMPLE (in a seperate terminal):"
echo "curl -X POST -H \"Content-Type: application/json\" -d '{\"query\": \"What is the weather in San Francisco?\"}' http://localhost:8080/completion"
echo ""

go run main.go

