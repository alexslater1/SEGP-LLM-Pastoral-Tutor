import {
  CoreUserMessage,
  CoreMessage,
  type Message,
} from 'ai';

const API_FETCH_TIMEOUT_SECONDS: number = 120;

type queryID = string

export enum Status {
    PENDING = "pending",
    COMPLETED = "completed",
    FAILED = "failed",
}

export type StatusResponse = {
    type: Status,
    error?: string,
    answer?: string,
    current_action?: string,
}

function createQueryText(previousMessages: Message[], userMessage: CoreUserMessage): string {
    let query = "Previous messages (oldest to most recent): " + previousMessages.slice(0, -1).map(message => {
            if (message.role == "user") {
                return "User: " + message.content;
            } else {
                return "Assistant: " + message.content;
            }
        }).join("\n");
    query += "\nMost recent user message to respond and answer to now: " + userMessage.content;
    return query
}

export async function makeV1Completion(previousMessages: Message[], userMessage: CoreUserMessage): Promise<string> {
    type CompletionResponse = {
        response: string;
        reason: string;
    }
    const completionEndpoint = "/completion"
    try {
        let query = createQueryText(previousMessages, userMessage);
        if (!process.env.BACKEND_URL) throw new Error()
        const response = await fetch(process.env.BACKEND_URL + completionEndpoint, {
            method: 'POST',
            signal: AbortSignal.timeout(1000 * API_FETCH_TIMEOUT_SECONDS),
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                query: query
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return (await response.json() as CompletionResponse).response;
    } catch (error) {
        console.error('Error:', error);
        return "Error fetching from backend";
    }
}

export async function makeV2InitialQuery(previousMessages: Message[], userMessage: CoreUserMessage): Promise<queryID> {
    try {
        const completionEndpoint = "/completion/v2"
        let query = createQueryText(previousMessages, userMessage);
        const response = await fetch(process.env.BACKEND_URL + completionEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                query: query
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return (await response.json()).request_id;
    } catch (error) {
        console.error('Error:', error);
        return "Error fetching from backend";
    }
}

export async function makeV2StatusQuery(queryID: queryID): Promise<StatusResponse> {
    try {
        const completionEndpoint = "/completion/v2/status/" + queryID
        const response = await fetch(process.env.BACKEND_URL + completionEndpoint, {
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return (await response.json());
    } catch (error) {
        console.error('Error:', error);
        return {
            type: Status.FAILED,
            error: "Error fetching status from backend",
        };
    }
}
