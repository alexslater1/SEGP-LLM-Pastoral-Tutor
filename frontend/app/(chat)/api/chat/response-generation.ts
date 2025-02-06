enum MessageRole {
    USER = "user",
    ASSISTANT = "assistant",
    SYSTEM = "system",
    DATA = "data",
}

const TOKEN_MILLISECOND_DELAY = 30

export class ResponseData {
    private data: Object[] = []
    private annotations: Object[] = []
    private message: string = ""
    private messageRole: string;
    private controller: ReadableStreamDefaultController<string>;
    private stream: ReadableStream;

    constructor(messageType: MessageRole) {
        let controller: ReadableStreamDefaultController<string>|null = null;
        this.stream = new ReadableStream({
            start(controllerArg) {
                controller = controllerArg;
            },
        });
        this.controller = controller as unknown as ReadableStreamDefaultController<string>;
        this.messageRole = messageType
    }

    safeEnqueue(data: string) {
        try {
          this.controller.enqueue(data);
        } catch (error) {
          // suppress errors when the stream has been closed
        }
      }

    addData(data: Object) {
        this.data.push(data)
        this.safeEnqueue("2:" + JSON.stringify([data]) + "\n")
    }

    addAnnotation(annotation: Object) {
        this.annotations.push(annotation)
        this.safeEnqueue("8:" + JSON.stringify([annotation]) + "\n")
    }

    // The UI will continue to show "Thinking..." until some message content is pushed to
    // the stream. This is a workaround to ensure the UI will update if an annotation was sent
    updateUIMessage() {
        this.safeEnqueue("0:\" \"\n")
    }

    async setMessage(message: string, delayed: boolean = true) {
        if (!message) {
            message = "Error: server provided no answer"
        }
        this.message = message;
        let tokens = this.message.split(" ")
        for (let token of tokens) {
            this.safeEnqueue("0:" + JSON.stringify(token + " ") + "\n")
            delayed && await new Promise(resolve => setTimeout(resolve, TOKEN_MILLISECOND_DELAY));
        }
    }

    closeStream() {
        this.controller.close()
    }

    getMessageRole(): string {
        return this.messageRole
    }

    getMessage(): string {
        return this.message
    }

    getAnnotations(): string { 
        return JSON.stringify(this.annotations)
    }

    toResponse() {
        // Not sure exactly what these two are needed for they appear to be required for indicating 
        // the response stream is complete for the AI SDK
        this.safeEnqueue('e:{"finishReason":"stop","usage":{"promptTokens":0,' +
            '"completionTokens":0},"isContinued":false}\n')
        this.safeEnqueue('d:{"finishReason":"stop","usage":{"promptTokens":0,' +
            '"completionTokens":0}}\n')

        return new Response(this.stream)
    }
}

async function executeAsyncInOrder(execute: (response: ResponseData) => Promise<void>, 
                                   onFinish: (response: ResponseData) => Promise<void>,
                                   response: ResponseData): Promise<void> {
    await execute(response)
    console.log("Finished executing")
    await onFinish(response)
    response.closeStream()
}

export function createResponse({execute, onFinish}: 
    {execute: (response: ResponseData) => Promise<void>,
     onFinish: (response: ResponseData) => Promise<void>}): Response {
    let responseData = new ResponseData(MessageRole.ASSISTANT)
    executeAsyncInOrder(execute, onFinish, responseData)
    return responseData.toResponse()
}